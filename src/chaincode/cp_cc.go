/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// Address Record
type Location   struct {
	Street   	string   `json:"Street"`
	Unit   		string   `json:"Unit"`
	City 		string   `json:"City"`
	State 		string   `json:"State"`
	Zip 		string   `json:"Zip"`
}

// Check Record
type DataCheck   struct {
	Result   	string   `json:"Status"`
	Message   	string   `json:"Message"`
}


// Employee Record
type Employee   struct {
	Name   		string   `json:"Name"`
	Address 	Location  `json:"Address"`
	Email 		string   `json:"Email"`
	Phone 		string   `json:"Phone"`
	DOB  		 string   `json:"DOB"`
	Gender    	string   `json:"Gender"`
	EmployerID  string   `json:"EmployerID"`
	EmployeeID	string   `json:"EmployeeID"`
	Type		string   `json:"Type"`
	Status		string   `json:"Status"`
	StartDate	string   `json:"StartDate"`
	EndDate		string   `json:"EndDate"`
}


type Member struct {
	MemberName string `json:"MemberName"`
	MemberID  string   `json:"MemberID"`
	MemberDOB  string   `json:"MemberDOB"`
	SubscriberID string `json:"SubscriberID"`
}

//Coverage Record
type Coverage struct {
		CoverageName string  `json:"CoverageName"`
		CoverageType string `json:"CoverageType"`
		CarrierID string `json:"CarrierID"`
		GroupNum string `json:"GroupNum"`
		PlanCode string `json:"PlanCode"`
		SubscriberID string `json:"SubscriberID"`
		SubscriberName string `json:"SubscriberName"`
		SubscriberDOB string `json:"SubscriberDOB"`
		IsPrimary string `json:"IsPrimary"`
		StartDate string   `json:"StartDate"`
		EndDate string   `json:"EndDate"`
		AnnualDeductible int `json:"AnnualDeductible"`
		AnnualBenefitMaximum int `json:"AnnualBenefitMaximum"`
		LifetimeBenefitMaximum string `json:"LifetimeBenefitMaximum"`
		PreventiveCare  string `json:"PreventiveCare"`
		MinorRestorativeCare string `json:"MinorRestorativeCare"`
		MajorRestorativeCare string `json:"MajorRestorativeCare"`
		OrthodonticTreatment string `json:"OrthodonticTreatment"`
		OrthodonticLifetimeBenefitMaximum string `json:"OrthodonticLifetimeBenefitMaximum"`
		AnnualDeductibleBal int `json:"AnnualDeductibleBal"`
		AnnualBenefitMaximumBal int `json:"AnnualBenefitMaximumBal"`
		EmployeeID string `json:"EmployeeID"`
		MemberID string  `json:"MemberID"`
		EmployerID string   `json:"EmployerID"`
		Dependents []Member `json:"Dependents"`
		Premium string `json:"Premium"`
		}
//Array for storing all coverages
type AllCoverages struct{
	Coverages []Coverage `json:"Coverages"`
}

type ALLMembers struct{
	Members []Member `json:"Members"`
}


var cpPrefix = "cp:"
var accountPrefix = "acct:"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func generateCUSIPSuffix(issueDate string, days int) (string, error) {

	t, err := msToTime(issueDate)
	if err != nil {
		return "", err
	}

	maturityDate := t.AddDate(0, 0, days)
	month := int(maturityDate.Month())
	day := maturityDate.Day()

	suffix := seventhDigit[month] + eigthDigit[day]
	return suffix, nil

}

const (
	millisPerSecond = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt / millisPerSecond,
		(msInt % millisPerSecond) * nanosPerMillisecond), nil
}

type Owner struct {
	Company  string    `json:"company"`
	Quantity int      `json:"quantity"`
}

type CP struct {
	CUSIP     string  `json:"cusip"`
	Ticker    string  `json:"ticker"`
	Par       float64 `json:"par"`
	Qty       int     `json:"qty"`
	Discount  float64 `json:"discount"`
	Maturity  int     `json:"maturity"`
	Owners    []Owner `json:"owner"`
	Issuer    string  `json:"issuer"`
	IssueDate string  `json:"issueDate"`
}

type Account struct {
	ID          string  `json:"id"`
	Prefix      string  `json:"prefix"`
	CashBalance float64 `json:"cashBalance"`
	AssetsIds   []string `json:"assetIds"`
}

