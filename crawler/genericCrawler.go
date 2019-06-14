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
	"errors"
	"fmt"
	config2 "github.com/micro/go-config"
	"io/ioutil"
	"net/http"
	"strings"
)

const valueUnknown = "UnknowN"
const ruleUrl = "url"
const rulePrice = "rule_price"
const ruleVolume = "rule_volume"
const ruleValue = "rule_value"
const rulePe = "rule_pe"
const ruleDividendYield = "rule_dividend_yield"
const rulePb = "rule_pb"
const ruleValueFlow = "rule_value_flow"
const ruleTurnoverRate = "rule_turnover_rate"
const ruleHandPerShare = "rule_hand_per_share"

type StructGenericCrawler struct {
	// inject the stock module config rules (map)
	StockModuleConfig map[string]config.StructStockModuleConfig
}

// constructor for Generic Crawler
func NewStructGenericCrawler(config map[string]config.StructStockModuleConfig) (pCrawler *StructGenericCrawler)  {
	pCrawler = new(StructGenericCrawler)
	if config != nil {
		pCrawler.StockModuleConfig = config
	}
	return
}

// 1) read config file for rule(s) to crawl (at least url and patterns to match) based on moduleKey (stock_module-rule)
// 2) based on the key above, invoke the corresponding crawler's crawl method and scrap out the values
// 3) output the results to a repository (filestore by default or any datastorage tech e.g. elasticsearch)
func (s *StructGenericCrawler) Crawl(moduleKey string) (err error) {
	// break the moduleKey back the moduleName and stockName
	names := strings.Split(moduleKey, ".")
	if names != nil && len(names) == 2 {
		stockModuleConfig := s.StockModuleConfig[names[0]]
		// download the html content for crawl / scrap
		ruleConfig := stockModuleConfig.Rules
		url := ruleConfig.Get(names[1], ruleUrl).String(valueUnknown)
		if strings.Compare(url, valueUnknown) == 0 {
			err = errors.New("url is not available~ can NOT retrieve content for crawling")
			return
		}
		// forward url for content crawl / scrap
		urlContent, err2 := s.geUrlContent(url)
		if err2 != nil {
			err = err2
			return
		}

		valPrice, err2 := s.crawlForRule(ruleConfig, urlContent, names[1], rulePrice)
		if err2 != nil {
			err = err2
			return
		}
		fmt.Println(valPrice)


	} else {
		err = errors.New("invalid moduleKey, it should be [STOCKS_MODULE_NAME][STOCK_CODE_UNDER_THE_MODULE]")
	}
	return
}

// retrieve content from url
func (s *StructGenericCrawler) geUrlContent(url string) (contentInString string, err error)  {
	resp, err := http.Get(url)
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			fmt.Println("could not close http connection after fetching data from url")
			err = err2
		}
	}()
	if err != nil {
		return
	}
	bContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	contentInString = string(bContent)

	return
}

func (s *StructGenericCrawler) crawlForRule(ruleConfig config2.Config, content, stockCode, rule string) (value string, err error) {
	ruleDef := ruleConfig.Get(stockCode, rule).String(valueUnknown)
	if strings.Compare(ruleDef, valueUnknown) == 0 {
		err = errors.New(fmt.Sprintf("rule-defintion not found: %v", rule))
		return
	}
	value, err = s.getValueFromRule(content, ruleDef)

	return
}

func (s *StructGenericCrawler) getValueFromRule(content, ruleDef string) (value string, err error) {
	idx := strings.Index(content, "337")
	if idx == -1 {
		err = errors.New("content does NOT contain the rule-definition")
	}
	fmt.Println(content)
	fmt.Println(idx)
	fmt.Println(ruleDef)


	return
}
