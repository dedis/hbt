package com.epfl.dedis.hbt.ui.wallet

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.epfl.dedis.hbt.data.user.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class WalletViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _walletForm = MutableLiveData<WalletFormState>()
    val walletFormState: LiveData<WalletFormState> = _walletForm

    val user =
        userRepository.loggedInUser ?: throw IllegalStateException("User should be logged in")
    val wallet = userRepository.wallet

    fun logout() {
        userRepository.logout()
    }
}
