package com.epfl.dedis.hbt.ui.login

/**
 * Data validation state of the login form.
 */
data class LoginFormState(
    val usernameError: Int? = null,
    val pincodeError: Int? = null,
    val isDataValid: Boolean = false,
    val isUserRegistered: Boolean = false
)