type Transaction struct {
	CUSIP       string   `json:"cusip"`
	FromCompany string   `json:"fromCompany"`
	ToCompany   string   `json:"toCompany"`
	Quantity    int      `json:"quantity"`
	Discount    float64  `json:"discount"`
}

func (t *SimpleChaincode) createAccounts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating accounts")

	//  				0
	// "number of accounts to create"
	var err error
	numAccounts, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("error creating accounts with input")
		return nil, errors.New("createAccounts accepts a single integer argument")
	}
	//create a bunch of accounts
	var account Account
	counter := 1
	for counter <= numAccounts {
		var prefix string
		suffix := "000A"
		if counter < 10 {
			prefix = strconv.Itoa(counter) + "0" + suffix
		} else {
			prefix = strconv.Itoa(counter) + suffix
		}
		var assetIds []string
		account = Account{ID: "company" + strconv.Itoa(counter), Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
		accountBytes, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("error creating account" + account.ID)
			return nil, errors.New("Error creating account " + account.ID)
		}
		err = stub.PutState(accountPrefix + account.ID, accountBytes)
		counter++
		fmt.Println("created account" + accountPrefix + account.ID)
	}

	fmt.Println("Accounts created")
	return nil, nil

}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating account")

	// Obtain the username to associate with the account
	if len(args) != 1 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAccount accepts a single username argument")
	}
	username := args[0]

	// Build an account object for the user
	var assetIds []string
	suffix := "000A"
	prefix := username + suffix
	var account = Account{ID: username, Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	}

	fmt.Println("Attempting to get state of any existing account for " + account.ID)
	existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {

		var company Account
		err = json.Unmarshal(existingBytes, &company)
		if err != nil {
			fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())

			if strings.Contains(err.Error(), "unexpected end") {
				fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
				err = stub.PutState(accountPrefix + account.ID, accountBytes)

				if err == nil {
					fmt.Println("created account" + accountPrefix + account.ID)
					return nil, nil
				} else {
					fmt.Println("failed to create initialize account for " + account.ID)
					return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
				}
			} else {
				return nil, errors.New("Error unmarshalling existing account " + account.ID)
			}
		} else {
			fmt.Println("Account already exists for " + account.ID + " " + company.ID)
			return nil, errors.New("Can't reinitialize existing user " + account.ID)
		}
	} else {

		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState(accountPrefix + account.ID, accountBytes)

		if err == nil {
			fmt.Println("created account" + accountPrefix + account.ID)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.ID)
			return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
		}

	}

}

func getBlockchainRecord(stub shim.ChaincodeStubInterface, recordKey string)([]byte, error){

	fmt.Println("Start getBlockchainRecord")
	fmt.Println("Looking for user with ID " + recordKey);
	var thisEmp Employee; 
	
	//get the User index
	fdAsBytes, err := stub.GetState(recordKey)
	if err != nil {
		return fdAsBytes, errors.New("Failed to get user account from blockchain")
	}
	
	
	
	err = json.Unmarshal(fdAsBytes, &thisEmp)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
	}

	return fdAsBytes, nil

}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init firing. Function will be ignored: " + function)

	
		// Create an example employee
	// Employee Record
	var aliceSheen Employee;
	var dental Coverage;

	aliceSheen.Name   			=	"Alice Sheen";
	aliceSheen.Address.Street 	=   "451 Indian Rocks Rd S";
	aliceSheen.Address.City 	=   "Largo";
	aliceSheen.Address.State 	=   "FL";
	aliceSheen.Address.Zip 		=   "33770";

	aliceSheen.Email 			=	"alicesheen@gmail.com";
	aliceSheen.Phone 			=	"727-223-5432";
	aliceSheen.DOB   			=	"08/16/1970";
	aliceSheen.Gender    		=	"Female";
	aliceSheen.EmployerID  		=	"Global Industries";
	aliceSheen.EmployeeID		=	"294048";
	aliceSheen.Type				= 	"Full Time";
	aliceSheen.Status 			= 	"Active";
	aliceSheen.StartDate		= 	"10/14/2008";
	aliceSheen.EndDate			= 	"NA";

	jsonAsBytes, _ := json.Marshal(aliceSheen)
	err := stub.PutState(aliceSheen.EmployeeID, jsonAsBytes)
	if err != nil {
		fmt.Println("Error Creating Bank user account")
		return nil, err
	}
	
	
	//create an array for storing all coverages , and store the array on the blockchain
