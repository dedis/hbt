package com.epfl.dedis.hbt.service.passport

import com.epfl.dedis.hbt.service.passport.mrz.MRZInfo
import org.jmrtd.lds.SODFile
import org.jmrtd.lds.icao.DG11File

data class Passport(
    val mrzInfo: MRZInfo,
    val sodFile: SODFile,
    val dg11File: DG11File?
)