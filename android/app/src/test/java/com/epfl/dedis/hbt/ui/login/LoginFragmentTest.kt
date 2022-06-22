package com.epfl.dedis.hbt.ui.login

import androidx.test.espresso.action.ViewActions.typeText
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.isEnabled
import androidx.test.espresso.matcher.ViewMatchers.isNotEnabled
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.UserRepository
import com.epfl.dedis.hbt.data.model.User
import com.epfl.dedis.hbt.test.fragment.FragmentScenarioRule
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.loginButton
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.pincodeInput
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.usernameInput
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
    private val user = User("Jon Smith", 12345)
    private var currentRepoResult: Result<User> = Result.Success(user)

    // Called before the fragment in initialized
    // allowing us to setup the dependencies
    private fun setup() {
        userRepo = mock {
            on { login(any(), any()) } doAnswer { currentRepoResult }
        }
    }

    @Test
    fun invalidPincodeDisableLoginButton() {
        // Input a valid username
        usernameInput().perform(typeText(user.name))
        loginButton().check(matches(isNotEnabled()))
        // Input only 3 digits (4 necessary)
        pincodeInput().perform(typeText("123"))
        loginButton().check(matches(isNotEnabled()))
        // Input 3 more, the button should be valid
        pincodeInput().perform(typeText("123"))
        loginButton().check(matches(isEnabled()))
        // Input 4 more, total of 10 which is too much
        pincodeInput().perform(typeText("1234"))
        loginButton().check(matches(isNotEnabled()))
    }
}
