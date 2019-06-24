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
	SetPrefix(prefixPattern string)

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