// var coverages AllCoverages
// jsonAsBytes, _ = json.Marshal(coverages)
// err = stub.PutState("allCvgs", jsonAsBytes)
// if err != nil {
// 	return nil, err
// }

// //create an array for storing all members and store array on the blockchain
//
// var members ALLMembers
// jsonAsBytes, _ = json.Marshal(members)
// err = stub.PutState("allMembrs", jsonAsBytes)
// if err != nil {
// 	return nil, err
// }
dental.EmployeeID="294048"
dental.MemberID="M-01"
dental.SubscriberID="ba2345"
var dep1 Member;
var dep2 Member;
dep1.MemberName ="Megan Sheen";
dep1.MemberID="M-03";
dep1.MemberDOB="08/20/1990";
dep1.SubscriberID="ba2345";
// jsonAsBytes, _ = json.Marshal(dep1);
// err= stub.PutState(dep1.MemberID, jsonAsBytes);
// if err != nil {
// 	fmt.Println("Error Creating dependents")
// 	return nil, err


dep2.MemberName ="Wade Sheen";
dep2.MemberID="M-02";
dep2.MemberDOB="08/20/1961";
dep2.SubscriberID="ba2345";
//jsonAsBytes, _ = json.Marshal(dep2);
//err= stub.PutState(dep2.MemberID, jsonAsBytes);
//if err != nil {
//	fmt.Println("Error Creating dependents")
//	return nil, err
//}
	dental.Dependents=append(dental.Dependents,dep1);
	dental.Dependents=append(dental.Dependents,dep2);

	dentalAsBytes, _ := json.Marshal(dental)
	fmt.Println("Test before");

	err = stub.PutState(dental.SubscriberID, dentalAsBytes)
	if err != nil {
		fmt.Println("Error adding Coverage")
		return nil, err
	}
	
	
	// Initialize the collection of commercial paper keys
	fmt.Println("Initializing paper keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err = stub.PutState("PaperKeys", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize paper key collection")
	}

	fmt.Println("Initialization complete")
	return nil, nil
}

func (t *SimpleChaincode) issueCommercialPaper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating commercial paper")

	/*		0
		json
	  	{
			"ticker":  "string",
			"par": 0.00,
			"qty": 10,
			"discount": 7.5,
			"maturity": 30,
			"owners": [ // This one is not required
				{
					"company": "company1",
					"quantity": 5
				},
				{
					"company": "company3",
					"quantity": 3
				},
				{
					"company": "company4",
					"quantity": 2
				}
			],				
			"issuer":"company2",
			"issueDate":"1456161763790"  (current time in milliseconds as a string)

		}
	*/
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}

	var cp CP
	var err error
	var account Account

	fmt.Println("Unmarshalling CP")
	err = json.Unmarshal([]byte(args[0]), &cp)
	if err != nil {
		fmt.Println("error invalid paper issue")
		return nil, errors.New("Invalid commercial paper issue")
	}

	//generate the CUSIP
	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
	accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}

	account.AssetsIds = append(account.AssetsIds, cp.CUSIP)

	// Set the issuer to be the owner of all quantity
	var owner Owner
	owner.Company = cp.Issuer
	owner.Quantity = cp.Qty

	cp.Owners = append(cp.Owners, owner)

	suffix, err := generateCUSIPSuffix(cp.IssueDate, cp.Maturity)
	if err != nil {
		fmt.Println("Error generating cusip")
		return nil, errors.New("Error generating CUSIP")
	}

	fmt.Println("Marshalling CP bytes")
	cp.CUSIP = account.Prefix + suffix

	fmt.Println("Getting State on CP " + cp.CUSIP)
	cpRxBytes, err := stub.GetState(cpPrefix + cp.CUSIP)
	if cpRxBytes == nil {
		fmt.Println("CUSIP does not exist, creating it")
		cpBytes, err := json.Marshal(&cp)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix + cp.CUSIP, cpBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(accountPrefix + cp.Issuer, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing commercial paper")
		}


		// Update the paper keys by adding the new key
		fmt.Println("Getting Paper Keys")
		keysBytes, err := stub.GetState("PaperKeys")
		if err != nil {
			fmt.Println("Error retrieving paper keys")
			return nil, errors.New("Error retrieving paper keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling paper keys ")
		}

		fmt.Println("Appending the new key to Paper Keys")
		foundKey := false
		for _, key := range keys {
			if key == cpPrefix + cp.CUSIP {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, cpPrefix + cp.CUSIP)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on PaperKeys")
			err = stub.PutState("PaperKeys", keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue commercial paper %+v\n", cp)
		return nil, nil
	} else {
		fmt.Println("CUSIP exists")

		var cprx CP
		fmt.Println("Unmarshalling CP " + cp.CUSIP)
		err = json.Unmarshal(cpRxBytes, &cprx)
		if err != nil {
			fmt.Println("Error unmarshalling cp " + cp.CUSIP)
			return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
		}

		cprx.Qty = cprx.Qty + cp.Qty

		for key, val := range cprx.Owners {
			if val.Company == cp.Issuer {
				cprx.Owners[key].Quantity += cp.Qty
				break
			}
		}

		cpWriteBytes, err := json.Marshal(&cprx)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix + cp.CUSIP, cpWriteBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Updated commercial paper %+v\n", cprx)
		return nil, nil
	}
}

func GetAllCPs(stub shim.ChaincodeStubInterface) ([]CP, error) {

	var allCPs []CP

	// Get list of all the keys
	keysBytes, err := stub.GetState("PaperKeys")
	if err != nil {
		fmt.Println("Error retrieving paper keys")
		return nil, errors.New("Error retrieving paper keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling paper keys")
		return nil, errors.New("Error unmarshalling paper keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)

		var cp CP
		err = json.Unmarshal(cpBytes, &cp)
		if err != nil {
			fmt.Println("Error retrieving cp " + value)
			return nil, errors.New("Error retrieving cp " + value)
		}

		fmt.Println("Appending CP" + value)
		allCPs = append(allCPs, cp)
	}

	return allCPs, nil
}

func GetCP(cpid string, stub shim.ChaincodeStubInterface) (CP, error) {
	var cp CP

	cpBytes, err := stub.GetState(cpid)
	if err != nil {
		fmt.Println("Error retrieving cp " + cpid)
		return cp, errors.New("Error retrieving cp " + cpid)
	}

	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + cpid)
		return cp, errors.New("Error unmarshalling cp " + cpid)
	}

	return cp, nil
}


