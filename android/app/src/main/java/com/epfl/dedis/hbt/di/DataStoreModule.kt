package com.epfl.dedis.hbt.di

import android.content.Context
import android.content.SharedPreferences
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object DataStoreModule {

    /**
     * If for some reason another shared preference file was to be created, one should use
     * https://developer.android.com/training/dependency-injection/hilt-android?hl=fr#multiple-bindings
     */
    @Provides
    @Singleton
    fun provideUsersPreferences(@ApplicationContext context: Context): SharedPreferences =
        context.getSharedPreferences("users", Context.MODE_PRIVATE)
}