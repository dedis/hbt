package com.epfl.dedis.hbt.integration

import androidx.test.espresso.action.ViewActions.*
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.*
import androidx.test.ext.junit.rules.ActivityScenarioRule
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.filters.LargeTest
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.currentFragment
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage
import com.epfl.dedis.hbt.test.ui.page.login.RegisterFragmentPage
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage
import com.epfl.dedis.hbt.ui.MainActivity
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.hamcrest.Matchers.allOf
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

/** Integration tests of various scenarios using the login service */
@LargeTest
@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class LoginTest {

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Test
    fun registerLogoutLoginScenario() {
        // Register
        LoginFragmentPage.registerButton()
            .check(matches(allOf(isDisplayed(), isEnabled(), withText(R.string.action_register))))
            .perform(click())
        // Assert Register fragment is opened
        currentFragment().check(matches(withId(RegisterFragmentPage.registerFragmentId())))

        // Enter register information
        RegisterFragmentPage.usernameInput()
            .check(matches(isDisplayed()))
            .perform(replaceText("Jon Smith"))

        RegisterFragmentPage.pincodeInput()
            .check(matches(isDisplayed()))
            .perform(replaceText("12345"))

        RegisterFragmentPage.passportInput()
            .check(matches(isDisplayed()))
            .perform(replaceText("ABCDEFGHI"), closeSoftKeyboard())

        // Register user
        RegisterFragmentPage.registerButton()
            .check(matches(allOf(isDisplayed(), isEnabled(), withText(R.string.action_register))))
            .perform(click())

        // Assert the wallet fragment is opened and has correct values
        currentFragment().check(matches(withId(WalletFragmentPage.walletFragmentId())))
        WalletFragmentPage.username()
            .check(matches(allOf(isDisplayed(), withText("Jon Smith"))))
        WalletFragmentPage.role()
            .check(matches(allOf(isDisplayed(), withText(R.string.role_beneficiary))))

        // Logout
        WalletFragmentPage.logout()
            .check(matches(allOf(isDisplayed(), withText(R.string.wallet_button_logout))))
            .perform(click())

        // Make sure the login fragment is reopened
        currentFragment().check(matches(withId(LoginFragmentPage.loginFragmentId())))

        // Enter login information
        LoginFragmentPage.usernameInput()
            .check(matches(isDisplayed()))
            .perform(replaceText("Jon Smith"))

        LoginFragmentPage.pincodeInput()
            .check(matches(isDisplayed()))
            .perform(replaceText("12345"))

        // Login
        LoginFragmentPage.loginButton()
            .check(matches(allOf(isDisplayed(), withText(R.string.action_sign_in))))
            .perform(click())

        // Assertions
        currentFragment().check(matches(withId(WalletFragmentPage.walletFragmentId())))
        WalletFragmentPage.username()
            .check(matches(allOf(isDisplayed(), withText("Jon Smith"))))
        WalletFragmentPage.role()
            .check(matches(allOf(isDisplayed(), withText(R.string.role_beneficiary))))
    }
}