func getEmployeeRecord(stub shim.ChaincodeStubInterface, employeeId string)([]byte, error){

	fmt.Println("Start getEmployeeRecord")
	fmt.Println("Looking for Employee with ID " + employeeId);

	//get the User index
	fdAsBytes, err := stub.GetState(employeeId)
	if err != nil {
		return nil, errors.New("Failed to get Employee account from blockchain")
	}

	return fdAsBytes, nil

}
func getCoverages(stub shim.ChaincodeStubInterface, subscriberID string)([]byte, error){

	fmt.Println("Start Get Coverage")
	fmt.Println("Looking for Coverage for SubscriberID" + subscriberID);

	coverageAsBytes, err := stub.GetState(subscriberID)
	if err != nil {
		return nil, errors.New("Failed to get coverage from blockchain")
	}
	return coverageAsBytes, nil
}
//Update Coverage

func (t *SimpleChaincode) updateCoverage(stub shim.ChaincodeStubInterface, args []string)([]byte, error){

			var dentalcoverage Coverage
			var subscriberID string
  		var  benefitMaximumBal int
			var deductibleBal int
			subscriberID=args[0];
  	  deductibleBal, err := strconv.Atoi(args[1]);
			benefitMaximumBal, err =strconv.Atoi(args[2]);
			coverageAsBytes, err := stub.GetState(subscriberID)
			fmt.Println("In Update coverage");
			if err != nil {
				return nil, errors.New("Failed to get Coverage from blopckchain")
			}
		json.Unmarshal(coverageAsBytes,&dentalcoverage)
		dentalcoverage.SubscriberID=subscriberID;
		dentalcoverage.AnnualDeductibleBal=deductibleBal;
		dentalcoverage.AnnualBenefitMaximumBal=benefitMaximumBal;
		dentalcvgAsBytes, _ := json.Marshal(dentalcoverage)
		err = stub.PutState(dentalcoverage.SubscriberID,dentalcvgAsBytes)
		if err != nil {
		fmt.Println("Error updating coverage")
		return nil, err
		}
			return nil,nil

		}
