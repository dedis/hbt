package com.epfl.dedis.hbt.service.json

enum class JsonType(val schemaPath: String) {
    COMPLETE_TRANSACTION("protocol/complete_transaction.json"),
    PENDING_TRANSACTION("protocol/pending_transaction.json"),
    USER_DATA("store/user_data.json")
}