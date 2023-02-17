The Examples are critical to have in working order and up to date.

We'll put runnable source files in the `examples` directory of the respective SDKs as well as in PyTeal and Beaker.

Goal and TEAL source location tbd

## REST APIs

We should show how to initialize a new client, set key header and do some very
simple calls against them. Refer to SDK docs or REST API docs for full functionality 

### Algod Client:

[Reference Page for Python](./docs/sdks/python/index.md) 
[Reference Page for JS](./docs/sdks/javascript/index.md) 
[Reference Page for Go](./docs/sdks/go/index.md) 
[Reference Page for Java](./docs/sdks/java/index.md) 
    
    CREATE_CLIENT - show how to init client and connect to algod and pass api keys in header
    USE_CLIENT - show how to get status/suggested params

### Indexer Client:

[Reference page](./docs/get-details/indexer.md)

    CREATE_INDEXER_CLIENT - show how to init client and connect to indexer and pass api keys in header
    USE_INDEXER_CLIENT - show how to get status/transactions 
    INDEXER_CLIENT_PAGINAGE - show how to paginate through results

### KMD:

    CREATE_KMD_CLIENT - Show how to init client and connect to kmd
    USE_KMD_CLIENT - Show how to get a wallet handle and grab accounts using the handle


## Accounts

[Reference Page](./docs/get-details/accounts/create.md)

    GENERATE_ACCOUNT - generate a new account with the SDK
    RECOVER_ACCOUNT_MNEMONIC - Create an account from a known mnemonic and print pubkey/privatekey
    RECOVER_ACCOUNT_PRIVATE_KEY - Create an account from a known private key and print pubkey/mnemonic

## ASAs

[Reference Page](./docs/get-details/asa.md)

### All the things

    CREATE_ASSET - Create an Asset from acct1 
    OPTIN_ASSET - Opt acct2 in to the newly created asset
    XFER_ASSET - Send asset from acct1 to acct2 
    FREEZE_ASSET - Freeze the asset in acct2 and then unfreeze the asset in acct2 
    CONFIGURE_ASSET - Reconfigure Asset to remove the freeze address 
    CLAWBACK_ASSET - Clawback Asset from acct2
    DELETE_ASSET - Delete the asset


## Atomic Transactions

### No ATC

[Reference Page](./docs/get-details/atomic-transfers.md)

    MANUAL_ATOMIC_GROUP - 
        Create a couple Payment transactions, Group them,
        print the group ids and txids, sign them, send them
    MANUAL_GROUP_TXIDS - 
        Show how to get _all_ the results back (not just the first txid)


### With ATC

[Reference Page](./docs/get-details/atc.md)

    ATC_CONTRACT_INIT - Create a contract object
    ATC_GROUP_TRANSACTION - Create a couple TransactionWithSigners (1 payment, 1 method call using the contract object)
    ATC_GROUP_TXIDS - Show how you can get the transaction group and print the group ids and txids
    ATC_RESULTS - Execute the group, Show how to get the results (esp for non method calls)

## Encoding/Decoding

We should show common encoding/decoding on [this page](./docs/get-details/encoding.md)

    CODEC_ADDRESS - Address
    CODEC_UINT64 - Uint64
    CODEC_STATE - Global/local state - how to decode, what the types mean 
    CODEC_ABI_TYPES - ABI types (using something like ABIType.from("(uint64,address)"))
    CODEC_TRANSACTION- Transactions (Signed and unsigned)
    CODEC_BLOCKS - Blocks

## 