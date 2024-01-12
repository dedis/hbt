package com.epfl.dedis.hbt.service.passport

import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.service.passport.mrz.MRZInfo
import org.jmrtd.lds.SODFile
import org.jmrtd.lds.icao.DG11File

data class Passport(
    val mrzInfo: MRZInfo,
    val sodFile: SODFile,
    val portrait: Portrait,
    val dg11File: DG11File?
)