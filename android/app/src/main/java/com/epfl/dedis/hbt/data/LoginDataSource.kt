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


    fun isRegistered(username: String): Boolean {
        return users.containsKey(username)
    }

    private fun register(username: String, pincode: Int): Result<User> {
        if (isRegistered(username)) return Result.Error(Exception("Already registered"))

        val user = User(username, pincode)
        users[username] = user

        return Result.Success(user)
    }

    fun login(username: String, pincode: Int): Result<User> {
        val user = users[username] ?: return Result.Error(Exception("The user does not exist"))

        return if (user.pincode == pincode)
            Result.Success(user)
        else
            Result.Error(Exception("Invalid pincode"))
    }
}