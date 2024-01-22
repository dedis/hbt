package com.epfl.dedis.hbt.service.document

import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.document.Document
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.data.user.User
import okhttp3.MediaType
import okhttp3.RequestBody
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class DocumentService @Inject constructor(private val endpoint: DocumentEndpoint) {

    suspend fun create(user: User, portrait: Portrait, registered: Boolean): Result<Document> =
        endpoint.create(
            user.name,
            user.passport,
            user.role.ordinal,
            RequestBody.create(MediaType.parse(portrait.type), portrait.data),
            registered
        )
}