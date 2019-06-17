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
	"Stockbinator/config"
	"Stockbinator/crawler"
	"Stockbinator/store"
	"flag"
	"fmt"
	config2 "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"os"
	"testing"
)

const stockModuleKey = "stock_aastocks.700_tencent"
const stockModuleName = "stock_aastocks"

// flag(s)
var (
	pFlagCrawler = flag.Bool("crawler", false, "run all crawler test")
	pFlagAAStocksCrawler = flag.Bool("crawler.aastocks", false, "run ONLY aastocks crawler test")
	pFlagGenericCrawler = flag.Bool("crawler.generic", false, "run ONLY generic crawler test")

	pFlagCommonUtil = flag.Bool("util.common", false, "run ONLY common-util test")
	pFlagCrawlerUtil = flag.Bool("util.crawler", false, "run ONLY crawler-util test")

	pFlagFilestore = flag.Bool("store.file", false, "run ONLY filestore test")

	// flag indicating logging feature
	pFlagLog = flag.Bool("log", false, "display logs about the test")
)

// Testing method
func TestMain(m *testing.M) {
	fmt.Println("************** setup in progress **************")

	var err error
	// parse command line flags (plus arguments if any)
	flag.Parse()

	// crawler test related
	if *pFlagCrawler {
		err = setupStockModuleConfig()
		handlerCommonError(err)

		err = setupCrawlerTestObjects()
		handlerCommonError(err)

		setupCrawlerAAStocks()
		setupCrawlerGeneric()
	} else if *pFlagAAStocksCrawler {
		err = setupStockModuleConfig()
		handlerCommonError(err)

		err = setupCrawlerTestObjects()
		handlerCommonError(err)

		setupCrawlerAAStocks()
	} else if *pFlagGenericCrawler {
		err = setupStockModuleConfig()
		handlerCommonError(err)

		err = setupCrawlerTestObjects()
		handlerCommonError(err)

		setupCrawlerGeneric()
	}

	// util series

	if *pFlagCommonUtil {
		err = setupStockModuleConfig()
		handlerCommonError(err)

	}
	if *pFlagCrawlerUtil {
		err = setupStockModuleConfig()
		handlerCommonError(err)

		err = setupHoliday2018Config()
		handlerCommonError(err)
	}

	// filestore series
	if *pFlagFilestore {
		err = setupSharableStore()
		handlerCommonError(err)

		setupFilestore()
	}

	// TODO: add other test setup

	code := m.Run()

	if *pFlagCrawler || *pFlagAAStocksCrawler || *pFlagGenericCrawler {
		err = teardownCrawlerTestObject()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	// TODO: add other test teardown

	os.Exit(code)
}

// handler common error (e.g. print out error message and Exit process with -1 code)
func handlerCommonError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// common shared objects or variables

// Crawler related test object (e.g pointers pointing the actual crawler)
type StructCrawlerTestObjects struct {
	ConfigMap map[string]config.StructStockModuleConfig

	// crawlers come into different fashions, each crawler maintains its rules on crawling / scrapping
	pCrawlerAAStocks *crawler.StructAAStocksCrawler
	pCrawlerGenric *crawler.StructGenericCrawler
}
var instanceStructCrawlerTestObjects *StructCrawlerTestObjects

var SharableStockModuleConfig config.StructStockModuleConfig

var Holidays2018Config config2.Config

var SharableStoreConfig config2.Config
var FileStore store.IStore

// corresponding setupXXX methods

func setupSharableStore() (err error) {
	// setup the config
	SharableStoreConfig = config2.NewConfig()
	cfgPath := file.NewSource(file.WithPath("../config/app_test.toml"))
	err = SharableStoreConfig.Load(cfgPath)
	handlerCommonError(err)

	return
}

func setupFilestore()  {
	// using the default store's filename...
	FileStore = store.NewStructFilestore(SharableStoreConfig, common.StoreDefaultDateFilename)
}

func setupStockModuleConfig() (err error) {
	// config object (could mock in the future)
	pRuleConfig := config2.NewConfig()
	fSrc := file.NewSource(file.WithPath("../config/rules.toml"))
	err = pRuleConfig.Load(fSrc)
	if err != nil {
		fmt.Printf("could not load the rules.toml, %v\n", err)
		return
	}

	pHoliday := config2.NewConfig()
	fSrc = file.NewSource(file.WithPath("../config/holiday.toml"))
	err = pHoliday.Load(fSrc)
	if err != nil {
		fmt.Printf("could not load the holiday.toml, %v\n", err)
		return
	}
	// stocks rule
	pStockModuleConfig := new(config.StructStockModuleConfig)
	pStockModuleConfig.Rules = pRuleConfig
	pStockModuleConfig.Holidays = pHoliday

	SharableStockModuleConfig = *pStockModuleConfig
	return
}

func setupHoliday2018Config() (err error) {
	Holidays2018Config = config2.NewConfig()
	fSrc := file.NewSource(file.WithPath("../config/holidays_2018.toml"))
	err = Holidays2018Config.Load(fSrc)
	if err != nil {
		fmt.Printf("could not load the holidays_2018.toml, %v\n", err)
		return
	}
	return
}

func setupCrawlerTestObjects() (err error)  {
	// add setup code here
	instanceStructCrawlerTestObjects = new(StructCrawlerTestObjects)
	instanceStructCrawlerTestObjects.ConfigMap = make(map[string]config.StructStockModuleConfig)
	instanceStructCrawlerTestObjects.ConfigMap[stockModuleName] = SharableStockModuleConfig

	return
}
func setupCrawlerAAStocks() {
	// crawler
	instanceStructCrawlerTestObjects.pCrawlerAAStocks = crawler.NewStructAAStocksCrawler(instanceStructCrawlerTestObjects.ConfigMap)
}
func setupCrawlerGeneric()  {
	instanceStructCrawlerTestObjects.pCrawlerGenric = crawler.NewStructGenericCrawler(instanceStructCrawlerTestObjects.ConfigMap)
}


func teardownCrawlerTestObject() (err error) {
	return
}


// * **************** *
// * common functions *
// * **************** *

// loggoing function, only logs when the "log" flag is passed
func LogTestOutput(testName, message string) {
	if *pFlagLog {
		fmt.Printf("[%v] %v\n", testName, message)
	}
}
