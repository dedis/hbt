package com.epfl.dedis.hbt.data

import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.service.passport.mrz.MRZExtractor
import com.epfl.dedis.hbt.service.passport.mrz.MRZInfo
import org.hamcrest.CoreMatchers.`is`
import org.hamcrest.MatcherAssert.assertThat
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class MRZInfoTest {

    @Test
    fun canParseValidPassport() {
        val text = """PXABCDUMMY<<BOB<SMITH<<<<<<<<<<<<<<<<<<<<<<<
                      A1234567<6ABC0102030X0405063<<<<<<<<<<<<<<<0"""
        assertThat(
            MRZExtractor.match(text),
            `is`(
                Result.Success(
                    MRZInfo(
                        "A1234567",
                        "010203",
                        "040506",
                        "ABC",
                        "DUMMY",
                        "BOB SMITH"
                    )
                )
            )
        )
    }
}