package com.epfl.dedis.hbt.test.ui.page

import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.matcher.ViewMatchers.withId
import com.epfl.dedis.hbt.R

object MainActivityPage {

    const val loginScreenIndex = 0

    fun tabLayout() = onView(withId(R.id.main_tab_layout))

}