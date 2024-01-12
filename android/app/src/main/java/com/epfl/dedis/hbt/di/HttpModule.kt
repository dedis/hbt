package com.epfl.dedis.hbt.di

import com.epfl.dedis.hbt.service.document.DocumentService
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
    fun provideBaseURL() = "http://localhost:80"

    @Provides
    @Singleton
    fun provideRetrofit(@BaseURL baseUrl: String, mapper: ObjectMapper): Retrofit =
        Retrofit.Builder()
            .baseUrl(baseUrl)
            .addConverterFactory(JacksonConverterFactory.create(mapper))
            .build()

    @Provides
    @Singleton
    fun provideDocumentService(retrofit: Retrofit): DocumentService =
        retrofit.create(DocumentService::class.java)

    @Qualifier
    @Retention(AnnotationRetention.BINARY)
    annotation class BaseURL
}
