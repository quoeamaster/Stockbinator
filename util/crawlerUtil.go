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
	"strings"
	"time"
)

func GetCurrentYearHolidays(pHolidayConfig *config.Config) (holidaySlice []string, err error) {
	holidaySlice = make([]string, 30)
	// catching runtime errors when translating the config keys' value to []string
	defer func() {
		r := recover()
		if r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if pHolidayConfig != nil {
		year := string(fmt.Sprintf("%v", time.Now().Year()))
		holidaySlice = (*pHolidayConfig).Get(year, common.ConfigKeyHolidays).StringSlice(holidaySlice)
	}
	// check if the contents are valid or not ("" is non valid)
	// return an empty slice (length of 0) if non valid contents available
	if len(holidaySlice) ==0 || 0 == strings.Compare("", strings.Trim(holidaySlice[0], "")) {
		holidaySlice = make([]string, 0)
	}
	return
}
