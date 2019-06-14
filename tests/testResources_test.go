package tests

import (
	"Stockbinator/config"
	"Stockbinator/crawler"
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
)

// Testing method
func TestMain(m *testing.M) {
	fmt.Println("************** setup in progress **************")

	var err error
	// parse command line flags (plus arguments if any)
	flag.Parse()

	// crawler test related
	if *pFlagCrawler {
		err = setupCrawlerTestObjects()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		setupCrawlerAAStocks()
		setupCrawlerGeneric()
	} else if *pFlagAAStocksCrawler {
		err = setupCrawlerTestObjects()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		setupCrawlerAAStocks()
	} else if *pFlagGenericCrawler {
		err = setupCrawlerTestObjects()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		setupCrawlerGeneric()
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



// common shared objects or variables

// Crawler related test object (e.g pointers pointing the actual crawler)
type StructCrawlerTestObjects struct {
	ConfigMap map[string]config.StructStockModuleConfig

	// crawlers come into different fashions, each crawler maintains its rules on crawling / scrapping
	pCrawlerAAStocks *crawler.StructAAStocksCrawler
	pCrawlerGenric *crawler.StructGenericCrawler
}
var instanceStructCrawlerTestObjects *StructCrawlerTestObjects

// corresponding setupXXX methods

func setupCrawlerTestObjects() (err error)  {
	// add setup code here
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

	instanceStructCrawlerTestObjects = new(StructCrawlerTestObjects)
	instanceStructCrawlerTestObjects.ConfigMap = make(map[string]config.StructStockModuleConfig)
	instanceStructCrawlerTestObjects.ConfigMap[stockModuleName] = *pStockModuleConfig

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