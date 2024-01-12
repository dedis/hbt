package com.epfl.dedis.hbt.service.document

import com.epfl.dedis.hbt.data.document.Document
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.data.user.User
import okhttp3.MediaType
import okhttp3.RequestBody
import retrofit2.Call
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.Part

interface DocumentService {

    @Multipart
    @POST("document")
    fun create(
        @Part("name") name: String,
        @Part("passport") passport: String,
        @Part("role") role: Int,
        @Part("image") image: RequestBody,
        @Part("registered") registered: Boolean
    ): Call<Document>

    fun create(user: User, portrait: Portrait, registered: Boolean): Call<Document> =
        create(
            user.name,
            user.passport,
            user.role.ordinal,
            RequestBody.create(MediaType.parse(portrait.type), portrait.data),
            registered
        )
}