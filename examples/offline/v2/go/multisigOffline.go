package main

import (
	"context"
	"crypto/ed25519"
	json "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// UPDATE THESE VALUES
// const algodAddress = "Your ADDRESS"
// const algodToken = "Your TOKEN"

// sandbox
const algodAddress = "http://localhost:4001"
const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// const algodAddress = "http://hackathon.algodev.network:9100"
// const algodToken = "ef920e2e7e002953f4b29a8af720efe8e4ecc75ff102b165e0472834b25832c1"

// Accounts to be used through examples
func loadAccounts() (map[int][]byte, map[int]string) {
	// Shown for demonstration purposes. NEVER reveal secret mnemonics in practice.
	// Change these values to use the accounts created previously.
	// var pks = map[int]string{
	//     1: "Account Address 1 ",
	//     2: "Account Address 2 ",
	//     3: "Account Address 3 ",
	// }
	// Paste in mnemonic phrases for all three accounts
	// mnemonic1 := "PASTE phrase for account 1"
	// mnemonic2 := "PASTE phrase for account 2"
	// mnemonic3 := "PASTE phrase for account 3"

	mnemonic1 := "predict mandate aware dizzy limit match hazard fantasy victory auto fortune hello public dragon ostrich happy blue spray parrot island odor actress only ability hurry"
	mnemonic2 := "moon grid random garlic effort faculty fence gym write skin they joke govern home huge there claw skin way bid fit bean damp able only"
	mnemonic3 := "mirror zone together remind rural impose balcony position minimum quick manage climb quit draft lion device pluck rug siege robust spirit fine luggage ability actual"

	mnemonics := []string{mnemonic1, mnemonic2, mnemonic3}
	pks := map[int]string{1: "", 2: "", 3: ""}
	var sks = make(map[int][]byte)

	for i, m := range mnemonics {
		var err error
		sk, err := mnemonic.ToPrivateKey(m)
		sks[i+1] = sk
		if err != nil {
			fmt.Printf("Issue with account %d private key conversion.", i+1)
		} else {
			fmt.Printf("Loaded Key %d: %s\n", i+1, pks[i+1])
		}
		// derive public address from Secret Key.
		pk := sk.Public()
		var a types.Address
		cpk := pk.(ed25519.PublicKey)
		copy(a[:], cpk[:])
		pks[i+1] = a.String()
		fmt.Printf("Loaded Key %d: %s\n", i+1, pks[i+1])
	}
	return sks, pks
}

// Function that waits for a given txId to be confirmed by the network
func waitForConfirmation(txID string, client *algod.Client, timeout uint64) (models.PendingTransactionInfoResponse, error) {
	pt := new(models.PendingTransactionInfoResponse)
	if client == nil || txID == "" || timeout < 0 {
		fmt.Printf("Bad arguments for waitForConfirmation")
		var msg = errors.New("Bad arguments for waitForConfirmation")
		return *pt, msg

	}

	status, err := client.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		var msg = errors.New(strings.Join([]string{"error getting algod status: "}, err.Error()))
		return *pt, msg
	}
	startRound := status.LastRound + 1
	currentRound := startRound

	for currentRound < (startRound + timeout) {

		fmt.Printf("Round checked: %d\n", currentRound)

		*pt, _, err = client.PendingTransactionInformation(txID).Do(context.Background())
		if err != nil {
			fmt.Printf("error getting pending transaction: %s\n", err)
			var msg = errors.New(strings.Join([]string{"error getting pending transaction: "}, err.Error()))
			return *pt, msg
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+txID+" confirmed in round %d\n", pt.ConfirmedRound)
			return *pt, nil
		}
		if pt.PoolError != "" {
			fmt.Printf("There was a pool error, then the transaction has been rejected!")
			var msg = errors.New("There was a pool error, then the transaction has been rejected")
			return *pt, msg
		}
		fmt.Printf("waiting for confirmation\n")
		status, err = client.StatusAfterBlock(currentRound).Do(context.Background())
		currentRound++
	}
	msg := errors.New("Tx not found in round range")
	return *pt, msg
}

// PrettyPrint prints Go structs
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func saveUnsignedMultisigTransaction() {

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}

	_, pks := loadAccounts()

	addr1, _ := types.DecodeAddress(pks[1])
	addr2, _ := types.DecodeAddress(pks[2])
	addr3, _ := types.DecodeAddress(pks[3])

	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})

	if err != nil {
		panic("invalid multisig parameters")
	}
	// Construct the transaction
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	// comment out the next two (2) lines to use suggested fees
	txParams.FlatFee = true
	txParams.Fee = 1000
	fromAddr, _ := ma.Address()
	//fromAddr := addr1.String()
	toAddr := addr2.String()
	var amount uint64 = 1000000
	var minFee uint64 = 1000
	note := []byte("Hello World")
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	txn, err := transaction.MakePaymentTxnWithFlatFee(fromAddr.String(), toAddr, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	unsignedTx := types.SignedTxn{
		Txn: txn,
	}

	// save unsigned transaction to file
	err = ioutil.WriteFile("./unsigned.txn", msgpack.Encode(unsignedTx), 0644)
	if err != nil {
		fmt.Printf("Failed in saving trx to file, error %s\n", err)
		return
	}
	fmt.Printf("Saved unsigned transaction to file\n")

	err = ioutil.WriteFile("./ma.txn", msgpack.Encode(ma), 0644)
	if err != nil {
		fmt.Printf("Failed in saving ma to file, error %s\n", err)
		return
	}
	fmt.Printf("Saved ma to file\n")
	return

}

