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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
	"io"
	"strings"
	"time"
)

type StructConsoleLogger struct {
	prefix string
	isPrefixSet bool
}

// ctor
func NewStructConsoleLogger() (pLogger *StructConsoleLogger) {
	pLogger = new(StructConsoleLogger)
	_ = pLogger.SetupLogger(nil)
	return
}

// setup for the logger, requires a Struct of logger-config and can provide an optional params map object
func (s *StructConsoleLogger) SetupLogger(config *StructLoggerConfig, params... map[string]interface{}) (err error) {
	// additional configuration if any (for now... no)
	return
}
// updates the Writer / output (e.g. set it to a valid FileWriter instance)
func (s *StructConsoleLogger) SetOutput(writer io.Writer) (err error) {
	// nothing here, use fmt.println directly
	return
}
// implementations might update the given prefix with the actual values to substitute
func (s *StructConsoleLogger) SetPrefix(prefixPattern string) (pLogger ILogger) {
	if !util.IsEmptyString(prefixPattern) {
		s.prefix = prefixPattern
		s.isPrefixSet = true
	}
	pLogger = s
	return
}

func (s *StructConsoleLogger) printPrefix() (finalLine string) {
	if s.isPrefixSet {
		ctfmt.Print(ct.Green, true, fmt.Sprintf("[%v]", s.prefix))
		ctfmt.Print(ct.Green, true, fmt.Sprintf("[%v] ", time.Now().UTC().Format(util.CommonDateFormat)))
		//finalLine = fmt.Sprintf("[%v] %v", s.prefix, line)
	}
	return
}

// print out values provided
func (s *StructConsoleLogger) Print(v... interface{}) {
	var bContents bytes.Buffer
	if v != nil {
		for _, interfaceValue := range v {
			bContents.WriteString(fmt.Sprintf("%v ", interfaceValue))
		}
		s.printPrefix()
		ctfmt.Printf(ct.White, true, bContents.String())
	}
}
// print out values provided by substituting them into the format
func (s *StructConsoleLogger) Printf(format string, v... interface{}) {
	if v == nil || len(v) == 0 {
		s.printPrefix()
		ctfmt.Printf(ct.White, true, format)
	} else {
		line := format
		for _, interfaceValue := range v {
			line = strings.Replace(line, "%v", interfaceValue.(string), 1)
		}
		s.printPrefix()
		ctfmt.Printf(ct.White, true, line)
	}
}
// print out values provided in a line format
func (s *StructConsoleLogger) Println(v... interface{}) {
	var bContents bytes.Buffer
	if v != nil {
		for _, interfaceValue := range v {
			bContents.WriteString(fmt.Sprintf("%v ", interfaceValue))
		}
		s.printPrefix()
		ctfmt.Println(ct.White, true, bContents.String())
	}
}

// close the logger, release resources involved
func (s *StructConsoleLogger) Close() (err error) {
	//if s.out != nil {
		//err = s.out.Close()
		//if err != nil {
		//	return
		//}
		//s.out = nil
	//}
	return
}

// return information describing the logger's implementation or config
func (s *StructConsoleLogger) Info() (info string) {
	iMap := make(map[string]interface{})
	iMap[common.LoggerInfoKeyLoggerName] = common.LoggerTypeConsoleLogger
	iMap[common.LoggerInfoKeyPrefix] = s.prefix

	bContents, err := json.Marshal(iMap)
	if err != nil {
		fmt.Println(err)
	}
	info = string(bContents)
	return
}
