package com.epfl.dedis.hbt.data

import androidx.test.espresso.matcher.ViewMatchers.assertThat
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.document.Document
import com.epfl.dedis.hbt.data.user.Role
import com.epfl.dedis.hbt.data.user.User
import com.epfl.dedis.hbt.data.user.UserDataSource
import com.epfl.dedis.hbt.di.JsonModule.provideObjectMapper
import com.epfl.dedis.hbt.service.document.DocumentService
import com.epfl.dedis.hbt.service.document.DocumentServiceTest
import com.epfl.dedis.hbt.service.json.JsonService
import com.epfl.dedis.hbt.test.MockSharedPreferences
import org.hamcrest.CoreMatchers.instanceOf
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.kotlin.any
import org.mockito.kotlin.doReturn
import org.mockito.kotlin.mock
import retrofit2.Call
import retrofit2.Response
import org.hamcrest.CoreMatchers.`is` as eq

/**
 * Test for the UserDataSource class
 */
@RunWith(AndroidJUnit4::class)
class UserDataSourceTest {

    private val jsonService = JsonService(provideObjectMapper())
    private val docService = mock<DocumentService> {
        on { create(any(), any(), any(), any(), any()) } doReturn mockCall(Document("ID"))
    }

    private fun <T> mockCall(body: T): Call<T> = mock {
        on { execute() } doReturn Response.success(body)
    }

    private val preferences = MockSharedPreferences()

    private val mockPortrait = DocumentServiceTest.getMockPortrait()

    private val alice = User("Alice", 12345, "XX4130X3")
    private val bob = User("Bob", 67890, "54X62C3", Role.MERCHANT)

    @Before
    fun setup() {
        preferences.reset()
    }

    @Test
    fun userDataSourceRegistrationTest() {
        val dataSource = UserDataSource(preferences, jsonService, docService)

        assertThat(
            dataSource.login(alice.name, alice.pincode),
            instanceOf(Result.Error::class.java)
        )
        assertThat(dataSource.login(bob.name, bob.pincode), instanceOf(Result.Error::class.java))

        dataSource.register(alice.name, alice.pincode, alice.passport, alice.role, mockPortrait)
        dataSource.register(bob.name, bob.pincode, bob.passport, bob.role, mockPortrait)

        assertThat(dataSource.login(alice.name, alice.pincode), eq(Result.Success(alice)))
        assertThat(dataSource.login(bob.name, bob.pincode), eq(Result.Success(bob)))
    }

    @Test
    fun userDataSourceStoresUsers() {
        val dataSource = UserDataSource(preferences, jsonService, docService)

        dataSource.register(alice.name, alice.pincode, alice.passport, alice.role, mockPortrait)
        dataSource.register(bob.name, bob.pincode, bob.passport, bob.role, mockPortrait)

        val dataSourceLoaded = UserDataSource(preferences, jsonService, docService)

        assertThat(dataSourceLoaded.login(alice.name, alice.pincode), eq(Result.Success(alice)))
        assertThat(dataSourceLoaded.login(bob.name, bob.pincode), eq(Result.Success(bob)))
    }
}