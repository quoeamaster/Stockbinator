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
package store

import (
	"Stockbinator/common"
	"Stockbinator/util"
	"errors"
	"fmt"
	"github.com/micro/go-config"
	"os"
	"strings"
)

type StructFilestore struct {
	// storing the application level settings
	AppConfig config.Config
	// target file to operate with
	Filename  string

	// the actual filepath to store and load the data
	filepath string
}

// creator / ctor method
func NewStructFilestore(config config.Config, filename string) (pStore *StructFilestore) {
	pStore = new(StructFilestore)
	if config != nil {
		pStore.AppConfig = config
	}
	pStore.Filename = filename
	// init and find the full filepath
	err := pStore.init()
	if err != nil {
		// log down the error but try to proceed
		fmt.Println(err)
	}
	return
}

func (s *StructFilestore) Persist(data map[string]StructStoreValue) (response StructStoreResponse, err error) {
	return
}

// read all contents from the store (might be an issue when the content size is HUGE
func (s *StructFilestore) ReadAll() (content string, err error) {
	return
}
// read only the content associated by the KEY, PARAMS contains additional information for the read operation
func (s *StructFilestore) ReadByKey(key string, params interface{}) (value StructStoreValue, err error) {
	return
}

// modify the value associated with the key
func (s *StructFilestore) ModifyByKey(key string, value StructStoreValue) (response StructStoreResponse, err error) {
	return
}

// remove value associated with the key
func (s *StructFilestore) RemoveByKey(key string) (valueRemoved StructStoreValue, err error) {
	return
}
// remove all data in the store, be careful~
func (s *StructFilestore) RemoveAll() (err error) {
	
	return
}


// * **************** *
// * common functions *
// * **************** *

func (s *StructFilestore) init() (err error) {
	repo := s.AppConfig.Get(common.ConfigKeyStoreFile, common.ConfigKeyRepo).String("")
	arrEnvVars, arrMatchIndices, err := util.ParseEnvVar(repo)
	// need to replace env var? (might not be the case, since could also hard code the path)
	if arrEnvVars != nil && len(arrEnvVars) > 0 {
		for i := len(arrMatchIndices)-1; i >= 0; i-- {
			matchIndices := arrMatchIndices[i]
			envVar := arrEnvVars[i][matchIndices[0]+1 : matchIndices[1]-1]

			envVarVal := os.Getenv(envVar)
			if util.IsEmptyString(envVarVal) {
				err = errors.New(fmt.Sprintf("env variable [%v] is NOT available.", envVar))

				repo = ""
				// break and use and empty "path"
				break
			}
			// replace it with the repo path
			repo = strings.Replace(repo, arrEnvVars[i], envVarVal, matchIndices[0]-1)
		} // end -- for (reverse loop on indices)
	}
	// append filename to the repo path resolved
	s.filepath = fmt.Sprintf("%v%v", repo, s.Filename)
	return
}

