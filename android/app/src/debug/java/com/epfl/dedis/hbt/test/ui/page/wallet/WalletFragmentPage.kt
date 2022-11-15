package com.epfl.dedis.hbt.test.ui.page.wallet

import androidx.annotation.IdRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object WalletFragmentPage {

    fun logout(): ViewInteraction = onView(withId(R.id.walletButtonLogout))

    fun receive(): ViewInteraction = onView(withId(R.id.rxAmountOk))

    fun send(): ViewInteraction = onView(withId(R.id.walletButtonSend))

    fun username(): ViewInteraction = onView(withId(R.id.walletName))

    fun role(): ViewInteraction = onView(withId(R.id.walletRole))

    @IdRes
    fun walletFragmentId() = R.id.walletFragment
}