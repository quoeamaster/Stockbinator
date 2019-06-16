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
package tests

import (
	"Stockbinator/util"
	"fmt"
	"testing"
)

func TestGetCurrentYearHolidays(t *testing.T) {
	if !*pFlagCrawlerUtil {
		t.SkipNow()
	}
	LogTestOutput("TestGetCurrentYearHolidays", "*** start test ***")

	results := []struct {
		year 		string
		holidayLen 	int
	}{
		{ "2019", 18 },
		// since you should not get any as this year is NOT 2018 (though in the holidays_2018.toml there are 19 entries)
		{ "2018", 0 },
	}

	// test on current year (e.g. 2019)
	holidaySlice, err := util.GetCurrentYearHolidays(&SharableStockModuleConfig.Holidays)
	if err != nil {
		t.Fatal(err)
	}
	if holidaySlice != nil && len(holidaySlice) != results[0].holidayLen {
		t.Fatal(fmt.Sprintf("expected holidays [current year] available and num-of-holidays should be the same, expected %v but %v",
			len(holidaySlice), results[0].holidayLen))
	}
	LogTestOutput("TestGetCurrentYearHolidays", "current year validation pass")

	holidayConfig2018 := Holidays2018Config

	// test on previous year (e.g. 2018)
	holidaySlice, err = util.GetCurrentYearHolidays(&holidayConfig2018)
	if err != nil {
		t.Fatal(err)
	}
	if holidaySlice != nil && len(holidaySlice) != results[1].holidayLen {
		LogTestOutput("TestGetCurrentYearHolidays", "2018 validation failed; is empty slice which means no holidays info for 2018")
		t.Fatal(fmt.Sprintf("expected holidays [2018] available and num-of-holidays should be the same, expected %v but %v",
			len(holidaySlice), results[1].holidayLen))
	}

	LogTestOutput("TestGetCurrentYearHolidays", "*** end test ***\n")
}
