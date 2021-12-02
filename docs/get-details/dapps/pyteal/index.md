title: PyTeal

[PyTeal](https://github.com/algorand/pyteal) is a python language binding for Algorand Smart Contracts (ASC) that abstracts away the complexities of writing smart contracts. PyTeal allows smart contracts and smart signatures to be written in Python and then compiled to TEAL. Note that the TEAL code is not automatically compiled to byte code, but that can be done with any of the SDKs, including Python. The TEAL code can also be compiled using the `goal` command-line tool or submitted to the blockchain to allow the node to compile it.
Complete installation instructions and developer guides are available in the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/). 

To quickly get PyTeal installed and running, see the [Getting started tutorial](../../../get-started/dapps/pyteal.md) for PyTeal.

!!! note
	This document refers to stateful smart contracts as smart contracts and stateless smart contracts as smart signatures.

# PyTeal overview
This section assumes the reader is familiar with Smart Contracts and Smart Signatures.
When building a dApp that makes use of smart contracts or smart signatures (smartsigs), PyTeal makes implementing these much simpler than writing the TEAL manually. Generally, developers install PyTeal, write the Python contract or smartsignature code using their preferred editor, and then use PyTeal’s `compileProgram` method to produce the TEAL code. The TEAL code can then be deployed to the blockchain.

In most applications, these contracts and smartsigs will only be a portion of the dApp’s architecture. Generally, developers will build functionality in the dApp that resides on the blockchain and some sort of front end to interact with the backend smart contracts, transactions, accounts, and assets. PyTeal can be used to produce both smart contracts and smart signatures. These can be pre-built and used by the dApp or the building of the contracts and smartsigs can be integrated into the dApp front end logic and then deployed as part of the normal operation of the application. For example, a front end may provide something like an exchange that allows limit orders to be created based on a template and then deployed once a user opens an order. In this case, the complete limit order may be implemented as part of a smart contract that is deployed when the order is opened by the dApp. For more information on deployment models see [What is a dApp](../../../get-started/dapps/index.md) in the developer documentation.

PyTeal supports both smart signatures and smart contracts. The following sections explain both of these scenarios.

# Building PyTeal smart contracts
On Algorand, smart contracts are small programs that contain logic that is evaluated when the contract is deployed or called. These contracts have a set of functions (opcodes) that can be called that make use of on-chain data (such as balances), additional arguments, additional transactions, and stored values to evaluate to either true or false. Additionally, as part of the logic on-chain variables can be stored on a per contract or per account basis. If the contract runs successfully these variables will be updated per the logic. If the contract fails the variable changes will be rolled back. For example, you may have a voting contract that stores whether a specific account has voted or not and a vote total for all the candidates in an election. The user may call the smart contract to vote and the logic may check to see if the specific account has voted already and if not allow the vote to be processed and increment the vote totals. The voting account would also be marked as having voted. If the account had already voted the call would fail. 

Smart contracts also have an Algorand address, which allows the contract to hold Algorand assets(ASAs) and Algos. Any account in the network can send ASAs or Algos to this address unimpeded. The only way either of these forms of value is allowed to leave the contract is when the logic within the contract performs a payment or asset transaction and the logic returns true. For more information on smart contracts, see the [smart contract documentation](../smart-contracts/apps/index.md).

Smart contracts are often referred to as applications on Algorand because they usually contain the bulk of blockchain logic that a distributed application will require. Once a contract is written it is deployed using an application creation transaction. This process will be described in the next section. 

When building smart contracts in PyTeal it is important to realize that a smart contract actually consists of two programs. These are called the approval and the clear programs. In PyTeal both of these programs are generally created in the same Python file. So the beginning of a  PyTeal program will contain logic similar to the following:

!!!info
       The following sample builds a simple counter smart contract that either adds or deducts one from a global counter based on how the contract is called. 


```python
#samplecontract.py
from pyteal import *
 
"""Basic Counter Application"""
 
def approval_program():
   program = Return(Int(1))
   # Mode.Application specifies that this is a smart contract
   return compileTeal(program, Mode.Application, version=5)

def clear_state_program():
   program = Return(Int(1))
   # Mode.Application specifies that this is a smart contract
   return compileTeal(program, Mode.Application, version=5)
``` 

In the above example, a function is defined to return each of the two programs. The `compileTeal` function compiles the program as defined by the program variable. In this case, we are just returning a one for both programs. The `compileTeal` method also sets the Mode.Application. This lets PyTeal know this is for a smart contract and not a smart signature. The version parameter instructs PyTeal on which version of TEAL to produce when compiling. 

All communication with Algorand smart contracts is achieved through a special transaction type, called an application transaction. Application transactions have six subtypes. The specific subtype will determine which of the two programs to call. 

<center>![Stateful Smart Contract](../../../imgs/sccalltypes.png)</center>
<center>*Application Transaction Types*</center>

Most smart contract logic will be implemented with a NoOp transaction type. The other subtypes are primarily less used or once-only transaction types. For example, Algorand allows smart contracts to be updated. Code must be implemented in the contract to prevent this if it is an unwanted feature. Smart contracts can also be deleted, although this can be disabled as well. The Optin transaction type is submitted by an account that wants to opt into the smart contract. Note that this is only required by contracts that store per account values. It is also possible to have smart contracts that store values for certain accounts but not others. In this case, if logic is encountered in the contract that attempts to store values for a particular account. It will fail unless the account has opted into the contract. The CloseOut application transaction is used to gracefully exit a smart contract. It is primarily an opt-out type of operation that allows the smart contract to do cleanup when an account wishes to leave the contract. This transaction can fail based on the logic, which would lock the user into the contract forever. To circumvent this issue, Algorand also has a Clear application transaction type. This type of transaction allows an ungraceful exit from the contract. This transaction may still fail but the blockchain will still clear any data associated with the contract from the account. Only the Clear transaction will call the clear program. All others will call the approval program.  