func (t *SimpleChaincode) addCoverage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
fmt.Println("In addcoverage chaincode");
	//var dental Coverage


	var dentalcvg Coverage
	coverageAsBytes, err := stub.GetState(args[5])
	fmt.Println("In Update coverage");
	if err != nil {
		return nil, errors.New("Failed to get Coverage from blockchain")
	}
	json.Unmarshal(coverageAsBytes,&dentalcvg)


	var AnnualBenefitMaximum int
	dentalcvg.CoverageName=args[0]
	dentalcvg.CoverageType=args[1]
 	dentalcvg.CarrierID=args[2]
 	dentalcvg.GroupNum=args[3]
	dentalcvg.PlanCode=args[4]
 	//dentalcvg.SubscriberID=args[5]
 	dentalcvg.SubscriberName=args[6]
	dentalcvg.SubscriberDOB=args[7]
	dentalcvg.IsPrimary=args[8]
	dentalcvg.EndDate=args[9]
	dentalcvg.StartDate=args[10]
  AnnualDeductible, err := strconv.Atoi(args[11])
	dentalcvg.AnnualDeductible=AnnualDeductible
	AnnualBenefitMaximum, err= strconv.Atoi(args[12])
	dentalcvg.AnnualBenefitMaximum=AnnualBenefitMaximum
	dentalcvg.LifetimeBenefitMaximum=args[13]
	dentalcvg.PreventiveCare =args[14]
	dentalcvg.MinorRestorativeCare=args[15]
	dentalcvg.MajorRestorativeCare=args[16]
	dentalcvg.OrthodonticTreatment=args[17]
	dentalcvg.OrthodonticLifetimeBenefitMaximum=args[18]
	dentalcvg.AnnualDeductibleBal=dentalcvg.AnnualDeductible
	dentalcvg.AnnualBenefitMaximumBal=AnnualBenefitMaximum
	dentalcvg.EmployerID=args[19]
	dentalcvg.Premium=args[20]
//new code for single struct
dentalAsBytes, _ := json.Marshal(dentalcvg)
err = stub.PutState(dentalcvg.SubscriberID, dentalAsBytes)
if err != nil {
	fmt.Println("Error adding Coverage")
	return nil, err
}

//new code for single struct
	return nil,nil


}
	// end add coverage
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}



func GetCompany(companyID string, stub shim.ChaincodeStubInterface) (Account, error) {
	var company Account
	companyBytes, err := stub.GetState(accountPrefix + companyID)
	if err != nil {
		fmt.Println("Account not found " + companyID)
		return company, errors.New("Account not found " + companyID)
	}

	err = json.Unmarshal(companyBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling account " + companyID + "\n err:" + err.Error())
		return company, errors.New("Error unmarshalling account " + companyID)
	}

	return company, nil
}

func verifyEmployment(stub shim.ChaincodeStubInterface, subscriberId string)([]byte, error){

	// Get the insurance coverage record
	coverageAsBytes, err := stub.GetState(subscriberId)
	if err != nil {
		return nil, errors.New("Failed to get coverage from blockchain")
	}
	var insRecord Coverage
	json.Unmarshal(coverageAsBytes, &insRecord)


	// Get the employee record
	employeeAsBytes, err := stub.GetState(insRecord.EmployeeID)
	if err != nil {
		return nil, errors.New("Failed to get user account from blockchain")
	}

	// Unmarshall the employee record from the blockchain
	var empRecord Employee
	json.Unmarshal(employeeAsBytes, &empRecord)


	// Check that the employment is valid

	var checkResults DataCheck
	checkResults.Result 	= "Passed";
	checkResults.Message 	= "Primary policy holder is a valid employee."

	if empRecord.Type != "Full Time" {
		checkResults.Result 	= "Failed";
		checkResults.Message 	= "Primary policy holder is not a full-time employee."
	}

	if empRecord.Status != "Active" {
		checkResults.Result 	= "Failed";
		checkResults.Message 	= "Primary policy holder is not an active employee."
	}


	resAsBytes, _ := json.Marshal(checkResults)
	return resAsBytes, nil

}


