package com.epfl.dedis.hbt.data.model

/**
 * Data class that captures an incomplete transaction of tokens created by the receiver
 */
data class PendingTransaction(
    val destination: String,
    val amount: Float,
    val datetime: Long
)