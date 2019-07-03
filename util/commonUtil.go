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
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const DefaultConfigValString = "NIL"
// common date format for parsing => yyyy-MM-ddThh:mm:ss-07:00
const CommonDateFormat = "2006-01-02T15:04:05-07:00"

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

func IsFolderExists(filePath string) (exists bool, err error) {
	exists = false

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	// non directory? great
	if fileInfo.IsDir() == true {
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
		var parts []string
		if strings.HasSuffix(finalPath, GetPathSeparator()) {
			parts = strings.Split(finalPath[0: len(finalPath)-1], GetPathSeparator())
		} else {
			parts = strings.Split(finalPath, GetPathSeparator())
		}
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


// validates if the given time part is valid (e.g. hour24 = 0~23, min = 0~59, sec = 0~59)
func IsValidTimePart(value int, partName string) (valid bool) {
	valid = true

	if value >= 0 {
		switch partName {
		case "hour24":
			if value > 23 {
				valid = false
			}
			break
		case "min":
		case "sec":
			if value > 59 {
				valid = false
			}
			break
		default:
			// unsupported time parts...
			valid = false
		}
	} else {
		valid = false
	}
	return
}


var timezoneRegexp = regexp.MustCompile(`[+-]?[0-1]?[0-9]:[0-3]0`)

// method to check if the given string is a valid timezone format =>
// 1) +08:00 OR
// 2) +8:00 OR
// 3) -07:00
func IsValidTimezone(value string) (valid bool) {
	return timezoneRegexp.MatchString(value)
}


// method to parse and return the request body contents in []byte
func GetRequestBodyInBytes(pBody *io.ReadCloser) (bContent []byte, err error) {
	bContent, err = ioutil.ReadAll(*pBody)
	return
}

// method to check if the supplied string is empty or not
func IsEmptyString(value string) bool {
	if strings.Compare(strings.Trim(value, " "), "") == 0 {
		return true
	}
	return false
}

// TODO ...
// parse string to Date object...

// 2019-05-28T07:30:00+0000
func CreateTodayTargetTimeByHourMinTimezone(hour, min int, timezone string) (todayTargetTime string) {
	// skip validation as assume isValidTimePart() has been called earlier
	bToday := ""
	now := time.Now()
	bToday = fmt.Sprintf("%v-", now.Year())

	dPart := int64(now.Month())
	if dPart < 10 {
		//bToday.WriteString(fmt.Sprintf("0%v-", dPart))
		bToday = fmt.Sprintf("%v0%v-", bToday, dPart)
	} else {
		//bToday.WriteString(fmt.Sprintf("%v-", dPart))
		bToday = fmt.Sprintf("%v%v-", bToday, dPart)
	}

	dPart2 := now.Day()
	if dPart2 < 10 {
		bToday = fmt.Sprintf("%v0%vT", bToday, dPart2)
	} else {
		bToday = fmt.Sprintf("%v%vT", bToday, dPart2)
	}

	// TODO: hh:mm:ssTZ
	if hour < 10 {
		bToday = fmt.Sprintf("%v0%v:", bToday, hour)
		//bToday.WriteString(fmt.Sprintf("0%v:", hour))
	} else {
		bToday = fmt.Sprintf("%v%v:", bToday, hour)
	}

	if min < 10 {
		bToday = fmt.Sprintf("%v0%v:", bToday, min)
		//bToday.WriteString(fmt.Sprintf("0%v:", min))
	} else {
		bToday = fmt.Sprintf("%v%v:", bToday, min)
	}
	bToday = fmt.Sprintf("%v00%v", bToday, timezone)

	return bToday
}


// helper method to create the date-time based on today's date associated with
// the hours plus minutes plus timezone provided
func ParseStringDateToTodayUTC(hh24, mm int, timezone string) (pDate time.Time, err error) {
	sDate := CreateTodayTargetTimeByHourMinTimezone(hh24, mm, timezone)

	pDate, err = time.Parse(CommonDateFormat, sDate)
	pDate = pDate.In(time.UTC)

	return
}

func GetTimeTruncatedDate(givenDate *time.Time) (date time.Time, err error) {
	if givenDate == nil {
		date = time.Now()
	} else {
		date = *givenDate
	}
	// start truncate
	sDate := fmt.Sprintf("%v", date.Year())
	iDatePart := int(date.Month())
	if iDatePart < 10 {
		sDate = fmt.Sprintf("%v-0%v", sDate, iDatePart)
	} else {
		sDate = fmt.Sprintf("%v-%v", sDate, iDatePart)
	}
	iDatePart = int(date.Day())
	if iDatePart < 10 {
		sDate = fmt.Sprintf("%v-0%v", sDate, iDatePart)
	} else {
		sDate = fmt.Sprintf("%v-%v", sDate, iDatePart)
	}
	sDate = fmt.Sprintf("%vT00:00:00", sDate)
	// timezone handling
	_, zoneDiff := date.Zone()
	// 3600 sec = 1 hour
	// TODO: handling => (assume all perfect timezone in hours and no "half-hour" such as +05:30)
	zoneDiffHour := zoneDiff / 3600
	if zoneDiffHour < 0 {
		sDate = fmt.Sprintf("%v-", sDate)
		// absolute it
		zoneDiffHour = -1 * zoneDiffHour
	} else {
		sDate = fmt.Sprintf("%v+", sDate)
	}
	if zoneDiffHour < 10 {
		sDate = fmt.Sprintf("%v0%v:00", sDate, zoneDiffHour)
	} else {
		sDate = fmt.Sprintf("%v%v:00", sDate, zoneDiffHour)
	}
	// convert to the final result
	date, err = time.Parse(CommonDateFormat, sDate)
	if err != nil {
		return
	}
	return
}


// check if the given date is on weekend (day-of-week 0 = sunday, 6 = saturday)
func IsWeekend(date time.Time) (isWeekend bool) {
	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 || dayOfWeek == 6 {
		isWeekend = true
	} else {
		isWeekend = false
	}
	return
}

// check if the given date is one of the holidays given
func IsHoliday(date *time.Time, dateInString *string, holidays []string) (isHoliday bool, err error) {
	var targetDate time.Time
	if date == nil && dateInString == nil {
		err = errors.New("both date and dateInstring parameter is invalid => nil\n")
		return
	} else if dateInString != nil {
		targetDate, err = time.Parse(CommonDateFormat, *dateInString)
		if err != nil {
			return
		}
	} else {
		targetDate = *date
	}
	// truncate / trim to date level for comparison
	targetDate = targetDate.UTC().Truncate(time.Hour * 24)

	// check against the holiday[]
	if holidays == nil || len(holidays) == 0 {
		err = errors.New("invalid holidays array, it is either nil or empty")
		return
	}
	for _, holiday := range holidays {
		hDate, err2 := time.Parse(CommonDateFormat, holiday)
		if err2 != nil {
			err = err2
			return
		}
		hDate = hDate.UTC().Truncate(time.Hour * 24)
		// compare
		// a) same = holiday (return true)
		// b) different, but hDate is already after the targetDate which this is a FUTURE holiday
		// 	comparing with targetDate and should skip the check (return false)
		if targetDate.Equal(hDate) {
			isHoliday = true
			return
		} else if targetDate.Before(hDate) {
			isHoliday = false
			return
		} // end -- if (targetDate vs hDate)
	}
	return
}

// parse and find out the env variable keys (e.g. {ENV_VAR_KEY} ) plus the indices involved per match
func ParseEnvVar(key string) (matches []string, matchesIndex [][]int,  err error) {
	regMatcher, err := regexp.Compile(`\{[a-z|A-Z|-|_]+\}`)
	matches = regMatcher.FindAllString(key, -1)
	matchesIndex = regMatcher.FindAllStringIndex(key, -1)
	return
}

// helper method to replace the time pattern string with the given time's values
func PrepareTimebasedPatternWithGivenTime(pattern string, givenTime time.Time) (v string) {
	v = strings.Replace(pattern, "yyyy", fmt.Sprintf("%v", givenTime.Year()), -1)
	if strings.Index(v, "yy") != -1 {
		v = strings.Replace(v, "yy", fmt.Sprintf("%v", givenTime.Year())[2:], -1)
	}
	mm := int(givenTime.Month())
	if mm < 10 {
		v = strings.Replace(v, "mm", fmt.Sprintf("0%v", mm), -1)
	} else {
		v = strings.Replace(v, "mm", fmt.Sprintf("%v", mm), -1)
	}
	dd := givenTime.Day()
	if dd < 10 {
		v = strings.Replace(v, "dd", fmt.Sprintf("0%v", dd), -1)
	} else {
		v = strings.Replace(v, "dd", fmt.Sprintf("%v", dd), -1)
	}
	return
}
