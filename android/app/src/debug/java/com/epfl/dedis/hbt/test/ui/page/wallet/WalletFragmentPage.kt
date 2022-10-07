package com.epfl.dedis.hbt.test.ui.page.wallet

import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object WalletFragmentPage {

    fun logout(): ViewInteraction = onView(withId(R.id.walletButtonLogout))
}