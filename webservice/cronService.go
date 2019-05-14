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
package webservice

import (
	"Stockbinator/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
	"github.com/emicklei/go-restful"
	"strings"
	"time"
)

const moduleWSCron = "cronService"

type StructCron struct {
	// instance of the jsonParser for this webservice module
	pJsonParser *util.StructJsonParser
}

// creation method for StructCron
func NewStructCron() (cron *StructCron) {
	cron = new(StructCron)
	cron.pJsonParser = util.NewStructJsonParser()
	return
}

// method to update or insert a cron schedule and its corresponding method.
// However this cron service only handles hour, min, sec and timezone and excludes year, month and date.
func (c *StructCron) UpsertTimeCron( hour24, min, sec int, timezone, stocksModuleRule string ) (inserted bool, err error) {
	// validation
	valid := false
	if util.IsValidTimePart(hour24, "hour24") &&
		util.IsValidTimePart(min, "min") &&
		util.IsValidTimePart(sec, "sec") &&
		util.IsValidTimezone(timezone) &&
		!util.IsEmptyString(stocksModuleRule) {
		valid = true
	}
	//fmt.Println("** heya => ", hour24, " ", min, " ", sec, " ", timezone, " ", stocksModuleRule)
	// by default, it should be just update on the cron; unless the entry doesn't exists
	inserted = false
	if valid == false {
		err = errors.New(fmt.Sprintf(`exception! parameters provided are not correct for creating a 
time-cron entry => hour24[%v], min[%v], sec[%v], 
timezone[%v], stockModuleRule[%v]`, hour24, min, sec, timezone, stocksModuleRule))
	} else {
		// add / update the cron
// TODO create a data structue to encapsulate the cron schedules
//  etc... map[string][]StructCronEntry;
//   string(key) = hour24:min:sec
//   StructCronEntry contains the stocksModuleRule => can look up the rule(s)
//    for scrapping etc (e.g. url, rule_price)
	}
	return
}


// #############################
// # webservice implementation #
// #############################

func (c *StructCron) CreateWebservice() *restful.WebService {
	pWs := new(restful.WebService)
	pWs.Path("/cron").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	// routes under "cron" endpoint (API)
	pWs.Route(pWs.POST("upsert").To(c.upsertTimeCronAPI))
	pWs.Route(pWs.GET("list").To(c.listTimeCronAPI))

	return pWs
}

func (c *StructCron) upsertTimeCronAPI(pReq *restful.Request, pRes *restful.Response) {
	// extract data for calling the corresponding endpoint "UpsertTimeCron"
	defer func() {
		pReq.Request.Body.Close()
	}()

	var bInfoMsg bytes.Buffer
	var pRO *util.StructCommonResponse

	bArr, err := util.GetRequestBodyInBytes(&pReq.Request.Body)
	if err != nil {
		panic(err)
	}
	// prepare the parameters (method oriented)
	hour24 := 0
	min := 0
	sec := 0
	timezone := ""
	stockModuleRule := ""
	c.pJsonParser.ResetIgnoreError().IgnoreError(util.ParseErrKeyPathNotFound)

	paramVal, err := c.pJsonParser.Get(bArr, "hour24")
	if err != nil {
		panic(err)
	} else {
		//if paramVal.Default(-1).IntValue() == -1 {
		//	// means key path not found...
		//	fmt.Println("key path not found??? hour24")
		//} else {
		//	fmt.Printf("%v = %v~\n", "hour24", paramVal.IntValue())
		//}
		hour24 = paramVal.Default(0).IntValue()
	}
	paramVal, err = c.pJsonParser.Get(bArr, "min")
	if err != nil {
		panic(err)
	} else {
		min = paramVal.Default(0).IntValue()
	}
	paramVal, err = c.pJsonParser.Get(bArr, "sec")
	if err != nil {
		panic(err)
	} else {
		sec = paramVal.Default(0).IntValue()
	}
	paramVal, err = c.pJsonParser.Get(bArr, "timezone")
	if err != nil {
		panic(err)
	} else {
		timezone = paramVal.Default("+00:00").StringValue()
	}
	paramVal, err = c.pJsonParser.Get(bArr, "stockModuleRule")
	if err != nil {
		panic(err)
	} else {
		stockModuleRule = paramVal.Default("").StringValue()
		if util.IsEmptyString(stockModuleRule) {
			// exception => stockModuleRule is a MUST parameter
			bInfoMsg.WriteString(`invalid parameters: 'stockModuleRule' is a MUST parameter. 
Optional parameters included: 
hour24 (default 0), min (default 0), sec (default 0), timezone (default "+00:00")`)
		}
	}
	//fmt.Printf("** params => hh:mm:ss Z = %v:%v:%v %v\n", hour24, min, sec, timezone)
	if bInfoMsg.Len() > 0 {
		pRO = util.NewStructCommonResponse(400, bInfoMsg.String())
	} else {
		bInserted, err := c.UpsertTimeCron(hour24, min, sec, timezone, stockModuleRule)
		if err != nil {
			pRO = util.NewStructCommonResponse(500, err.Error())
		} else {
			switch bInserted {
			case true:
				pRO = util.NewStructCommonResponse(201, fmt.Sprintf("a new time-cron entry has been scheduled, stockModuleRule => %v", stockModuleRule))
				break
			default:
				pRO = util.NewStructCommonResponse(200, fmt.Sprintf("an existing time-cron entry has been re-scheduled, stockModuleRule => %v", stockModuleRule))
			}
		}
	}
	err = pRes.WriteAsJson(*pRO)
	if err != nil {
		// just log and continue to serve (sometimes it is a disconnection which could be re-covered)
		c.logError("upsertTimeCronAPI", err.Error())
	}
}

