package com.epfl.dedis.hbt.utility.json

import com.epfl.dedis.hbt.data.model.CompleteTransaction
import com.epfl.dedis.hbt.data.model.PendingTransaction
import kotlin.reflect.KClass

sealed class JsonType<T : Any>(val schemaPath: String, val type: KClass<T>) {
    object CompleteTransactionType :
        JsonType<CompleteTransaction>(
            "protocol/complete_transaction.json",
            CompleteTransaction::class
        )

    object PendingTransactionType :
        JsonType<PendingTransaction>(
            "protocol/pending_transaction.json",
            PendingTransaction::class
        )

    companion object {
        val TYPES: List<JsonType<*>> =
            JsonType::class.sealedSubclasses.mapNotNull { it.objectInstance }
    }
}