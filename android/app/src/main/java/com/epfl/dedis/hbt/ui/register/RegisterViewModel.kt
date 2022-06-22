package com.epfl.dedis.hbt.ui.register

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.Result.Success
import com.epfl.dedis.hbt.data.UserRepository
import com.epfl.dedis.hbt.data.model.Role
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class RegisterViewModel @Inject constructor(private val userRepository: UserRepository) :
    ViewModel() {

    private val _registerForm = MutableLiveData<RegisterFormState>()
    val registerFormState: LiveData<RegisterFormState> = _registerForm

    private val _registerResult = MutableLiveData<RegisterResult>()
    val registerResult: LiveData<RegisterResult> = _registerResult

    fun register(username: String, pincode: String, passport: String, role: Role) {
        // can be launched in a separate asynchronous job
        val result = userRepository.register(username, pincode, passport, role)

        if (result is Success) {
            _registerResult.value = RegisterResult(error = null)
        } else {
            _registerResult.value = RegisterResult(error = R.string.login_failed)
        }
    }

    fun registerDataChanged(username: String, pincode: String, passport: String) {
        if (!isUserNameValid(username)) {
            _registerForm.value = RegisterFormState(usernameError = R.string.invalid_username)
        } else if (!isPincodeValid(pincode)) {
            _registerForm.value = RegisterFormState(pincodeError = R.string.invalid_pin_code)
        } else if (!isPassportValid(passport)) {
            _registerForm.value = RegisterFormState(passportError = R.string.invalid_passport)
        } else {
            _registerForm.value = RegisterFormState(isDataValid = true)
        }
    }

    // Validate username
    private fun isUserNameValid(username: String): Boolean {
        return username.length >= 2
    }

    // Validate pincode
    private fun isPincodeValid(pincode: String): Boolean {
        return pincode.length in 4..9
    }

    // Validate passport number
    private fun isPassportValid(passport: String): Boolean {
        return passport.length >= 9
    }
}