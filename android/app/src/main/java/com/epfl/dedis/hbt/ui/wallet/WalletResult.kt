package com.epfl.dedis.hbt.ui.wallet

/**
 * Wallet transaction result : success or error message.
 */
data class WalletResult(
    val error: Int? = null
)