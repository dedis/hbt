package com.epfl.dedis.hbt.ui.wallet

import androidx.test.espresso.action.ViewActions.click
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.withId
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.user.User
import com.epfl.dedis.hbt.data.user.UserRepository
import com.epfl.dedis.hbt.data.user.Wallet
import com.epfl.dedis.hbt.test.fragment.FragmentScenarioRule
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.currentFragment
import com.epfl.dedis.hbt.test.ui.page.login.LoginFragmentPage.loginFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.logout
import dagger.hilt.android.testing.BindValue
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.junit.Rule
import org.junit.Test
import org.junit.rules.ExternalResource
import org.junit.runner.RunWith
import org.mockito.kotlin.*

@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class WalletFragmentTest {

    @BindValue
    lateinit var userRepo: UserRepository

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val setupRule = object : ExternalResource() {
        override fun before() = setup()
    }

    @get:Rule(order = 2)
    val fragmentRule = FragmentScenarioRule.launch(WalletFragment::class.java)

    private val user = User("Jon Smith", 12345, "passport")
    private val wallet = Wallet()

    private fun setup() {
        userRepo = mock {
            on { loggedInUser } doAnswer { user }

            on { isLoggedIn } doAnswer { true }

            on { wallet } doAnswer { wallet }

            on { isRegistered(any()) } doReturn true
        }
    }

    @Test
    fun logoutCallsLogoutAndGoesToLogin() {
        logout().perform(click())

        currentFragment().check(matches(withId(loginFragmentId())))
        verify(userRepo).logout()
    }
}