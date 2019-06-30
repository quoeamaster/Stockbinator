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
	"fmt"
	"io"
	"strings"
)

const (
	loggerTimebasedPattern = "yyyy-mm-dd"
)

// structure containing the config information for the logger implementations; if the implementation is non file based,
// then there is no need to pass in this structure
type StructLoggerConfig struct {
	// the filepath for the log file (simply directory)
	Filepath string
	// the basic filename for the log file
	Filename string

	// is the filename timebased ?? If not, would be number based, starting from "1",
	// increment by 1 every time a rolled file is needed
	IsFilenameTimebased bool
	// if timebased => pattern would be ??
	FilenameTimebasedPattern string
}

// interface for common-logger
type ILogger interface {
	// setup for the logger, requires a Struct of logger-config and can provide an optional params map object
	SetupLogger(config *StructLoggerConfig, params... map[string]interface{}) (err error)
	// updates the Writer / output (e.g. set it to a valid FileWriter instance)
	SetOutput(writer io.Writer) (err error)
	// implementations might update the given prefix with the actual values to substitute
	SetPrefix(prefixPattern string) (pLogger ILogger)

	// print out values provided
	Print(v... interface{})
	// print out values provided by substituting them into the format
	Printf(format string, v... interface{})
	// print out values provided in a line format
	Println(v... interface{})

	// close the logger, release resources involved
	Close() (err error)

	// return information describing the logger's implementation or config
	Info() (info string)
}

// ctor.
func NewStructLoggerConfig(filepath, filename string, isFilenameTimebased bool, filenameTimebasedPattern string) (cfg *StructLoggerConfig) {
	cfg = new(StructLoggerConfig)
	cfg.Filepath = filepath
	cfg.Filename = filename
	cfg.IsFilenameTimebased = isFilenameTimebased
	if cfg.IsFilenameTimebased {
		if strings.Compare(strings.Trim(filenameTimebasedPattern, ""), "") == 0 {
			cfg.FilenameTimebasedPattern = loggerTimebasedPattern
		} else {
			cfg.FilenameTimebasedPattern = filenameTimebasedPattern
		}
	}
	return
}

// ******************************
// *** infrastructure related ***
// ******************************

// repository of ILogger instances
type structRepositoryLoggers struct {
	// the map of ILogger instance(s)
	repo map[string]ILogger
	// default ILogger id / name
	defaultILogger string

	// flag to indicate setup in progress (prevent next getLogger call to re-init the initialization process)
	isSetupInProgess bool
	// flag to indicate the repo is READY to serve
	isReady bool
}
// actual private instance of the repo of ILogger(s)
var pRepositoryLoggers = NewStructRepositoryLoggers()

func NewStructRepositoryLoggers() (pRepo *structRepositoryLoggers) {
	pRepo = new(structRepositoryLoggers)
	pRepo.repo = make(map[string]ILogger)
	err := pRepo.Init()
	if err != nil {
		fmt.Println(err)
		pRepo = nil
	}
	return
}

func (s *structRepositoryLoggers) Init() (err error) {
	if s.isReady {
		return
	}
	if s.isSetupInProgess {
		return
	}
	s.isReady = false
	s.isSetupInProgess = true

	// a. read config file from env var or from a default config path
	repoPath, err := util.GetConfigFolderPath()
	if err != nil {
		return
	}
	cfg, err := util.LoadConfig(fmt.Sprintf("%v%v", repoPath, common.ConfigFileLoggerToml))
	sliceLoggerDefinition := cfg.Map()[common.ConfigKeyLoggers].([]interface{})

	// b. init all ILogger defined in the toml file
	for _, loggerDefinition := range sliceLoggerDefinition {
		loggerCfgMap := loggerDefinition.(map[string]interface{})
		// must have a "name" attribute
		name := loggerCfgMap["name"].(string)
		switch name {
		case common.LoggerTypeFileLogger:
			err2 := s.loadFileLogger(loggerCfgMap)
			if err2 != nil {
				err = err2
				return
			}
			break
		case common.LoggerTypeConsoleLogger:
			err2 := s.loadConsoleLogger(loggerCfgMap)
			if err2 != nil {
				err = err2
				return
			}
			break
		default:
			fmt.Printf("non supported logger type yet ... [%v]", name)
		}
	}
	return
}