func (c *StructCron) listTimeCronAPI(pReq *restful.Request, pRes *restful.Response) {
	now := time.Now()
	tzName, tzOffsetInMin := now.Zone()
	// can construct a new Time... based on data parts
	now = time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 00, 00, now.Location())

	// use Parse
	tFormat := "2006-01-02 15:04:05 MST"
	// tFormat = "2006-01-02T15:04:05-07:00"
	utcLoc, err := time.LoadLocation("UTC")

	// TODO: use this logic...
	// a) format a current-time
	utcFormatTime := now.Format(tFormat)
	fmt.Println(utcFormatTime)
	// b) form back into a Time object
	utcTime, err := time.Parse(tFormat, utcFormatTime)
	fmt.Println(utcTime)
	// c) convert to UTC
	utcTime = utcTime.In(utcLoc)
	fmt.Println(utcTime)

	// TODO: try another approach (every sec parse the targeted cronTime... not efficient)
	fmt.Println("")
	gmtTimeObj, err := time.Parse("2006-01-02T15:04:05-0700", "2019-05-19T15:30:00+0800")
	tzName, _ = gmtTimeObj.Zone()
	fmt.Println(gmtTimeObj)
	//fmt.Println(tzName)
	//fmt.Println(gmtTimeObj.Local())

	utcTimeObj := gmtTimeObj.In(utcLoc)
	tzName, _ = utcTimeObj.Zone()
	fmt.Println(utcTimeObj)
	//fmt.Println(tzName)
	//fmt.Println(utcTimeObj.Local())

	// TODO: another approach (1. create cronTime in UTC format {string} - conversion, 2. convert now to utc {string}, 3. compare the values)
	// RECOMMENDED
	fmt.Println("\nutc string comparison")
	fmt.Println("a) target cronTime converted to UTC {string} once : need to update every midnight")
	// create everyday just after midnight... etc
	gmtTimeObj, err = time.Parse("2006-01-02T15:04:05-0700", "2019-05-28T15:30:00+0800")
	strGmtTime := gmtTimeObj.In(time.UTC).Format("2006-01-02T15:04:05-0700")
	fmt.Println("a => ", strGmtTime)
	fmt.Println("b) convert current time to utc as well")
	strUtcTime := time.Now().In(time.UTC).Format("2006-01-02T15:04:05-0700")
	fmt.Println("b => ", strUtcTime)
	fmt.Println("c) compare timestamp")
	fmt.Println(strings.Compare(strGmtTime, strUtcTime))










	//utcTime, err := time.ParseInLocation(tFormat, "2019-05-14T15:30:00+08:00", utcLoc)
	//utcTime = now.In(utcLoc)


// TODO do xxx to UTC conversion here
	err = pRes.WriteAsJson(fmt.Sprintf("current time: %v, timezone: %v, diff in min: %v; utc time => %v\n", now, tzName, tzOffsetInMin, utcTime))
	if err != nil {
		// just log and continue to serve (sometimes it is a disconnection which could be re-covered)
		c.logError("listTimeCronAPI", err.Error())
	}
}

func (c *StructCron) logError(funcName string, msg string) {
	ctfmt.Print(ct.Red, true, fmt.Sprintf("[%v%v] ", moduleWSCron, funcName))
	ctfmt.Println(ct.White, true, msg)
}