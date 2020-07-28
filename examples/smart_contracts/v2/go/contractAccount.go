package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	// "encoding/json"
	"io/ioutil"
	"log"
	"os"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
)

// Function that waits for a given txId to be confirmed by the network
func waitForConfirmation(txID string, client *algod.Client) {
	status, err := client.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}
	lastRound := status.LastRound
	for {
		pt, _, err := client.PendingTransactionInformation(txID).Do(context.Background())
		if err != nil {
			fmt.Printf("error getting pending transaction: %s\n", err)
			return
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+txID+" confirmed in round %d\n", pt.ConfirmedRound)
			break
		}
		fmt.Printf("waiting for confirmation\n")
		lastRound++
		status, err = client.StatusAfterBlock(lastRound).Do(context.Background())
	}
}
func main() {

	// const algodToken = "algod-token<PLACEHOLDER>"
	// const algodAddress = "algod-address<PLACEHOLDER>"

	// sandbox
	const algodAddress = "http://localhost:4001"
	const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}
	// Create logic signature
	var sk ed25519.PrivateKey
    var ma crypto.MultisigAccount

    // samplearg.teal  ... This code is meant for learning purposes only
    // It should not be used in production
    // arg_0
    // btoi
    // int 123
    // == 
	file, err := os.Open("./samplearg.teal")
	// file, err := os.Open("<PLACEHOLDER>")
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()
    tealFile, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Printf("failed to read file: %s\n", err)
		return}
		
    response, err := algodClient.TealCompile(tealFile).Do(context.Background())
    fmt.Printf("Hash = %s\n", response.Hash)
    fmt.Printf("Result = %s\n", response.Result)
    
    program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    // if no args use these two lines
    // var args [][]byte
    // lsig, err := crypto.MakeLogicSig(program, args, sk, ma)

    // string parameter
    // args := make([][]byte, 1)
    // args[0] = []byte("my string")
    // lsig, err := crypto.MakeLogicSig(program, args, sk, ma)
    
    // integer args parameter
    args := make([][]byte, 1)
    var buf [8]byte
    binary.BigEndian.PutUint64(buf[:], 123)
    args[0] = buf[:]
    lsig, err := crypto.MakeLogicSig(program, args, sk, ma)

    addr := crypto.LogicSigAddress(lsig).String()
	fmt.Printf("Escrow Address: %s\n" , addr )

	// Get suggested params for the transaction
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	// comment out the next two (2) lines to use suggested fees
	txParams.FlatFee = true
	txParams.Fee = 1000

	// Make transaction
	// const receiver = "transaction-receiver"<PLACEHOLDER>
	// const fee = fee<PLACEHOLDER>
	// const amount = amount<PLACEHOLDER>
	const receiver = "QUDVUXBX4Q3Y2H5K2AG3QWEOMY374WO62YNJFFGUTMOJ7FB74CMBKY6LPQ"
	const fee = 1000
	const amount = 100000
	var minFee uint64 = 1000
	note := []byte("Hello World")
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)
	tx, err := transaction.MakePaymentTxnWithFlatFee(
		addr, receiver, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)

	txID, stx, err := crypto.SignLogicsigTransaction(lsig, tx)
	if err != nil {
		fmt.Printf("Signing failed with %v", err)
		return
	}
	fmt.Printf("Signed tx: %v\n", txID)
	// logic signature transaction can be written to a file
	// f, err := os.Create("simple.stxn")

	// defer f.Close()
	// if _, err := f.Write(stx); err != nil {
	//     // handle
	// }
	// if err := f.Sync(); err != nil {
	//     // handle
	// }
		
	// Submit the raw transaction to network

	transactionID, err := algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		fmt.Printf("Sending failed with %v\n", err)
	}
    // Wait for transaction to be confirmed
    waitForConfirmation(txID, algodClient)
    fmt.Printf("Transaction ID: %v\n", transactionID)
}
