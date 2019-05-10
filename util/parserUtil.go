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
	"github.com/buger/jsonparser"
	"strconv"
	"strings"
)

const ParseErrKeyPathNotFound = "Key path not found"


// important concept about usage of StructJsonParser.
// scenario 1: do not set ignore-error list
// - basically, if there were any error occurred during the Get(), all kinds of error would exposed;
//   however there are error(s) that you might want to ignore -> "key path not found"
//   especially when certain parameters might be optional. In this case the developer would need
//   to handle all error checking and ignore error(s) that might do no harm to your biz logic.
// scenario 2: set ignore-error list + using DefaultValue()
// - if we would like to ignore "key path not found" error; then how could we
//   know if the value returned is valid or just the default value for casting failure?
//   We can deal with this by setting the DefaultValue() and perform a bool check
//   against the value returned by a key versus the provided DefaultValue; if both values are identical that means
//   the casting failed (probably because of "key path not found" in this scenario)
//
// code snippet:
//   c.pJsonParser.ResetIgnoreError().IgnoreError(util.ParseErrKeyPathNotFound)
//
//   paramVal, err := c.pJsonParser.Get(bArr, "hour24")
//   if err != nil { // really is an unexpected error
//     fmt.Println(err)
//   } else {
//     if paramVal.Default(-1).IntValue() == -1 { // means key path not found...
//       fmt.Println("key path not found??? hour24")
//     } else {
//       fmt.Printf("%v = %v~\n", "hour24", paramVal.IntValue())
//     }
//   }
type StructJsonParser struct {
	// a list of ignored errors based on the key-words description
	ignoreList []string
}

// creator method for the JsonParser struct
func NewStructJsonParser() (pParser *StructJsonParser) {
	pParser = new(StructJsonParser)
	pParser.ignoreList = make([]string, 0)
	return
}

// updates the ignoreList with the given description(s);
// only new description(s) would be appended to the ignore-list
func (j *StructJsonParser) IgnoreError(description ...string) {
	for _, desc := range description {
		found := false
		for _, ignoreDesc := range j.ignoreList {
			if strings.Compare(desc, ignoreDesc) == 0 {
				// found
				found = true
				break
			}
		}	// end -- if (ignoreList range)
		if found == false {
			j.ignoreList = append(j.ignoreList, desc)
		}
	}	// end -- if (description range)
}

// reset the ignore-error descriptions
func (j *StructJsonParser) ResetIgnoreError() *StructJsonParser {
	// redundant, but kind of play safe to set to nil first and release previous slice's memory
	j.ignoreList = nil
	j.ignoreList = make([]string, 0)
	return j
}


func (j *StructJsonParser) Get(byteContent []byte, keys ...string) (value StructJsonValue, err error) {
	pValue := new(StructJsonValue)
	bVal, _, _, err := jsonparser.Get(byteContent, keys...)
	if err != nil {
		// check against ignore-error list
		if j.isErrorIgnorable(err) {
			err = nil
		}
		return
	}
	pValue.value = string(bVal)
	value = *pValue
	return
}

// method to check if the error string could be ignored (counter check with the ignore-error list)
func (j *StructJsonParser) isErrorIgnorable(err error) (canIgnore bool) {
	canIgnore = false
	errString := err.Error()
	for _, desc := range j.ignoreList {
		if strings.Contains(desc, errString) {
			canIgnore = true
			break
		}
	}
	return
}

// ###########################################
// # return value of the parser object above #
// ###########################################

type StructJsonValue struct {
	// the actual extracted value (default is string datatype, cast to different datatypes when necessary)
	value string
	// default value provided... just in case
	defaultValue interface{}
	err error
}

// set a default value for the underneath value property
func (v *StructJsonValue) Default(value interface{}) *StructJsonValue {
	v.defaultValue = value
	return v
}

// cast the original (string) to an int, error would be encapsulate inside the StructJsonValue.error property instead
func (v *StructJsonValue) IntValue() (iValue int) {
	defer func() {
		if r := recover(); r != nil {
			v.err = r.(error)
			// return the default value set earlier if any
			if v.defaultValue != nil {
				iValue = v.defaultValue.(int)
			}
		}
	}()
	// totally empty string
	if IsEmptyString(v.value) {
		panic(errors.New("empty string value found which could not be parsed to integer"))
	}
	//if strings.Compare(strings.Trim(v.value, " "), "") == 0 {
	//	panic(errors.New("empty string value found which could not be parsed to integer"))
	//}
	iValue, err := strconv.Atoi(v.value)
	if err != nil {
		v.err = err
	}
	return
}

// super easy implementation since, the underneath value attribute is already a string. Simply return this value
func (v *StructJsonValue) StringValue() (sValue string) {
	if IsEmptyString(v.value) {
		if v.defaultValue != nil {
			return v.defaultValue.(string)
		}
	}
	return v.value
}
