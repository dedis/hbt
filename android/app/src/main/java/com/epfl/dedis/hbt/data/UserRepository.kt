package com.epfl.dedis.hbt.data

import com.epfl.dedis.hbt.data.model.User
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Class that requests authentication and user information from the remote data source and
 * maintains an in-memory cache of login status and user credentials information.
 */

@Singleton
class UserRepository @Inject constructor(private val dataSource: LoginDataSource) {

    // in-memory cache of the loggedInUser object
    var user: User? = null
        private set

    val isLoggedIn: Boolean
        get() = user != null

    init {
        // If user credentials will be cached in local storage, it is recommended it be encrypted
        // @see https://developer.android.com/training/articles/keystore
        user = null
    }

    fun isRegistered(username: String): Boolean {
        //TODO implement actual check whether user is registered
        return false
    }

    fun logout() {
        user = null
        dataSource.logout()
    }

    fun login(username: String, pincode: String): Result<User> {
        // handle login
        val pin = pincode.toIntOrNull() ?: return Result.Error(NumberFormatException())
        val result = dataSource.login(username, pin)

        if (result is Result.Success) {
            setLoggedInUser(result.data)
        }

        return result
    }

    private fun setLoggedInUser(loggedInUser: User) {
        this.user = loggedInUser
        // If user credentials will be cached in local storage, it is recommended it be encrypted
        // @see https://developer.android.com/training/articles/keystore
    }
}