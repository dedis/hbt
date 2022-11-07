package com.epfl.dedis.hbt.utility.json

enum class JsonType(val schemaPath: String) {
    COMPLETE_TRANSACTION("protocol/complete_transaction.json"),
    PENDING_TRANSACTION("protocol/pending_transaction.json")
}