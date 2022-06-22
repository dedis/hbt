package com.epfl.dedis.hbt.test.ui.page.login

import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object LoginFragmentPage {

    fun usernameInput(): ViewInteraction = onView(withId(R.id.loginUsername))

    fun pincodeInput(): ViewInteraction = onView(withId(R.id.loginPincode))

    fun loginButton(): ViewInteraction = onView(withId(R.id.loginSignin))
}