func verifyCoverage(stub shim.ChaincodeStubInterface, subscriberId string, memberId string)([]byte, error){

	dateForDemo, _ := time.Parse("01/02/2006",  "01/15/2017")

	// Get the insurance coverage record
	coverageAsBytes, err := stub.GetState(subscriberId)
	if err != nil {
		return nil, errors.New("Failed to get coverage from blockchain")
	}
	var insRecord Coverage
	json.Unmarshal(coverageAsBytes, &insRecord)

	/// Add some magic to get the right dependent here


	var checkResults DataCheck
	checkResults.Result 	= "Passed";
	checkResults.Message 	= "Coverage is valid."


	// Check if age is > 26 ONLY if the person is NOT the primary member
	if	insRecord.MemberID != memberId {

		numDependents := len(insRecord.Dependents)

		for index := numDependents - 1; index >= 0; index-- {

			if insRecord.Dependents[index].MemberID == memberId{

				memberDobDate, _ := time.Parse("01/02/2006", insRecord.Dependents[index].MemberDOB)
				ageInYears, _, _, _, _, _ 	  := diff(memberDobDate, dateForDemo)

				if ageInYears >= 26 {
					checkResults.Result 	= "Failed";
					checkResults.Message 	= "Dependent is above legal maximum coverage age. DOB: " + insRecord.Dependents[index].MemberDOB;
				}

			}

		}
	}

	startDate,_ 	:= time.Parse("2006-01-02",  insRecord.StartDate)
	endDate,_	:= time.Parse("2006-01-02",  insRecord.EndDate)

	// Check if plan is active
	if (!(dateForDemo.After(startDate) && dateForDemo.Before(endDate))) {
		checkResults.Result 	= "Failed";
		checkResults.Message 	= "Policy is not active."
	}

	// Check if annual benefit max has not yet been reached
	if (insRecord.AnnualBenefitMaximumBal  < 1 ) {
		checkResults.Result 	= "Failed";
		checkResults.Message 	= "Annual maximum has already been met."
	}

	resAsBytes, _ := json.Marshal(checkResults)
	return resAsBytes, nil

}

