package com.epfl.dedis.hbt.test.ui.page.wallet

import androidx.annotation.IdRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object ShowQrFragmentPage {

    fun showOk(): ViewInteraction = onView(withId(R.id.walletButtonOk))

    @IdRes
    fun showFragmentId() = R.id.showQRWalletFragment
}