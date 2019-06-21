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
	"Stockbinator/store"
	"testing"
)

func TestAAStocksCrawlerCrawl(t *testing.T)  {
	// check if test needs to be run
	if !*pFlagAAStocksCrawler && !*pFlagCrawler {
		t.SkipNow()
	}
	// table test for several different stock-code(s)
	var stockCodes = []struct{
		moduleKey string
	}{
		{ "stock_aastocks.700_tencent" },
		{ "stock_aastocks.857_petrol_china_oil" },
		{ "stock_aastocks.1299_aia" },
		{ "stock_aastocks.939_construction_bank_cn" },
	}

	storeList := make([]store.IStore, 0)
	if FileStore != nil {
		storeList = append(storeList, FileStore)
	}
	for _, sCode := range stockCodes {
		err := instanceStructCrawlerTestObjects.pCrawlerAAStocks.Crawl(sCode.moduleKey, storeList)
		if err != nil {
			t.Errorf("[TestAAStocksCrawlerCrawl] exception: %v", err)
		}
	}
}