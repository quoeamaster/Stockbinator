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
package common

// * ******************************************************************
// * the package is storing common / sharable constants or objects
// * that would not involve circular import issues.
// * ******************************************************************


const (
	// default folder name for storing the CONFIG files
	ConfigDefaultFoldername = "config"
	// default folder name for storing the LOG files
	LogDefaultFoldername = "log"

	ConfigEnvPathConfig = "config"
	ConfigEnvPathPath = "path"
	ConfigEnvPrefix = "SB"

	// config entry / key => "holidays" (used in holiday.toml)
	ConfigKeyHolidays = "holidays"

	// config entry / key => "filestore" (app.toml)
	ConfigKeyStoreFile = "filestore"
	// config entry / key => "datastore" (app.toml)
	ConfigKeyStoreData = "datastore"

	// config entry / key => "repo" (app.toml)
	ConfigKeyRepo = "repo"

	// default logger config file -> logger.toml
	ConfigFileLoggerToml = "logger.toml"
	ConfigKeyLoggers = "loggers"
	ConfigKeyFileLoggerFilename = "filename"
	ConfigKeyFileLoggerTimebased = "timebased"
	ConfigKeyFileLoggerTimebasePattern = "timebasePattern"
	ConfigKeyFileLoggerFilepath = "filepath"
	// is the logger a default; if no logger id / name specified the logger "defaulted" will be employed
	ConfigKeyLoggerPathDefault = "default"

	// store -> filestore's default file name if not provided
	StoreDefaultDateFilename = "default.data"
	StoreKeyDefaultDateFilename = "filestore.default.filename"

	// common file status
	FileStatusAvailable    = 200
	FileStatusNotAvailable = 404
	FileStatusLocked       = 401
	FileStatusUnknown      = 500

	// filename for the file-logger (for the "Info" method)
	LoggerFileInfoKeyFilename = "filename"
	// full filepath for the file-logger
	LoggerFileInfoKeyFilepath = "filepath"
	LoggerInfoKeyLoggerName   = "name"
	LoggerInfoKeyPrefix   = "prefix"

	LoggerTypeFileLogger = "fileLogger"
	LoggerTypeConsoleLogger = "consoleLogger"
)
