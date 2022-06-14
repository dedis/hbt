package com.epfl.dedis.hbt.data

import com.epfl.dedis.hbt.data.model.User
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Class that handles authentication w/ login credentials and retrieves user information.
 */
@Singleton
class LoginDataSource @Inject constructor() {

    private val users: MutableMap<String, User> = mutableMapOf()

    private fun register(username: String, pincode: Int): Result<User> {
        if (users.containsKey(username)) return Result.Error(Exception("Already registered"))

        val user = User(username, pincode)
        users[username] = user

        return Result.Success(user)
    }

    fun login(username: String, pincode: Int): Result<User> {
        val user = users[username]

        // TODO return Result.Error(Exception("Not registered"))
        return if (user == null) register(username, pincode)
        else Result.Success(user)
    }

    fun logout() {
        // TODO: revoke authentication
    }
}