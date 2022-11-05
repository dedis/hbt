package com.epfl.dedis.hbt.ui.wallet

import com.epfl.dedis.hbt.data.model.PendingTransaction
import com.epfl.dedis.hbt.data.model.Transaction

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
    data class ReceiverShow(val amount: Float, val datetime: Long) : TransactionState()
    object SenderRead : TransactionState()
    data class SenderShow(val transaction: Transaction) : TransactionState()
    data class ReceiverRead(val expected: PendingTransaction) : TransactionState()
}
