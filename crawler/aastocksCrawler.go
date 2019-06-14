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
package crawler

import (
	"Stockbinator/config"
	"Stockbinator/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type StructAAStocksCrawler struct {
	// inject the stock module config rules (map)
	StockModuleConfig map[string]config.StructStockModuleConfig
}

// constructor for Generic Crawler
func NewStructAAStocksCrawler(config map[string]config.StructStockModuleConfig) (pCrawler *StructAAStocksCrawler)  {
	pCrawler = new(StructAAStocksCrawler)
	if config != nil {
		pCrawler.StockModuleConfig = config
	}
	return
}

func (s *StructAAStocksCrawler) Crawl(moduleKey string) (err error) {
	names := strings.Split(moduleKey, ".")
	if names != nil && len(names) == 2 {
		stockModuleConfig := s.StockModuleConfig[names[0]]
// TODO: also check if today is a holiday based on the rule-config
// stockModuleConfig.Holidays
fmt.Println(stockModuleConfig.Holidays.Map())

		ruleConfig := stockModuleConfig.Rules
		url := ruleConfig.Get(names[1], ruleUrl).String(valueUnknown)
		if strings.Compare(url, valueUnknown) == 0 {
			err = errors.New("url is not available~ can NOT retrieve content for crawling")
			return
		}
		// forward url for content crawl / scrap
		urlContent, err2 := util.GetContentFromUrl(url)
		if err2 != nil {
			err = err2
			return
		}

		valPrice, valPriceFluctuations, valTrxAmount, err2 := s.crawlForMetrics(urlContent)
		if err2 != nil {
			err = err2
			return
		}
		fmt.Println(valPrice)
		fmt.Println(valPriceFluctuations)
		fmt.Println(valTrxAmount)


	} else {
		err = errors.New("invalid moduleKey, it should be [STOCKS_MODULE_NAME][STOCK_CODE_UNDER_THE_MODULE]")
	}
	return
}

func (s *StructAAStocksCrawler) crawlForMetrics(content string) (price, priceFluctuations, trxAmount string, err error) {
	//fmt.Println(content)
	pFeedList := make([]structAAStocksFeed, 1)

	err = json.Unmarshal([]byte(content), &pFeedList)
	if err != nil {
		return
	}
	if len(pFeedList) > 0 {
		price = pFeedList[0].Price
		priceFluctuations = pFeedList[0].PriceFluctuation
		trxAmount = pFeedList[0].TrxAmount
	}
	return
}

type structAAStocksFeed struct {
	// a sample => [{"a": "330.000", "b": "<span class='neg'>-4.200(1.257%)</span>", "c": "329.200-336.800", "d": "46.61億", "e": "2019/06/14 16:08"}]
	Price string `json:"a"`
	PriceFluctuation string `json:"c"`
	// need to parse the float value though
	TrxAmount string `json:"d"`
}
