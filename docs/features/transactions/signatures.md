title: Authorization

In the [Transactions Section](../transactions/index.md), you learned how transactions are composed. In this section you will learn how to authorize them. 

Before a transaction is sent to the network, it must first be authorized by the [sender](../../reference/transactions.md#sender). Authorization occurs through the addition of a **signature** to the transaction object. Specifically, a transaction object, when signed, is wrapped in a [`SignedTxn`](../../reference/transactions.md#signed-transaction) object that includes the [transaction](../../reference/transactions.md#txn) and a type of [signature](../../reference/transactions.md#sig). 

There are three types of signatures:

- [Single Signatures](#single-signatures)
- [Multisignatures](#multisignatures)
- [Logic Signatures](#logic-signatures)



# Single Signatures
A single signature corresponds to a signature from the private key of an [Algorand public/private key pair](../accounts/index.md#keys-and-addresses).

This is an example of a transaction signed by an Algorand private key displayed with `goal clerk inspect` command:

```json
{
  "sig": "ynA5Hmq+qtMhRVx63pTO2RpDrYiY1wzF/9Rnnlms6NvEQ1ezJI/Ir9nPAT6+u+K8BQ32pplVrj5NTEMZQqy9Dw==",
  "txn": {
    "amt": 10000000,
    "fee": 1000,
    "fv": 4694301,
    "gen": "testnet-v1.0",
    "gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
    "lv": 4695301,
    "rcv": "QC7XT7QU7X6IHNRJZBR67RBMKCAPH67PCSX4LYH4QKVSQ7DQZ32PG5HSVQ",
    "snd": "EW64GC6F24M7NDSC5R3ES4YUVE3ZXXNMARJHDCCCLIHZU6TBEOC7XRSBG4",
    "type": "pay"
  }
}
```
This transaction sends 1 Algo from `"EW64GC..."` to `"QC7XT7..."` on TestNet. The transaction was signed with the private key that corresponds to the `"snd"` address of `"EW64GC..."`. The base64 encoded signature is shown as the value of the [`"sig"`](../../reference/transactions.md#sig) field.

**Related How-To**

- [Sign a transaction with your private key](../../build-apps/hello_world.md#sign-the-transaction)

# Multisignatures

When the [sender](../../reference/transactions.md#sender) of a transaction is the address of a [multisignature account](../accounts/create.md#multisignature) then authorization requires a subset of signatures, _equal to or greater than the threshold value_, from the associated private keys of the addresses that multisignature account is composed of. See [Multisignature Accounts](../accounts/create.md#multisignature) for details on how to configure a multisignature account.

!!! important
	Upon signing, either the signing agent or the transaction needs to know the composition of the multisignature account, i.e. the ordered addresses, threshold, and version. 

Here is what the same transaction above would look like if sent from a 2/3 multisig account.

```json
{
  "msig": {
    "subsig": [
      {
        "pk": "SYGHTA2DR5DYFWJE6D4T34P4AWGCG7JTNMY4VI6EDUVRMX7NG4KTA2WMDA"
      },
      {
        "pk": "VBDMPQACQCH5M6SBXKQXRWQIL7QSR4FH2UI6EYI4RCJSB2T2ZYF2JDHZ2Q"
      },
      {
        "pk": "W3KONPXCGFNUGXGDCOCQYVD64KZOLUMHZ7BNM2ZBK5FSSARRDEXINLYHPI"
      }
    ],
    "thr": 2,
    "v": 1
  },
  "txn": {
    "amt": 10000000,
    "fee": 1000,
    "fv": 4694301,
    "gen": "testnet-v1.0",
    "gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
    "lv": 4695301,
    "rcv": "QC7XT7QU7X6IHNRJZBR67RBMKCAPH67PCSX4LYH4QKVSQ7DQZ32PG5HSVQ",
    "snd": "GQ3QPLJL4VKVGQCHPXT5UZTNZIJAGVJPXUHCJLRWQMFRVL4REVW7LJ3FGY",
    "type": "pay"
  }
}
```
The difference between this transaction and the one above is the form of its signature component. For multisignature accounts, an [`"msig"`](../../reference/transactions.md#msig) struct is added which contains the 3 public addresses (`"pk"`), the threshold value (`"thr"`) and the multisig version `"v"`. This transaction is still unsigned but the addition of the correct `"msig"` struct is confirmation that the transaction is "aware" of the fact that the sender is multisig and will have no trouble accepting sub-signatures from single keys even if the signing agent does not contain information about its multisignature properties.

!!! tip
	Adding the `"msig"` template to make the transaction "aware" of its multisig sender is highly recommended, particularly in cases where the transaction is signed by multiple parties or offline. Without it, the signing agent would need to have its own knowledge of the multisignature account. For example, `goal` can sign a multisig transaction that does not contain an `"msig"` template _if_ the multisig address was created within its wallet. On signing, it will add the `"msig"` template. 


Sub-signatures can be added to the transaction one at a time, cumulatively, or merged together from multiple transactions. Here is the same transaction above, fully authorized:

```json hl_lines="6 13"
{
  "msig": {
    "subsig": [
      {
        "pk": "SYGHTA2DR5DYFWJE6D4T34P4AWGCG7JTNMY4VI6EDUVRMX7NG4KTA2WMDA",
        "s": "xoQkPyyqCPEhodngmOTP2930Y2GgdmhU/YRQaxQXOwh775gyVSlb1NWn70KFRZvZU96cMtq6TXW+r4sK/lXBCQ=="
      },
      {
        "pk": "VBDMPQACQCH5M6SBXKQXRWQIL7QSR4FH2UI6EYI4RCJSB2T2ZYF2JDHZ2Q"
      },
      {
        "pk": "W3KONPXCGFNUGXGDCOCQYVD64KZOLUMHZ7BNM2ZBK5FSSARRDEXINLYHPI",
        "s": "p1ynP9+LZSOZCBcrFwt5JZB2F+zqw3qpLMY5vJBN83A+55cXDYp5uz/0b+vC0VKEKw+j+bL2TzKSL6aTESlDDw=="
      }
    ],
    "thr": 2,
    "v": 1
  },
  "txn": {
    "amt": 10000000,
    "fee": 1000,
    "fv": 4694301,
    "gen": "testnet-v1.0",
    "gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
    "lv": 4695301,
    "rcv": "QC7XT7QU7X6IHNRJZBR67RBMKCAPH67PCSX4LYH4QKVSQ7DQZ32PG5HSVQ",
    "snd": "GQ3QPLJL4VKVGQCHPXT5UZTNZIJAGVJPXUHCJLRWQMFRVL4REVW7LJ3FGY",
    "type": "pay"
  }
}
```

The two signatures are added underneath their respective addresses. Since 2 meets the threshold, this transaction is now fully authorized and can be sent to the network.

!!! info
	Adding more sub-signatures than the threshold requires is unnecessary but perfectly valid.

**How-To**

Extend the example from the [Multisignature Account](../accounts/create.md#multisignature) section by creating, signing, and sending a transaction from a multisig account on TestNet.

```javascript tab="JavaScript"
const algosdk = require('algosdk');

const waitForConfirmation = async function (algodclient, txId, timeout) {
    // ...covered in a previous example
};

// enter token, server, and port
const token = "<algod-token>";
const server = "<algod-address>";
const port = "<algod-port>";

const main = async () => {
    // mnemonic phrases here for each of your test accounts
    const account1_mnemonic = "< mnemonic 1 >";
    const account2_mnemonic = "< mnemonic 2 >";
    const account3_mnemonic = "< mnemonic 3 >"

    const account1 = algosdk.mnemonicToSecretKey(account1_mnemonic);
    const account2 = algosdk.mnemonicToSecretKey(account2_mnemonic);
    const account3 = algosdk.mnemonicToSecretKey(account3_mnemonic);

    // Setup the parameters for the multisig account
    const mparams = {
        version: 1,
        threshold: 2,
        addrs: [
            account1.addr,
            account2.addr,
            account3.addr,
        ],
    };

    const multsigaddr = algosdk.multisigAddress(mparams);
    console.log("Multisig Address: " + multsigaddr);

    const algodClient = new algosdk.Algodv2(token, server, port);

    // Get the relevant params from the algod
    let params = await algodClient.getTransactionParams().do();
    // comment out the next two lines to use suggested fee
    params.fee = 1000;
    params.flatFee = true;

    const receiver = account3.addr;
    const names = '{"firstName":"John", "lastName":"Doe"}';
    const enc = new TextEncoder();
    const note = enc.encode(names);

    const txn = algosdk.makePaymentTxnWithSuggestedParams(multsigaddr, receiver, 1000000, undefined, note, params);
    const txId = txn.txID().toString();
    
    // Sign with first signature
    const rawSignedTxn = algosdk.signMultisigTransaction(txn, mparams, account1.sk).blob;
    //sign with second account
    const twosigs = algosdk.appendSignMultisigTransaction(rawSignedTxn, mparams, account2.sk).blob;
    
    //submit the transaction
    await algodclient.sendRawTransaction(twosigs).do();

    // Wait for confirmation
    const confirmedTxn = await waitForConfirmation(algodclient, txId, 4);
    console.log("Transaction " + txId + " confirmed in round " + confirmedTxn["confirmed-round"]);
    
    const mytxInfo = JSON.stringify(confirmedTxn.txn.txn, undefined, 2);
    console.log("Transaction information: %o", mytxInfo);

    const noteString = new TextDecoder().decode(confirmedTxn.txn.txn.note);
    console.log("Note field: ", noteString);

    const obj = JSON.parse(noteString);
    console.log("Note first name: %s", obj.firstName);
}

main().then().catch(e => {
    const error = e.body && e.body.message ? e.body.message : e;
    console.log(error);
});
```

```python tab="Python"
import json
from algosdk.v2client import algod
from algosdk import account, encoding, mnemonic
from algosdk.future.transaction import Multisig, PaymentTxn, MultisigTransaction
import base64

# Change these values with mnemonics
mnemonic1 = "PASTE phrase for account 1"
mnemonic2 = "PASTE phrase for account 2"
mnemonic3 = "PASTE phrase for account 3"


# For ease of reference, add account public and private keys to
# an accounts dictionary.

private_key_1 = mnemonic.to_private_key(mnemonic1)
account_1 = mnemonic.to_public_key(mnemonic1)

private_key_2 = mnemonic.to_private_key(mnemonic2)
account_2 = mnemonic.to_public_key(mnemonic2)

private_key_3 = mnemonic.to_private_key(mnemonic3)
account_3 = mnemonic.to_public_key(mnemonic3)

# utility for waiting on a transaction confirmation
def wait_for_confirmation(client, transaction_id, timeout):
    """
    Wait until the transaction is confirmed or rejected, or until 'timeout'
    number of rounds have passed.
    Args:
        transaction_id (str): the transaction to wait for
        timeout (int): maximum number of rounds to wait    
    Returns:
        dict: pending transaction information, or throws an error if the transaction
            is not confirmed or rejected in the next timeout rounds
    """
    start_round = client.status()["last-round"] + 1
    current_round = start_round

    while current_round < start_round + timeout:
        try:
            pending_txn = client.pending_transaction_info(transaction_id)
        except Exception:
            return 
        if pending_txn.get("confirmed-round", 0) > 0:
            return pending_txn
        elif pending_txn["pool-error"]:  
            raise Exception(
                'pool error: {}'.format(pending_txn["pool-error"]))
        client.status_after_block(current_round)                   
        current_round += 1
    raise Exception(
        'pending tx not found in timeout rounds, timeout value = : {}'.format(timeout))


# create a multisig account
version = 1  # multisig version
threshold = 2  # how many signatures are necessary
msig = Multisig(version, threshold, [account_1, account_2])

print("Multisig Address: ", msig.address())
print("Please go to: https://bank.testnet.algorand.network/ to fund multisig account.", msig.address())
# input("Please go to: https://bank.testnet.algorand.network/ to fund multisig account." + '\n' + "Press Enter to continue...")

# Specify your node address and token. This must be updated.
algod_address = ""  # ADD ADDRESS
algod_token = ""  # ADD TOKEN


# Initialize an algod client
algod_client = algod.AlgodClient(algod_token, algod_address)

# get suggested parameters
params = algod_client.suggested_params()
# comment out the next two (2) lines to use suggested fees
params.flat_fee = True
params.fee = 1000


# create a transaction
sender = msig.address()
recipient = account_3
amount = 10000
note = "Hello Multisig".encode()
txn = PaymentTxn(sender, params, recipient, amount, None, note, None)


# create a SignedTransaction object
mtx = MultisigTransaction(txn, msig)

# sign the transaction
mtx.sign(private_key_1)
mtx.sign(private_key_2)

# send the transaction
txid = algod_client.send_raw_transaction(
    encoding.msgpack_encode(mtx))
    # wait for confirmation	
try:
    confirmed_txn = wait_for_confirmation(algod_client, txid, 4)  
    print("Transaction information: {}".format(
        json.dumps(confirmed_txn, indent=4)))
    print("Decoded note: {}".format(base64.b64decode(
        confirmed_txn["txn"]["txn"]["note"]).decode()))
except Exception as err:
    print(err)
```

```java tab="Java"
package com.algorand.javatest.multisig;

import java.io.Console;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

import com.algorand.algosdk.account.Account;
import com.algorand.algosdk.v2.client.common.AlgodClient;
import com.algorand.algosdk.v2.client.common.Response;
import com.algorand.algosdk.v2.client.model.PendingTransactionResponse;
import com.algorand.algosdk.v2.client.model.NodeStatusResponse;
import com.algorand.algosdk.v2.client.model.TransactionParametersResponse;
import com.algorand.algosdk.algod.client.ApiException;
import com.algorand.algosdk.crypto.Address;
import com.algorand.algosdk.crypto.Ed25519PublicKey;
import com.algorand.algosdk.crypto.MultisigAddress;
import com.algorand.algosdk.transaction.SignedTransaction;
import com.algorand.algosdk.transaction.Transaction;
import com.algorand.algosdk.util.Encoder;
import org.json.JSONObject;
import com.algorand.algosdk.v2.client.model.PostTransactionsResponse;

/**
 * Test Multisignature
 *
 */
public class Multisig {

    public AlgodClient client = null;
    public static void waitForEnter(String message) {
        Console c = System.console();
        if (c != null) {
            // printf-like arguments
            if (message != null)
                c.format(message);
            c.format("\nPress ENTER to proceed.\n");
            c.readLine();
        }
    }
    
    // utility function to connect to a node
    private AlgodClient connectToNetwork() {

        // Initialize an algod client
        final String ALGOD_API_ADDR = "algod-address<PLACEHOLDER>";
        final String ALGOD_API_TOKEN = "algod-token<PLACEHOLDER>";
        final Integer ALGOD_PORT = port<PLACEHOLDER>;


        AlgodClient client = new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);
        return client;

    }

    /**
     * utility function to wait on a transaction to be confirmed
     * the timeout parameter indicates how many rounds do you wish to check pending transactions for
     */
    public PendingTransactionResponse waitForConfirmation(AlgodClient myclient, String txID, Integer timeout)
    throws Exception {
        if (myclient == null || txID == null || timeout < 0) {
            throw new IllegalArgumentException("Bad arguments for waitForConfirmation.");
        }
        Response < NodeStatusResponse > resp = myclient.GetStatus().execute();
        if (!resp.isSuccessful()) {
            throw new Exception(resp.message());
        }
        NodeStatusResponse nodeStatusResponse = resp.body();
        Long startRound = nodeStatusResponse.lastRound + 1;
        Long currentRound = startRound;
        while (currentRound < (startRound + timeout)) {
            // Check the pending transactions                 
            Response < PendingTransactionResponse > resp2 = myclient.PendingTransactionInformation(txID).execute();
            if (resp2.isSuccessful()) {
                PendingTransactionResponse pendingInfo = resp2.body();
                if (pendingInfo != null) {
                    if (pendingInfo.confirmedRound != null && pendingInfo.confirmedRound > 0) {
                        // Got the completed Transaction
                        return pendingInfo;
                    }
                    if (pendingInfo.poolError != null && pendingInfo.poolError.length() > 0) {
                        // If there was a pool error, then the transaction has been rejected!
                        throw new Exception("The transaction has been rejected with a pool error: " + pendingInfo.poolError);
                    }
                }
            }
            resp = myclient.WaitForBlock(currentRound).execute();
            if (!resp.isSuccessful()) {
                throw new Exception(resp.message());
            }
            currentRound++;
        }
        throw new Exception("Transaction not confirmed after " + timeout + " rounds!");
    }
    

    public void multisigExample() throws Exception {
        if (client == null)
            this.client = connectToNetwork();

        final String account1_mnemonic = <var>your-25-word-mnemonic</var>
        final String account2_mnemonic = <var>your-25-word-mnemonic</var>
        final String account3_mnemonic = <var>your-25-word-mnemonic</var>

        Account act1 = new Account(account1_mnemonic);
        Account act2 = new Account(account2_mnemonic);
        Account act3 = new Account(account3_mnemonic);
        System.out.println("Account1: " + act1.getAddress());
        System.out.println("Account2: " + act2.getAddress());
        System.out.println("Account3: " + act3.getAddress());

        final String DEST_ADDR = act3.getAddress().toString();

        // List for Pks for multisig account
        List<Ed25519PublicKey> publicKeys = new ArrayList<>();
        publicKeys.add(act1.getEd25519PublicKey());
        publicKeys.add(act2.getEd25519PublicKey());
        publicKeys.add(act3.getEd25519PublicKey());

        // Instantiate the Multisig Account
        MultisigAddress msa = new MultisigAddress(1, 2, publicKeys);

        System.out.println("Multisignature Address: " + msa.toString());
        // waitForEnter("Use TestNet Dispenser to add funds, wait for the transaction to be finalized and press enter");

        // setup transaction   
        try {
            Response < TransactionParametersResponse > resp = client.TransactionParams().execute();
            if (!resp.isSuccessful()) {
                throw new Exception(resp.message());
            }
            TransactionParametersResponse params = resp.body();
            if (params == null) {
                throw new Exception("Params retrieval error");
            }                      
            BigInteger amount = BigInteger.valueOf(1000000); // microAlgos
            // add some notes to the transaction
            byte[] notes = "These are some notes encoded in some way!".getBytes();
            // Setup Transaction
            Address sender = new Address(msa.toString());

            Transaction tx = Transaction.PaymentTransactionBuilder()
                    .sender(sender)
                    .amount(amount)
                    .receiver(DEST_ADDR)
                    .note(notes)
                    .suggestedParams(params).build();
            // Sign the Transaction for two accounts
            SignedTransaction signedTx = act1.signMultisigTransaction(msa, tx);
            SignedTransaction completeTx = act2.appendMultisigTransaction(msa, signedTx);
            // Msgpack encode the signed transaction
            byte[] encodedTxBytes = Encoder.encodeToMsgPack(completeTx);

            // Submit the transaction to the network
            Response < PostTransactionsResponse > rawtxresponse = client.RawTransaction().rawtxn(encodedTxBytes).execute();
            if (!rawtxresponse.isSuccessful()) {
                throw new Exception(rawtxresponse.message());
            }
            String id = rawtxresponse.body().txId;
            System.out.println("Successfully sent tx with ID: " + id);
            // Wait for transaction confirmation
            PendingTransactionResponse pTrx = waitForConfirmation(client, id, 4);

            System.out.println("Transaction " + id + " confirmed in round " + pTrx.confirmedRound);
            // Read the transaction
            JSONObject jsonObj = new JSONObject(pTrx.toString());
            System.out.println("Transaction information (with notes): " + jsonObj.toString(2));
            System.out.println("Decoded note: " + new String(pTrx.txn.tx.note));


        } catch (ApiException e) {
            // This is generally expected, but should give us an informative error message.
            System.err.println("Exception when calling algod#rawTransaction: " + e.getResponseBody());
        }
    }
  
    public static void main(String args[]) throws Exception {
        Multisig t = new Multisig();
        t.multisigExample();
    }
}
```

```go tab="Go"
package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"errors"
	"strings"
	json "encoding/json"
	// "io/ioutil"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// UPDATE THESE VALUES
const algodAddress = "Your ADDRESS"
const algodToken = "Your TOKEN"

// Accounts to be used through examples
func loadAccounts() (map[int][]byte, map[int]string) {
	// Shown for demonstration purposes. NEVER reveal secret mnemonics in practice.


	// Paste in mnemonic phrases for all three accounts
	mnemonic1 := "PASTE phrase for account 1"
	mnemonic2 := "PASTE phrase for account 2"
	mnemonic3 := "PASTE phrase for account 3"

	mnemonics := []string{mnemonic1, mnemonic2, mnemonic3}
	pks := map[int]string{1: "", 2: "", 3: ""}
	var sks = make(map[int][]byte)

	for i, m := range mnemonics {
		var err error
		sk, err := mnemonic.ToPrivateKey(m)
		sks[i+1] = sk
		if err != nil {
			fmt.Printf("Issue with account %d private key conversion.", i+1)
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

func main() {

    // Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}
	// Get network-related transaction parameters and assign
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}
	// comment out the next two (2) lines to use suggested fees
	txParams.FlatFee = true
	txParams.Fee = 1000
	// Get pre-defined set of keys for example
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

	fromAddr, _ := ma.Address()
	// Print multisig account
    fmt.Printf("Here is your multisig address : %s \n", fromAddr.String())
    fmt.Println("Please go to: https://bank.testnet.algorand.network/ to fund your multisig account.")


    toAddr := addr3.String()
    var amount uint64 = 10000
    note := []byte("Hello World")
    genID := txParams.GenesisID
    genHash := txParams.GenesisHash
    firstValidRound := uint64(txParams.FirstRoundValid)
    lastValidRound := uint64(txParams.LastRoundValid)
	var minFee uint64 = 1000
	txn, err := transaction.MakePaymentTxn(
		fromAddr.String(),
		toAddr,
		minFee,     // fee per byte
		amount,  // amount
		firstValidRound, // first valid round
		lastValidRound, // last valid round
		note,    // note
		"",     // closeRemainderTo
		genID,     // genesisHash
		genHash,     // genesisHash
	)

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

	// Broadcast the transaction to the network
	txid, err = algodClient.SendRawTransaction(twoOfThreeTxBytes).Do(context.Background())

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
```

```zsh tab="goal"
# Sign cumulatively
$ goal clerk multisig sign -t multisig.txn -a $ADDRESS1
$ goal clerk multisig sign -t multisig.txn -a $ADDRESS2

# Or sign two separate files and merge
$ goal clerk multisig sign -t multisig1.txn -a $ADDRESS1
$ goal clerk multisig sign -t multisig2.txn -a $ADDRESS2
$ goal clerk multisig merge multisig1.txn multisig2.txn --out=merged.stxn
```

```zsh tab="algokey"
# algokey takes account-level mnemonics at the time of signing
# requires the transaction to include the msig struct before signing
$ algokey multisig --txfile=multisig1.txn --outfile=multisig1.stxn -m <25-word-mnemonic>
$ algokey multisig --txfile=multisig2.txn --outfile=multisig2.stxn -m <25-word-mnemonic>

# Use goal to merge the the *.stxn files.
```

# Logic Signatures

Logic Signatures (or LogicSigs) authorize transactions associated with an Algorand Smart Contract. Logic signatures are added to transactions to authorize spends from a [Contract Account](../asc1/stateless/modes.md#contract-account) or from a [Delegated Account](../asc1/stateless/modes.md#delegated-account).

A full explanation of Logic Signatures can be found in the [Algorand Smart Contract Usage Modes Guide](../asc1/stateless/modes.md#logic-signatures).

**Related How-To**

- [Use LogicSigs with the SDKs](../asc1/stateless/sdks.md)
- [Attach a LogicSig with `goal`](../asc1/stateless/walkthrough.md)

!!! info
    Full running code examples for each SDK are available within the GitHub repo for V1 and V2 at [/examples/multisig](https://github.com/algorand/docs/tree/master/examples/multisig) and for [download](https://github.com/algorand/docs/blob/master/examples/multisig/multisig.zip?raw=true) (.zip).
