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
	"io/ioutil"
	"net/http"
	"strings"
)

type StructCommonResponse struct {
	// http response code (only 200 or 201 is a Successful operation)
	ResponseCode int
	// optional string message to describe the operation
	Message string
}

// create method for the common-response object
func NewStructCommonResponse(responseCode int, message string) *StructCommonResponse {
	pInst := new(StructCommonResponse)
	pInst.ResponseCode = responseCode
	pInst.Message = message
	return pInst
}

// method to get contents from a URL
func GetContentFromUrl(url string) (content string, err error) {
	if len(url) == 0 || strings.Compare(strings.Trim(url, ""), "") == 0 {
		err = errors.New("url provided is invalid, probably EMPTY~")
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		resp.Body.Close()
	}()
	bContent, err := ioutil.ReadAll(resp.Body)
	content = string(bContent)

	return
}