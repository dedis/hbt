package com.epfl.dedis.hbt.data.transaction

/**
 * Data class that captures a transaction information to send or receive token
 */
data class CompleteTransaction(
    val source: String,
    val destination: String,
    val amount: Float,
    val datetime: Long
) {
    fun pendingTransaction() = PendingTransaction(destination, amount, datetime)
}