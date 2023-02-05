package com.epfl.dedis.hbt.data.transaction

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import com.epfl.dedis.hbt.data.transaction.TransactionState.*
import javax.inject.Inject
import javax.inject.Singleton

/**
 * This service manages the state of a transaction. It stores the current state of a transaction and
 * handles its transitions
 */
@Singleton
class TransactionStateManager @Inject constructor() {

    private val _currentState: MutableLiveData<TransactionState> = MutableLiveData(None)
    val currentState: LiveData<TransactionState> = _currentState

    /**
     * Start a new transaction for a sender
     *
     * Must be in the None state
     * */
    fun startSendingTransaction() {
        validateCurrentState(None)
        transitionTo(SenderRead)
    }

    /**
     * Start a new transaction for a receiver
     *
     * Must be in the None state
     *
     * @param pendingTrx transaction proposed by the receiver
     */
    fun startReceivingTransaction(pendingTrx: PendingTransaction) {
        validateCurrentState(None)
        transitionTo(ReceiverShow(pendingTrx))
    }

    /**
     * Move to the show state of the sender:
     * The sender will now show the completed transaction built from the receiver's pending
     * transaction with its identifier
     *
     * Must be called in the SenderRead state
     *
     * @param completeTrx The complete transaction created by the sender
     */
    fun showCompleteTransaction(completeTrx: CompleteTransaction) {
        validateCurrentState(SenderRead)
        transitionTo(SenderShow(completeTrx))
    }

    /**
     * Move to the read state of the receiver:
     * The reader will now scan the sender's QRCode to learn the source of the transaction
     *
     * Must be called in the ReceiverShow state
     *
     * @param pendingTrx The pending transaction created by the receiver
     */
    fun readCompleteTransaction(pendingTrx: PendingTransaction) {
        validateCurrentState(ReceiverShow(pendingTrx))
        transitionTo(ReceiverRead(pendingTrx))
    }

    /**
     * Complete the transaction from the sender's end
     *
     * Must be called in the SenderShow state
     *
     * @param completeTrx The complete transaction as showed to the receiver
     */
    fun completeSending(completeTrx: CompleteTransaction) {
        validateCurrentState(SenderShow(completeTrx))
        transitionTo(None)

        send(completeTrx)
    }

    /**
     * Complete the transaction from the receiver's end
     *
     * Must be called in the ReceiverRead state
     *
     * @param completeTrx The complete transaction as read from the sender.
     *        Its data must be coherent with the one stored in ReceiverRead state
     */
    fun completeReceiving(completeTrx: CompleteTransaction) {
        validateCurrentState(ReceiverRead(completeTrx.pendingTransaction()))
        transitionTo(None)

        receive(completeTrx)
    }

    /** Cancel the current transaction and move to the None state */
    fun cancelTransaction() {
        transitionTo(None)
    }

    private fun validateCurrentState(expected: TransactionState) {
        if (currentState.value != expected) {
            throw IllegalStateException("The current state is ${currentState.value} but expected $expected")
        }
    }

    private fun transitionTo(newState: TransactionState) {
        if (_currentState.value != newState) {
            _currentState.postValue(newState)
        }
    }

    private fun send(transaction: CompleteTransaction) {
        // TODO Probably use a transaction repo to store and dispatch transactions
        Log.i("Wallet VM", "Sending $transaction")
    }

    private fun receive(transaction: CompleteTransaction) {
        // TODO Probably use a transaction repo to store and dispatch transactions
        Log.i("Wallet VM", "Receiving $transaction")
    }
}