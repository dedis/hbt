package com.epfl.dedis.hbt.ui.wallet

/**
 * Data state of the wallet form.
 */
data class WalletFormState(
    val username: String? = null,
    val balance: Float = 0F
)