package com.epfl.dedis.hbt.di

import android.content.Context
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import org.chromium.net.CronetEngine
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
class HttpModule {
    @Provides
    @Singleton
    fun provideCronet(@ApplicationContext context: Context): CronetEngine =
        CronetEngine.Builder(context).build()
}
