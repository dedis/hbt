package com.epfl.dedis.hbt.test.ui.page.login

import androidx.annotation.IdRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object RegisterFragmentPage {

    fun usernameInput(): ViewInteraction = onView(withId(R.id.registerUsername))

    fun pincodeInput(): ViewInteraction = onView(withId(R.id.registerPincode))

    fun passportInput(): ViewInteraction = onView(withId(R.id.passport_number))

    fun registerButton(): ViewInteraction = onView(withId(R.id.registerRegister))

    @IdRes
    fun registerFragmentId(): Int = R.id.registerFragment
}