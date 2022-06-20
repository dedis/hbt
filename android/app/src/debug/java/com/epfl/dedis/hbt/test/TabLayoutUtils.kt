package com.epfl.dedis.hbt.test

import android.view.View
import androidx.test.espresso.matcher.BoundedMatcher
import com.google.android.material.tabs.TabLayout
import org.hamcrest.Description

class IsTabSelected(private val index: Int) :
    BoundedMatcher<View, TabLayout>(TabLayout::class.java) {

    override fun describeTo(description: Description) {
        description.appendText("selected tab is")
            .appendValue(index)
    }

    override fun describeMismatch(item: Any?, description: Description) {
        if (super.matches(item)) {
            description.appendText("selected tab was")
                .appendValue((item as TabLayout).selectedTabPosition)
        } else {
            super.describeMismatch(item, description)
        }
    }

    override fun matchesSafely(item: TabLayout): Boolean {
        return item.selectedTabPosition == index
    }
}