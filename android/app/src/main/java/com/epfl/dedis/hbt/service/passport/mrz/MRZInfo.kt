package com.epfl.dedis.hbt.service.passport.mrz

import org.jmrtd.lds.icao.MRZInfo as MRZ

data class MRZInfo(
    val bacData: BACData,
    val country: String,
    val surname: String,
    val name: String
) : BACData by bacData {
    constructor(mrz: MRZ) : this(
        mrz.documentNumber,
        mrz.dateOfBirth,
        mrz.dateOfExpiry,
        mrz.issuingState,
        mrz.primaryIdentifier,
        mrz.secondaryIdentifier
    )

    constructor(
        passNumber: String,
        dateOfBirth: String,
        expiration: String,
        country: String,
        surname: String,
        name: String
    ) : this(
        BACData.create(
            passNumber,
            dateOfBirth,
            expiration
        ),
        country,
        surname,
        name
    )
}
