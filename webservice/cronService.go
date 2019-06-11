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
	"Stockbinator/config"
	"Stockbinator/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
	"github.com/emicklei/go-restful"
	"time"
)

const moduleWSCron = "cronService"
// the time / date layout
const cronTimeLayout = "2006-01-02T15:04:05-0700"

type StructCron struct {
	// instance of the jsonParser for this webservice module
	pJsonParser *util.StructJsonParser
	// map of StructCronEntry entries; each of these represent a time for running a crawler job
	cronTimeEntries map[string]*StructCronEntry
	// the config information for the crawler job
	pCfg *config.StructConfig
}

// creation method for StructCron
func NewStructCron(pCfg *config.StructConfig) (cron *StructCron) {
	cron = new(StructCron)
	cron.pJsonParser = util.NewStructJsonParser()
	cron.cronTimeEntries = make(map[string]*StructCronEntry)
	cron.pCfg = pCfg
	return
}

// structure contains the display-name of the cron-time entry;
// plus a UTC converted time object,
// finally a list of stocks-module-rule(s) associated
type StructCronEntry struct {
	// display name for the cron-time; it could be in any timezone
	DisplayName string
	// UTC converted time / date
	UTCTime time.Time
	// list of stocksModuleRule under this cron-time entry (usually size of 1)
	StocksModuleRuleList []string
}

func NewStructCronEntry() (entry *StructCronEntry) {
	entry = new(StructCronEntry)
	entry.StocksModuleRuleList = make([]string, 0)
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
		// a) prepare the cron-time for today
		cronTime := util.CreateTodayTargetTimeByHourMinTimezone(hour24, min, timezone)

		// b) check if the entry is already there or not
		if c.cronTimeEntries[cronTime] != nil {
			// means existed
			pEntry := c.cronTimeEntries[cronTime]
			pEntry.StocksModuleRuleList = append(pEntry.StocksModuleRuleList, stocksModuleRule)

		} else {
			// create a new entry
			pDate, err2 := util.ParseStringDateToTodayUTC(hour24, min, timezone)
			if err2 != nil {
				return
			} else {
				pEntry := NewStructCronEntry()
				pEntry.DisplayName = cronTime
				pEntry.UTCTime = pDate
				pEntry.StocksModuleRuleList = append(pEntry.StocksModuleRuleList, stocksModuleRule)
				c.cronTimeEntries[cronTime] = pEntry
			}
		} // end -- if (cronTimeEntries exists check)
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
// TODO do xxx to UTC conversion here
	err := pRes.WriteAsJson(c.cronTimeEntries)
	if err != nil {
		// just log and continue to serve (sometimes it is a disconnection which could be re-covered)
		c.logError("listTimeCronAPI", err.Error())
	}
}

func (c *StructCron) logError(funcName string, msg string) {
	ctfmt.Print(ct.Red, true, fmt.Sprintf("[%v%v] ", moduleWSCron, funcName))
	ctfmt.Println(ct.White, true, msg)
}