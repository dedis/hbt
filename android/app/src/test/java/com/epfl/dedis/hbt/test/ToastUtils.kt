package com.epfl.dedis.hbt.test

import android.content.Context
import androidx.annotation.StringRes
import androidx.test.core.app.ApplicationProvider
import org.hamcrest.MatcherAssert
import org.hamcrest.Matchers
import org.junit.Assert
import org.robolectric.shadows.ShadowToast

/** This class holds utility functions when retrieving particular elements of a view in a test  */
object ToastUtils {
    /**
     * Assert that the latest toast was shown with the expected text
     *
     * @param resId resource of the text
     * @param args arguments to the resource
     */
    fun assertToastIsDisplayedWithText(@StringRes resId: Int, vararg args: Any?) {
        MatcherAssert.assertThat(
            "No toast was displayed",
            ShadowToast.getLatestToast(),
            Matchers.notNullValue()
        )
        val expected = ApplicationProvider.getApplicationContext<Context>().getString(resId, *args)
        Assert.assertEquals(expected, ShadowToast.getTextOfLatestToast())
    }
}