package com.epfl.dedis.hbt.data

import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.data.model.User
import com.epfl.dedis.hbt.data.model.Wallet
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Class that requests authentication and user information from the remote data source and
 * maintains an in-memory cache of login status and user credentials information.
 */

@Singleton
class UserRepository @Inject constructor(private val dataSource: UserDataSource) {

    // in-memory cache of the loggedInUser object
    var loggedInUser: User? = null
        private set

    var wallet: Wallet? = null
        private set

    val isLoggedIn: Boolean
        get() = loggedInUser != null

    init {
        // If user credentials will be cached in local storage, it is recommended it be encrypted
        // @see https://developer.android.com/training/articles/keystore
        loggedInUser = null
    }

    fun isRegistered(username: String): Boolean {
        return dataSource.isRegistered(username)
    }

    fun logout() {
        loggedInUser = null
    }

    fun register(username: String, pincode: String, passport: String, role: Role): Result<User> {
        val pin = pincode.toIntOrNull() ?: return Result.Error(NumberFormatException())
        val result = dataSource.register(username, pin, passport, role)

        if (result is Result.Success) {
            setLoggedInUser(result.data)
        }

        return result
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
        // If user credentials will be cached in local storage, it is recommended it be encrypted
        // @see https://developer.android.com/training/articles/keystore
        this.loggedInUser = loggedInUser

        val result = dataSource.getWallet(loggedInUser)
        if (result is Result.Success) {
            this.wallet = result.data
        }
    }
}