func (s *structRepositoryLoggers) loadFileLogger(loggerCfgMap map[string]interface{}) (err error) {
	// load all related fields and cast them to string or bool
	loggerName := common.LoggerTypeFileLogger
	loggerFilename := loggerCfgMap[common.ConfigKeyFileLoggerFilename].(string)
	loggerIsTimebased := false
	iValue := loggerCfgMap[common.ConfigKeyFileLoggerTimebased]
	if iValue != nil {
		loggerIsTimebased = iValue.(bool)
	}
	loggerTimebasedPattern := ""
	if loggerIsTimebased {
		iValue = loggerCfgMap[common.ConfigKeyFileLoggerTimebasePattern]
		if iValue != nil {
			loggerTimebasedPattern = iValue.(string)
		}
	}
	loggerFilepath := loggerCfgMap[common.ConfigKeyFileLoggerFilepath].(string)
	loggerFilepath, err = util.GetLoggerFolderPath(loggerFilepath)
	if err != nil {
		return
	}
	pLoggerCfg := NewStructLoggerConfig(loggerFilepath, loggerFilename, loggerIsTimebased, loggerTimebasedPattern)
	// call ctor
	pFileLogger := NewStructFileLogger()
	err = pFileLogger.SetupLogger(pLoggerCfg)
	if err != nil {
		return
	}
	// set key and value pair
	s.repo[loggerName] = pFileLogger
	// is default?
	if s.isDefaultLogger(loggerCfgMap) && util.IsEmptyString(s.defaultILogger) {
		s.defaultILogger = loggerName
	}
	// fmt.Println(fmt.Sprintf("filename: %v; filepath: %v; isTimebased: %v; timebasePattern: %v\n", loggerFilename, loggerFilepath, loggerIsTimebased, loggerTimebasedPattern))
	return
}

func (s *structRepositoryLoggers) loadConsoleLogger(loggerCfgMap map[string]interface{}) (err error) {
	loggerName := common.LoggerTypeConsoleLogger
	pConfg := NewStructLoggerConfig("", "", false, "")
	pLogger := NewStructConsoleLogger()
	err = pLogger.SetupLogger(pConfg)
	if err != nil {
		return
	}
	s.repo[loggerName] = pLogger
	// is default?
	if s.isDefaultLogger(loggerCfgMap) && util.IsEmptyString(s.defaultILogger) {
		s.defaultILogger = loggerName
	}
	return
}

// check if the logger config contains a "default" key or not;
// if available check the bool value to determine if this logger instance would be the default logger of the system
func (s *structRepositoryLoggers) isDefaultLogger(loggerCfgMap map[string]interface{}) (isDefault bool) {
	defValue := loggerCfgMap[common.ConfigKeyLoggerPathDefault]
	if defValue != nil && defValue.(bool) {
		isDefault = true
	}
	return
}

// return the logger. Logic as follow:
// 1. if loggerName is given, try to get a logger instance based on the given logger-name
// 2. final check on availability of the logger instance;
// 	if nil -> check if any default logger is available and grab the default logger out
func GetLogger(loggerName...string) (pLogger ILogger) {
	if loggerName != nil && !util.IsEmptyString(loggerName[0]) {
		pLogger = pRepositoryLoggers.repo[loggerName[0]]
	}
	if pLogger == nil && !util.IsEmptyString(pRepositoryLoggers.defaultILogger) {
		// return default logger instead (would not let the logging be ignored)
		pLogger = pRepositoryLoggers.repo[pRepositoryLoggers.defaultILogger]
	}
	return
}

func CloseAllLoggers() (err error) {
	if pRepositoryLoggers != nil {
		for _, pLogger := range pRepositoryLoggers.repo {
			if pLogger != nil {
				err2 := pLogger.Close()
				if err2 != nil {
					fmt.Println("unexpected error on closing logger", err2)
					err = err2
					// continue though
				}
			}
		} // end -- for (logger repository loop)
	}

	return
}





