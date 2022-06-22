package com.epfl.dedis.hbt.ui.register

/**
 * Data validation state of the register form.
 */
data class RegisterFormState(
    val usernameError: Int? = null,
    val pincodeError: Int? = null,
    val passportError: Int? = null,
    val isDataValid: Boolean = false
)