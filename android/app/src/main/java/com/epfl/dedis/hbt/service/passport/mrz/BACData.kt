package com.epfl.dedis.hbt.service.passport.mrz

import org.jmrtd.BACKey
import java.io.Serializable

interface BACData : Serializable {
    val number: String
    val dateOfBirth: String
    val expiration: String

    val bacKey: BACKey

    companion object {
        fun create(documentNumber: String, dateOfBirth: String, dateOfExpiry: String): BACData =
            BACDataImpl(documentNumber, dateOfBirth, dateOfExpiry)
    }

    data class BACDataImpl(
        override val number: String,
        override val dateOfBirth: String,
        override val expiration: String
    ) : BACData {
        override val bacKey: BACKey
            get() = BACKey(number, dateOfBirth, expiration)
    }
}
