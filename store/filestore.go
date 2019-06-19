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
	"bytes"
	"errors"
	"fmt"
	"github.com/micro/go-config"
	"os"
	"reflect"
	"strings"
	"time"
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


// save all data into the store
func (s *StructFilestore) Persist(data map[string]StructStoreValue) (response StructStoreResponse, err error) {
	response = s.newStructStoreResponse(CodeSuccess, "")
	pFile, err := os.OpenFile(s.filepath, os.O_APPEND|os.O_WRONLY, 0666)
	// dtor
	defer func() {
		err2 := pFile.Close()
		// do not shadow the original error
		if err == nil && err2 != nil {
			err = err2
			s.handleCommonErrorForResponse(&response, err)
			return
		}
	}()
	if err != nil {
		s.handleCommonErrorForResponse(&response, err)
		return
	}
	// write the values (json in 1 line)
	// build json content by code...
	jsonValue, err := s.toJson(data)
	if err != nil {
		return
	}
	// write to file
	_, err = pFile.WriteString(jsonValue)
	if err != nil {
		fmt.Println("haha... expected")
		return
	}
	_, err = pFile.Seek(int64(len(jsonValue) + 1), 0)
	if err != nil {
		return
	}
	return
}

// read all contents from the store (might be an issue when the content size is HUGE
func (s *StructFilestore) ReadAll() (response StructStoreResponse, content string, err error) {
	return
}
// read only the content associated by the KEY, PARAMS contains additional information for the read operation
func (s *StructFilestore) ReadByKey(key string, params interface{}) (response StructStoreResponse, value StructStoreValue, err error) {
	return
}

// modify the value associated with the key
func (s *StructFilestore) ModifyByKey(key string, value StructStoreValue) (response StructStoreResponse, err error) {
	return
}

// remove value associated with the key
func (s *StructFilestore) RemoveByKey(key string) (response StructStoreResponse, valueRemoved StructStoreValue, err error) {
	return
}
// remove all data in the store, be careful~ For file-store case, the file would not be removed,
// ONLY contents would be truncated
func (s *StructFilestore) RemoveAll() (response StructStoreResponse, err error) {
	response = s.newStructStoreResponse(CodeSuccess, "")

	exists, err := util.IsFileExists(s.filepath)
	if exists {
		// open file and truncate (must be RW mode)
		pFile, err2 := os.OpenFile(s.filepath, os.O_RDWR, 0666)
		// dtor
		defer func() {
			err3 := pFile.Close()
			// do not shadow the original error
			if err == nil && err3 != nil {
				err = err3
				s.handleCommonErrorForResponse(&response, err)
				return
			}
		}()
		if err2 != nil {
			err = err2
			s.handleCommonErrorForResponse(&response, err)
			return
		}
		err2 = pFile.Truncate(0)
		if err2 != nil {
			err = err2
			s.handleCommonErrorForResponse(&response, err)
			return
		}
		// change pointer back to the 1st character space (Beginning-Of-File)
		_, err2 = pFile.Seek(0 ,0)
		if err2 != nil {
			err = err2
			s.handleCommonErrorForResponse(&response, err)
			return
		}
	}
	return
}

func (s *StructFilestore) newStructStoreResponse(code int, message string) (response StructStoreResponse) {
	response = *new(StructStoreResponse)
	response.Code = code
	response.Message = message
	return
}

func (s *StructFilestore) handleCommonErrorForResponse(pResponseStruct *StructStoreResponse, err error) {
	if pResponseStruct != nil && err != nil {
		pResponseStruct.Code = CodeFailure
		pResponseStruct.Message = err.Error()
	}
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

// method to transform the given map into json value
func (s *StructFilestore) toJson(pValueMap map[string]StructStoreValue) (jsonValue string, err error) {
	var bContent bytes.Buffer
	// default value is empty string
	jsonValue = ""
	if pValueMap != nil {
		bContent.WriteString("{")
		// loop all key-value pairs
		for fieldName, storeValue := range pValueMap {
			if !util.IsEmptyString(storeValue.Key) {
				fieldName = storeValue.Key
			}
			bContent.WriteString(fmt.Sprintf("\"%v\": ", fieldName))
			// handle array values
			if storeValue.IsArray {
				fmt.Println(reflect.TypeOf(storeValue))
				err = errors.New("array type not YET implemented")
			} else if storeValue.IsObject {
				fmt.Println(reflect.TypeOf(storeValue))
				err = errors.New("object type not YET implemented")
			} else {
				// normal singular field
				switch storeValue.Type {
				case TypeString:
					bContent.WriteString(fmt.Sprintf("\"%v\", ", storeValue.Value.(string)))
				case TypeInteger:
					bContent.WriteString(fmt.Sprintf("%v, ", storeValue.Value.(int)))
				case TypeFloat:
					bContent.WriteString(fmt.Sprintf("%v, ", storeValue.Value.(float64)))
				case TypeBool:
					bContent.WriteString(fmt.Sprintf("%v, ", storeValue.Value.(bool)))
				case TypeDate:
					typeVal := reflect.TypeOf(storeValue.Value)
					switch typeVal.String() {
					case "string":
						sDate := storeValue.Value.(string)
						_, err2 := time.Parse(util.CommonDateFormat, sDate)
						if err2 != nil {
							err = err2
							// invalid date format; hence parsing failed
							return
						}
						bContent.WriteString(fmt.Sprintf("\"%v\", ", sDate))
					case "time.Time":
						if storeValue.Value == nil {
							err = errors.New(fmt.Sprintf("field [%v] is configured to be a nil value", fieldName))
							return
						}
						dDate := storeValue.Value.(time.Time)
						sDate := dDate.Format(util.CommonDateFormat)
						bContent.WriteString(fmt.Sprintf("\"%v\", ", sDate))
					}
				} // end -- switch (store.Type)
			} // end -- if (singular, array OR object)
		} // end -- for (k v pair for the fields to be jsonfy)
		// remove the last "," if length of bContent is > 1
		if bContent.Len() > 1 {
			jsonValue = bContent.String()
			jsonValue = jsonValue[0:len(jsonValue)-2]
			jsonValue = fmt.Sprintf("%v }\n", jsonValue)
		} else {
			bContent.WriteString("}\n")
			jsonValue = bContent.String()
		} // end -- if (have real content or not)
	}
	return
}
