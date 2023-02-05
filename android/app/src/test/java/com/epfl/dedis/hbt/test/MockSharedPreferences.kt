package com.epfl.dedis.hbt.test

import android.content.SharedPreferences
import android.content.SharedPreferences.OnSharedPreferenceChangeListener
import android.os.Bundle

/**
 * A straight-forward implementation of the SharedPreferences used to test services where it is used
 *
 * It does not support change listener as it is not used in the application
 */
class MockSharedPreferences : SharedPreferences {

    private val store: Bundle = Bundle()

    fun reset() = store.clear()

    @Suppress("DEPRECATION") // We need to use the untyped API to get Any value
    override fun getAll(): MutableMap<String, *> =
        store.keySet().associateWith { store.get(it) }.toMutableMap()

    override fun getString(key: String, default: String?): String? =
        store.getString(key, default)

    override fun getStringSet(key: String, default: MutableSet<String>?): MutableSet<String>? =
        store.getStringArray(key)?.toMutableSet() ?: default

    override fun getInt(key: String, default: Int): Int =
        store.getInt(key, default)

    override fun getLong(key: String, default: Long): Long =
        store.getLong(key, default)

    override fun getFloat(key: String, default: Float): Float =
        store.getFloat(key, default)

    override fun getBoolean(key: String, default: Boolean): Boolean =
        store.getBoolean(key, default)

    override fun contains(key: String): Boolean =
        store.containsKey(key)

    override fun edit(): SharedPreferences.Editor =
        object : SharedPreferences.Editor {
            override fun putString(key: String, value: String?): SharedPreferences.Editor =
                this.apply {
                    store.putString(key, value)
                }

            override fun putStringSet(
                key: String,
                values: MutableSet<String>?
            ): SharedPreferences.Editor =
                this.apply {
                    store.putStringArray(key, values?.toTypedArray())
                }

            override fun putInt(key: String, value: Int): SharedPreferences.Editor =
                this.apply {
                    store.putInt(key, value)
                }

            override fun putLong(key: String, value: Long): SharedPreferences.Editor =
                this.apply {
                    store.putLong(key, value)
                }

            override fun putFloat(key: String, value: Float): SharedPreferences.Editor =
                this.apply {
                    store.putFloat(key, value)
                }

            override fun putBoolean(key: String, value: Boolean): SharedPreferences.Editor =
                this.apply {
                    store.putBoolean(key, value)
                }

            override fun remove(key: String): SharedPreferences.Editor =
                this.apply {
                    store.remove(key)
                }

            override fun clear(): SharedPreferences.Editor =
                this.apply {
                    store.clear()
                }

            override fun commit(): Boolean = true

            override fun apply() {}
        }

    override fun registerOnSharedPreferenceChangeListener(listener: OnSharedPreferenceChangeListener?) =
        throw UnsupportedOperationException()

    override fun unregisterOnSharedPreferenceChangeListener(listener: OnSharedPreferenceChangeListener?) =
        throw UnsupportedOperationException()
}
