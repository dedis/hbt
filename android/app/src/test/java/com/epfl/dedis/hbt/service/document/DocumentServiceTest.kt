package com.epfl.dedis.hbt.service.document

import androidx.test.espresso.matcher.ViewMatchers.assertThat
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.data.user.Role
import com.epfl.dedis.hbt.data.user.User
import com.epfl.dedis.hbt.di.HttpModule
import com.epfl.dedis.hbt.di.JsonModule
import org.junit.Ignore
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.kotlin.isNull
import org.mockito.kotlin.notNull
import java.io.FileNotFoundException

@RunWith(AndroidJUnit4::class)
class DocumentServiceTest {


    companion object {
        private val retrofit =
            HttpModule.provideRetrofit(
                "http://localhost:3000",
                JsonModule.provideObjectMapper()
            )

        fun getMockPortrait(): Portrait {
            val stream = DocumentServiceTest::class.java.getResourceAsStream("/mock-portrait.jpeg")
                ?: throw FileNotFoundException()

            stream.use {
                return Portrait("image/jpeg", it.readBytes())
            }
        }
    }

    @Test
    @Ignore("This is merely a PoC")
    fun createSimpleDocumentSucceed() {
        val service = HttpModule.provideDocumentService(retrofit)

        val user = User(
            name = "Jon Smith",
            pincode = 24256,
            passport = "ABCDEFGHI",
            role = Role.BENEFICIARY
        )

        val call = service.create(user, getMockPortrait(), false)

        val response = call.execute()

        assertThat(response.errorBody(), isNull())
        assertThat(response.body(), notNull())
    }
}