package com.epfl.dedis.hbt.service.passport.mrz

import org.jmrtd.BACKey
import java.io.Serializable
import org.jmrtd.lds.icao.MRZInfo as MRZ

data class MRZInfo(
    val number: String,
    val dateOfBirth: String,
    val expiration: String,
    val country: String? = null,
    val surname: String? = null,
    val name: String? = null
) : Serializable {
    constructor(mrz: MRZ) : this(
        mrz.documentNumber,
        mrz.dateOfBirth,
        mrz.dateOfExpiry,
        mrz.issuingState,
        mrz.primaryIdentifier,
        mrz.secondaryIdentifier
    )

    val bacKey: BACKey
        get() = BACKey(number, dateOfBirth, expiration)
}
