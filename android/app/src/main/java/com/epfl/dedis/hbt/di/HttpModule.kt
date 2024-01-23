package com.epfl.dedis.hbt.di

import com.epfl.dedis.hbt.service.document.DocumentEndpoint
import com.epfl.dedis.hbt.service.http.ResultCallAdapterFactory
import com.fasterxml.jackson.databind.ObjectMapper
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import retrofit2.Retrofit
import retrofit2.converter.jackson.JacksonConverterFactory
import javax.inject.Qualifier
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object HttpModule {

    @BaseURL
    @Provides
    @Singleton
    fun provideBaseURL() = "http://10.0.2.2:3000"

    @Provides
    @Singleton
    fun provideRetrofit(@BaseURL baseUrl: String, mapper: ObjectMapper): Retrofit =
        Retrofit.Builder()
            .baseUrl(baseUrl)
            .addConverterFactory(JacksonConverterFactory.create(mapper))
            .addCallAdapterFactory(ResultCallAdapterFactory())
            .build()

    @Provides
    @Singleton
    fun provideDocumentService(retrofit: Retrofit): DocumentEndpoint =
        retrofit.create(DocumentEndpoint::class.java)

    @Qualifier
    @Retention(AnnotationRetention.BINARY)
    annotation class BaseURL
}
