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
package config

import (
	"Stockbinator/util"
	"fmt"
	"github.com/micro/go-config"
)



// config file description:
// a) app.conf - application (Stockbinator) related configurations; applied to all stock sources
// b) xxx (folder) - a set of config files related to stock source "xxx" (e.g. aastocks)
//		1) holiday.conf => sets of holidays where stock data should not be available
//		2) rules.conf => sets of rules for capturing a specific metric about a particular stock + a list of targeted stocks url (html containing the metrics)


// env variable holding the config repository's location (could be nil)
const appToml = "app.toml"
const stockModuleFolderPrefix = "stock_"
const holidayToml = "holiday.toml"
const rulesToml = "rules.toml"

type StructConfig struct {
	// where the config files are located
	RepositoryPath string

	// the actual config(s) loaded
	AppConfig config.Config
	// map of stock module config(s)
	ModuleConfigs map[string]StructStockModuleConfig
}

// creation method for StructConfig instance
func NewStructConfig() (*StructConfig, error) {
	pInst := new(StructConfig)
	pInst.ModuleConfigs = make(map[string]StructStockModuleConfig)

	err := pInst.LoadConfigs()
	if err != nil {
		return nil, err
	}
	return pInst, nil
}

// logic =>
// 1) check if env variable "SB_CONFIG_PATH" exists; if exists => use that folder as repositories for the config file(s)
// 2) if 1) is not available, use the current executable path's config folder =>
// 	e.g. /exec/config (if exec is where the Stockbinator executable lives)
func (s *StructConfig) LoadConfigs() (err error) {
	repoPath, err := util.GetConfigFolderPath()
	if err != nil {
		return
	}
	s.RepositoryPath = repoPath

	// 2a) load app.conf
	err = s.loadAppConf()
	if err != nil {
		return
	}
	// 2b) load stock module conf
	err = s.loadStockModuleConf()

	return
}

// load the application's config(s)
func (s *StructConfig) loadAppConf() (err error) {
	s.AppConfig, err = util.LoadConfig(fmt.Sprintf("%v%v", s.RepositoryPath, appToml))
	return
}


type StructStockModuleConfig struct {
	// stock module's name (might be redundant in this case)
	Name string
	// config about holidays of this stock module
	// (diff stock modules should be originated in diff country and hence should have its own set of public holidays)
	Holidays config.Config
	// rules for the metrics collection; remember could have multiple rule(s) on the same stock modules
	// (e.g. ALPHABET and IBM are the target metrics on Nasdaq)
	Rules config.Config
}


// load the stock module's config(s) - e.g. holiday.toml and rules.toml
func (s *StructConfig) loadStockModuleConf() (err error) {
	// scan for the "modules"; all modules start with a folder prefix "stock_"
	// (check constant stockModuleFolderPrefix)
	arrPaths, err := util.ListMatchingFoldersUnderPath(s.RepositoryPath, stockModuleFolderPrefix)
	if err != nil {
		return
	}
	// loop through all modules and read the corresponding config / toml etc
	for _, modulePath := range arrPaths {
		pCfgInst := new(StructStockModuleConfig)
		pCfgInst.Name, err = util.GetLastPathElement(modulePath)
		if err != nil {
			return
		}
		// holiday.toml
		pCfgInst.Holidays, err = util.LoadConfig(util.ConcatenateFilePath(modulePath, holidayToml))
		if err != nil {
			return
		}
		//fmt.Println( pCfgInst.Holidays.Get("2019", "holidays").StringSlice([]string{}))

		// rulesToml
		pCfgInst.Rules, err = util.LoadConfig(util.ConcatenateFilePath(modulePath, rulesToml))
		if err != nil {
			return
		}
		//fmt.Println(pCfgInst.Rules.Get("700_tencent", "rule_dividend_yield").String("something's wrong..."))

		s.ModuleConfigs[pCfgInst.Name] = *pCfgInst
	}
	return
}


