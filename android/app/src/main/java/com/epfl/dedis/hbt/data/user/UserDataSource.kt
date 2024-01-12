package com.epfl.dedis.hbt.data.user

import android.content.SharedPreferences
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.service.document.DocumentService
import com.epfl.dedis.hbt.service.json.JsonService
import com.epfl.dedis.hbt.service.json.JsonType.USER_DATA
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Class that handles authentication w/ login credentials and retrieves user information.
 */
@Singleton
class UserDataSource @Inject constructor(
    private val sharedPref: SharedPreferences,
    private val jsonService: JsonService,
    private val documentService: DocumentService
) {

    private val usernamesKey: String = "users"

    private val users: MutableMap<String, User> = mutableMapOf()
    private val wallets: MutableMap<User, Wallet> = mutableMapOf()

    init {
        // Load data from preferences
        val usernames = sharedPref.getStringSet(usernamesKey, setOf())!!
        usernames.forEach {
            val user = getUserData(it, sharedPref, jsonService)
            users[it] = user
            // TODO : Wallet is currently not stored
            wallets[user] = Wallet.newInstance()
        }
    }

    private fun getUserData(
        name: String,
        sharedPref: SharedPreferences,
        jsonService: JsonService
    ): User =
        jsonService.fromJson(
            USER_DATA,
            sharedPref.getString(name, null)
                ?: throw IllegalStateException("The user $name is present is the user list but has no data")
        )

    fun isRegistered(username: String): Boolean {
        return username == usernamesKey || users.containsKey(username)
    }

    fun register(
        username: String,
        pincode: Int,
        passport: String,
        role: Role,
        portrait: Portrait
    ): Result<User> {
        if (isRegistered(username)) return Result.Error(Exception("Already registered"))

        // create user
        val user = User(username, pincode, passport, role)
        val call = documentService.create(user, portrait, false)
        val response = call.execute()
        if (response.errorBody() != null) {
            return Result.Error(Exception("Failed to register : " + response.message()))
        }

        users[username] = user

        //create wallet
        wallets[user] = Wallet.newInstance()

        createUserInStore(user)

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

    private fun createUserInStore(user: User) {
        with(sharedPref.edit()) {
            putString(user.name, jsonService.toJson(USER_DATA, user))
            putStringSet(usernamesKey, users.keys)
            apply()
        }
    }
}