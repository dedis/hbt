package com.epfl.dedis.hbt.ui

import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.ext.junit.rules.ActivityScenarioRule
import com.epfl.dedis.hbt.test.IsTabSelected
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.loginScreenIndex
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.tabLayout
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.junit.Rule
import org.junit.Test

@HiltAndroidTest
class MainActivityTest {

    @get:Rule(order = 0)
    var hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    var activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Test
    fun mainActivityStartsOnLoginScreen() {
        tabLayout().check(matches(IsTabSelected(loginScreenIndex)))
    }
}