Within PyTeal, a developer can switch on the type of transaction. This is the preferred way of building a PyTeal contract. Sections of code should be created that handle any of the transaction types they may encounter. The above example's approval program can be changed to the following to handle the different application transaction types.

```python
def approval_program():
   # Mode.Application specifies that this is a smart contract

   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   return compileTeal(program, Mode.Application, version=5)

```

The `program` variable is changed to use the [PyTeal `Cond` expression](https://pyteal.readthedocs.io/en/latest/control_structures.html?highlight=Cond#chaining-tests-cond). This expression allows several conditions to be chained and the first to return true will then evaluate its second parameter, known as the condition body. If none of the conditions are true the smart contract will return an `err` and fail. The conditions are the first parameter to each condition and the body is the second parameter. In the above example, most of the conditions check the transaction type using the `on_completion` transaction field. This field will map to one of the subtypes of application transactions described above. The body for each condition here will point to another location in your PyTeal contract. 

When a smart contract is deployed, it actually is a special case of a `NoOp` transaction type. The Algorand SDKs provide functions to create this special type of transaction. For example, the Python SDK uses the function named `ApplicationCreateTxn`. When the contract is first created, the contract’s ID will be equal to 0. Once created, all smart contracts will have a unique ID and a unique Algorand address. This is why the first condition checked in the example above is checking the ID of the current smart contract. If this value is 0, the logic will know this is the first execution of the contract. This is a perfect place to add code to the contract to initialize any global variables the contract will use. More details on deploying a contract will be covered in the next section of this guide.

The `Txn` object is used to access transaction fields for the current transaction. Many fields can be examined. See the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/accessing_transaction_field.html) for a complete list. 

The example also uses the arithmetic expression `==` to check equality. See the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/arithmetic_expression.html) for additional arithmetic expressions.

```python

def approval_program():
   # Mode.Application specifies that this is a smart contract

   handle_creation = Seq([
       App.globalPut(Bytes("Count"), Int(0)),
       Return(Int(1))
   ])

   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   return compileTeal(program, Mode.Application, version=5)
```

