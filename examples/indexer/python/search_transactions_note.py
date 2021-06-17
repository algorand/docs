# search_transactions_note.py
import base64
import json
# requires Python SDK version 1.3 or higher
from algosdk.v2client import indexer

# instantiate indexer client
myindexer = indexer.IndexerClient(indexer_token="", indexer_address="http://localhost:8980")

import base64

note_prefix = 'showing prefix'.encode()
address = "IAMIRIFW3ERXIMR5LWNYHNK7KRTESUGS4QHOPKF2GL3CLHWWGW32XWB7OI"

response = myindexer.search_transactions(
    address=address,
    note_prefix=note_prefix)

print("note_prefix = " +
      json.dumps(response, indent=2, sort_keys=True))
