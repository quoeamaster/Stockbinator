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
package store

import (
	"Stockbinator/common"
	"errors"
	"fmt"
	"github.com/micro/go-config"
	"strings"
)

// internal storage (map) for the store implementation(s)
var storeCache map[string]IStore

// get store implementation by key (e.g. filestore);
// if the store instance is not available, try its best to create the store
// and save to cache by the given key.
// Parameter config is the config read from app.toml
// Optional parameter params which is a map of object(s)
func GetStoreByKey(key string, config config.Config, params... map[string]interface{}) (store IStore, err error)  {
	// TODO singleton??? thunder herd should not happen normally
	// https://news.ycombinator.com/item?id=1722213
	if storeCache == nil {
		storeCache = make(map[string]IStore)
	}
	// key MUST be... starting with either "filestore.x.y" or "datastore.x.y"
	store = storeCache[key]
	if store == nil {
		// try to create store that are recognizable
		isFilestore := strings.Index(key, common.ConfigKeyStoreFile) == 0
		isDatastore := strings.Index(key, common.ConfigKeyStoreData) == 0

		if isFilestore {
			filename := common.StoreDefaultDateFilename
			filename = key[len(common.ConfigKeyStoreFile)+1:]
			// PS. might need to use params... but in this nope
			//if params != nil {
			//	interfaceFilename := params[0][common.StoreKeyDefaultDateFilename]
			//	if interfaceFilename != nil {
			//		filename = interfaceFilename.(string)
			//	} // end -- if (interface-filename from map is non nil)
			//}
fmt.Println("*** filename", filename)
fmt.Println("*** key", key)
			store = NewStructFilestore(config, filename)
			storeCache[key] = store
		} else if isDatastore {
			// TODO: add other store(s) like datastore
		} else {
			store = nil
			err = errors.New(fmt.Sprintf("unknown Store type => %v", key))
		}
	} // end -- if (store is nil, create one if recognized)
	return
}