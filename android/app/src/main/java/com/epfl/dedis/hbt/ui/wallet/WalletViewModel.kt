package com.epfl.dedis.hbt.ui.wallet

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class WalletViewModel @Inject constructor(private val userRepository: UserRepository) :
    ViewModel() {

    private val _walletForm = MutableLiveData<WalletFormState>()
    val walletFormState: LiveData<WalletFormState> = _walletForm

    private val _walletResult = MutableLiveData<WalletResult>()
    val walletResult: LiveData<WalletResult> = _walletResult

    private val _transactionState = MutableLiveData<TransactionState>(TransactionState.None)
    val transactionState: LiveData<TransactionState> = _transactionState

    val user =
        userRepository.loggedInUser ?: throw IllegalStateException("User should be logged in")
    val wallet = userRepository.wallet

    fun send(amount: Float) {
        // can be launched in a separate asynchronous job
        if (wallet?.send(amount, "destination") == true) {
            _walletResult.value = WalletResult(error = null)
        } else {
            _walletResult.value = WalletResult(error = R.string.wallet_send_failed)
        }
    }

    fun receive(amount: Float) {
        // can be launched in a separate asynchronous job
        if (wallet?.receive("source", amount) == true) {
            _walletResult.value = WalletResult(error = null)
        } else {
            _walletResult.value = WalletResult(error = R.string.wallet_receive_failed)
        }
    }

    fun logout() {
        userRepository.logout()
    }

    fun transitionTo(newState: TransactionState) {
        if (_transactionState.value != newState) {
            _transactionState.postValue(newState)
        }
    }
}
