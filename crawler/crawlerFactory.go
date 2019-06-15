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
	"strings"
)

// interface defining crawler behavior.
type InterfaceCrawler interface {
	Crawl(moduleKey string) (err error)
}

// * ************************************* *
// * cache for Interface-crawler instances *
// * ************************************* *

var cacheCrawlersMap map[string]InterfaceCrawler

// return the cached Crawler instance
func getCrawlerByCacheKey(key string) (pCrawler InterfaceCrawler) {
	pCrawler = nil
	if cacheCrawlersMap != nil && len(cacheCrawlersMap) > 0 {
		pCrawler = cacheCrawlersMap[key]
	} else if cacheCrawlersMap == nil {
		cacheCrawlersMap = make(map[string]InterfaceCrawler)
	}
	return
}


// * *********************** *
// * crawler-factory related *
// * *********************** *

func GetCrawler(key string, config map[string]config.StructStockModuleConfig) (pCrawler InterfaceCrawler) {
	pCrawler = getCrawlerByCacheKey(key)
	if pCrawler == nil {
		// try to create an instance if the key is recognizable
		if strings.Index(key, crawlerPrefixAAStocks) != -1 {
			cacheCrawlersMap[key] = NewStructAAStocksCrawler(config)
			pCrawler = cacheCrawlersMap[key]
		}
		// TODO: add other crawler implementations
	}
	return
}


// * *************************** *
// * constants (internal private *
// * *************************** *

const crawlerPrefixAAStocks = "stock_aastocks"
