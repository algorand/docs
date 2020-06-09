title: Install the Indexer 🆕

The Algorand Indexer is a feature that enables searching the blockchain for transactions, assets, accounts, and blocks with various criteria. Currently, Algorand has a V1 and V2 Indexer. The V1 Indexer is deprecated and users should now use the V2 Indexer. The V2 Indexer runs as an independent process that must connect to a [PostgreSQL](https://www.postgresql.org/) compatible database that contains the ledger data. The PostgeSQL database is populated by the indexer which connects to an Algorand node and processes all the ledger data and loads the database. The node the Indexer connects to must be an archival node to get all the ledger data. Alternatively, the Indexer can just connect to a PostgresSQL database that is populated by another instance of Indexer. This allows reader instances to be set up that provide the REST APIs for searching the database and another Indexer to be responsible for loading the ledger data.


The V2 Indexer is network agnostic, meaning it can point at BetaNet, TestNet, or MainNet. 

The source code for the Indexer is provided on [github](https://github.com/algorand/indexer).

For details on Indexer usage, read the [Searching the Blockchain](../../features/indexer.md) feature guide. See [Indexer README](https://github.com/algorand/indexer) for more details on running the Indexer.

# Indexer V2
To Install the new follow the instructions below. The Indexer binaries are available on [github](https://github.com/algorand/indexer/releases).

!!! info
    Additional install methods will be available in the near future.

## Download the Indexer Binaries
Download the Indexer binaries for specific operating system.

* Linux binaries ([AMD64](https://github.com/algorand/indexer/releases/download/2.0.0/algorand-indexer_linux_arm64_2.0.0.tar.bz2), [ARM64](https://github.com/algorand/indexer/releases/download/2.0.0/algorand-indexer_linux_arm64_2.0.0.tar.bz2), [ARM32](https://github.com/algorand/indexer/releases/download/2.0.0/algorand-indexer_linux_arm_2.0.0.tar.bz2))
* [Mac](https://github.com/algorand/indexer/releases/download/2.0.0/algorand-indexer_darwin_amd64_2.0.0.tar.bz2) binaries.
  
## Extract the binaries to a specific directory. 
The location does not matter. In these instructions, an indexer folder is used which is located in the current accounts home directory.

```bash
$ mkdir ~/indexer
$ cd /location-of-downloaded-tar
$ tar -xf <your-os-tarfile> -C ~/indexer
```

## Run the Indexer
The Indexer primarily provides two services, loading a PostgreSQL database with ledger data and supplying a REST API to search this ledger data. You can set the Indexer to point at a database that was loaded by another instance of the Indexer. The database does not have to be on the current node. In fact, you can have one Indexer that loads the database and many Indexers that share this data through their REST APIs. How the Indexer operates is determined with parameters that are passed to the Indexer as it is started.

The Indexer has many options which can be seen using the -h option when running the Indexer binary.

```bash
$ ./algorand-indexer daemon -h
run indexer daemon. Serve api on HTTP.

Usage:
  indexer daemon [flags]

Flags:
  -d, --algod string         path to algod data dir, or $ALGORAND_DATA
      --algod-net string     host:port of algod
      --algod-token string   api access token for algod
  -c, --config string        path to 'key: value' config file, keys are same as command line options
      --dev-mode             allow performance intensive operations like searching for accounts at a particular round
  -g, --genesis string       path to genesis.json (defaults to genesis.json in algod data dir if that was set)
  -h, --help                 help for daemon
      --no-algod             disable connecting to algod for block following
  -S, --server string        host:port to serve API on (default :8980) (default ":8980")
  -t, --token string         an optional auth token, when set REST calls must use this token in a bearer format, or in a 'X-Indexer-API-Token' header

Global Flags:
      --cpuprofile string   file to record cpu profile to
  -n, --dummydb             use dummy indexer db
      --pidfile string      file to write daemon's process id to
  -P, --postgres string     connection string for postgres database
```

To start the Indexer as a reader (ie not connecting to an Algorand node), supply the the `--postgres` or `-P` option when running the indexer. The value should be a valid connection string for a postgres database.

```bash
$ ./algorand-indexer daemon -P “host=[your-host] port=[your-port] user=[uname] password=[password] dbname=[ledgerdb] sslmode=disable”  --no-algod
```

To start the Indexer so it populates the PostgreSQL database, supply the Algorand Archival node connection details. This can be done by either specifying the data directory (`--algod`), if the node is on the same machine as the Indexer, or by supplying the algod network host and port string (`--algod-net`) and the proper API token (`--token`). The database needs to be created and running prior to starting the Indexer.

```bash
# start with local data directory
$ ./algorand-indexer daemon -P “host=[your-host] port=[your-port] user=[uname] password=[password] dbname=[ledgerdb] sslmode=disable” --algod=~/node/data

# start with networked Algorand node
$ ./algorand-indexer daemon -P “host=[your-host] port=[your-port] user=[uname] password=[password] dbname=[ledgerdb] sslmode=disable” --algod-net="http://[your-host]:[your-port]" ---algod-token="[your-api-token]

```

!!! info
    The initial loading of the Indexer Database will take a considerable amount of time.


## REST API Token and Server
When starting the Indexer, a REST API is exposed. To control access to this API you can you use the `--token` parameter, which allows specifying any desired token. REST API clients will be required to pass this token in their calls in order to return successful searches. The REST API defaults to serving on port 8980. This can be changed by supply a [host:port] value to the Indexer with the `--server` option.

# Indexer V1

!!! info
     This section explains using V1 of the indexer and should be considered deprecated.

All transaction searching by default is limited to a 1000 round range. If a node is configured in archival mode, an additional configuration option can be used to turn on a node indexer and remove this restriction. With the indexer turned on, searching for specific transactions will be quicker. Two additional REST calls are also made available for more refined searching. 

The two Additional REST calls are:

```
/v1/transaction/{txid}
```
This call allows quickly locating a transaction using the txid
See [REST API Reference](../../reference/rest-apis/algod/v1.md#get-v1transactiontxid) for more details.

```
/v1/account/{account}/transactions?fromDate=YYYY-MM-DD&toDate=YYYY-MM-DD) 
```

This call allows locating all transactions within a date range. Date parameters support RFC3339 (ie 2006-01-02T15:04:05Z07:00).
See [REST API Reference](../../reference/rest-apis/algod/v1.md#get-v1accountaddresstransactions) for more details.

To turn on indexing for a node, the `isIndexerActive` configuration parameter must be set to `true`. The [Node Configuration](../../reference/node/config.md) guide describes setting node configuration properties.

!!! warning
     Turning on indexing with a node will increase the disk space required by the node.

!!! info
    Indexing on a node is only allowed with nodes that have archival mode turned on.
