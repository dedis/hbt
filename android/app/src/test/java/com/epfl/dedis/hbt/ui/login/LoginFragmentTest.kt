package com.epfl.dedis.hbt.ui.login

import androidx.test.espresso.action.ViewActions.*
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.*
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.user.User
import com.epfl.dedis.hbt.data.user.UserRepository
import com.epfl.dedis.hbt.test.ToastUtils
import com.epfl.dedis.hbt.test.fragment.FragmentScenarioRule
import com.epfl.dedis.hbt.test.typeNumbers
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.currentFragment
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.loginButton
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.pincodeInput
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.registerButton
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.usernameInput
import com.epfl.dedis.hbt.test.ui.page.register.PassportScanFragmentPage.scanPassportFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.walletFragmentId
import dagger.hilt.android.testing.BindValue
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.junit.Rule
import org.junit.Test
import org.junit.rules.ExternalResource
import org.junit.runner.RunWith
import org.mockito.kotlin.any
import org.mockito.kotlin.doAnswer
import org.mockito.kotlin.mock
import org.mockito.kotlin.whenever

@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class LoginFragmentTest {

    @BindValue
    lateinit var userRepo: UserRepository

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val setupRule = object : ExternalResource() {
        override fun before() = setup()
    }

    @get:Rule(order = 2)
    val fragmentRule = FragmentScenarioRule.launch(LoginFragment::class.java)

    // Test data
    private val user = User("Jon Smith", 12345, "passport")

    private var currentUser: User? = null
    private var currentRepoResult: Result<User> = Result.Success(user)
    private var currentRegistered = true

    // Called before the fragment in initialized
    // allowing us to setup the dependencies
    private fun setup() {
        // Reset
        currentRepoResult = Result.Success(user)
        currentRegistered = true
        currentUser = null

        // Create mock
        userRepo = mock {
            on { login(any(), any()) } doAnswer { currentRepoResult }

            on { isRegistered(any()) } doAnswer { currentRegistered }

            on { loggedInUser } doAnswer { currentUser }

            on { isLoggedIn } doAnswer { currentUser != null }
        }
    }

    @Test
    fun loginOpensWallet() {
        // Modify the userRepo mock such that it sets the current user in the repo as user when login is called
        whenever(userRepo.login(any(), any())).thenAnswer {
            currentUser = user
            currentRepoResult
        }

        // Login procedure
        usernameInput().perform(replaceText(user.name))
        pincodeInput().perform(typeNumbers(user.pincode.toString()))
        loginButton().perform(click())

        currentFragment().check(matches(withId(walletFragmentId())))
    }


    @Test
    fun registerKeepsAlreadyGivenInput() {
        currentRegistered = false

        // Login procedure
        usernameInput().perform(replaceText(user.name))
        pincodeInput().perform(typeNumbers(user.pincode.toString()))
        registerButton().perform(click())

        // make sure the scan passport has been opened
        currentFragment().check(matches(withId(scanPassportFragmentId())))
    }

    @Test
    fun doneButtonOnKeyboardActsAsLogin() {
        // Modify the userRepo mock such that it sets the current user in the repo as user when login is called
        whenever(userRepo.login(any(), any())).thenAnswer {
            currentUser = user
            currentRepoResult
        }

        // Login procedure
        usernameInput().perform(replaceText(user.name))
        pincodeInput().perform(typeNumbers(user.pincode.toString()), pressImeActionButton())

        currentFragment().check(matches(withId(walletFragmentId())))
    }

    @Test
    fun failedLoginShowsError() {
        currentRepoResult = Result.Error(Exception())
        usernameInput().perform(typeText(user.name))
        pincodeInput().perform(typeNumbers(user.pincode.toString()))

        loginButton().check(matches(isEnabled())).perform(click())

        ToastUtils.assertToastIsDisplayedWithText(R.string.login_failed)
    }

    @Test
    fun invalidPincodeDisableLoginButton() {
        // Input a valid username
        usernameInput().perform(typeText(user.name))
        loginButton().check(matches(isNotEnabled()))
        // Input only 3 digits (4 necessary)
        pincodeInput().perform(typeNumbers("123"))
        loginButton().check(matches(isNotEnabled()))
        // Input 3 more, the button should be valid
        pincodeInput().perform(typeNumbers("123"))
        loginButton().check(matches(isEnabled()))
        // Input 4 more, total of 10 which is too much
        pincodeInput().perform(typeNumbers("1234"))
        loginButton().check(matches(isNotEnabled()))
    }
}
