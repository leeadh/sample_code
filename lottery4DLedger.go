package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Lottery4DLedger struct {
}

// We can assume that the game name is always 4D
type Lottery4DBet struct {
	AccountNumber  string `json: "account_number"`
	BetID          string `json: "bet_id"`
	BetReceipt     string `json: "bet_receipt"'`
	BetDate        string `json: "bet_date"`
	BetslipID      string `json: "betslip_id"`
	BetslipReceipt string `json: "betslip_receipt"`
	Pick           string `json: "pick"`
	Source         string `json: "betting_source"`
	Game           string `json: "game"`
	DrawID         string `json: "draw_id"`
	TotalStake     string `json: "total_stake"`
	BigStake       string `json: "big_stake"`
	SmallStake     string `json: "small_stake"`
	Currency       string `json: "currency"`
	BetType        string `json: "bet_type"`
	StakePerLine   string `json : "stake_per_line"`
}

const (

	// Transaction Types
	TX_PLACE_BET       = "PlaceBet"
	TX_PREFIX          = "tx_"
	KEY_PREFIX         = "Game~BetId~"
	ERROR_WRONG_FORMAT = "{\"code\":301, \"reason\": \"Command format is wrong. Insufficient arguments\"}"
)

var logger = shim.NewLogger("Lottery4DLedger")

// ===========================================================
// Init with seed data
// =============================================

func (t *Lottery4DLedger) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Initializing Chaincode")
	return shim.Success(nil)
}

func (t *Lottery4DLedger) initBets(stub shim.ChaincodeStubInterface) pb.Response {
	// Seed data for the chaincode
	fmt.Printf("Invoking initBets()")
	logger.Info("Init Bets")

	bet1AsString := `{""account_number": "010432899", "bet_id": "324636", "bet_receipt": "L/0040544/0000001", "bet_date": "2018-09-07 10:14:35", "betslip_id": "36006", "betslip_receipt": "B/0040544/0000001", "pick": "1111", "betting_source": "internet", "game": "4D", "draw_id": "556578", "total_stake": "10.00", "big_stake": "5.00", "small_stake": "5.00", "currency": "SGD", "bet_type": "4D", "stake_per_line": "1.00" }`
	bet2AsString := `{"account_number": "010432899", "bet_id": "324637", "bet_receipt": "L/0040544/0000002", "bet_date": "2018-09-07 10:15:35", "betslip_id": "36007", "betslip_receipt": "B/0040544/0000002", "pick": "0087", "betting_source": "internet", "game": "4D", "draw_id": "556578", "total_stake": "24.00", "big_stake": "1.00", "small_stake": "1.00", "currency": "SGD", "bet_type": "4D1PG", "stake_per_line": "2.00" }`
	bet3AsString := `{""account_number": "010432899", "bet_id": "324638", "bet_receipt": "L/0040544/0000003", "bet_date": "2018-09-07 10:15:21", "betslip_id": "36008", "betslip_receipt": "B/0040544/0000003", "pick": "8700", "betting_source": "internet", "game": "4D", "draw_id": "556578", "total_stake": "2.00", "big_stake": "1.00", "small_stake": "1.00", "currency": "SGD", "bet_type": "4D1iG", "stake_per_line": "2.00" }`
	bet4AsString := `{""account_number": "010432899", "bet_id": "324639", "bet_receipt": "L/0040544/0000004", "bet_date": "2018-09-07 10:15:48", "betslip_id": "36009", "betslip_receipt": "B/0040544/0000004", "pick": "123R", "betting_source": "internet", "game": "4D", "draw_id": "556578", "total_stake": "10.00", "big_stake": "1.00", "small_stake": "-", "currency": "SGD", "bet_type": "4DR4B",   "stake_per_line": "1.00" }`

	var err error
	var compositeKey string

	compositeKey, err = stub.CreateCompositeKey(KEY_PREFIX, []string{"4D", "324636"})
	logger.Info("Seeding Data - Composite Key : ", compositeKey)
	err = stub.PutState(compositeKey, []byte(bet1AsString))
	if err != nil {
		return shim.Error(err.Error())
	}

	compositeKey, err = stub.CreateCompositeKey(KEY_PREFIX, []string{"4D", "324637"})
	logger.Info("Seeding Data - Composite Key : ", compositeKey)

	err = stub.PutState(compositeKey, []byte(bet2AsString))
	if err != nil {
		return shim.Error(err.Error())
	}

	compositeKey, err = stub.CreateCompositeKey(KEY_PREFIX, []string{"4D", "324638"})
	logger.Info("Seeding Data - Composite Key : ", compositeKey)

	err = stub.PutState(compositeKey, []byte(bet3AsString))
	if err != nil {
		return shim.Error(err.Error())
	}

	compositeKey, err = stub.CreateCompositeKey(KEY_PREFIX, []string{"4D", "324638"})
	logger.Info("Seeding Data - Composite Key : ", compositeKey)

	err = stub.PutState(compositeKey, []byte(bet4AsString))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===========================================================
// Router function
// ===========================================================
func (t *Lottery4DLedger) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoking Function : ", function)

	switch function {
	case "init":
		return t.initBets(stub)
	case "placeBet":
		return t.placeBet(stub, args)
	case "queryBet":
		return t.queryBet(stub, args)
	case "queryByDrawID":
		return t.queryByDrawID(stub, args)
	case "queryByBetID":
		return t.queryByBetID(stub, args)
	case "queryAllBets":
		return t.queryAllBet(stub)
		//	case "queryAllBets":
		//		return t.queryAllBets(stub, args)
	default:
		return shim.Error("[ERROR] Not a valid function : " + function)
	}
}

