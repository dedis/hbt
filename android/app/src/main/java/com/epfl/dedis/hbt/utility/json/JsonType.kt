package com.epfl.dedis.hbt.utility.json

enum class JsonType(val schemaPath: String) {
    TRANSACTION("protocol/transaction.json"),
    PENDING_TRANSACTION("protocol/pending_transaction.json")
}