package com.epfl.dedis.hbt.data

import com.epfl.dedis.hbt.data.model.User
import org.hamcrest.CoreMatchers.nullValue
import org.hamcrest.MatcherAssert.assertThat
import org.junit.Test
import org.mockito.kotlin.*
import org.hamcrest.CoreMatchers.`is` as eq

class UserRepositoryTest {

    private val username = "Jon Smith"
    private val pincode = 12345
    private val passport = "12jdjpdwa"

    private val invalidPincode = "abcdef"

    private fun mockDataSource(result: Result<User>): LoginDataSource =
        mock {
            on { login(any(), any()) } doReturn result
        }

    @Test
    fun loginCacheUser() {
        val user = User(username, pincode, passport)
        val dataSource = mockDataSource(Result.Success(user))
        val repo = UserRepository(dataSource)

        assertThat(repo.user, nullValue())
        assertThat(repo.isLoggedIn, eq(false))

        val result = repo.login(username, pincode.toString())

        assertThat(result, eq(Result.Success(user)))
        assertThat(repo.user, eq(user))
        assertThat(repo.isLoggedIn, eq(true))

        verify(dataSource, times(1)).login(username, pincode)
    }

    @Test
    fun invalidPincodeFails() {
        val user = User(username, pincode, passport)
        val dataSource = mockDataSource(Result.Success(user))
        val repo = UserRepository(dataSource)

        val result = repo.login(username, invalidPincode)
        assertThat("The login result is not an error", result is Result.Error)
        verify(dataSource, never()).login(username, pincode)
    }

    @Test
    fun logoutUserChangesState() {
        val user = User(username, pincode, passport)
        val dataSource = mockDataSource(Result.Success(user))
        val repo = UserRepository(dataSource)
        repo.login(username, pincode.toString())
        repo.logout()
        assertThat(repo.user, nullValue())
        assertThat(repo.isLoggedIn, eq(false))
        verify(dataSource, times(1)).login(username, pincode)
    }

    @Test
    fun loginFailsWhenSourceFails() {
        val dataSource = mockDataSource(Result.Error(Exception()))
        val repo = UserRepository(dataSource)

        val result = repo.login(username, pincode.toString())

        assertThat("The login result is not an error", result is Result.Error)
        assertThat(repo.user, nullValue())
        assertThat(repo.isLoggedIn, eq(false))
        verify(dataSource, times(1)).login(username, pincode)
    }
}