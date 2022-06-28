package com.epfl.dedis.hbt.data

import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.data.model.User
import com.epfl.dedis.hbt.data.model.Wallet
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Class that handles authentication w/ login credentials and retrieves user information.
 */
@Singleton
class UserDataSource @Inject constructor() {

    private val users: MutableMap<String, User> = mutableMapOf()
    private val wallets: MutableMap<User, Wallet> = mutableMapOf()

    fun isRegistered(username: String): Boolean {
        return users.containsKey(username)
    }

    fun register(username: String, pincode: Int, passport: String, role: Role): Result<User> {
        if (isRegistered(username)) return Result.Error(Exception("Already registered"))

        // create user
        val user = User(username, pincode, passport, role)
        users[username] = user

        //create wallet
        val wallet = Wallet()
        wallets[user] = wallet

        return Result.Success(user)
    }

    fun login(username: String, pincode: Int): Result<User> {
        val user = users[username] ?: return Result.Error(Exception("The user does not exist"))

        return if (user.pincode == pincode)
            Result.Success(user)
        else
            Result.Error(Exception("Invalid pincode"))
    }

    fun getWallet(user: User): Result<Wallet> {
        val wallet = wallets[user] ?: return Result.Error(Exception("The wallet does not exist"))

        return Result.Success(wallet)
    }
}