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
package main

import (
	"Stockbinator/server"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/daviddengcn/go-colortext/fmt"
)

const moduleMain = "main."

func main() {
	pSvr := new(server.Server)
	err := pSvr.Start()
	if err != nil {
		logInfo("bootstrap", fmt.Sprintf("exception!!!! => %v", err))
		logInfo("bootstrap", "server startup failure and exited successfully")
	}
}


// logging function for info level
func logInfo(funcName string, msg string) {
	ctfmt.Print(ct.Red, true, fmt.Sprintf("[%v%v] ", moduleMain, funcName))
	ctfmt.Println(ct.White, true, msg)
}

