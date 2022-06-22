package com.epfl.dedis.hbt.test

import android.view.View
import android.widget.EditText
import androidx.test.espresso.UiController
import androidx.test.espresso.ViewAction
import androidx.test.espresso.matcher.ViewMatchers.isEnabled
import org.hamcrest.CoreMatchers.allOf
import org.hamcrest.Matcher

/**
 * For an EditText with inputType="number", Espresso's typeText() doesn't have any impact
 * when run on robolectric.
 *
 * Found here : https://github.com/robolectric/robolectric/issues/5110
 */
fun typeNumbers(text: String): ViewAction {
    return object : ViewAction {
        override fun getDescription(): String {
            return "force type text"
        }

        override fun getConstraints(): Matcher<View> {
            return allOf(isEnabled())
        }

        override fun perform(uiController: UiController?, view: View?) {
            (view as? EditText)?.append(text)
            uiController?.loopMainThreadUntilIdle()
        }
    }
}