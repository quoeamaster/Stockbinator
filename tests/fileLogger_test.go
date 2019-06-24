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
	"Stockbinator/common"
	"Stockbinator/logger"
	"Stockbinator/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFileLoggerFlow_01(t *testing.T) {
	if !*pFlagFileLogger {
		t.SkipNow()
	}
	LogTestOutput("TestFileLoggerFlow_01", "** start test **")
	// a. create config struct
	filename := SharableStoreConfig.Get(common.ConfigKeyStoreFile, common.ConfigKeyRepo).String("")
	envMatches, _, err := util.ParseEnvVar(filename)
	if err != nil {
		t.Fatal(fmt.Sprintf("trying to parse the env vars whilst got exception => %v", err))
	}
	if len(envMatches) > 0 {
		for idx := len(envMatches) - 1; idx >= 0; idx-- {
			envMatch := envMatches[idx]
			envMatchEnv := envMatch[1:len(envMatch)-1]
			filename = strings.Replace(filename, envMatch, os.Getenv(envMatchEnv), 1)
		}
	}
	LogTestOutput("TestFileLoggerFlow_01", fmt.Sprintf("final filepath => %v", filename))

	pCfg := logger.NewStructLoggerConfig(
		filename, "file.log",true, "yy-mm-dd")
	// b. create the logger
	pLogger := logger.NewStructFileLogger()
	defer func() {
		if pLogger != nil {
			pLogger.Close()
		}
	}()
	err = pLogger.SetupLogger(pCfg, nil)
	if err != nil {
		t.Fatal(fmt.Sprintf("cannot setup logger => %v", err))
	}
	// c. write logs (without prefix)
	pLogger.Println("this is demo line of log ($#$%^&**) with special characters")
	pLogger.Print("1st part of the Log-line")
	pLogger.Printf(", 2nd part of the log-line [%v] {%v}\n", "some-DATA", "some-MORE")

	infoMap := make(map[string]string)
	err = json.Unmarshal([]byte(pLogger.Info()), &infoMap)
	if err != nil {
		t.Fatal(err)
	}
	LogTestOutput("TestFileLoggerFlow_01", fmt.Sprintf("%v => filepath = %v", infoMap, infoMap[common.LoggerFileInfoKeyFilepath]))
	pLogFile, err := os.Open(infoMap[common.LoggerFileInfoKeyFilepath])
	defer func() {
		if pLogFile != nil {
			pLogFile.Close()
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	bContents, err := ioutil.ReadAll(pLogFile)
	if err != nil {
		t.Fatal(err)
	}
	sArrContents := string(bContents)
	if len(sArrContents) <= 0 {
		LogTestOutput("TestFileLoggerFlow_01", fmt.Sprintf("actual content => [%v]", sArrContents))
		t.Fatal("should HAVE some contents....")
	}

	// d. cleanup the contents...
	pLogFile2, err := os.OpenFile(infoMap[common.LoggerFileInfoKeyFilepath], os.O_WRONLY|os.O_TRUNC, 0660)
	defer func() {
		if pLogFile2 != nil {
			pLogFile2.Close()
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	err = pLogFile2.Truncate(0)
	if err != nil {
		t.Fatal(err)
	}


	LogTestOutput("TestFileLoggerFlow_01", "** end test **\n")
}