The [`Seq` expression](https://pyteal.readthedocs.io/en/latest/control_structures.html?highlight=seq#chaining-expressions-seq) is used here when the contract is created. This expression is used to provide a sequence of expressions. In this example, the first expression stores a global variable named Count, and its value is set to 0. This expression is followed by the `Return` expression which exits the program with the return value. In this case, it returns a value of 1, indicating success. So when this specific smart contract is first deployed it will store a global variable named Count with a value of 0 and immediately return success. The `Seq` expression is a flow control expression much like the `Cond` or `Return` expression. PyTeal contains many different flow expressions which are explained in the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/control_structures.html?highlight=seq#). 

As stated earlier, global variables can be stored for a contract. In addition, for any account that opts into the contract, local variables can be stored for each account. To read more about how to do this with PyTeal see the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/control_structures.html?highlight=seq#chaining-expressions-seq).

Opting into this contract is not required as no local variables are stored. The logic can reject the opt-in application transaction. Also, if no opt-in is allowed, close-out application transactions are not useful either, so those can be rejected as well. Finally, the contract can deny application transactions that attempt to delete or update the contract.

```python

def approval_program():
   # Mode.Application specifies that this is a smart contract

   handle_creation = Seq([
       App.globalPut(Bytes("Count"), Int(0)),
       Return(Int(1))
   ])

   handle_optin = Return(Int(0))
 
   handle_closeout = Return(Int(0))
 
   handle_updateapp = Return(Int(0))
 
   handle_deleteapp = Return(Int(0))


   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   return compileTeal(program, Mode.Application, version=5)

```

All four of these transaction types simply return a 0, which will cause the transactions to fail. This contract now handles all application transactions but the standard NoOp type. Do remember that deploying the contract for the first time is actually a NoOp transaction type, but this case is accounted for in the `handle_creation` function. The NoOp transaction type is the primary location where application logic will be implemented in most smart contracts. This example requires an add and a deduct function, to increment and decrement the counter respectively, to be handled for NoOp application transactions. Which of these two methods is executed will depend on the first parameter to the stateful smart contract. In addition, we want to verify that application transactions are not grouped with any other transactions.

For more information on passing parameters to smart contracts, see the [smart contract documentation](../smart-contracts/apps/index.md). 

```python
from pyteal import *
 
"""Basic Counter Application"""
 
def approval_program():
   # Mode.Application specifies that this is a smart contract

   handle_creation = Seq([
       App.globalPut(Bytes("Count"), Int(0)),
       Return(Int(1))
   ])

   handle_optin = Return(Int(0))
 
   handle_closeout = Return(Int(0))
 
   handle_updateapp = Return(Int(0))
 
   handle_deleteapp = Return(Int(0))

   handle_noop = Cond(
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Add")
       ), add],
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Deduct")
       ), deduct],
   )

   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   return compileTeal(program, Mode.Application, version=5)

```

In the example, another `Cond` expression is used to handle the NoOp transaction. The condition expression only has two options. The expression for each option is actually created using a logical `And` arithmetic expression. The `And` expression first checks the global group size variable to verify that is set to 1, indicating the application call is not grouped with any other transactions. This is logically Anded with an expression that checks that the first argument passed with the application transaction is either the string Add or Deduct. The body of each condition will be added next. PyTeal has access to many global variables like group size. See the [PyTeal documentation](https://pyteal.readthedocs.io/en/latest/accessing_transaction_field.html?highlight=global#global-parameters) for more information.

The final step for the approval program is to implement the add and deduct functions for the smart contract.

```python
def approval_program():
   # Mode.Application specifies that this is a smart contract

   handle_creation = Seq([
       App.globalPut(Bytes("Count"), Int(0)),
       Return(Int(1))
   ])

   handle_optin = Return(Int(0))
 
   handle_closeout = Return(Int(0))
 
   handle_updateapp = Return(Int(0))
 
   handle_deleteapp = Return(Int(0))

   scratchCount = ScratchVar(TealType.uint64)
 
   add = Seq([
       scratchCount.store(App.globalGet(Bytes("Count"))),
       App.globalPut(Bytes("Count"), scratchCount.load() + Int(1)),
       Return(Int(1))
   ])
 
    deduct = Seq([
       scratchCount.store(App.globalGet(Bytes("Count"))),
        If(scratchCount.load() > Int(0),
            App.globalPut(Bytes("Count"), scratchCount.load() - Int(1)),
        ),
        Return(Int(1))
   ])


   handle_noop = Cond(
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Add")
       ), add],
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Deduct")
       ), deduct],
   )

   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   return compileTeal(program, Mode.Application, version=5)

```

First, the contract is modified to create a temporary variable in scratch space. Smart contracts can store up to 256 temporary variables in scratch space. PyTeal provides a class that can be used to interface with these temporary variables. The scratch variable in this example happens to be an integer, byte arrays can also be stored. The implementation of the add and deduct functions rely on the `Seq` expression, guaranteeing an order of operations. First, the current value of the global variable Count is read for the contract and placed in scratch space. Next, the contract either increments this number or decrements and then stores the resultant back into the contract’s global variable. On the deduct, an additional `If` expression is used to verify the current global variable is above 0. Finally, both methods then exit the smart contract call, returning a 1, indicating that the transaction was successful.

This sample application requires no local variables and user’s do not have to opt into the contract to call either of its methods. This means that if a clear transaction is submitted to the contract, it will execute the clear program and simply return 1, indicating success. The full example is presented below. Additionally, print commands are added to the contract to illustrate the output.

```python
#samplecontract.py
from pyteal import *
 
"""Basic Counter Application"""
 
def approval_program():
   handle_creation = Seq([
       App.globalPut(Bytes("Count"), Int(0)),
       Return(Int(1))
   ])
 
   handle_optin = Return(Int(0))
 
   handle_closeout = Return(Int(0))
 
   handle_updateapp = Return(Int(0))
 
   handle_deleteapp = Return(Int(0))
 
   scratchCount = ScratchVar(TealType.uint64)
 
   add = Seq([
       scratchCount.store(App.globalGet(Bytes("Count"))),
       App.globalPut(Bytes("Count"), scratchCount.load() + Int(1)),
       Return(Int(1))
   ])
 
   deduct = Seq([
       scratchCount.store(App.globalGet(Bytes("Count"))),
        If(scratchCount.load() > Int(0),
            App.globalPut(Bytes("Count"), scratchCount.load() - Int(1)),
        ),
        Return(Int(1))
   ])
 
   handle_noop = Cond(
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Add")
       ), add],
       [And(
           Global.group_size() == Int(1),
           Txn.application_args[0] == Bytes("Deduct")
       ), deduct],
   )

 
   program = Cond(
       [Txn.application_id() == Int(0), handle_creation],
       [Txn.on_completion() == OnComplete.OptIn, handle_optin],
       [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
       [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
       [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
       [Txn.on_completion() == OnComplete.NoOp, handle_noop]
   )
   # Mode.Application specifies that this is a smart contract
   return compileTeal(program, Mode.Application, version=5)
 
 
def clear_state_program():
   program = Return(Int(1))
   # Mode.Application specifies that this is a smart contract
   return compileTeal(program, Mode.Application, version=5)

# print out the results
print(approval_program())
print(clear_state_program())
```

This program can be executed to illustrate compiling the PyTeal and printing out the resultant TEAL code.
 
```bash
$ python3 samplecontract.py
```

This example of a smart contract is very simple. Using PyTeal, more sophisticated contracts can be created. To learn more about what can be done in smart contracts, see the [smart contract documentation](../smart-contracts/apps/index.md). The documentation also contains many PyTeal code snippets that can be used within smart contracts.

# Deploying and calling the smart contract
This section explains how to deploy and call the smart contract developed in the previous section.

## Deploying the contract
In the previous section, the development of a simple smart contract was explained. This smart contract can be deployed in many different ways, but generally, this will be done using one of the Algorand SDKs ([Python](../../../sdks/python/index.md), [JavaScript](../../../sdks/javascript/index.md), [Go](../../../sdks/go/index.md), and [Java](../../../sdks/java/index.md)). This section will add additional code to the previous section’s example using the Python SDK to illustrate deploying the example contract. 

!!! note
	This example expects the developer to use the sandbox install. Additionally, one account should be set up and funded. See the [Python SDK](../../../sdks/python/index.md) getting started guide for more details. 

Before getting into the details of deploying the contract, a couple of global variables must be added to the PyTeal Python example. 

```python
# user declared account mnemonics
creator_mnemonic = "REPLACE WITH YOUR OWN MNEMONIC"
# user declared algod connection parameters. 
# Node must have EnableDeveloperAPI set to true in its config
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
```

The first is a creator mnemonic. This mnemonic is used to recover the private key for the funded account that will own and create the smart contract. Placing a mnemonic like this in code should never be done in production. Typically applications will link to some protected wallet to sign transactions. Some examples of wallets are the Algorand mobile wallet, AlgoSigner, MyAlgo Wallet, and Aikon ORE. When using the Algorand mobile wallet, transactions can be signed using the [Wallet Connect API](../../walletconnect/index.md). The mnemonic is used here for learning purposes only.

The algod_address and algod_token values are the default values to connect to a sandbox installed node. Also note that in this example, the sandbox node is connected to the Algorand TestNet network (eg  `./sandbox up testnet`). 

In this example, the TEAL for the smart contract will be compiled programmatically by the node. The `EnableDeveloperAPI` configuration parameter must be set to `true` in the node’s configuration to allow this operation. For more information on changing node configuration parameters, see the [developer documentation](../../../run-a-node/reference/config.md). When using the sandbox install, this value is automatically set to true.

Next, a few helper functions need to be added to the sample.

```python
# helper function to compile program source
def compile_program(client, source_code):
    compile_response = client.compile(source_code)
    return base64.b64decode(compile_response['result'])
    
# helper function that converts a mnemonic passphrase into a private signing key
def get_private_key_from_mnemonic(mn) :
    private_key = mnemonic.to_private_key(mn)
    return private_key

# helper function that waits for a given txid to be confirmed by the network
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

```

The `compile_program` function is a utility function that allows passing the generated TEAL code to a node that will compile and return the byte code. This returned byte code will be used with the application creation transaction (deploying the contract) later.

The ‘get_private_key_from_mnemonic` function is a utility function that takes a mnemonic (account backup phrase) and returns the private key of the specific account. This will be used in this sample to recover the private key of the funded account of the smart contract creator.

The `wait_for_confirmation` function is a utility function that when called will wait until a specific transaction is confirmed on the Algorand blockchain. This will be used to confirm that the application creation transaction is successful and the smart contract is actively deployed.

As the sample smart contract manipulates global variables, a couple of helper functions are needed to display the contents of these values.

```python
# helper function that formats global state for printing
def format_state(state):
    formatted = {}
    for item in state:
        key = item['key']
        value = item['value']
        formatted_key = base64.b64decode(key).decode('utf-8')
        if value['type'] == 1:
            # byte string
            if formatted_key == 'voted':
                formatted_value = base64.b64decode(value['bytes']).decode('utf-8')
            else:
                formatted_value = value['bytes']
            formatted[formatted_key] = formatted_value
        else:
            # integer
            formatted[formatted_key] = value['uint']
    return formatted

# helper function to read app global state
def read_global_state(client, addr, app_id):
    results = client.account_info(addr)
    apps_created = results['created-apps']
    for app in apps_created:
        if app['id'] == app_id:
            return format_state(app['params']['global-state'])
    return {}

```

Global variables for smart contracts are actually stored in the creator account’s ledger entry on the blockchain. The location is referred to as global state and the SDKs provide a function to retrieve the account’s record. In this example, the function `read_global_state` uses the Python SDK function `acount_info` to connect to the Algorand node and retrieve the account information. The function then locates the created application within this record. The `format_state` function takes the application data and formats the values for display. For more information on global and local state see the [smart contract documentation](../smart-contracts/apps/index.md).

As covered earlier in this guide, to deploy the contract an application creation transaction must be created and submitted to the blockchain. The SDKs provide a method for creating this transaction. The following code illustrates creating and submitting this transaction.

```python
# create new application
def create_app(client, private_key, approval_program, clear_program, global_schema, local_schema):
    # define sender as creator
    sender = account.address_from_private_key(private_key)

    # declare on_complete as NoOp
    on_complete = transaction.OnComplete.NoOpOC.real

    # get node suggested parameters
    params = client.suggested_params()
   
    # create unsigned transaction
    txn = transaction.ApplicationCreateTxn(sender, params, on_complete, \
                                            approval_program, clear_program, \
                                            global_schema, local_schema)

    # sign transaction
    signed_txn = txn.sign(private_key)
    tx_id = signed_txn.transaction.get_txid()

    # send transaction
    client.send_transactions([signed_txn])

    # await confirmation
    wait_for_confirmation(client, tx_id, 5)

    # display results
    transaction_response = client.pending_transaction_info(tx_id)
    app_id = transaction_response['application-index']
    print("Created new app-id:", app_id)

    return app_id
```

This function is a simple example of creating an application creation transaction, which when submitted will deploy a smart contract. This example is very generic and can be used to deploy any smart contract. First, the creator’s address is resolved from the private key passed to the function, the transaction type is set to a NoOp application transaction, and the blockchain suggested parameters are retrieved from the connected node. These suggested parameters provide the default values that are required to submit a transaction, such as the expected fee for the transaction.

The Python SDK’s `ApplicationCreateTxn` function is called to create the transaction.  This function takes the creator’s address, the approval and clear programs byte code, and a declaration of how much global and local state the smart contract will reserve. When creating a smart contract, the creation transaction has to specify how much state will be reserved. A contract can store up to 64 key-value pairs in global state and up to 16 key-value pairs per user who opts into the contract. Once these values are set, they can never be changed. The key is limited to 64 bytes. The key plus the value is limited to 128 bytes total. Using smaller keys to have more storage available for the value is possible. The keys are stored as byte slices (byte-array value) and the values are stored as either byte slices (byte-array value) or uint64s.  More information on state values can be found in the [smart contract documentation](../smart-contracts/apps/index.md#modifying-state-in-smart-contract).

The passed-in private key is then used to sign the transaction and the ID of the transaction is retrieved. This ID is unique and can be used to look up the transaction later.

The transaction is then submitted to the connected node and the `wait_for_confirmation` function is called to wait for the blockchain to process the transaction. Once the blockchain processes the transaction, a unique ID, called application ID, is returned for the smart contract. This can be used later to issue calls against the smart contract.

Now that all required functions are implemented, the main function can be created to deploy the contract.

```python
def main() :
    # initialize an algodClient
    algod_client = algod.AlgodClient(algod_token, algod_address)

    # define private keys
    creator_private_key = get_private_key_from_mnemonic(creator_mnemonic)

    # declare application state storage (immutable)
    local_ints = 0
    local_bytes = 0
    global_ints = 1 
    global_bytes = 0
    global_schema = transaction.StateSchema(global_ints, global_bytes)
    local_schema = transaction.StateSchema(local_ints, local_bytes)

    # compile program to TEAL assembly
    with open("./approval.teal", "w") as f:
        approval_program_teal = approval_program()
        f.write(approval_program_teal)

    # compile program to TEAL assembly
    with open("./clear.teal", "w") as f:
        clear_state_program_teal = clear_state_program()
        f.write(clear_state_program_teal)

    # compile program to binary
    approval_program_compiled = compile_program(algod_client, approval_program_teal)

    # compile program to binary
    clear_state_program_compiled = compile_program(algod_client, clear_state_program_teal)

    print("--------------------------------------------")
    print("Deploying Counter application......")
    
    # create new application
    app_id = create_app(algod_client, creator_private_key, approval_program_compiled, clear_state_program_compiled, global_schema, local_schema)

    # read global state of application
    print("Global state:", read_global_state(algod_client, account.address_from_private_key(creator_private_key), app_id))
```

First, a connection to the sandbox node is established. This is followed by recovering the account of the creator. Next, the amount of state to be used is defined. In this example, only one global integer is specified. 

The SDK is then used to first convert the approval and clear programs to TEAL using the PyTeal library and both are written to local files. Each is then complied to byte code by the connected node. Finally, the smart contract is deployed using the `create_app` function created earlier and the current global state is then printed out for the contract. On deployment, this value will be set to 0.

## Calling the deployed smart contract
Now that the contract is deployed, the Add or Deduct functions can be called using a standard NoOp application transaction. The example created throughout this guide can be further modified to illustrate making a call to the smart contract. 

To begin with, a function can be added to support calling the smart contract.

```python
# call application
def call_app(client, private_key, index, app_args) : 
    # declare sender
    sender = account.address_from_private_key(private_key)

    # get node suggested parameters
    params = client.suggested_params()

    # create unsigned transaction
    txn = transaction.ApplicationNoOpTxn(sender, params, index, app_args)

    # sign transaction
    signed_txn = txn.sign(private_key)
    tx_id = signed_txn.transaction.get_txid()

    # send transaction
    client.send_transactions([signed_txn])

    # await confirmation
    wait_for_confirmation(client, tx_id, 5)

    print("Application called")
```

This function operates similarly to the `create_app` function we defined earlier. In this case, we use the Python SDK’s `ApplicationNoOpTxn` function to create a standard NoOp application transaction. The address of the account sending the call is specified, followed by the network suggested parameters, the application id of the smart contract, and any arguments to the call. The arguments will be used to specify either the Add or Deduct methods.  

The `main` function can then be modified to call the smart contract after deploying by adding the following to the bottom of the `main` function.

```python
    print("--------------------------------------------")
    print("Calling Counter application......")
    app_args = ["Add"]
    call_app(algod_client, creator_private_key, app_id, app_args)

    # read global state of application
    print("Global state:", read_global_state(algod_client, account.address_from_private_key(creator_private_key), app_id))
```

In this example, the Add string is added to the application arguments array and the smart contract is called. The updated global state is then printed out. The value should now be set to 1.

The complete example is shown below.

```python
import base64

from algosdk.future import transaction
from algosdk import account, mnemonic, logic
from algosdk.v2client import algod
from pyteal import *

# user declared account mnemonics
creator_mnemonic = "finger rigid hat room course salmon say detect avocado assault awake sea public curious exit valve donkey tired escape dash drink diagram section absent cruise"
# user declared algod connection parameters. Node must have EnableDeveloperAPI set to true in its config
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# helper function to compile program source
def compile_program(client, source_code):
    compile_response = client.compile(source_code)
    return base64.b64decode(compile_response['result'])
    
# helper function that converts a mnemonic passphrase into a private signing key
def get_private_key_from_mnemonic(mn) :
    private_key = mnemonic.to_private_key(mn)
    return private_key

# helper function that waits for a given txid to be confirmed by the network
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

# helper function that formats global state for printing
def format_state(state):
    formatted = {}
    for item in state:
        key = item['key']
        value = item['value']
        formatted_key = base64.b64decode(key).decode('utf-8')
        if value['type'] == 1:
            # byte string
            if formatted_key == 'voted':
                formatted_value = base64.b64decode(value['bytes']).decode('utf-8')
            else:
                formatted_value = value['bytes']
            formatted[formatted_key] = formatted_value
        else:
            # integer
            formatted[formatted_key] = value['uint']
    return formatted

# helper function to read app global state
def read_global_state(client, addr, app_id):
    results = client.account_info(addr)
    apps_created = results['created-apps']
    for app in apps_created:
        if app['id'] == app_id:
            return format_state(app['params']['global-state'])
    return {}


"""Basic Counter Application in PyTeal"""

def approval_program():
    on_creation = Seq([
        App.globalPut(Bytes("Count"), Int(0)),
        Return(Int(1))
    ])

    handle_optin = Return(Int(0))

    handle_closeout = Return(Int(0))

    handle_updateapp = Return(Int(0))

    handle_deleteapp = Return(Int(0))

    scratchCount = ScratchVar(TealType.uint64)

    add = Seq([ 
        scratchCount.store(App.globalGet(Bytes("Count"))),
        App.globalPut(Bytes("Count"), scratchCount.load() + Int(1)),
        Return(Int(1))
    ])

    deduct = Seq([
       scratchCount.store(App.globalGet(Bytes("Count"))),
        If(scratchCount.load() > Int(0),
            App.globalPut(Bytes("Count"), scratchCount.load() - Int(1)),
        ),
        Return(Int(1))
   ])

    handle_noop = Cond(
        [And(
            Global.group_size() == Int(1),
            Txn.application_args[0] == Bytes("Add")
        ), add],
        [And(
            Global.group_size() == Int(1),
            Txn.application_args[0] == Bytes("Deduct")
        ), deduct],
    )

    program = Cond(
        [Txn.application_id() == Int(0), on_creation],
        [Txn.on_completion() == OnComplete.OptIn, handle_optin],
        [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
        [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
        [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp],
        [Txn.on_completion() == OnComplete.NoOp, handle_noop]
    )
    # Mode.Application specifies that this is a smart contract
    return compileTeal(program, Mode.Application, version=5)

def clear_state_program():
    program = Return(Int(1))
    # Mode.Application specifies that this is a smart contract
    return compileTeal(program, Mode.Application, version=5)

    
# create new application
def create_app(client, private_key, approval_program, clear_program, global_schema, local_schema):
    # define sender as creator
    sender = account.address_from_private_key(private_key)

    # declare on_complete as NoOp
    on_complete = transaction.OnComplete.NoOpOC.real

    # get node suggested parameters
    params = client.suggested_params()

    # create unsigned transaction
    txn = transaction.ApplicationCreateTxn(sender, params, on_complete, \
                                            approval_program, clear_program, \
                                            global_schema, local_schema)

    # sign transaction
    signed_txn = txn.sign(private_key)
    tx_id = signed_txn.transaction.get_txid()

    # send transaction
    client.send_transactions([signed_txn])

    # await confirmation
    wait_for_confirmation(client, tx_id, 5)

    # display results
    transaction_response = client.pending_transaction_info(tx_id)
    app_id = transaction_response['application-index']
    print("Created new app-id:", app_id)

    return app_id


# call application
def call_app(client, private_key, index, app_args) : 
    # declare sender
    sender = account.address_from_private_key(private_key)

    # get node suggested parameters
    params = client.suggested_params()

    # create unsigned transaction
    txn = transaction.ApplicationNoOpTxn(sender, params, index, app_args)

    # sign transaction
    signed_txn = txn.sign(private_key)
    tx_id = signed_txn.transaction.get_txid()

    # send transaction
    client.send_transactions([signed_txn])

    # await confirmation
    wait_for_confirmation(client, tx_id, 5)

    print("Application called")

def main() :
    # initialize an algodClient
    algod_client = algod.AlgodClient(algod_token, algod_address)

    # define private keys
    creator_private_key = get_private_key_from_mnemonic(creator_mnemonic)

    # declare application state storage (immutable)
    local_ints = 0
    local_bytes = 0
    global_ints = 1 
    global_bytes = 0
    global_schema = transaction.StateSchema(global_ints, global_bytes)
    local_schema = transaction.StateSchema(local_ints, local_bytes)

    # compile program to TEAL assembly
    with open("./approval.teal", "w") as f:
        approval_program_teal = approval_program()
        f.write(approval_program_teal)


    # compile program to TEAL assembly
    with open("./clear.teal", "w") as f:
        clear_state_program_teal = clear_state_program()
        f.write(clear_state_program_teal)
        
    # compile program to binary
    approval_program_compiled = compile_program(algod_client, approval_program_teal)
            
    # compile program to binary
    clear_state_program_compiled = compile_program(algod_client, clear_state_program_teal)

    print("--------------------------------------------")
    print("Deploying Counter application......")
    
    # create new application
    app_id = create_app(algod_client, creator_private_key, approval_program_compiled, clear_state_program_compiled, global_schema, local_schema)

    # read global state of application
    print("Global state:", read_global_state(algod_client, account.address_from_private_key(creator_private_key), app_id))

    print("--------------------------------------------")
    print("Calling Counter application......")
    app_args = ["Add"]
    call_app(algod_client, creator_private_key, app_id, app_args)

    # read global state of application
    print("Global state:", read_global_state(algod_client, account.address_from_private_key(creator_private_key), app_id))

main()

```

For more information on using the SDKs to deploy and interact with smart contracts see the [developer documentation](../smart-contracts/frontend/smartsigs.md).

# Building PyTeal smart signatures
Smart signatures are small programs that are submitted as part of a transaction and evaluated at submission time. These types of signatures can be used as an escrow-type of account or can be used to delegate a portion of the authority for a specific account. 

When used as an escrow, they can hold Algos or Algorand assets (ASAs). When used this way any transaction can send Algos or ASAs to the escrow but the logic in the signature determines when value leaves the escrow. In this respect, they act very similarly to smart contracts, but the logic must be supplied with every transaction.  

When used as a delegate, the logic can be signed by a specific account. The logic is then evaluated when a transaction is submitted from the signing account that is signed by the logic and not the private key of the sender. This is often used to allow restricted access to an account. For example, a mortgage company may provide logic to an account to remove a certain number of Algos from the account once a month. The user then signs this logic and once a month the mortgage company can submit a transaction from the signing account, but the transaction is signed by the smart signature and not the private key of the account.

Any time a smart signature is used the complete logic must be submitted as part of the transaction where the logic is used. The logic is recorded as part of the transaction but this is after the fact.

PyTeal supports building smart signatures in Python. For example, assume an escrow account is needed. This escrow can be funded by anyone but only a specific account is the beneficiary of the escrow and that account can withdraw funds at any time.  

```python
#sample_smart_sig.py
from pyteal import *

"""Basic Donation Escrow"""
def donation_escrow(benefactor):
    Fee = Int(1000)

    #Only the benefactor account can withdraw from this escrow
    program = And(
        Txn.type_enum() == TxnType.Payment,
        Txn.fee() <= Fee,
        Txn.receiver() == Addr(benefactor),
        Global.group_size() == Int(1),
        Txn.rekey_to() == Global.zero_address()        
    )
    # Mode.Signature specifies that this is a smart signature
    return compileTeal(program, Mode.Signature, version=5)
```

This is a very simplistic smart signature. The code for the complete signature is defined in the `donatation_escrow` function. This function takes an Algorand address as a parameter. This address represents the beneficiary of the escrow. The entire program is a set of conditions anded together using the [`And` logical expression](https://pyteal.readthedocs.io/en/latest/api.html#pyteal.And). This expression takes two or more arguments that are logically anded and produces a 0 (logically false) or 1 (logically true). In this sample, a set of transaction fields are compared to expected values. The transaction type is first verified to be a payment transaction, the transaction fee is compared to make sure it is less than 1000 microAlgos, the transaction receiver is compared to the benefactor’s address, the group size is verified to guarantee that this transaction is not submitted with other transactions in a group, and the rekey field of the transaction is verified to be the zero address. The zero address is used to verify that the rekey field is not set. This prevents the escrow from being rekeyed to another account. This sample uses transaction fields and global properties. See the PyTeal documentation for additional [transaction fields](https://pyteal.readthedocs.io/en/latest/accessing_transaction_field.html?highlight=global#id1) and [global properties](https://pyteal.readthedocs.io/en/latest/accessing_transaction_field.html?highlight=global#global-parameters). The entire program is compiled to TEAL using the `compileTeal` PyTeal function. This function compiles the program as defined by the program variable. The `compileTeal` method also sets the Mode.Signature. This lets PyTeal know this is for a smart signature and not a smart contract. The version parameter instructs PyTeal on which version of TEAL to produce when compiling. 

To test this sample, a sample address can be defined and a print command calling the `donation_escrow` function can be added to the sample.

```python
#sample_smart_sig.py
from pyteal import *

"""Basic Donation Escrow"""
def donation_escrow(benefactor):
    Fee = Int(1000)

    #Only the benefactor account can withdraw from this escrow
    program = And(
        Txn.type_enum() == TxnType.Payment,
        Txn.fee() <= Fee,
        Txn.receiver() == Addr(benefactor),
        Global.group_size() == Int(1),
        Txn.rekey_to() == Global.zero_address()
    )
    # Mode.Signature specifies that this is a smart signature
    return compileTeal(program, Mode.Signature, version=5)

test_benefactor = "CZHGG36RBYTTK36N3ZC7MENGFOL3R6D4NNEJQU3G43U5GH457SU34ZGRLY"
print( donation_escrow(test_benefactor))
```

This sample can be executed using the following command.

```bash
$ python3 sample_smart_sig.py
```

This will print out the compiled TEAL. The Algorand address of the escrow can be retrieved by first saving the produced TEAL to a file and then compiled to byte code using the `goal` command-line tool. In the next section,  using this smart signature with a transaction will be demonstrated.

```bash
$ python3 smart_sig.py > test.teal
$ ./goal clerk compile test.teal
test.teal: ZNJNTBMZKTCSO2RF4AJ3TLVFCZ5ZTHKAUBGR5AHJ23IHRFGK6GRIUVH2MU
```

# Deploying the smart signature
As stated in the previous section, smart signatures are used in conjunction with a transaction submission. In the previous section, a sample escrow was created. With escrows, any account can fund these accounts. These funds can not leave the escrow unless the logic evaluates to true. Once you have the escrow address, simple payment transactions can be used to fund it. To remove funds, a payment transaction can also be used but the transaction needs to be signed with the logic, not a private key. The following example illustrates:

* Compiling the escrow
* Funding the escrow with a simple payment transaction
* Dispensing funds using a payment transaction to the beneficiary signed with the logic

A few global variables are created and some utility functions are added to the previous section’s sample. The `benefactor_mnemonic` is the backup phrase for the address of the benefactor and the `sender_mnemonic` represents the account that will fund the escrow. Mnemonics should never be included in the source of a production environment. It is done here for learning purposes only. Key management should be handled by a proper wallet.

```python
import base64

from algosdk.future import transaction
from algosdk import mnemonic
from algosdk.v2client import algod
from pyteal import *

# user declared account mnemonics
benefactor_mnemonic = "REPLACE WITH YOUR OWN MNEMONIC"
sender_mnemonic = "REPLACE WITH YOUR OWN MNEMONIC"

# user declared algod connection parameters. Node must have EnableDeveloperAPI set to true in its config
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# helper function to compile program source
def compile_smart_signature(client, source_code):
    compile_response = client.compile(source_code)
    return compile_response['result'], compile_response['hash']
    
# helper function that converts a mnemonic passphrase into a private signing key
def get_private_key_from_mnemonic(mn) :
    private_key = mnemonic.to_private_key(mn)
    return private_key

# helper function that waits for a given txid to be confirmed by the network
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
```

The  `compile_smart_contract`, `get_private_key_from_mnemonic`, and `wait_for_confirmation` functions are explained in [Deploying and calling a smart contract](#deploying-the-contract).

A utility function is then added to create and submit a simple payment transaction.

```python
def payment_transaction(creator_mnemonic, amt, rcv, algod_client)->dict:
    params = algod_client.suggested_params()
    add = mnemonic.to_public_key(creator_mnemonic)
    key = mnemonic.to_private_key(creator_mnemonic)
    unsigned_txn = transaction.PaymentTxn(add, params, rcv, amt)
    signed = unsigned_txn.sign(key)
    txid = algod_client.send_transaction(signed)
    pmtx = wait_for_confirmation(algod_client, txid , 5)
    return pmtx
```

This function takes a creator mnemonic of the address that is creating the payment transaction as the first parameter. The amount to send and the receiver of the payment transaction are the next two parameters. The final parameter is a connection to a valid Algorand node. In this example, the sandbox installed node is used. 

In this function, the blockchain suggested parameters are retrieved from the connected node. These suggested parameters provide the default values that are required to submit a transaction, such as the expected fee for the transaction. The creator of the transaction’s address and private key are resolved from the mnemonic. The unsigned payment transaction is created using the Python SDK’s `PaymentTxn` method. This transaction is then signed with the recovered private key. As noted earlier, in a production application, the transaction should be signed by a valid wallet provider. The signed transaction is submitted to the node and the `wait_for_confirmation` utility function is called, which will return when the transaction is finalized on the blockchain. 

Another utility function is also added to create a payment transaction that is signed by the escrow logic. This function is very similar to the previous function.

```python
def lsig_payment_txn(escrowProg, escrow_address, amt, rcv, algod_client):
    params = algod_client.suggested_params()
    unsigned_txn = transaction.PaymentTxn(escrow_address, params, rcv, amt)
    encodedProg = escrowProg.encode()
    program = base64.decodebytes(encodedProg)
    lsig = transaction.LogicSig(program)
    stxn = transaction.LogicSigTransaction(unsigned_txn, lsig)
    tx_id = algod_client.send_transaction(stxn)
    pmtx = wait_for_confirmation(algod_client, tx_id, 10)
    return pmtx 
```

The primary difference is that the function is passed the base64 encoded string of the compiled bytecode for the smart signature and the escrow’s Algorand address. The program is then converted to a byte array and the Python SDK’s `LogicSig` function is used to create a logic signature from the program bytes. The payment transaction is then signed with the logic using the SDKs `LogicSigTransaction` function. For more information on Logic Signatures and smart signatures see the [smart signatures documentation](../smart-contracts/smartsigs/index.md).

The solution can be completed by adding a main function to put the utility functions to use.

```python
def main() :
    # initialize an algodClient
    algod_client = algod.AlgodClient(algod_token, algod_address)

    # define private keys
    receiver_public_key = mnemonic.to_public_key(benefactor_mnemonic)

    print("--------------------------------------------")
    print("Compiling Donation Smart Signature......")

    stateless_program_teal = donation_escrow(receiver_public_key)
    escrow_result, escrow_address= compile_smart_signature(algod_client, stateless_program_teal)

    print("Program:", escrow_result)
    print("hash: ", escrow_address)

    print("--------------------------------------------")
    print("Activating Donation Smart Signature......")

    # Activate escrow contract by sending 2 algo and 1000 microalgo for transaction fee from creator
    amt = 2001000
    payment_transaction(sender_mnemonic, amt, escrow_address, algod_client)

    print("--------------------------------------------")
    print("Withdraw from Donation Smart Signature......")

    # Withdraws 1 ALGO from smart signature using logic signature.
    withdrawal_amt = 1000000
    lsig_payment_txn(escrow_result, escrow_address, withdrawal_amt, receiver_public_key, algod_client)

```

The main function first makes a connection to the sandbox installed node, then the benefactor’s address is recovered. The `donation_escrow` built in the previous section is called to produce the TEAL for the smart signature.  This TEAL is then compiled, returning both the base64 encoded bytes of the program and the address of the escrow.

A simple payment transaction is then created to fund the escrow with a little over 2 Algos. Finally, 1 Algo is dispensed from the escrow to the benefactor using a payment transaction signed by the smart signature. The complete example is shown below.

```python
import base64

from algosdk.future import transaction
from algosdk import mnemonic
from algosdk.v2client import algod
from pyteal import *

# user declared account mnemonics
#benefactor_mnemonic = "REPLACE WITH YOUR OWN MNEMONIC"
#sender_mnemonic = "REPLACE WITH YOUR OWN MNEMONIC"
benefactor_mnemonic = "finger dizzy engage favorite purpose blade hybrid fun calm rely pink oven make calm gaze absorb book hood floor observe venue cancel question abstract army"
sender_mnemonic = "ecology average boost pony around voice daring story host brother elephant cargo drift fiber crystal bracket vivid lumber liar inquiry sketch phrase fade abstract exotic"

# user declared algod connection parameters. Node must have EnableDeveloperAPI set to true in its config
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# helper function to compile program source
def compile_smart_signature(client, source_code):
    compile_response = client.compile(source_code)
    return compile_response['result'], compile_response['hash']
    
# helper function that converts a mnemonic passphrase into a private signing key
def get_private_key_from_mnemonic(mn) :
    private_key = mnemonic.to_private_key(mn)
    return private_key

# helper function that waits for a given txid to be confirmed by the network
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

def payment_transaction(creator_mnemonic, amt, rcv, algod_client)->dict:
    params = algod_client.suggested_params()
    add = mnemonic.to_public_key(creator_mnemonic)
    key = mnemonic.to_private_key(creator_mnemonic)
    unsigned_txn = transaction.PaymentTxn(add, params, rcv, amt)
    signed = unsigned_txn.sign(key)
    txid = algod_client.send_transaction(signed)
    pmtx = wait_for_confirmation(algod_client, txid , 5)
    return pmtx

def lsig_payment_txn(escrowProg, escrow_address, amt, rcv, algod_client):
    params = algod_client.suggested_params()
    unsigned_txn = transaction.PaymentTxn(escrow_address, params, rcv, amt)
    encodedProg = escrowProg.encode()
    program = base64.decodebytes(encodedProg)
    lsig = transaction.LogicSig(program)
    stxn = transaction.LogicSigTransaction(unsigned_txn, lsig)
    tx_id = algod_client.send_transaction(stxn)
    pmtx = wait_for_confirmation(algod_client, tx_id, 10)
    return pmtx 

"""Basic Donation Escrow"""

def donation_escrow(benefactor):
    Fee = Int(1000)

    #Only the benefactor account can withdraw from this escrow
    program = And(
        Txn.type_enum() == TxnType.Payment,
        Txn.fee() <= Fee,
        Txn.receiver() == Addr(benefactor),
        Global.group_size() == Int(1),
        Txn.rekey_to() == Global.zero_address()
    )

    # Mode.Signature specifies that this is a smart signature
    return compileTeal(program, Mode.Signature, version=5)

def main() :
    # initialize an algodClient
    algod_client = algod.AlgodClient(algod_token, algod_address)

    # define private keys
    receiver_public_key = mnemonic.to_public_key(benefactor_mnemonic)

    print("--------------------------------------------")
    print("Compiling Donation Smart Signature......")

    stateless_program_teal = donation_escrow(receiver_public_key)
    escrow_result, escrow_address= compile_smart_signature(algod_client, stateless_program_teal)

    print("Program:", escrow_result)
    print("hash: ", escrow_address)

    print("--------------------------------------------")
    print("Activating Donation Smart Signature......")

    # Activate escrow contract by sending 2 algo and 1000 microalgo for transaction fee from creator
    amt = 2001000
    payment_transaction(sender_mnemonic, amt, escrow_address, algod_client)

    print("--------------------------------------------")
    print("Withdraw from Donation Smart Signature......")

    # Withdraws 1 ALGO from smart signature using logic signature.
    withdrawal_amt = 1000000
    lsig_payment_txn(escrow_result, escrow_address, withdrawal_amt, receiver_public_key, algod_client)

main()

```

For more information on smart signatures, see the [developer documentation](../smart-contracts/smartsigs/index.md).