/// ===========================================================
// Query by Bet ID 
// ===========================================================
func (t *Lottery4DLedger) queryByBetID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var bet_id, jsonResp string
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting BetID number of the vehicle to query")
	}

	bet_id = args[0]
	valAsbytes, err := stub.GetState(bet_id) //get the bet from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + bet_id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Bet does not exist: " + bet_id + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)

}

/// ===========================================================
// Query by draw ID
// ===========================================================
func (t *Lottery4DLedger) queryByDrawID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	drawId := args[0]
	//queryString := fmt.Sprintf("{ 'selector' : { 'drawId' : '%s' }''}", drawId)
	//queryString := fmt.Sprintf("{\"selector\":{\"drawId\":\"%s\"}}", drawId)

	queryString := "{\"selector\":{\"draw_id\":\"" + drawId + "\"}}"

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===========================================================
// start placing a bet
// ===========================================================

func (t *Lottery4DLedger) placeBet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("Placing Bet")

	if len(args) != 16 {
		return shim.Error("Incorrect number of arguments. Expecting 16")
	}

	// ==== Create new bet  and marshal to JSON ====
	new4DBetJSON := &Lottery4DBet{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15]}
	new4DBetJSONasBytes, err := json.Marshal(new4DBetJSON)
	logger.Info("new4DBetJSON        : ", new4DBetJSON)
	logger.Info("new4DBetJSONasBytes : ", new4DBetJSONasBytes)
	err = stub.PutState(new4DBetJSON.BetID, new4DBetJSONasBytes)
	if err != nil {
		shim.Error("Error. Failed to place bet : " + err.Error())
	}

	return shim.Success(nil)
}

// ===========================================================
// query all bets
// ===========================================================
func (t *Lottery4DLedger) queryAllBet(stub shim.ChaincodeStubInterface) pb.Response {

	logger.Info("Querying All Bets")

	resultsIterator, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Info(" queryAllBet ", buffer.String())
	fmt.Printf("- queryAllBet queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

// ===========================================================
// query a single bet
// ===========================================================
func (t *Lottery4DLedger) queryBet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	betKey := args[0]

	logger.Info("Querying Bet : ", betKey)
	w, err := stub.GetState(betKey)
	if err != nil {
		logger.Info("queryBet: Error query: ", err)
		return shim.Error("Unable to get the state")
	}
	logger.Info("Results : ", w)
	return shim.Success(w)
}

// =========================================================================================
// private  - rich query - getQueryResultForQueryString executes the passed in query string.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}




func main() {

	// configure the logger

	if err := shim.Start(new(Lottery4DLedger)); err != nil {
		fmt.Printf("Error starting Lottery4DLedger chaincode: %s", err)
	}
}
