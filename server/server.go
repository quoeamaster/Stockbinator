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
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
	"github.com/emicklei/go-restful"
	"net/http"
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

// TODO
func (s *Server) setupCrons() (err error) {
	// get back all stock_module's rules' collection time
	for _, stockModuleObj := range s.pCfg.ModuleConfigs {
		fmt.Println(stockModuleObj.Name, "%%%")
		mapToplevelRules := stockModuleObj.Rules.Map()
		for keySub, ruleVal := range mapToplevelRules {
			mapSubRules := ruleVal.(map[string]interface{})
			collectTime := mapSubRules["collect_time"].(string)
			// direct call the api and not through http
			// e.g. 21:00T+08:00 => hour24:21, min:00, sec:00, timezone:+08:00, stockModuleRule:
			ruleKey := fmt.Sprintf("%v.%v", stockModuleObj.Name, keySub)

			// TODO break the 21:00T+08:00 into hour24, min, sec, timezone....

			fmt.Println(ruleKey, " ", collectTime)
		}
	}

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
	s.pCronSrv = webservice.NewStructCron()
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