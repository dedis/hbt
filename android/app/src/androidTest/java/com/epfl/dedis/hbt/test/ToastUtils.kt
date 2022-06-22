package com.epfl.dedis.hbt.test

import android.os.IBinder
import android.view.WindowManager
import androidx.annotation.StringRes
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.Root
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.isDisplayed
import androidx.test.espresso.matcher.ViewMatchers.withText
import androidx.test.platform.app.InstrumentationRegistry
import org.hamcrest.Description
import org.hamcrest.TypeSafeMatcher


/** This class holds utility functions when retrieving particular elements of a view in a test  */
object ToastUtils {

    /**
     * Assert that the latest toast was shown with the expected text
     *
     * @param resId resource of the text
     * @param args arguments to the resource
     */
    fun assertToastIsDisplayedWithText(@StringRes resId: Int, vararg args: Any?) {
        val text = InstrumentationRegistry.getInstrumentation().targetContext.getString(resId, args)
        onView(withText(text)).inRoot(ToastMatcher()).check(matches(isDisplayed()))
    }

    /**
     * This matcher matches the root of a Toast
     *
     * https://stackoverflow.com/questions/28390574/checking-toast-message-in-android-espresso
     */
    class ToastMatcher : TypeSafeMatcher<Root>() {

        override fun describeTo(description: Description) {
            description.appendText("is toast")
        }

        override fun matchesSafely(root: Root): Boolean {
            val type: Int = root.windowLayoutParams.get().type
            if (type == WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY) {
                val windowToken: IBinder = root.decorView.windowToken
                val appToken: IBinder = root.decorView.applicationWindowToken
                if (windowToken === appToken) {
                    // windowToken == appToken means this window isn't contained by any other windows.
                    // if it was a window for an activity, it would have TYPE_BASE_APPLICATION.
                    return true
                }
            }
            return false
        }
    }
}

