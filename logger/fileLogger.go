/*
 *  Copyright Project - Stockbinator, Author - quoeamaster, (C) 2019
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
package logger

import (
	"Stockbinator/common"
	"Stockbinator/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type StructFileLogger struct {
	// actual file handler....
	fileHandle *os.File
	prefix string
	isPrefixSet bool

	currentLogFilename string
	currentLogFilenameFullpath string
}

// ctor
func NewStructFileLogger() (logger *StructFileLogger) {
	logger = new(StructFileLogger)
	return
}

// setup for the logger, requires a Struct of logger-config and can provide an optional params map object
func (l *StructFileLogger) SetupLogger(config *StructLoggerConfig, params... map[string]interface{}) (err error) {
	if config != nil {
		// a. check if filepath exists
		exists, err2 := util.IsFolderExists(config.Filepath)
		if err2 != nil {
			err = err2
			return
		}
		if !exists {
			err = errors.New(fmt.Sprintf("filepath %v does NOT exists", config.Filepath))
			return
		}
		// b. get filename prepared
		finalFilename, err2  := l.prepareFilename(config)
		if err2 != nil {
			err = err2
			return
		}
		l.currentLogFilename = finalFilename
		l.currentLogFilenameFullpath = fmt.Sprintf("%v%v", config.Filepath, finalFilename)
		// c. create the writer, file
		l.fileHandle, err = os.OpenFile(l.currentLogFilenameFullpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

	} else {
		err = errors.New("a structLoggerConfig is required to setup the logger")
		return
	}
	return
}

func (l *StructFileLogger) prepareFilename(config *StructLoggerConfig) (name string, err error) {
	// assume filepath all good
	// check filename
	if util.IsEmptyString(config.Filename) {
		err = errors.New("filename for logging should NOT be empty")
		return
	}
	// if is time-based but filepattern is not provided... use the default pattern
	finalPattern := loggerTimebasedPattern
	if config.IsFilenameTimebased && !util.IsEmptyString(config.FilenameTimebasedPattern) {
		finalPattern = config.FilenameTimebasedPattern
	}
	// substitute the pattern with real date parts
	finalPattern = l.prepareTimebasedPattern(finalPattern)
	// concat the filename with the pattern
	name = fmt.Sprintf("%v-%v", config.Filename, finalPattern)
	return
}

func (l *StructFileLogger) prepareTimebasedPattern(pattern string) (v string) {
	now := time.Now()

	v = strings.Replace(pattern, "yyyy", fmt.Sprintf("%v", now.Year()), -1)
	if strings.Index(v, "yy") != -1 {
		v = strings.Replace(v, "yy", fmt.Sprintf("%v", now.Year())[2:], -1)
	}
	mm := int(now.Month())
	if mm < 10 {
		v = strings.Replace(v, "mm", fmt.Sprintf("0%v", mm), -1)
	} else {
		v = strings.Replace(v, "mm", fmt.Sprintf("%v", mm), -1)
	}
	dd := now.Day()
	if dd < 10 {
		v = strings.Replace(v, "dd", fmt.Sprintf("0%v", dd), -1)
	} else {
		v = strings.Replace(v, "dd", fmt.Sprintf("%v", dd), -1)
	}
	return
}

// updates the Writer / output (e.g. set it to a valid FileWriter instance) (non implemented though)
func (l *StructFileLogger) SetOutput(writer io.Writer) (err error) {
	// non implemented
	return
}
// implementations might update the given prefix with the actual values to substitute
func (l *StructFileLogger) SetPrefix(prefixPattern string) {
	if strings.Compare(strings.Trim(prefixPattern, ""), "") != 0 {
		l.prefix = prefixPattern
		l.isPrefixSet = true
	}
}

// print out values provided
func (l *StructFileLogger) Print(v... interface{}) {
	if v != nil {
		sList := make([]string, 0)

		for _, s := range v {
			sList = append(sList, s.(string))
		}
		sLine := strings.Join(sList, " ")
		// append prefix???
		sLine = l.prependPrefixToLine(sLine)
		_, err := l.fileHandle.WriteString(sLine)
		if err != nil {
			fmt.Println(err)
		}
	}
}
// print out values provided by substituting them into the format
func (l *StructFileLogger) Printf(format string, v... interface{}) {
	if v != nil {
		sLine := format
		for _, s := range v {
			sVal := s.(string)
			sLine = strings.Replace(sLine, "%v", sVal, 1)
		}
		// append prefix???
		sLine = l.prependPrefixToLine(sLine)
		_, err := l.fileHandle.WriteString(sLine)
		if err != nil {
			fmt.Println(err)
		}
	}
}
// print out values provided in a line format
func (l *StructFileLogger) Println(v... interface{}) {
	if v != nil {
		sList := make([]string, 0)

		for _, s := range v {
			sList = append(sList, s.(string))
		}
		sLine := strings.Join(sList, " ")
		// append prefix???
		sLine = l.prependPrefixToLine(sLine)
		_, err := l.fileHandle.WriteString(fmt.Sprintf("%v\n", sLine))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (l *StructFileLogger) prependPrefixToLine(line string) (finalLine string) {
	finalLine = line
	if l.isPrefixSet {
		finalLine = fmt.Sprintf("[%v] %v", l.prefix, line)
	}
	return
}

// close the logger, release resources involved
func (l *StructFileLogger) Close() (err error) {
	if l.fileHandle != nil {
		err = l.fileHandle.Close()
	}
	return
}

func (l *StructFileLogger) Info() (info string) {
	iMap := make(map[string]interface{})

	iMap[common.LoggerFileInfoKeyFilename] = l.currentLogFilename
	iMap[common.LoggerFileInfoKeyFilepath] = l.currentLogFilenameFullpath

	bContents, err := json.Marshal(iMap)
	if err != nil {
		panic(err)
	}
	info = string(bContents)
	return
}