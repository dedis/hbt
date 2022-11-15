package com.epfl.dedis.hbt.data.model

/**
 * Data class that captures an incomplete transaction of tokens created by the receiver
 */
data class PendingTransaction(
    val destination: String,
    val amount: Float,
    val datetime: Long
) {
    fun withSource(source: String): CompleteTransaction =
        CompleteTransaction(source, destination, amount, datetime)
}