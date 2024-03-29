package com.epfl.dedis.hbt.ui.register

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.Result.Error
import com.epfl.dedis.hbt.data.Result.Success
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.data.user.Role
import com.epfl.dedis.hbt.data.user.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class RegisterViewModel @Inject constructor(private val userRepository: UserRepository) :
    ViewModel() {

    private val _registerForm = MutableLiveData<RegisterFormState>()
    val registerFormState: LiveData<RegisterFormState> = _registerForm

    private val _registerResult = MutableLiveData<RegisterResult>()
    val registerResult: LiveData<RegisterResult> = _registerResult

    fun register(
        username: String,
        pincode: String,
        passport: String,
        portrait: Portrait,
        checksum: ByteArray,
        role: Role
    ) {
        CoroutineScope(Dispatchers.IO).launch {
            // can be launched in a separate asynchronous job
            val result = userRepository.register(username, pincode, passport, role, portrait)

            CoroutineScope(Dispatchers.Main).launch {
                when (result) {
                    is Success -> _registerResult.value = RegisterResult(error = null)
                    is Error -> {
                        Log.e("Register", "Failed to register", result.exception)
                        _registerResult.value = RegisterResult(error = R.string.login_failed)
                    }
                }
            }
        }
    }

    fun registerDataChanged(username: String, pincode: String) {
        var isValid = true
        var userName: Int? = null
        var pinError: Int? = null

        if (!isUserNameValid(username)) {
            isValid = false
            userName = R.string.invalid_username
        }

        if (!isPincodeValid(pincode)) {
            isValid = false
            pinError = R.string.invalid_pin_code
        }

        _registerForm.value = RegisterFormState(
            usernameError = userName,
            pincodeError = pinError,
            isDataValid = isValid
        )
    }

    // Validate username
    private fun isUserNameValid(username: String): Boolean {
        return username.length >= 2
    }

    // Validate pincode
    private fun isPincodeValid(pincode: String): Boolean {
        return pincode.length in 4..9
    }
}