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
	"strings"
	"testing"
	"time"
)

// class for testing all common-util functions

// Test on isWeekend

func TestIsWeekend(t *testing.T)  {
	if !*pFlagCommonUtil {
		t.SkipNow()
	}
	LogTestOutput("TestIsWeekend", "** start test **")

	knownWeek, err := helperTestIsWeekend()
	if err != nil {
		t.Fatal(err)
	}
	for _, dateStruct := range knownWeek {
		isWeekend := util.IsWeekend(dateStruct.date)
		switch dateStruct.dayOfWeek {
		case 0:
			// 0 and 6 are weekends, let it fall-through
			fallthrough
		case 6:
			LogTestOutput("TestIsWeekend", fmt.Sprintf("weekend met => %v; %v", dateStruct.date, dateStruct.dayOfWeek))
			if !isWeekend {
				t.Fatal(fmt.Sprintf("should be weekend but interpreted as weekday instead: %v (%v)", dateStruct.date, dateStruct.dayOfWeek))
			}
		default:
			if isWeekend {
				t.Fatal(fmt.Sprintf("should be weekday but interpreted as weekend instead: %v (%v)", dateStruct.date, dateStruct.dayOfWeek))
			}
		}
	}
	LogTestOutput("TestIsWeekend", "** end test **\n")
}

type structDateForTest struct {
	date      	time.Time
	dayOfWeek 	int
	isHoliday 	bool
	displayTime string
}

// creation of 1 week of data which starts with 2019-06-02 which is a sunday (day-of-week == 0)
func helperTestIsWeekend() (knownWeek []structDateForTest, err error) {
	knownWeek = make([]structDateForTest, 7)

	date, err := time.Parse(util.CommonDateFormat, "2019-06-02T00:00:00+08:00")
	if err != nil {
		return
	}
	// date = date.UTC()
	for idx := range knownWeek {
		// LogTestOutput("helperTestIsWeekend", fmt.Sprintf("%v %v", date, int(date.Weekday())))
		pDateStruct := new(structDateForTest)
		pDateStruct.date = date
		pDateStruct.dayOfWeek = int(date.Weekday())
		knownWeek[idx] = *pDateStruct
		// append 1 more day to the date
		date = date.Add(time.Hour * 24)
	}
	return
}

// Test on isHoliday

// note that the holidays are targeted for the year 2019 of HK's public holidays
// https://www.gov.hk/en/about/abouthk/holiday/2019.htm
func TestIsHoliday(t *testing.T) {
	if !*pFlagCommonUtil {
		t.SkipNow()
	}
	LogTestOutput("TestIsHoliday", "** start test **")

	// normally a country should not have more than 30 statuary holidays
	holidaySlice := make([]string, 30)
	holidaySlice = SharableStockModuleConfig.Holidays.Get("2019", "holidays").StringSlice(holidaySlice)
	// for debug only
	//for _, h := range holidaySlice {
	//	LogTestOutput("TestIsHoliday", h)
	//}
	// compare
	targetDates := helperTestIsHoliday()
	for _, tDate := range targetDates {
		isH, err := util.IsHoliday(&tDate.date, nil, holidaySlice)
		if err != nil {
			t.Fatal(err)
		}
		if isH != tDate.isHoliday {
			t.Fatal(fmt.Sprintf("a Holiday was expected but somehow not interpreted correctly => %v (%v)", tDate.date, tDate.isHoliday))
		}
		LogTestOutput("TestIsHoliday", fmt.Sprintf("%v is-holiday? %v", tDate.date, tDate.isHoliday))
	} // end -- for (targetDates)


	LogTestOutput("TestIsHoliday", "** end test **\n")
}

func helperTestIsHoliday() (dates []structDateForTest) {
	dates = make([]structDateForTest, 0)

	// 2019-01-01 true
	d, err := time.Parse(util.CommonDateFormat, "2019-01-01T00:00:00+08:00")
	if err != nil {
		return
	}
	dates = append(dates, structDateForTest{ date: d, isHoliday: true })
	// 2019-01-01 false
	d, err = time.Parse(util.CommonDateFormat, "2019-01-02T00:00:00+08:00")
	if err != nil {
		return
	}
	dates = append(dates, structDateForTest{ date: d, isHoliday: false })
	// 2019-02-07 true
	d, err = time.Parse(util.CommonDateFormat, "2019-02-07T00:00:00+08:00")
	if err != nil {
		return
	}
	dates = append(dates, structDateForTest{ date: d, isHoliday: true })
	// 2019-09-14 true
	d, err = time.Parse(util.CommonDateFormat, "2019-09-14T00:00:00+08:00")
	if err != nil {
		return
	}
	dates = append(dates, structDateForTest{ date: d, isHoliday: true })

	return
}

// Test on truncated time date

