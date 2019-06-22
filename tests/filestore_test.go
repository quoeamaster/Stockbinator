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
	"Stockbinator/util"
	"fmt"
	"strings"
	"testing"
)


func TestFilestorePersistFlow(t *testing.T) {
	if !*pFlagFilestore {
		t.SkipNow()
	}
	LogTestOutput("TestFilestorePersistFlow", "** start test **")
	LogTestOutput("TestFilestorePersistFlow", "a. remove test file contents")
	resp, err := FileStore.RemoveAll()
	if err != nil {
		t.Fatal(fmt.Sprintf("could not remove the entries from filestore instance: %v", err))
	}
	helperFilestoreFlowsCommonResponseHandler(resp, t)

	LogTestOutput("TestFilestorePersistFlow", "b. add a few entries to the test file")
	// write 1st 4 entries; leave the last entry for modification
	for i := 0; i<len(TestStoreEntriesList)-1; i++ {
		entryMap := TestStoreEntriesList[i]
		// LogTestOutput("TestFilestorePersistFlow", fmt.Sprintf("index %v => %v", i, entryMap))
		resp, err = FileStore.Persist(entryMap)
		if err != nil {
			LogTestOutput("TestFilestorePersistFlow", fmt.Sprintf("persist entry failed at index %v", i))
			t.Fatal(err.Error())
		}
		helperFilestoreFlowsCommonResponseHandler(resp, t)
	}

	LogTestOutput("TestFilestorePersistFlow", "c. retrieve all contents just written")
	resp, contents, err := FileStore.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Code != store.CodeSuccess {
		LogTestOutput("TestFilestorePersistFlow",
			fmt.Sprintf("read operation returned non success code - %v; msg > %v", resp.Code, resp.Message))
		t.Fatal(fmt.Sprintf("read operation returned non success code - %v; msg > %v", resp.Code, resp.Message))
	}
	if util.IsEmptyString(contents) {
		t.Fatal("expected contents READ should be non EMPTY!")
	}
	// kind of weird... there is 1 additional line feed... (???)
	lines := strings.Split(contents, "\n")
	if lines == nil || len(lines) != (4 + 1) {
		for i, line := range lines {
			LogTestOutput("TestFilestorePersistFlow", fmt.Sprintf("%v - %v", i, line))
		}
		t.Fatal(fmt.Sprintf("expecting contents to be in 4 lines of data BUT got %v", len(lines)))
	}






	LogTestOutput("TestFilestorePersistFlow", "** end test **\n")
}

func helperFilestoreFlowsCommonResponseHandler(r store.StructStoreResponse, t *testing.T) {
	if r.Code != store.CodeSuccess {
		t.Fatal(fmt.Sprintf("no Error previously HOWEVER there is a non SUCCESS code and message from response =>[%v] %v", r.Code, r.Message))
	}
}

