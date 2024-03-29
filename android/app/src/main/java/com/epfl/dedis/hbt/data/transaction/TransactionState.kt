package com.epfl.dedis.hbt.data.transaction

/**
 * Describe the five states of a transaction process :
 *  - None : there is no transaction being made
 *  - ReceiverShow : the receiver shows the pending transaction
 *  - SenderRead : the sender reads the pending transaction
 *  - SenderShow : the sender shows the complete transaction
 *  - ReceiverRead : the receiver reads and validate the complete transaction
 */
sealed class TransactionState {
    object None : TransactionState()
    data class SenderShow(val transaction: CompleteTransaction) : TransactionState()
    object SenderRead : TransactionState()
    data class ReceiverRead(val expected: PendingTransaction) : TransactionState()
    data class ReceiverShow(val transaction: PendingTransaction) : TransactionState()
}