func TestGetTimeTruncatedDate(t *testing.T)  {
	if !*pFlagCommonUtil {
		t.SkipNow()
	}
	LogTestOutput("TestGetTimeTruncatedDate", "** start test **")

	results := helperTestGetTimeTruncatedDate()
	for idx, target := range results {
		dDate, err := util.GetTimeTruncatedDate(&target.date)
		if err != nil {
			t.Fatal(err)
		}
		sDate := dDate.Format(util.CommonDateFormat)
		if strings.Compare(target.displayTime, sDate) != 0 {
			t.Fatal(fmt.Sprintf("expected date to be [%v] but got [%v]", target.displayTime, sDate))
		}
		LogTestOutput("TestGetTimeTruncatedDate", fmt.Sprintf("round %v passed", idx))
	}

	LogTestOutput("TestGetTimeTruncatedDate", "** end test **\n")
}

func helperTestGetTimeTruncatedDate() (targets []structDateForTest) {
	targets = make([]structDateForTest, 5)

	// create several dates with timezone
	// 0. current time in current timezone
	date := time.Now()
	dStruct := new(structDateForTest)
	dStruct.date = date
	dStruct.displayTime = util.CreateTodayTargetTimeByHourMinTimezone(0,0,"+08:00")
	targets[0] = *dStruct

	// 1. 2017-12-25T13:12:34-06:00
	date, _ = time.Parse(util.CommonDateFormat, "2017-12-25T13:12:34-06:00")
	dStruct = new(structDateForTest)
	dStruct.date = date
	dStruct.displayTime = "2017-12-25T00:00:00-06:00"
	targets[1] = *dStruct

	// 2. 2021-09-02T00:23:00+08:00
	date, _ = time.Parse(util.CommonDateFormat, "2021-09-02T00:23:00+08:00")
	dStruct = new(structDateForTest)
	dStruct.date = date
	dStruct.displayTime = "2021-09-02T00:00:00+08:00"
	targets[2] = *dStruct

	// 3. 2021-10-30T00:00:56+00:00
	date, _ = time.Parse(util.CommonDateFormat, "2021-10-30T00:00:56+00:00")
	dStruct = new(structDateForTest)
	dStruct.date = date
	dStruct.displayTime = "2021-10-30T00:00:00+00:00"
	targets[3] = *dStruct

	// 4. 2021-10-30T23:00:00-01:00
	date, _ = time.Parse(util.CommonDateFormat, "2021-10-30T23:00:00-01:00")
	dStruct = new(structDateForTest)
	dStruct.date = date
	dStruct.displayTime = "2021-10-30T00:00:00-01:00"
	targets[4] = *dStruct

	return
}

// test env parsing
func TestParseEnvVar(t *testing.T)  {
	LogTestOutput("TestParseEnvVar", "** start test **")
	if !*pFlagCommonUtil {
		t.SkipNow()
	}

	results := []struct {
		key          string
		matchList    []string
		matchIndices [][]int
	}{
		{ key: "/{testing_id}/movie", matchList: []string{ "{testing_id}" }, matchIndices: [][]int{ { 1, 13 } } },
		{ key: "hd0/{DB_PATH}/backup", matchList: []string{ "{DB_PATH}" }, matchIndices: [][]int{ { 4, 13 } } },
		{ key: "hd0/{twitter}/{gifs}", matchList: []string{ "{twitter}", "{gifs}" }, matchIndices: [][]int{ { 4, 13 }, { 14, 20 } } },
	}
	for _, result := range results {
		matches, matchIndices, err := util.ParseEnvVar(result.key)
		handlerCommonError(err)
		// compare matches with matchList
		if matches == nil || len(matches) != len(result.matchList) {
			LogTestOutput("TestParseEnvVar", fmt.Sprintf("matches returned => %v VS expected matchList => %v",
				matches, result.matchList))
			t.Fatal(fmt.Sprintf("either returned matches are nil (%v) OR length of matches VS expected matches are different => %v vs %v",
				matches == nil, len(matches), len(result.matchList)))
		} else {
			for i2, match := range matches {
				if strings.Compare(match, result.matchList[i2]) != 0 {
					t.Fatal(fmt.Sprintf("expected %v BUT got %v", match, result.matchList[i2]))
				}
			} // end -- for (matches elements loop)
		}

		// check also the indices (start, end idx per match)
		if matchIndices == nil || len(matchIndices) != len(result.matchIndices) {
			t.Fatal(fmt.Sprintf("either indices returned is nil (%v) OR expected length is %v BUT got %v",
				matchIndices==nil, len(matchIndices), len(result.matchIndices)))
		} else {
			for i3, idx1stLv := range matchIndices {
				// length check again
				if len(idx1stLv) != len(result.matchIndices[i3]) {
					t.Fatal(fmt.Sprintf("on indices level check, element %v, lenght expected is %v BUT got %v",
						i3, len(idx1stLv), len(result.matchIndices)))
				}
				expectedIndices := result.matchIndices[i3]
				for i4, idxValue := range idx1stLv {
					if idxValue != expectedIndices[i4] {
						t.Fatal(fmt.Sprintf("expected %v BUT got %v", idxValue, expectedIndices[i4]))
					}
				}
			} // end - for (1st level indices match)
		}
	}
	LogTestOutput("TestParseEnvVar", "ALL env value validation passed")
	LogTestOutput("TestParseEnvVar", "ALL matched indices level check passed")

	LogTestOutput("TestParseEnvVar", "** end test **\n")
}

