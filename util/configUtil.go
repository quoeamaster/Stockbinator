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
package util

import (
	"Stockbinator/common"
	"errors"
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"os"
	"strings"
)

// common method to get the config-folder's path:
// 1. check if env variable SB_CONFIG_PATH is set or not, if so, use it as repo path OR
// 2. get the current executable path and append "/config" to this path
//
// PS. there is no guarantee that the path decided is available or not
func GetConfigFolderPath() (repoPath string, err error) {
	repoPath = DefaultConfigValString

	// 1) check env variable
	cfg := config.NewConfig()
	// https://github.com/micro/go-micro/tree/master/config/source/env => stripped prefix (reducing 1 level of hierarchy)
	eSrc := env.NewSource(env.WithStrippedPrefix(common.ConfigEnvPrefix))
	err = cfg.Load(eSrc)
	if err != nil {
		return
	}
	repoPath = cfg.Get(common.ConfigEnvPathConfig, common.ConfigEnvPathPath).String(DefaultConfigValString)
	if strings.Compare(repoPath, DefaultConfigValString) == 0 {
		// means no env path available, get the executable's path
		repoPath, err = GetExecutablePath(true)
		if err != nil {
			return
		}
		repoPath = fmt.Sprintf("%v%v%v/", repoPath, GetPathSeparator(), common.ConfigDefaultFoldername)
	}
	return
}

// similar to GetConfigFolderPath(), this method parses the env variable OR
// the hard-coded path provided
// when hard-coded path is provided; first check if the path is empty string or not;
// if empty => get the executable's path append "log" to it as final path OR
// just use the given hard-coded path if non empty
//
// PS. there is no guarantee on the path is available or not
func GetLoggerFolderPath(initialPath string) (repoPath string, err error) {
	matches, _, err := ParseEnvVar(initialPath)
	if err != nil {
		return
	}
	if matches != nil && len(matches) > 0 {
		repoPath = initialPath
		for i := len(matches)-1; i >= 0; i-- {
			envMatch := matches[i]
			envMatchKey := os.Getenv(envMatch[1:len(envMatch)-1])
			repoPath = strings.Replace(repoPath, envMatch, envMatchKey, 1)
		}
	} else {
		// just a hard-coded path
		if strings.Compare("", strings.Trim(initialPath, "")) == 0 {
			// use the executable path + "log/"
			repoPath, err = GetExecutablePath(true)
			if err != nil {
				return
			}
			repoPath = fmt.Sprintf("%v%v%v/", repoPath, GetPathSeparator(), common.LogDefaultFoldername)
		} else {
			repoPath = initialPath
		}
	}
	return
}

// method to load a given file into a Config object
func LoadConfig(filePath string) (configObject config.Config, err error) {
	exists, err := IsFileExists(filePath)
	if !exists {
		err = errors.New(fmt.Sprintf("The supplied config file path DOES NOT exists => %v", filePath))
		return
	}
	if err != nil {
		return
	}

	// continue parsing
	configObject = config.NewConfig()
	fSrc := file.NewSource(file.WithPath(filePath))
	err = configObject.Load(fSrc)
	if err != nil {
		return
	}
	return
}
