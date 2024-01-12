package com.epfl.dedis.hbt.data.document

import java.io.Serializable

data class Portrait(
    val type: String,
    val data: ByteArray
) : Serializable {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as Portrait

        if (type != other.type) return false
        return data.contentEquals(other.data)
    }

    override fun hashCode(): Int {
        var result = type.hashCode()
        result = 31 * result + data.contentHashCode()
        return result
    }
}
