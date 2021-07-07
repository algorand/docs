title: Using the SDKs

This guide covers using stateless smart contracts with the Algorand SDKs. Stateless smart contracts can be used to create contract accounts or to handle account delegation which is described in the [Modes](modes.md) documentation. In either case, the contracts are written in [Transaction Execution Approval Language (TEAL)](../teal/index.md) or with Python using the [PyTeal](../teal/pyteal.md) library.

Each SDK's install process is discussed in the [SDK Reference](../../../reference/sdks/index.md) documentation.

!!! info
    The example code snippets are provided throughout this page and are abbreviated for conciseness and clarity. Full running code examples for each SDK are available within the GitHub repo for V1 and V2 at [/examples/smart_contracts](https://github.com/algorand/docs/tree/master/examples/smart_contracts) and for [download](https://github.com/algorand/docs/blob/master/examples/smart_contracts/smart_contracts.zip?raw=true) (.zip).


# Compiling TEAL program from SDKs
Before a TEAL program can be used, it must be compiled. SDKs provide this capability. If using the `goal` tool see the [goal TEAL walkthrough](walkthrough.md) documentation for this process.  The examples in this section read a file called `sample.teal` which contains one line of TEAL code, `int 0` . This will always return `false`. So, any transactions that use this TEAL file will fail. 


To use the SDK compile command, the [config settings](../../../../reference/node/config/) may need to be modified to set a value for `EnableDeveloperAPI`, which should be set to `true`. The default is false. If using the [sandbox](../../../../build-apps/setup/#2-use-docker-sandbox), the following modification is already made. If [running your own node](../../../../build-apps/setup/#3-run-your-own-node), you may see an error similar to "compile was not enabled in the configuration file by setting the EnableDeveloperAPI to true". Make the following modification to the `config.json` file located in the node’s data directory. First, if there is not a `config.json`, make a copy of the `config.json.example` file.

```
$ goal node stop -d data
$ cd data
$ cp config.json.example config.json
```

Then edit the config.json file and replace `false` on the line

`"EnableDeveloperAPI": false,`

with `true`

`"EnableDeveloperAPI": true,`

Restart the node.

```
$ goal node start -d data
$ goal node status -d data
```

The following TEAL code is used in the examples below.

`sample.teal`
```
// This code is meant for learning purposes only
// It should not be used in production
int 0
```

!!! info
    The samples on this page use a [docker sandbox install](https://developer.algorand.org/docs/build-apps/setup/#2-use-docker-sandbox)

```javascript tab="JavaScript"
const fs = require('fs');
const path = require('path');
const algosdk = require('algosdk');
// We assume that testing is done off of sandbox, hence the settings below
const token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
const server = "http://localhost";
const port = 4001;

// create v2 client
const algodClient = new algosdk.Algodv2(token, server, port);

const main = async () => {
    // Read file for Teal code - int 0
    const filePath = path.join(__dirname, 'sample.teal');
    const data = fs.readFileSync(filePath);

    // Compile teal
    const results = await algodClient.compile(data).do();
    return results;
};

main().then((results) => {
    // Print results
    console.log("Hash = " + results.hash);
    console.log("Result = " + results.result);
}).catch(e => {
    const error = e.body && e.body.message ? e.body.message : e;
    console.log(error);
});

// output would be similar to this... 
// Hash : KI4DJG2OOFJGUERJGSWCYGFZWDNEU2KWTU56VRJHITP62PLJ5VYMBFDBFE
// Result : ASABACI=
```

```python tab="Python"
# compile teal code
from algosdk.v2client import algod

try:
    # create an algod client
    algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    algod_address = "http://localhost:4001"
    algod_client = algod.AlgodClient(algod_token, algod_address)

    # int 0 - sample.teal
    myprogram = "sample.teal"
    # read teal program
    data = open(myprogram, 'r').read()
    # compile teal program
    response = algod_client.compile(data)
    # print(response)
    print("Response Result = ",response['result'])
    print("Response Hash = ",response['hash'])
except Exception as e:
    print(e)

# results should look similar to this:
# Response Result = ASABACI=
# Response Hash = KI4DJG2OOFJGUERJGSWCYGFZWDNEU2KWTU56VRJHITP62PLJ5VYMBFDBFE
```

```java tab="Java"
package com.algorand.javatest.smart_contracts;

import com.algorand.algosdk.v2.client.common.AlgodClient;
import java.nio.file.Files;
import java.nio.file.Paths;
import com.algorand.algosdk.v2.client.model.CompileResponse;

public class CompileTeal {
// Utility function to update changing block parameters
public AlgodClient client = null;

// utility function to connect to a node
private AlgodClient connectToNetwork() {

    // Initialize an algod client
    final Integer ALGOD_PORT = 4001;
    final String ALGOD_API_ADDR = "localhost";
    final String ALGOD_API_TOKEN = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
    AlgodClient client = new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);
    return client;
}

public void compileTealSource() throws Exception {
    // Initialize an algod client
    if (client == null)
        this.client = connectToNetwork();

    // read file - int 0
    byte[] data = Files.readAllBytes(Paths.get("./sample.teal"));
    // compile
    CompileResponse response = client.TealCompile().source(data).execute().body();
    // print results
    System.out.println("response: " + response);
    System.out.println("Hash: " + response.hash); 
    System.out.println("Result: " + response.result); 
}

public static void main(final String args[]) throws Exception {
    CompileTeal t = new CompileTeal();
    t.compileTealSource();
}

}
// Output should look similar to this... 
// response:
// {"hash":"KI4DJG2OOFJGUERJGSWCYGFZWDNEU2KWTU56VRJHITP62PLJ5VYMBFDBFE","result":"ASABACI="}
// Hash: KI4DJG2OOFJGUERJGSWCYGFZWDNEU2KWTU56VRJHITP62PLJ5VYMBFDBFE 
// Result: ASABACI=
```

```go tab="Go"
package main

import (

	"context"
	"io/ioutil"
	"log"
	"fmt"
	"os"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func main() {

	const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	const algodAddress = "http://localhost:4001"

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}
	// int 0 in sample.teal
	file, err := os.Open("./sample.teal")
    if err != nil {
        log.Fatal(err)
    }    
	defer file.Close()
	tealFile, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("failed to teal file: %s\n", err)
		return}
    // compile teal program
	response, err := algodClient.TealCompile(tealFile).Do(context.Background())
    // print response	
	fmt.Printf("Hash = %s\n", response.Hash)
	fmt.Printf("Result = %s\n", response.Result)
}
// results should look similar to
// Hash = KI4DJG2OOFJGUERJGSWCYGFZWDNEU2KWTU56VRJHITP62PLJ5VYMBFDBFE
// Result = ASABACI=
```

Once a TEAL program is compiled, the bytes of the program can be used as a parameter to the LogicSig method. Most of the SDKs support the bytes encoded in base64 or hexadecimal format. 

The binary bytes are used in the SDKs as shown below. If using the `goal` command-line tool to compile the TEAL code, these same bytes can be retrieved using the following commands. 


``` bash
//simple.teal contains int 0
//hexdump 
$ hexdump -C simple.teal.tok
00000000  01 20 01 00 22                                    |. .."|
00000005
//base64 format
$ cat simple.teal.tok | base64
ASABACI=
```

The response result from the TEAL `compile` command above is used to create the `program` variable. This variable can then be used as an input parameter to the function to make a logic signature.

```javascript tab="JavaScript"
    const program = new Uint8Array(Buffer.from(results.result , "base64"));
    const lsig = algosdk.makeLogicSig(program);   
```

```python tab="Python"
    import base64
    from algosdk.future.transaction import LogicSig

    programstr = response['result']
    t = programstr.encode()
    program = base64.decodebytes(t)
    lsig = LogicSig(program)
```

```java tab="Java"
    // byte[] program = {
    //     0x01, 0x20, 0x01, 0x00, 0x22  // int 0
    // };
    byte[] program = Base64.getDecoder().decode(response.result.toString());
    LogicsigSignature lsig = new LogicsigSignature(program, null);
```

```go tab="Go"
    // program, err :=  base64.StdEncoding.DecodeString("ASABACI=")
    program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    var sk ed25519.PrivateKey
    var ma crypto.MultisigAccount
    var args [][]byte
    lsig, err := crypto.MakeLogicSig(program, args, sk, ma)  
```

# Passing Parameters using the SDKs
The SDKs require that parameters to a stateless smart contract TEAL program be in byte arrays. This byte array is passed to the method that creates the logic signature. Currently, these parameters must be either unsigned integers or binary strings. If comparing a constant string in TEAL, the constant within the TEAL program must be encoded in hex or base64. See the TEAL tab below for a simple example of comparing the string argument used in the other examples. SDK native language functions can be used to encode the parameters to the TEAL program correctly. The example below illustrates both a string parameter and an integer.

!!! info
    The samples show setting parameters at the creation of the logic signature. These parameters can be changed at the time of submitting the transaction.

```javascript tab="JavaScript"
    // string parameter
    const args = [];
    args.push([...Buffer.from("my string")]);
    const lsig = algosdk.makeLogicSig(program, args);
    
    // integer parameter
    const args = [];
    args.push(algosdk.encodeUint64(123));
    const lsig = algosdk.makeLogicSig(program, args);
```

```python tab="Python"
    from algosdk.future.transaction import LogicSig

    # string parameter
    arg_str = "my string"
    arg1 = arg_str.encode()
    lsig = LogicSig(program, args=[arg1])
    # integer parameter
    arg1 = (123).to_bytes(8, 'big')
    lsig = LogicSig(program, args=[arg1])
```

```java tab="Java"
    // string parameter
    ArrayList<byte[]> teal_args = new ArrayList<byte[]>();
    String orig = "my string";
    teal_args.add(orig.getBytes());
    LogicsigSignature lsig = new LogicsigSignature(program, teal_args);    

    // integer parameter
    ArrayList<byte[]> teal_args = new ArrayList<byte[]>();
    byte[] arg1 = {123};
    teal_args.add(arg1);
    LogicsigSignature lsig = new LogicsigSignature(program, teal_args);
```

```go tab="Go"
    // string parameter
	args := make([][]byte, 1)
	args[0] = []byte("my string")
    lsig, err := crypto.MakeLogicSig(program, args, sk, ma)
    
    // integer parameter
    args := make([][]byte, 1)
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], 123)
	args[0] = buf[:]
```

```text tab="TEAL"
// Never use this code for a real transaction
// for educational purposes only
// string compare
arg_0
byte b64 bXkgc3RyaW5n
==
// integer compare
arg_0
btoi
int 123
==
```

# Contract Account SDK usage
ASC1 Contract accounts are used to allow TEAL logic to determine when outgoing account transactions are approved. The compiled TEAL program produces an Algorand Address, which is funded with Algos or Algorand Assets. As the receiver of a transaction, these accounts function as any other account. When the account is specified as the sender in a transaction, the TEAL logic is evaluated and determines if the transaction is approved. The [ASC1 Usage Modes](./modes.md) documentation explains ASC1 modes in more detail. 

TEAL contract account transactions where the sender is set to the contract account, function much in the same way as normal Algorand [transactions](../../transactions/index.md). The major difference is that instead of the transaction being signed with a private key, the transaction is signed with a [logic signature](./modes.md#logic-signatures). See [Transaction](../../transactions/index.md) documentation for details on setting up a payment transaction.

Contract Accounts are created by compiling the TEAL logic. Once the contract account is created, it can be used as any other address. To send tokens or assets from the account the transaction must be signed by a Logic Signature. From an SDK standpoint, the following process should be used.

* Load the Program Bytes into the SDK.
* Create a Logic Signature based on the program.
* Create the Transaction.
* Set the `from` transaction property to the contract address.
* Sign the Transaction with the Logic Signature.
* Send the Transaction to the network.

<center>![Transaction From Contract Account](../../../imgs/asc1_sdk_usage-1.png)</center>
<center>*Transaction From Contract Account*</center>

The following example illustrates compiling a TEAL program with one argument and signing a transaction with a created logic signature. The example TEAL program `samplearg.teal` takes one argument. Information on TEAL Opcodes can be [found here](https://developer.algorand.org/docs/reference/teal/opcodes/#opcodes). 

`samplearg.teal`

```
// samplearg.teal
// This code is meant for learning purposes only
// It should not be used in production
arg_0
btoi
int 123
==
```

```javascript tab="JavaScript"
const fs = require('fs');
const path = require('path');
const algosdk = require('algosdk');
// We assume that testing is done off of sandbox, hence the settings below
const token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
const server = "http://localhost";
const port = 4001;

// create v2 client
const algodClient = new algosdk.Algodv2(token, server, port);

// Function used to wait for a tx confirmation
const waitForConfirmation = async function (algodClient, txId) {
    const response = await algodClient.status().do();
    let lastround = response["last-round"];
    while (true) {
        const pendingInfo = await algodClient.pendingTransactionInformation(txId).do();
        if (pendingInfo["confirmed-round"] !== null && pendingInfo["confirmed-round"] > 0) {
            //Got the completed Transaction
            console.log("Transaction " + txId + " confirmed in round " + pendingInfo["confirmed-round"]);
            break;
        }
        lastround++;
        await algodClient.statusAfterBlock(lastround).do();
    }
};

const main = async () => {
    // get suggested parameters
    const params = await algodClient.getTransactionParams().do();
    // comment out the next two lines to use suggested fee
    params.fee = 1000;
    params.flatFee = true;

    const filePath = path.join(__dirname, 'samplearg.teal');
      
    const data = fs.readFileSync(filePath);
    const  results = await algodClient.compile(data).do();
    console.log("Hash = " + results.hash);
    console.log("Result = " + results.result);

    const program = new Uint8Array(Buffer.from(results.result, "base64"));
    // Use this if no args
    // const lsig = algosdk.makeLogicSig(program);

    // Initialize arguments array
    const args = [];

    // String parameter
    // args.push([...Buffer.from("my string")]);

    // Integer parameter
    args.push(algosdk.encodeUint64(123));


    const lsig = algosdk.makeLogicSig(program, args);
    console.log("lsig : " + lsig.address());   

    // create a transaction
    const sender = lsig.address();
    const receiver = "<receiver-address>";
    const amount = 10000;
    const closeToRemaninder = undefined;
    const note = undefined;
    const txn = algosdk.makePaymentTxnWithSuggestedParams(sender, receiver, amount, closeToRemaninder, note, params)

    const rawSignedTxn = algosdk.signLogicSigTransactionObject(txn, lsig);

    // send raw LogicSigTransaction to network
    const tx = await algodClient.sendRawTransaction(rawSignedTxn.blob).do();
    console.log("Transaction : " + tx.txId);   
    await waitForConfirmation(algodClient, tx.txId);
};

main().then().catch(e => {
    const error = e.body && e.body.message ? e.body.message : e;
    console.log(error);
});
```

```python tab="Python"
from algosdk.v2client import algod
from algosdk.future.transaction import PaymentTxn, LogicSig, LogicSigTransaction
import base64

def wait_for_confirmation(client, txid):
    """
    Utility function to wait until the transaction is
    confirmed before proceeding.
    """
    last_round = client.status().get('last-round')
    txinfo = client.pending_transaction_info(txid)
    while not (txinfo.get('confirmed-round') and txinfo.get('confirmed-round') > 0):
        print("Waiting for confirmation")
        last_round += 1
        client.status_after_block(last_round)
        txinfo = client.pending_transaction_info(txid)
    print("Transaction {} confirmed in round {}.".format(
        txid, txinfo.get('confirmed-round')))
    return txinfo

try:
    # Create an algod client
    algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" 
    algod_address = "http://localhost:4001" 
    algod_client = algod.AlgodClient(algod_token, algod_address)

    receiver = "<receiver-address>" 
   
    myprogram = "samplearg.teal"
    # Read TEAL program
    data = open(myprogram, 'r').read()
    # Compile TEAL program
    response = algod_client.compile(data)
    # Print(response)
    print("Response Result = ", response['result'])
    print("Response Hash = ", response['hash'])

    # Create logic sig
    programstr = response['result']
    t = programstr.encode()
    # program = b"hex-encoded-program"
    program = base64.decodebytes(t)
    print(program)

    # Create arg to pass if TEAL program requires an arg
    # if not, omit args param
    arg1 = (123).to_bytes(8, 'big')
    lsig = LogicSig(program, args=[arg1])
    print("lsig Address: " + lsig.address())
    sender = lsig.address()

    # Get suggested parameters
    params = algod_client.suggested_params()
    # Comment out the next two (2) lines to use suggested fees
    params.flat_fee = True
    params.fee = 1000
    
    # Build transaction  
    amount = 10000 
    closeremainderto = None

    # Create a transaction
    txn = PaymentTxn(
        sender, params, receiver, amount, closeremainderto)

    # Create the LogicSigTransaction with contract account LogicSig
    lstx = LogicSigTransaction(txn, lsig)

    # Send raw LogicSigTransaction to network
    txid = algod_client.send_transaction(lstx)
    print("Transaction ID: " + txid)    
    wait_for_confirmation(algod_client, txid) 

except Exception as e:
    print(e)
```

```java tab="Java"
package com.algorand.javatest.smart_contracts;

import com.algorand.algosdk.account.Account;
import com.algorand.algosdk.algod.client.ApiException;

import com.algorand.algosdk.crypto.Address;
import com.algorand.algosdk.crypto.LogicsigSignature;
import com.algorand.algosdk.transaction.SignedTransaction;
import com.algorand.algosdk.transaction.Transaction;
import com.algorand.algosdk.util.Encoder;
import com.algorand.algosdk.v2.client.common.AlgodClient;
import com.algorand.algosdk.v2.client.common.Response;
import com.algorand.algosdk.v2.client.model.PendingTransactionResponse;
import com.algorand.algosdk.v2.client.model.PostTransactionsResponse;
import com.algorand.algosdk.v2.client.model.TransactionParametersResponse;
import java.nio.file.Files;
import java.nio.file.Paths;
import org.json.JSONObject;
import java.util.ArrayList;
import java.util.Base64;
import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import com.algorand.algosdk.v2.client.model.CompileResponse;


public class ContractAccount {
    // Utility function to update changing block parameters
    public AlgodClient client = null;

    // utility function to connect to a node
    private AlgodClient connectToNetwork() {

        // Initialize an algod client
        final String ALGOD_API_ADDR = "localhost";
        final Integer ALGOD_PORT = 4001;
        final String ALGOD_API_TOKEN = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

        AlgodClient client = new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);
        return client;
    }

    // utility function to wait on a transaction to be confirmed

    public void waitForConfirmation(String txID) throws Exception {
        if (client == null)
            this.client = connectToNetwork();
        Long lastRound = client.GetStatus().execute().body().lastRound;
        while (true) {
            try {
                // Check the pending transactions
                Response<PendingTransactionResponse> pendingInfo = client.PendingTransactionInformation(txID).execute();
                if (pendingInfo.body().confirmedRound != null && pendingInfo.body().confirmedRound > 0) {
                    // Got the completed Transaction
                    System.out.println(
                            "Transaction " + txID + " confirmed in round " + pendingInfo.body().confirmedRound);
                    break;
                }
                lastRound++;
                client.WaitForBlock(lastRound).execute();
            } catch (Exception e) {
                throw (e);
            }
        }
    }

    public void contractAccountExample() throws Exception {
        // Initialize an algod client
        if (client == null)
            this.client = connectToNetwork();

        // Set the receiver
        final String RECEIVER = "<receiver-address>";

        // Read program from file samplearg.teal
        byte[] source = Files.readAllBytes(Paths.get("./samplearg.teal"));
        // compile
        CompileResponse response = client.TealCompile().source(source).execute().body();
        // print results
        System.out.println("response: " + response);
        System.out.println("Hash: " + response.hash);
        System.out.println("Result: " + response.result);
        byte[] program = Base64.getDecoder().decode(response.result.toString());

        // create logic sig
        // integer parameter
        ArrayList<byte[]> teal_args = new ArrayList<byte[]>();
        byte[] arg1 = { 123 };
        teal_args.add(arg1);
        LogicsigSignature lsig = new LogicsigSignature(program, teal_args);
        // For no args use null as second param
        // LogicsigSignature lsig = new LogicsigSignature(program, null);
        System.out.println("lsig address: " + lsig.toAddress());        
        TransactionParametersResponse params = client.TransactionParams().execute().body();
        // create a transaction
        String note = "Hello World";
        Transaction txn = Transaction.PaymentTransactionBuilder()
                .sender(lsig
                        .toAddress())
                .note(note.getBytes())
                .amount(100000)
                .receiver(new Address(RECEIVER))
                .suggestedParams(params)
                .build();   
        try {
            // create the LogicSigTransaction with contract account LogicSig
            SignedTransaction stx = Account.signLogicsigTransaction(lsig, txn);
            // send raw LogicSigTransaction to network
            byte[] encodedTxBytes = Encoder.encodeToMsgPack(stx);
            String id = client.RawTransaction().rawtxn(encodedTxBytes).execute().body().txId;
            // Wait for transaction confirmation
            waitForConfirmation(id);
            System.out.println("Successfully sent tx with id: " + id);
            // Read the transaction
            PendingTransactionResponse pTrx = client.PendingTransactionInformation(id).execute().body(); 
            JSONObject jsonObj = new JSONObject(pTrx.toString());
            System.out.println("Transaction information (with notes): " + jsonObj.toString(2)); // pretty print
            System.out.println("Decoded note: " + new String(pTrx.txn.tx.note));
        } catch (ApiException e) {
            System.err.println("Exception when calling algod#rawTransaction: " + e.getResponseBody());
        }
    }
    public static void main(final String args[]) throws Exception {
        ContractAccount t = new ContractAccount();
        t.contractAccountExample();
    }
}
```

```go tab="Go"
package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
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

	file, err := os.Open("./samplearg.teal")
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()
    tealFile, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Printf("failed to read file: %s\n", err)
		return}
		
    response, err := algodClient.TealCompile(tealFile).Do(context.Background())
    	if err != nil {
		fmt.Printf("Error on TEAL compile: %s\n", err)
		return
	}
    fmt.Printf("Hash = %s\n", response.Hash)
    fmt.Printf("Result = %s\n", response.Result)
    
    program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    
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
	const receiver = "receiver-address"
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

		
	// Submit the raw transaction to network
	transactionID, err := algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		fmt.Printf("Sending failed with %v\n", err)
	}
    // Wait for transaction to be confirmed
    waitForConfirmation(txID, algodClient)
    fmt.Printf("Transaction ID: %v\n", transactionID)
}
```

# Account Delegation SDK Usage
ASC1 allows TEAL logic to be used to delegate signature authority. This allows specific accounts or multi-signature accounts to sign logic that allows transactions from the account to be approved based on the TEAL logic. The [ASC1 Usage Modes](./modes.md) documentation explains ASC1 modes in more detail. 

Delegated transactions are special transactions where the `sender` also signs the logic and the transaction is then signed with the [logic signature](./modes.md#logic-signature). In all other aspects, the transaction functions as any other transaction. See [Transaction](../../transactions/index.md) documentation for details on setting up a payment transaction.

Delegated Logic Signatures require that the logic signature be signed from a specific account or a multi-signature account. The TEAL program is first loaded, then a Logic Signature is created and then the Logic Signature is signed by a specific account or multi-signature account. The transaction is created as normal. The transaction is then signed with the Logic Signature. From an SDK standpoint, the following process should be used.

* Load the Program Bytes into the SDK.
* Create a Logic Signature based on the program.
* Sign The Logic Signature with a specific account
* Create the Transaction.
* Set the `from` transaction property to the Address that signed the logic.
* Sign the Transaction with the Logic Signature.
* Send the Transaction to the network.

<center>![Delegated Signature Transaction](../../../imgs/asc1_sdk_usage-2.png)</center>
<center>*Delegated Signature Transaction*</center>

The following example illustrates signing a transaction with a created logic signature that is signed by a specific account.

```javascript tab="JavaScript"
const main = async () => {
    // get suggested parameters
    const params = await algodClient.getTransactionParams().do();
    params.fee = 1000;
    params.flatFee = true;

    const filePath = path.join(__dirname, 'samplearg.teal');
    const data = fs.readFileSync(filePath);
    const  results = await algodClient.compile(data).do();
    const program = new Uint8Array(Buffer.from(results.result, "base64"));

    const args = [];
    args.push(algosdk.encodeUint64(123));
   
    const lsig = algosdk.makeLogicSig(program, args); 
    
    // *** Begin account delegation changes ***
    // import your private key mnemonic
    const PASSPHRASE = "< mnemonic >";
    const myAccount = algosdk.mnemonicToSecretKey(PASSPHRASE);

    // sign the logic signature with an account sk
    lsig.sign(myAccount.sk);
    // *** End account delegation changes ***

    // Setup a transaction
    const sender = myAccount.addr;
    const receiver = "GVVIH6IE3ZAFLIQF6UFBOH4TI5ZOPHFRMDQUTXQRKZ6XWSRTVRDBQ4GZLY";
    const amount = 10000;
    const closeToRemaninder = undefined;
    const note = undefined;
    const txn = algosdk.makePaymentTxnWithSuggestedParams(sender, receiver, amount, closeToRemaninder, note, params)
    const rawSignedTxn = algosdk.signLogicSigTransactionObject(txn, lsig);
   
    const tx = await algodClient.sendRawTransaction(rawSignedTxn.blob).do();
    console.log("Transaction : " + tx.txId);    
    await waitForConfirmation(algodClient, tx.txId);
};
```

```python tab="Python"
from algosdk import account, mnemonic
from algosdk.v2client import algod
from algosdk.future.transaction import PaymentTxn, LogicSig, LogicSigTransaction
import base64

def wait_for_confirmation(client, txid):
    """
    Utility function to wait until the transaction is
    confirmed before proceeding.
    """
    last_round = client.status().get('last-round')
    txinfo = client.pending_transaction_info(txid)
    while not (txinfo.get('confirmed-round') and txinfo.get('confirmed-round') > 0):
        print("Waiting for confirmation")
        last_round += 1
        client.status_after_block(last_round)
        txinfo = client.pending_transaction_info(txid)
    print("Transaction {} confirmed in round {}.".format(
        txid, txinfo.get('confirmed-round')))
    return txinfo

try:
    # Create an algod client
    algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    algod_address = "http://localhost:4001"
    algod_client = algod.AlgodClient(algod_token, algod_address)
    receiver = "<receiver-address>" 

    myprogram = "samplearg.teal"

    # Read TEAL program
    data = open(myprogram, 'r').read()

    # Compile TEAL program
    response = algod_client.compile(data)
    # Print(response)
    print("Response Result = ", response['result'])
    print("Response Hash = ", response['hash'])

    # Create logic sig
    programstr = response['result']
    t = programstr.encode()
    # program = b"hex-encoded-program"   
    program = base64.decodebytes(t)
    print(program)

    # Create arg to pass if TEAL program requires an arg,
    # if not, omit args param
    # string parameter
    arg1 = (123).to_bytes(8, 'big')
    lsig = LogicSig(program, args=[arg1])

    # Recover the account that is wanting to delegate signature
    passphrase = "<25-word-mnemonic>"
    sk = mnemonic.to_private_key(passphrase)
    addr = account.address_from_private_key(sk)
    print("Address of Sender/Delegator: " + addr)

    # Sign the logic signature with an account sk
    lsig.sign(sk)
 
    # Get suggested parameters
    params = algod_client.suggested_params()
    # Comment out the next two (2) lines to use suggested fees
    params.flat_fee = True
    params.fee = 1000

    # Build transaction
    amount = 10000 
    closeremainderto = None
 
    # Create a transaction
    txn = PaymentTxn(
        addr, params, receiver, amount, closeremainderto)
    # Create the LogicSigTransaction with contract account LogicSig
    lstx = LogicSigTransaction(txn, lsig)
    # transaction.write_to_file([lstx], "simple.stxn")
    # Send raw LogicSigTransaction to network
    txid = algod_client.send_transaction(lstx)
    print("Transaction ID: " + txid)
    wait_for_confirmation(algod_client, txid)
except Exception as e:
    print(e)
```

```java tab="Java"
package com.algorand.javatest.smart_contracts;

import com.algorand.algosdk.account.Account;
import com.algorand.algosdk.algod.client.ApiException;

import com.algorand.algosdk.crypto.Address;
import com.algorand.algosdk.crypto.LogicsigSignature;
import com.algorand.algosdk.transaction.SignedTransaction;
import com.algorand.algosdk.transaction.Transaction;
import com.algorand.algosdk.util.Encoder;
import com.algorand.algosdk.v2.client.common.AlgodClient;
import com.algorand.algosdk.v2.client.common.Response;
import com.algorand.algosdk.v2.client.model.PendingTransactionResponse;
import com.algorand.algosdk.v2.client.model.PostTransactionsResponse;
import com.algorand.algosdk.v2.client.model.TransactionParametersResponse;
import java.util.Base64;
import org.json.JSONObject;
import java.util.ArrayList;
import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import com.algorand.algosdk.v2.client.model.CompileResponse;

public class AccountDelegation {
// Utility function to update changing block parameters
public AlgodClient client = null;

// utility function to connect to a node
private AlgodClient connectToNetwork() {

// Initialize an algod client
// sandbox
final String ALGOD_API_ADDR = "localhost";
final Integer ALGOD_PORT = 4001;
final String ALGOD_API_TOKEN = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";


AlgodClient client = new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);
return client;
}

// utility function to wait on a transaction to be confirmed
public void waitForConfirmation(String txID) throws Exception {
if (client == null)
    this.client = connectToNetwork();
Long lastRound = client.GetStatus().execute().body().lastRound;
while (true) {
    try {
        // Check the pending transactions
        Response<PendingTransactionResponse> pendingInfo = client.PendingTransactionInformation(txID).execute();
        if (pendingInfo.body().confirmedRound != null && pendingInfo.body().confirmedRound > 0) {
            // Got the completed Transaction
            System.out.println(
                    "Transaction " + txID + " confirmed in round " + pendingInfo.body().confirmedRound);
            break;
        }
        lastRound++;
        client.WaitForBlock(lastRound).execute();
    } catch (Exception e) {
        throw (e);
    }
}
}

public void accountDelegationExample() throws Exception {
// Initialize an algod client
if (client == null)
    this.client = connectToNetwork();
// import your private key mnemonic and address
final String SRC_ACCOUNT = "<25-word-mnemonic>";

Account src = new Account(SRC_ACCOUNT);
// Set the receiver
final String RECEIVER = "<receiver-address>";

// Read program from file samplearg.teal
byte[] source = Files.readAllBytes(Paths.get("./samplearg.teal"));


// compile
CompileResponse response = client.TealCompile().source(source).execute().body();
// print results
System.out.println("response: " + response);
System.out.println("Hash: " + response.hash);
System.out.println("Result: " + response.result);
byte[] program = Base64.getDecoder().decode(response.result.toString());

// create logic sig

// string parameter
// ArrayList<byte[]> teal_args = new ArrayList<byte[]>();
// String orig = "my string";
// teal_args.add(orig.getBytes());
// LogicsigSignature lsig = new LogicsigSignature(program, teal_args);

// integer parameter
ArrayList<byte[]> teal_args = new ArrayList<byte[]>();
byte[] arg1 = { 123 };
teal_args.add(arg1);
LogicsigSignature lsig = new LogicsigSignature(program, teal_args);
//    For no args use null as second param
//    LogicsigSignature lsig = new LogicsigSignature(program, null);
// sign the logic signature with an account sk
src.signLogicsig(lsig);
TransactionParametersResponse params = client.TransactionParams().execute().body();
// create a transaction
String note = "Hello World";
Transaction txn = Transaction.PaymentTransactionBuilder()
        .sender(src.getAddress())
        .note(note.getBytes())
        .amount(100000)
        .receiver(new Address(RECEIVER))
        .suggestedParams(params)
        .build();   
try {
    // create the LogicSigTransaction with contract account LogicSig
    SignedTransaction stx = Account.signLogicsigTransaction(lsig, txn);
    // send raw LogicSigTransaction to network
    byte[] encodedTxBytes = Encoder.encodeToMsgPack(stx);

    String id = client.RawTransaction().rawtxn(encodedTxBytes).execute().body().txId;
    // Wait for transaction confirmation
    waitForConfirmation(id);

    System.out.println("Successfully sent tx with id: " + id);
    // Read the transaction
    PendingTransactionResponse pTrx = client.PendingTransactionInformation(id).execute().body();

    JSONObject jsonObj = new JSONObject(pTrx.toString());
    System.out.println("Transaction information (with notes): " + jsonObj.toString(2)); // pretty print
    System.out.println("Decoded note: " + new String(pTrx.txn.tx.note));
} catch (ApiException e) {
    System.err.println("Exception when calling algod#rawTransaction: " + e.getResponseBody());
}

}
public static void main(final String args[]) throws Exception {

AccountDelegation t = new AccountDelegation();
t.accountDelegationExample();
}
}
```

```go tab="Go"
package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
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
    // sandbox
    const algodAddress = "http://localhost:4001"
    const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

    // Create an algod client
    algodClient, err := algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        fmt.Printf("failed to make algod client: %s\n", err)
        return
    }	
    // Get private key for sender address
    PASSPHRASE := "<25-word-mnemonic>"		
    sk, err := mnemonic.ToPrivateKey(PASSPHRASE)	
    pk := sk.Public()
    var a types.Address
    cpk := pk.(ed25519.PublicKey)
    copy(a[:], cpk[:])
    fmt.Printf("Address: %s\n", a.String())	
    sender := a.String()

    // Create logic signature
    var ma crypto.MultisigAccount
    // file, err := os.Open("<PLACEHOLDER>")
    file, err := os.Open("./samplearg.teal")
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
    // lsig, err := crypto.MakeLogicSig(program, args, sk, m a)

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

    
    // Construct the transaction
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
    note := []byte("Hello World")
    genID := txParams.GenesisID
    genHash := txParams.GenesisHash
    firstValidRound := uint64(txParams.FirstRoundValid)
    lastValidRound := uint64(txParams.LastRoundValid)
    tx, err := transaction.MakePaymentTxnWithFlatFee(
        sender, receiver, fee, amount, firstValidRound, lastValidRound,
        note, "", genID, genHash )

    txID, stx, err := crypto.SignLogicsigTransaction(lsig, tx)
    if err != nil {
        fmt.Printf("Signing failed with %v", err)
        return
    }
    fmt.Printf("Signed tx: %v\n", txID)

    // Submit the raw transaction as normal 
    transactionID, err := algodClient.SendRawTransaction(stx).Do(context.Background())
    if err != nil {
        fmt.Printf("Sending failed with %v\n", err)
    }
    // Wait for transaction to be confirmed
    waitForConfirmation(txID, algodClient)
    fmt.Printf("Transaction ID: %v\n", transactionID)
    }
```

!!! Notice
    The samplearg.teal file will compile to the address UVBYHRZIHUNUELDO6HWUAHOZF6G66W6T3JOXIIUSV3LDSBWVCFZ6LM6NCA, please fund this address with at least 11000 microALGO else executing the sample code as written will result in an overspend response from the network node.

!!! info
    The example code snippets are provided throughout this page and are abbreviated for conciseness and clarity. Full running code examples for each SDK are available within the GitHub repo for V1 and V2 at [/examples/smart_contracts](https://github.com/algorand/docs/tree/master/examples/smart_contracts) and for [download](https://github.com/algorand/docs/blob/master/examples/smart_contracts/smart_contracts.zip?raw=true) (.zip).


