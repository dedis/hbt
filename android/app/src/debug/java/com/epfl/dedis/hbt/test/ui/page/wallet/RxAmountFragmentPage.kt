package com.epfl.dedis.hbt.test.ui.page.wallet

import androidx.annotation.IdRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object RxAmountFragmentPage {

    fun rxAmount(): ViewInteraction = onView(withId(R.id.walletRxAmount))

    fun rxAmountOk(): ViewInteraction = onView(withId(R.id.rxAmountOk))

    @IdRes
    fun rxAmountFragmentId() = R.id.rxAmountWalletFragment
}