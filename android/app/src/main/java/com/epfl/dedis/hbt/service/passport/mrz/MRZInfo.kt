package com.epfl.dedis.hbt.service.passport.mrz

import org.jmrtd.BACKey
import java.io.Serializable
import org.jmrtd.lds.icao.MRZInfo as MRZ

data class MRZInfo(
    val country: String,
    val surname: String,
    val name: String,
    val number: String,
    val dateOfBirth: String,
    val expiration: String
) : Serializable {
    constructor(mrz: MRZ) : this(
        mrz.issuingState,
        mrz.primaryIdentifier,
        mrz.secondaryIdentifier,
        mrz.documentNumber,
        mrz.dateOfBirth,
        mrz.dateOfExpiry
    )

    val bacKey: BACKey
        get() = BACKey(number, dateOfBirth, expiration)
}