// Still working on this one
func (t *SimpleChaincode) transferPaper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Transferring Paper")
	/*		0
		json
	  	{
			  "CUSIP": "",
			  "fromCompany":"",
			  "toCompany":"",
			  "quantity": 1
		}
	*/
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}

	var tr Transaction

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid commercial paper issue")
	}

	fmt.Println("Getting State on CP " + tr.CUSIP)
	cpBytes, err := stub.GetState(cpPrefix + tr.CUSIP)
	if err != nil {
		fmt.Println("CUSIP not found")
		return nil, errors.New("CUSIP not found " + tr.CUSIP)
	}

	var cp CP
	fmt.Println("Unmarshalling CP " + tr.CUSIP)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + tr.CUSIP)
		return nil, errors.New("Error unmarshalling cp " + tr.CUSIP)
	}

	var fromCompany Account
	fmt.Println("Getting State on fromCompany " + tr.FromCompany)
	fromCompanyBytes, err := stub.GetState(accountPrefix + tr.FromCompany)
	if err != nil {
		fmt.Println("Account not found " + tr.FromCompany)
		return nil, errors.New("Account not found " + tr.FromCompany)
	}

	fmt.Println("Unmarshalling FromCompany ")
	err = json.Unmarshal(fromCompanyBytes, &fromCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.FromCompany)
		return nil, errors.New("Error unmarshalling account " + tr.FromCompany)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + tr.ToCompany)
	toCompanyBytes, err := stub.GetState(accountPrefix + tr.ToCompany)
	if err != nil {
		fmt.Println("Account not found " + tr.ToCompany)
		return nil, errors.New("Account not found " + tr.ToCompany)
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.ToCompany)
		return nil, errors.New("Error unmarshalling account " + tr.ToCompany)
	}

	// Check for all the possible errors
	ownerFound := false
	quantity := 0
	for _, owner := range cp.Owners {
		if owner.Company == tr.FromCompany {
			ownerFound = true
			quantity = owner.Quantity
		}
	}

	// If fromCompany doesn't own this paper
	if ownerFound == false {
		fmt.Println("The company " + tr.FromCompany + "doesn't own any of this paper")
		return nil, errors.New("The company " + tr.FromCompany + "doesn't own any of this paper")
	} else {
		fmt.Println("The FromCompany does own this paper")
	}

	// If fromCompany doesn't own enough quantity of this paper
	if quantity < tr.Quantity {
		fmt.Println("The company " + tr.FromCompany + "doesn't own enough of this paper")
		return nil, errors.New("The company " + tr.FromCompany + "doesn't own enough of this paper")
	} else {
		fmt.Println("The FromCompany owns enough of this paper")
	}

	amountToBeTransferred := float64(tr.Quantity) * cp.Par
	amountToBeTransferred -= (amountToBeTransferred) * (cp.Discount / 100.0) * (float64(cp.Maturity) / 360.0)

	// If toCompany doesn't have enough cash to buy the papers
	if toCompany.CashBalance < amountToBeTransferred {
		fmt.Println("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")
		return nil, errors.New("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")
	} else {
		fmt.Println("The ToCompany has enough money to be transferred for this paper")
	}

	toCompany.CashBalance -= amountToBeTransferred
	fromCompany.CashBalance += amountToBeTransferred

	toOwnerFound := false
	for key, owner := range cp.Owners {
		if owner.Company == tr.FromCompany {
			fmt.Println("Reducing Quantity from the FromCompany")
			cp.Owners[key].Quantity -= tr.Quantity
			//			owner.Quantity -= tr.Quantity
		}
		if owner.Company == tr.ToCompany {
			fmt.Println("Increasing Quantity from the ToCompany")
			toOwnerFound = true
			cp.Owners[key].Quantity += tr.Quantity
			//			owner.Quantity += tr.Quantity
		}
	}

	if toOwnerFound == false {
		var newOwner Owner
		fmt.Println("As ToOwner was not found, appending the owner to the CP")
		newOwner.Quantity = tr.Quantity
		newOwner.Company = tr.ToCompany
		cp.Owners = append(cp.Owners, newOwner)
	}

	fromCompany.AssetsIds = append(fromCompany.AssetsIds, tr.CUSIP)

	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the toCompany")
		return nil, errors.New("Error marshalling the toCompany")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix + tr.ToCompany, toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the toCompany back")
		return nil, errors.New("Error writing the toCompany back")
	}

	// From company
	fromCompanyBytesToWrite, err := json.Marshal(&fromCompany)
	if err != nil {
		fmt.Println("Error marshalling the fromCompany")
		return nil, errors.New("Error marshalling the fromCompany")
	}
	fmt.Println("Put state on fromCompany")
	err = stub.PutState(accountPrefix + tr.FromCompany, fromCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the fromCompany back")
		return nil, errors.New("Error writing the fromCompany back")
	}

	// cp
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on CP")
	err = stub.PutState(cpPrefix + tr.CUSIP, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the cp back")
		return nil, errors.New("Error writing the cp back")
	}

	fmt.Println("Successfully completed Invoke")
	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

	if function == "GetAllCPs" {
		fmt.Println("Getting all CPs")
		allCPs, err := GetAllCPs(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	} else if function == "GetCP" {
		fmt.Println("Getting particular cp")
		cp, err := GetCP(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return cpBytes, nil
		}
	
	} else if function == "getCoverages" {
	
		return getCoverages(stub, args[0])

	} else if function == "getBlockchainRecord" {
	
		return getBlockchainRecord(stub, args[0])

	}  else if function == "getEmployeeRecord" {
	
		return getEmployeeRecord(stub, args[0])

	}  else if function == "verifyEmployment" {
	
		return verifyEmployment(stub, args[0])

	}  else if function == "verifyCoverage" {
	
		return verifyCoverage(stub, args[0], args[1])

	}  else if function == "GetCompany" {
		fmt.Println("Getting the company")
		company, err := GetCompany(args[0], stub)
		if err != nil {
			fmt.Println("Error from getCompany")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the company")
				return nil, err1
			}
			fmt.Println("All success, returning the company")
			return companyBytes, nil
		}
	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

	if function == "issueCommercialPaper" {
		return t.issueCommercialPaper(stub, args)
	} else if function == "transferPaper" {
		return t.transferPaper(stub, args)
	} else if function == "createAccounts" {
		return t.createAccounts(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	} else if function == "addCoverage" {											//create a transaction
		return t.addCoverage(stub, args)
	}	else if function == "updateCoverage" {											//create a transaction
	return t.updateCoverage(stub, args) 
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}

//lookup tables for last two digits of CUSIP
var seventhDigit = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "P",
	15: "Q",
	16: "R",
	17: "S",
	18: "T",
	19: "U",
	20: "V",
	21: "W",
	22: "X",
	23: "Y",
	24: "Z",
}

var eigthDigit = map[int]string{
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "J",
	19: "K",
	20: "L",
	21: "M",
	22: "N",
	23: "P",
	24: "Q",
	25: "R",
	26: "S",
	27: "T",
	28: "U",
	29: "V",
	30: "W",
	31: "X",
}
