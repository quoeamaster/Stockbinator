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
}


// load server and stock module config(s)
func (s *Server) loadConfig() (err error) {
	s.pCfg, err = config.NewStructConfig()
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
	// start Http service
	s.logInfo("Start", "server started at port => 9000")
	err = http.ListenAndServe(":9000", nil)

	return
}

func (s *Server) loadWebServices() (err error) {
	// add webService module
	pWs := new(restful.WebService)
	pWs.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	pWs.Route(pWs.GET("").To(s.welcomeFunc))

	restful.DefaultContainer.Add(pWs)

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