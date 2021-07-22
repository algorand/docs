title: Your First Transaction

This section is a quick start guide for sending your first transaction on the Algorand TestNet network using the Go programming language. This guide installs the Go SDK, creates an account and submits a payment transaction. This guide also installs Algorand Sandbox, which provides required infrastructure for development and testing. 

!!! Info
    If you are a visual learner, try our [live demo](https://replit.com/@Algorand/Getting-Started-with-Go){target=_blank} or watch a [video walkthrough](https://youtu.be/rFG7Zo2JvIY?t=){target=_blank} explaining all the code in the steps below.
 
# Install Algorand Sandbox

Algorand Sandbox is developer-focused tool for quickly spinning up the Algorand infrastructure portion of your development environment. It uses Docker to provide an `algod` instance for connecting to the network of your choosing and an `indexer` instance for querying blockchain data. APIs are exposed by both instances for client access from the SDK. Read more about [Algorand networks](../../get-details/algorand-networks/index.md){target=_blank}, their capabilities and intended use.

!!! Prerequisites
    - Docker Compose ([install guide](https://docs.docker.com/compose/install/){target=_blank})
    - Git ([install guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git){target=_blank})

From a terminal window, install Algorand Sandbox connected to TestNet:

```bash
git clone https://github.com/algorand/sandbox.git
cd sandbox
./sandbox up testnet
```

!!! Warning
    The Algorand Sandbox installation may take a few minutes to complete in order to catch up to the current round on TestNet. To learn more about fast catchup, see [Sync Node Network using Fast Catchup](https://developer.algorand.org/docs/run-a-node/setup/install/#sync-node-network-using-fast-catchup){target=_blank}.

!!! Info
    The `indexer` is enabled **only** for _private networks_. Therefore, all blockchain queries in this guide will use the `algod` API.

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=18){target=_blank}
- [More Information](https://developer.algorand.org/articles/introducing-sandbox-20/){target=_blank}

# Install Go SDK

Algorand provides an SDK for Go. 

!!! Prerequisites
    - Go programming language ([install guide](https://golang.org/doc/install){target=_blank})

From a terminal window, install the Go SDK:

```bash
go get -u github.com/algorand/go-algorand-sdk/...
```

- [`Watch Video`](https://youtu.be/rFG7Zo2JvIY?t=88){target=_blank}
- [`More Information`](https://github.com/algorand/go-algorand-sdk){target=_blank}
 
The SDK is installed and can now interact with the running Algorand Sandbox environment, as configured above.

# Create account

In order to interact with the Algorand blockchain, you must have a funded account on the network. To quickly create an account on Algorand TestNet create a new file **yourFirstTransaction.go** and insert the following code:

```go linenums="1"
package main

import (
    "context"
    json "encoding/json"
    "errors"
    "fmt"
    "strings"

    "github.com/algorand/go-algorand-sdk/client/v2/algod"
    "github.com/algorand/go-algorand-sdk/client/v2/common/models"
    "github.com/algorand/go-algorand-sdk/crypto"
    "github.com/algorand/go-algorand-sdk/mnemonic"
    "github.com/algorand/go-algorand-sdk/transaction"
)

// TODO: insert additional utility functions here

func main() {
    // Create account
    account := crypto.GenerateAccount()
    passphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)
    myAddress := account.Address.String()

    if err != nil {
        fmt.Printf("Error creating transaction: %s\n", err)
    } else {
        fmt.Printf("My address: %s\n", myAddress)
        fmt.Printf("My passphrase: %s\n", passphrase)
        fmt.Println("--> Copy down your address and passpharse for future use.")
        fmt.Println("--> Once secured, press ENTER key to continue...")
        fmt.Scanln()
    }

    // TODO: insert additional codeblocks here
}
```

!!! Info 
    Lines 17 and 35 contain comments about inserting additional code. As you proceed with this guide, ensure the line numbers remain in sync.

!!! Tip
    Make sure to save the generated address and passphrase in a secure location, as they will be used later on.

!!! Warning 
    Never share your mnemonic passphrase or private keys. Production environments require stringent private key management. For more information on key management in community Wallets, click [here](https://developer.algorand.org/docs/community/#wallets). For the open source [Algorand Wallet](https://developer.algorand.org/articles/algorand-wallet-now-open-source/), click [here](https://github.com/algorand/algorand-wallet).

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=97){target=_blank}
- [More Information](https://developer.algorand.org/docs/features/accounts/create/#standalone){target=_blank}
 
# Fund account

The code below prompts to fund the newly generated account. Before sending transactions to the Algorand network, the account must be funded to cover the minimal transaction fees that exist on Algorand. To fund the account use the [Algorand TestNet faucet](https://dispenser.testnet.aws.algodev.network/){target=_blank}. 

```go linenums="35"
// Fund account
fmt.Println("Fund the created account using the Algorand TestNet faucet:\n--> https://dispenser.testnet.aws.algodev.network?account=" + myAddress)
fmt.Println("--> Once funded, press ENTER key to continue...")
fmt.Scanln()
```

!!! Info
    All Algorand accounts require a minimum balance to be registered in the ledger. To read more about Algorand minimum balance see [Account Overview](https://developer.algorand.org/docs/features/accounts/#minimum-balance){target=_blank}

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=138){target=_blank}

# Instantiate client

You must instantiate a client prior to making calls to the API endpoints. The Go SDK implements the client natively using the following code:

```go  linenums="40"
// instantiate algod client to Algorand Sandbox
const algodAddress = "http://localhost:4001"
const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

algodClient, err := algod.MakeClient(algodAddress, algodToken)
if err != nil {
    fmt.Printf("Issue with creating algod client: %s\n", err)
    return
}
```
 
!!! Info
    This guide provides values for `algodAddress` and `algodToken` as specified by Algorand Sandbox. If you want to connect to a third-party service provider, see [Purestake](https://developer.purestake.io/code-samples){target=_blank} or [AlgoExplorer Developer API](https://algoexplorer.io/api-dev/v2){target=_blank} and adjust these values accordingly.
 
- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=149){target=_blank}

# Check account balance

Before moving on to the next step, make sure your account has been funded by the faucet.
 
```go  linenums="50"
//Check account balance
fmt.Printf("My address: %s\n", myAddress)

accountInfo, err := algodClient.AccountInformation(myAddress).Do(context.Background())
if err != nil {
    fmt.Printf("Error getting account info: %s\n", err)
    return
}
fmt.Printf("Account balance: %d microAlgos\n", accountInfo.Amount)
fmt.Println("--> Ensure balance greater than 0, press ENTER key to continue...")
fmt.Scanln()
```

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=161){target=_blank}

# Build transaction

Communication with the Algorand network is performed using transactions. Create a payment transaction sending 9 ALGO from your account to the TestNet faucet address:

```go linenums="62"
// Construct the transaction
txParams, err := algodClient.SuggestedParams().Do(context.Background())
if err != nil {
    fmt.Printf("Error getting suggested tx params: %s\n", err)
    return
}
fromAddr := myAddress
toAddr := "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"
var amount uint64 = 9000000
var minFee uint64 = 1000
note := []byte("Hello World")
genID := txParams.GenesisID
genHash := txParams.GenesisHash
firstValidRound := uint64(txParams.FirstRoundValid)
lastValidRound := uint64(txParams.LastRoundValid)
txn, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
if err != nil {
    fmt.Printf("Error creating transaction: %s\n", err)
    return
}
```

!!! Info
    Algorand supports many transaction types. To see what types are supported see [Transactions](https://developer.algorand.org/docs/features/transactions/){target=_blank}.

[`Watch Video`](https://youtu.be/rFG7Zo2JvIY?t=178){target=_blank}

# Sign transaction

Before the transaction is considered valid, it must be signed by a private key. Use the following code to sign the transaction.

```go linenums="83"
// Sign the transaction
txID, signedTxn, err := crypto.SignTransaction(account.PrivateKey, txn)
if err != nil {
    fmt.Printf("Failed to sign transaction: %s\n", err)
    return
}
fmt.Printf("Signed txid: %s\n", txID)
```

!!! Info
    Algorand provides many ways to sign transactions. To see other ways see [Authorization](https://developer.algorand.org/docs/features/transactions/signatures/#single-signatures){target=_blank}.

[`Watch Video`](https://youtu.be/rFG7Zo2JvIY?t=204){target=_blank}

# Submit transaction

The signed transaction can now be broadcast to the network for validation and inclusion in a future block. The `waitForConfirmation` method polls the `algod` node for the transaction ID to ensure it succeeded.

```go linenums="91"
// Submit the transaction
sendResponse, err := algodClient.SendRawTransaction(signedTxn).Do(context.Background())
if err != nil {
    fmt.Printf("failed to send transaction: %s\n", err)
    return
}
fmt.Printf("Submitted transaction %s\n", sendResponse)

// Wait for confirmation
confirmedTxn, err := waitForConfirmation(txID, algodClient, 4)
if err != nil {
    fmt.Printf("Error waiting for confirmation on txID: %s\n", txID)
    return
}
```

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=216){target=_blank}

# Display completed transaction

Finally, we can query the blockchain for the committed transaction data and display in on the command line. 

```go linenums="106"
// Display completed transaction
txnJSON, err := json.MarshalIndent(confirmedTxn.Transaction.Txn, "", "\t")
if err != nil {
    fmt.Printf("Can not marshall txn data: %s\n", err)
}
fmt.Printf("Transaction information: %s\n", txnJSON)
fmt.Printf("Decoded note: %s\n", string(confirmedTxn.Transaction.Txn.Note))
```

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=232){target=_blank}
 
# Add utility functions

The utility function `waitFoConfirmation` should be inserted between your `imports` and `main()` code blocks:

```go linenums="17"
// Utility function that waits for a given txId to be confirmed by the network
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
```

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=241){target=_blank}
 
# Run the program
 
Save your file and execute the program:

```bash
$ go run yourFirstTransaction.go
```

!!! Warning
    In order for your transaction to be successful, you must fund the generated account during runtime.

!!! Info
	View the confirmed transaction in your web browser by clicking the link to these third-party block explorers and inserting the transactionID within their search bar:
	
	- [AlgoExplorer](https://testnet.algoexplorer.io/){target=_blank}
	- [Goal Seeker](https://goalseeker.purestake.io/algorand/testnet){target=_blank}

- [Watch Video](https://youtu.be/rFG7Zo2JvIY?t=232){target=_blank}
 
# Complete example

If you have any trouble compiling or running your program, please check the complete example below which details how to quickly submit your first transaction.
 
[Run Code](https://replit.com/@Algorand/Getting-Started-with-Go){target=_blank}

[Watch Video](https://youtu.be/rFG7Zo2JvIY?t=){target=_blank}

# Setting up your editor/framework

The Algorand community provides many editors, frameworks, and plugins that can be used to work with the Algorand Network. Tutorials have been created for configuring each of these for use with Algorand. Select your Editor preference below.

* [Setting Up VSCode](https://developer.algorand.org/tutorials/vs-code-go/){target=_blank}
* [Algorand VSCode Extension](https://developer.algorand.org/articles/intro-algorand-studio-algorand-vs-code-extension/){target=_blank}
* [Algo Studio](https://developer.algorand.org/articles/intro-algorand-studio-algorand-vs-code-extension/){target=_blank}
* [AlgoDEA InteliJ Plugin](https://developer.algorand.org/articles/making-development-easier-algodea-intellij-plugin/){target=_blank}
* [Algo Builder Framework](https://developer.algorand.org/articles/introducing-algorand-builder/){target=_blank}