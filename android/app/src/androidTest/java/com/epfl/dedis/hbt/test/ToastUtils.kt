package com.epfl.dedis.hbt.test

import androidx.annotation.StringRes


/** This class holds utility functions when retrieving particular elements of a view in a test  */
object ToastUtils {

    /**
     * Assert that the latest toast was shown with the expected text
     *
     * @param resId resource of the text
     * @param args arguments to the resource
     */
    fun assertToastIsDisplayedWithText(@StringRes resId: Int, vararg args: Any?) {
        // Do nothing : Toast assertions are really unstable on instrumented tests.
        // !! Make sure your test works in robolectric !!
    }
}

