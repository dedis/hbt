package com.epfl.dedis.hbt.test.ui.page

import androidx.annotation.IdRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.ViewInteraction
import androidx.test.espresso.matcher.ViewMatchers.withId
import androidx.test.espresso.matcher.ViewMatchers.withParent
import com.epfl.dedis.hbt.R

object MainActivityPage {

    fun currentFragment(): ViewInteraction = onView(withParent(withId(containerId())))

    @IdRes
    fun containerId(): Int = R.id.container
}