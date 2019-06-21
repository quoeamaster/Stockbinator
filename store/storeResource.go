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

// type constants
const (
	TypeInteger = iota
	TypeFloat
	TypeString
	TypeBool
	TypeDate
)

// code constants
const (
	CodeSuccess = iota
	CodeFailure
	CodeUnknownError
)

// structure to describe store-value for persistence
type StructStoreValue struct {
	// actual value
	Value    interface{}
	// data type
	Type     int
	// is it array
	IsArray  bool
	// is it object or map
	IsObject bool
	// key / fieldname representing the value (if not given AND was passed through a
	// map[string]StructStoreValue; then use that key from the map)
	Key      string
}

// helper method to construct a StructStoreValue object
func NewStructStoreValue(key string, value interface{}, valueType int, isArray, isObject bool) (val *StructStoreValue) {
	val = new(StructStoreValue)
	val.Key = key
	val.Value = value
	val.Type = valueType
	val.IsObject = isObject
	val.IsArray = isArray
	return
}

// structure to describe a response from store implementor
type StructStoreResponse struct {
	Code    int
	Message string
}

// interface for "store"
type IStore interface {
	// save all data into the store
	Persist(data map[string]StructStoreValue) (response StructStoreResponse, err error)

	// read all contents from the store (might be an issue when the content size is HUGE
	ReadAll() (response StructStoreResponse, content string, err error)
	// read only the content associated by the KEY, PARAMS contains additional information for the read operation
	ReadByKey(key string, params interface{}) (response StructStoreResponse, value StructStoreValue, err error)

	// modify the value associated with the key
	ModifyByKey(key string, value StructStoreValue) (response StructStoreResponse, err error)

	// remove value associated with the key
	RemoveByKey(key string) (response StructStoreResponse, valueRemoved StructStoreValue, err error)
	// remove all data in the store, be careful~
	RemoveAll() (response StructStoreResponse, err error)
}