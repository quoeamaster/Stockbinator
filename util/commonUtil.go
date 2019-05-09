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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const DefaultConfigValString = "NIL"

// return the path separator based on current OS
func GetPathSeparator() string {
	return string(os.PathSeparator)
}

// return the executable's fullpath; in "excludeExecutable" is true
// then the return path would exclude the executable's name
func GetExecutablePath(excludeExecutable bool) (ePath string, err error) {
	ePath, err = os.Executable()
	if err != nil {
		return
	}
	// need to exclude executable path??
	if excludeExecutable == true {
		ePath = path.Dir(ePath)
	}
	return
}

// check if the given filepath exists and valid or not
func IsFileExists(filePath string) (exists bool, err error) {
	exists = false

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	// non directory? great
	if fileInfo.IsDir() == false {
		exists = true
		return
	}
	// unknown situation; however treated as non-exists
	return
}

// get back only matched folders under a given path;
// matching criteria is the folder starts with the given folder-prefix
func ListMatchingFoldersUnderPath(repoPath string, folderPrefix string) (matchedFolderPaths []string, err error) {
	matchedFolderPaths = make([]string, 0)

	arrFileInfo, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return
	}
	for _, fileInfo := range arrFileInfo {
		if fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), folderPrefix) {
			matchedFolderPaths = append(matchedFolderPaths, fmt.Sprintf("%v%v/", repoPath, fileInfo.Name()))
		}
	}
	return
}

// way to concatenate strings together to form a final-path
func ConcatenateFilePath(paths ...string) (finalPath string) {
	finalPath = ""
	if len(paths) > 0 {
		for _, path := range paths {
			finalPath = fmt.Sprintf("%v%v", finalPath, path)
		}
	}
	return
}

// extract the final sub-path of the given path; the separator is based on the OS's path separator
func GetLastPathElement(finalPath string) (string, error) {
	if len(finalPath) > 0 {
		parts := strings.Split(finalPath, GetPathSeparator())
		partsLen := len(parts)
		if partsLen > 0 {
			return parts[partsLen - 1], nil
		} else if partsLen == 0 {
			return parts[0], nil
		} else {
			return "", errors.New(fmt.Sprintf("exception! The given PATH does not have a sub-path => [%v]", finalPath))
		}
	} else {
		return "", errors.New(fmt.Sprintf("exception! The given PATH is not valid => [%v]", finalPath))
	}
}


// TODO ...
// parse string to Date object...
