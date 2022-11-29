package com.epfl.dedis.hbt.data

import androidx.test.ext.junit.runners.AndroidJUnit4
import org.hamcrest.CoreMatchers.`is`
import org.hamcrest.MatcherAssert.assertThat
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class PassportTest {

    @Test
    fun canParseValidPassport() {
        val text = """PXABCDUMMY<<BOB<SMITH<<<<<<<<<<<<<<<<<<<<<<<
                      A1234567<6ABC0102030X0405063<<<<<<<<<<<<<<<0"""
        val passport = Passport.match(text)
        assertThat(
            passport,
            `is`(
                Passport(
                    "ABC",
                    "DUMMY",
                    "BOB SMITH",
                    "A1234567",
                    "010203",
                    "040506"
                )
            )
        )
    }
}