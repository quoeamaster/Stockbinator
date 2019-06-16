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
	"time"
)

// class for testing all common-util functions

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
	date      time.Time
	dayOfWeek int
	isHoliday bool
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



