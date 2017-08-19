/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//慈善智能合约

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"strings"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//新增走访记录信息
//type Ncharity struct{
//CharityHash string
//VisitInf Visit
//}
//走访机构信息
type Visit struct {
	Organization string `json:"Organization"` //慈善机构名称
	Result       string `json:"Result"`       //走访结果
	VTime        string `json:"VTime"`        //走访时间
	Comment      string `json:"Comment"`      //备注
}

//捐助情况
type Sum struct {
	SOrganization string `json:"SOrganization"`
	Money         string `json:"Money"`
	Reason        string `json:"Reason"`
	STime         string `json:"STime"`
}

//慈善信息结构
type ChariInf struct {
	CharityHash string  `json:"CharityHash"` //所有信息的hash
	Name        string  `json:"Name"`        //姓名
	//TotalSum    string  `json:"TotalSum"`    //慈善捐助总金额
	VisitInf    []Visit `json:"VisitInf"`    //走访信息
	ChSum       []Sum   `json:"ChSum"`       //慈善机构捐助具体信息
}

//输入参数：“init”，two function:add,update.[0]是操作人员ID
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	return nil, nil
}

//输入参数：function：“add(或update)”，[0]是接受慈善帮助人员的ID，[1]是json格式，[2]是操作人ID
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	IdentID := args[0]
	Hashval := args[1]
	//Hashval = strings.Replace(Hashval, "\\", "", -1)
	// Hashval := string(args[1])
	var chariInf ChariInf
	var visitInf Visit //输入的参数是[1]:{"Organization": "慈善机构", "Result": "捐助","VTime":"2017-01-05","Comment":"无"}
	var err error
	// Perform the execution

	// Write the state back to the ledger
	if function == "add" {
		// Perform the execution
		TempHashval, err := stub.GetState(IdentID)

		if TempHashval != nil {
			return nil, errors.New("This ID already exists")
		}
		// Write the state back to the ledger
		err = stub.PutState(IdentID, []byte(Hashval))
		if err != nil {
			return []byte(Hashval), err
		}
	}
	//更新信息
	if function == "update" {
		HashvalTemp, errs := stub.GetState(IdentID)

		if errs != nil {
			return nil, errors.New("list is not here")
		}
		if HashvalTemp == nil {
			return nil, errors.New("Entity not found")
		}

		err = stub.PutState(IdentID, []byte(Hashval))
		if err != nil {
			return nil, err
		}
	}
	//新增走访记录
	if function == "addVisit" {
		HashvalTemp, err := stub.GetState(IdentID)

		if err != nil {
			return nil, errors.New("list is not here")
		}
		if HashvalTemp == nil {
			return nil, errors.New("Entity not found")
		}

		// charT := Visit{
		// 	Organization: "china",    //慈善机构名称
		// 	Result:       "help",     //走访结果
		// 	VTime:        "20170101", //走访时间
		// 	Comment:      "no",       //备注
		// }

		json.Unmarshal(HashvalTemp, &chariInf)
		json.Unmarshal([]byte(Hashval), &visitInf)

		//修改哈希

		//chariInf.CharityHash = visitInf.CharityHash
		//新增走访记录
		chariInf.VisitInf = append(chariInf.VisitInf, visitInf)

		jsonchari, _ := json.Marshal(chariInf)
		err = stub.PutState(IdentID, []byte(jsonchari))
		if err != nil {
			return jsonchari, err
		}
	}

	return nil, nil

}

// 查询，输入参数:[0]是接受慈善对象的身份证号
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var IdentID string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	IdentID = args[0]

	// Get the state from the ledger
	Hashval, err := stub.GetState(IdentID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + IdentID + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Hashval == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + IdentID + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + IdentID + "\",\"Amount\":\"" + string(Hashval) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Hashval, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

