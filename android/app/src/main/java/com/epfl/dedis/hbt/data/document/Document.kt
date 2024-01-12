package com.epfl.dedis.hbt.data.document

import com.fasterxml.jackson.annotation.JsonProperty

data class Document(
    @JsonProperty("doc_id") val id: String
)
