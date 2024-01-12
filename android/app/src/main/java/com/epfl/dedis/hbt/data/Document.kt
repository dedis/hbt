package com.epfl.dedis.hbt.data

import com.fasterxml.jackson.annotation.JsonProperty

data class Document(
    @JsonProperty("doc_id") val id: String
)
