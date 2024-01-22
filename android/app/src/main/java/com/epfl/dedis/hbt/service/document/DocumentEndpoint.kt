package com.epfl.dedis.hbt.service.document

import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.document.Document
import okhttp3.RequestBody
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.Part

interface DocumentEndpoint {

    @Multipart
    @POST("document")
    suspend fun create(
        @Part("name") name: String,
        @Part("passport") passport: String,
        @Part("role") role: Int,
        @Part("image") image: RequestBody,
        @Part("registered") registered: Boolean
    ): Result<Document>
}