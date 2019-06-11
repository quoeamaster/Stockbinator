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
package server

import (
	"Stockbinator/config"
	"Stockbinator/webservice"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"strings"
)

const moduleServer = "server."

// Server struct that runs a REST api layer
type Server struct {
	// server config
	pCfg *config.StructConfig
	// Cron service
	pCronSrv *webservice.StructCron
}


// load server and stock module config(s)
func (s *Server) loadConfig() (err error) {
	s.pCfg, err = config.NewStructConfig()
	return
}

// setup the initial cron schedules based on the stock-module's config(s)
func (s *Server) setupCrons() (err error) {
	// get back all stock_module's rules' collection time
	for _, stockModuleObj := range s.pCfg.ModuleConfigs {
		mapToplevelRules := stockModuleObj.Rules.Map()
		for keySub, ruleVal := range mapToplevelRules {
			mapSubRules := ruleVal.(map[string]interface{})
			collectTime := mapSubRules["collect_time"].(string)
			// direct call the api and not through http
			ruleKey := fmt.Sprintf("%v.%v", stockModuleObj.Name, keySub)
			//fmt.Println(ruleKey, " ", collectTime)

			// e.g. 21:00T+08:00 => hour24:21, min:00, sec:00, timezone:+08:00, stockModuleRule:
			// break the 21:00T+08:00 into hour24, min, sec, timezone....
			hour24, min, sec, timezone, err := s.prepareCronScheduleParams(collectTime)
			if err != nil {
				panic(err)
			}
			_, err = s.pCronSrv.UpsertTimeCron(hour24, min, sec, timezone, ruleKey)
			if err != nil {
				panic(err)
			}
		}	// end -- for (sub level rule values e.g. [007_tencet][collect_time])
	}	// end -- for (top level rule values e.g. [stockmodule_aastocks])
	return
}

// method to break the 21:00T+08:00 into hour24, min, sec, timezone....
// since this parsing is application dependant, hence not available under the commonUtil.go
func (s *Server) prepareCronScheduleParams(cronVal string) (hours24, min, sec int, timezone string, err error) {
	topParts := strings.Split(cronVal, "T")
	if len(topParts) != 2 {
		err = errors.New(fmt.Sprintf("exception! Invalid cron value => %v, exepcted cron value is [21:00T+07:00]", cronVal))
		return
	}
	// time parts
	timeParts := strings.Split(topParts[0], ":")
	if len(timeParts) != 2 {
		err = errors.New(fmt.Sprintf("exception! Invalid time value => %v, exepcted time value is [21:00]", topParts[0]))
		return
	}
	hours24, err = strconv.Atoi(timeParts[0])
	if err != nil {
		return
	}
	min, err = strconv.Atoi(timeParts[1])
	if err != nil {
		return
	}
	sec = 0
	// timezone
	// regexp check??? (ignore for now)
	timezone = topParts[1]

	return
}

// Start the REST server
func (s *Server) Start() (err error) {
	// load config
	err = s.loadConfig()
	if err != nil {
		return
	}
	// load the webservice(s)
	err = s.loadWebServices()
	if err != nil {
		return
	}
	// setup cron(s)
	err = s.setupCrons()
	if err != nil {
		return
	}
	// start Http service
	s.logInfo("Start", "server started at port => 9000")
	err = http.ListenAndServe(":9000", nil)

	return
}

func (s *Server) loadWebServices() (err error) {
	// add webService module (dummy tester)
	pWs := new(restful.WebService)
	pWs.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	pWs.Route(pWs.GET("").To(s.welcomeFunc))

	restful.DefaultContainer.Add(pWs)

	// load CronService module
	s.pCronSrv = webservice.NewStructCron(s.pCfg)
	restful.DefaultContainer.Add(s.pCronSrv.CreateWebservice())

	return
}

// dummy welcome function for testing the server
func (s *Server) welcomeFunc(pReq *restful.Request, pRes *restful.Response) {
	err := pRes.WriteAsJson("welcome using the server")
	if err != nil {
		panic(err)
	}
}

// logging function for info level
func (s *Server) logInfo(funcName string, msg string) {
	ctfmt.Print(ct.Green, true, fmt.Sprintf("[%v%v] ", moduleServer, funcName))
	ctfmt.Println(ct.White, true, msg)
}