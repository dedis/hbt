package com.epfl.dedis.hbt.data.model

import com.epfl.dedis.hbt.data.model.Role.BENEFICIARY

/**
 * Data class that captures user information for logged in users retrieved from LoginRepository
 */
data class User(
    val name: String,
    val pincode: Int,
    val passport: String,
    val role: Role = BENEFICIARY
    )