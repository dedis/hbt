package com.epfl.dedis.hbt.data.model

/**
 * Data class that captures a transaction information to send or receive token
 */
data class PendingTransaction(
    val destination: String,
    val amount: Float,
    val datetime: Long
)