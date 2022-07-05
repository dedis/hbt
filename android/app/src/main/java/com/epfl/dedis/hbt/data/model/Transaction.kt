package com.epfl.dedis.hbt.data.model

/**
 * Data class that captures a transaction information to send or receive token
 */
data class Transaction(
    val sourcePk: String,
    val destinationPk: String,
    val amount: Float,
)