func readUnsignedMultisigTransaction() {

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}

	// read unsigned transaction from file
	dat, err := ioutil.ReadFile("./unsigned.txn")
	if err != nil {
		fmt.Printf("Error reading transaction from file: %s\n", err)
		return
	}
	var unsignedTxRaw types.SignedTxn
	var unsignedTxn types.Transaction

	msgpack.Decode(dat, &unsignedTxRaw)
	unsignedTxn = unsignedTxRaw.Txn
	sks, pks := loadAccounts()
	addr := pks[1]
	fmt.Printf("Address is: %s\n", addr)
	msgpack.Decode(dat, &unsignedTxRaw)
	unsignedTxn = unsignedTxRaw.Txn
	// read ma from file
	datma, err := ioutil.ReadFile("./ma.txn")
	if err != nil {
		fmt.Printf("Error reading ma from file: %s\n", err)
		return
	}
	var ma crypto.MultisigAccount
	msgpack.Decode(datma, &ma)

	txid, txBytes, err := crypto.SignMultisigTransaction(sks[1], ma, unsignedTxn)
	if err != nil {
		println(err.Error)
		panic("could not sign multisig transaction")
	}
	fmt.Printf("Made partially-signed multisig transaction with TxID %s: %x\n", txid, txBytes)
	txid, twoOfThreeTxBytes, err := crypto.AppendMultisigTransaction(sks[2], ma, txBytes)

	if err != nil {
		panic("could not append signature to multisig transaction")
	}
	fmt.Printf("Appended bytes %x\n", twoOfThreeTxBytes)

	fmt.Printf("Made 2-out-of-3 multisig transaction with TxID %s: %x\n", txid, twoOfThreeTxBytes)

	txid, err = algodClient.SendRawTransaction(twoOfThreeTxBytes).Do(context.Background())

	// Wait for confirmation
	confirmedTxn, err := waitForConfirmation(txid, algodClient, 5)
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txid)
		return
	}
	txnJSON, err := json.MarshalIndent(confirmedTxn.Transaction.Txn, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)

	fmt.Printf("Decoded note: %s\n", string(confirmedTxn.Transaction.Txn.Note))
}
func saveSignedMultisigTransaction() {

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}

	sks, pks := loadAccounts()

	addr1, _ := types.DecodeAddress(pks[1])
	addr2, _ := types.DecodeAddress(pks[2])
	addr3, _ := types.DecodeAddress(pks[3])

	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})

	if err != nil {
		panic("invalid multisig parameters")
	}
	// Check account balance
	// accountInfo, err := algodClient.AccountInformation(addr1.String()).Do(context.Background())
	// if err != nil {
	// 	fmt.Printf("Error getting account info: %s\n", err)
	// 	return
	// }
	// fmt.Printf("Account balance: %d microAlgos\n", accountInfo.Amount)

	// Construct the transaction
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	// comment out the next two (2) lines to use suggested fees
	txParams.FlatFee = true
	txParams.Fee = 1000
	fromAddr, _ := ma.Address()
	// Print multisig account
	fmt.Printf("Here is your multisig address : %s \n", fromAddr.String())
	fmt.Println("Please go to: https://bank.testnet.algorand.network/ to fund your multisig account.")
	toAddr := addr2.String()
	var amount uint64 = 1000000
	var minFee uint64 = 1000
	note := []byte("Hello World")
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	txn, err := transaction.MakePaymentTxnWithFlatFee(fromAddr.String(), toAddr, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	txid, txBytes, err := crypto.SignMultisigTransaction(sks[1], ma, txn)
	if err != nil {
		println(err.Error)
		panic("could not sign multisig transaction")
	}
	fmt.Printf("Made partially-signed multisig transaction with TxID %s: %x\n", txid, txBytes)
	txid, twoOfThreeTxBytes, err := crypto.AppendMultisigTransaction(sks[2], ma, txBytes)

	if err != nil {
		panic("could not append signature to multisig transaction")
	}
	fmt.Printf("Appended bytes %x\n", twoOfThreeTxBytes)

	fmt.Printf("Made 2-out-of-3 multisig transaction with TxID %s: %x\n", txid, twoOfThreeTxBytes)

	//Save the signed transaction to file
	err = ioutil.WriteFile("./signed.stxn", twoOfThreeTxBytes, 0644)
	if err != nil {

		fmt.Printf("Failed in saving stxn to file, error %s\n", err)
		return
	}
	fmt.Printf("Saved signed transaction to file\n")

	return
}

func readSignedMultisigTransaction() {

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}

	// read unsigned transaction from file
	dat, err := ioutil.ReadFile("./signed.stxn")
	if err != nil {
		fmt.Printf("Error reading signed transaction from file: %s\n", err)
		return
	}

	// Broadcast the transaction to the network
	txid, err := algodClient.SendRawTransaction(dat).Do(context.Background())

	// Wait for confirmation
	confirmedTxn, err := waitForConfirmation(txid, algodClient, 4)
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txid)
		return
	}
	txnJSON, err := json.MarshalIndent(confirmedTxn.Transaction.Txn, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)

	fmt.Printf("Decoded note: %s\n", string(confirmedTxn.Transaction.Txn.Note))
}

func main() {

	// saveUnsignedMultisigTransaction()
	// readUnsignedMultisigTransaction()

	saveSignedMultisigTransaction()
	readSignedMultisigTransaction()

}
