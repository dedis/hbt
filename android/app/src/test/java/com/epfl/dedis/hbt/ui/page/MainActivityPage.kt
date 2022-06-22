package com.epfl.dedis.hbt.ui.page

import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.matcher.ViewMatchers.withId

object MainActivityPage {

    const val loginScreenIndex = 0

    fun tabLayout() = onView(withId(androidx.fragment.R.id.fragment_container_view_tag))

}