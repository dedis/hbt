package com.epfl.dedis.hbt.service.passport.ncf

import android.content.Intent
import android.nfc.NfcAdapter
import android.nfc.Tag
import android.nfc.tech.IsoDep
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.Result.Error
import com.epfl.dedis.hbt.data.Result.Success
import com.epfl.dedis.hbt.service.passport.Passport
import com.epfl.dedis.hbt.service.passport.mrz.MRZInfo
import com.epfl.dedis.hbt.ui.MainActivity.Companion.getSafeParcelable
import net.sf.scuba.smartcards.CardService
import org.jmrtd.PassportService
import java.security.Security

object NFCUtils {

    fun readPassport(intent: Intent, mrzInfo: MRZInfo): Result<Passport> {
        val tag = intent.extras?.getSafeParcelable<Tag>(NfcAdapter.EXTRA_TAG) ?: return Error(
            UnsupportedOperationException("The nfc tag is not present in the provided intent")
        )

        val nfc = IsoDep.get(tag).apply { timeout = 5 * 1000 } //5 seconds timeout
            ?: return Error(UnsupportedOperationException("ISODep could not be created"))

        val ps = PassportService(
            CardService.getInstance(nfc),
            PassportService.NORMAL_MAX_TRANCEIVE_LENGTH,
            PassportService.DEFAULT_MAX_BLOCKSIZE,
            false,
            true
        )

        return try {
            ps.open()

            val passportNFC = PassportNFC(ps, mrzInfo)
            Success(
                Passport(
                    MRZInfo(passportNFC.dg1File!!.mrzInfo),
                    passportNFC.sodFile!!,
                    passportNFC.dg11File
                )
            )
        } catch (e: Exception) {
            Error(e)
        }.also {
            ps.close()
        }
    }

    init {
        Security.insertProviderAt(org.spongycastle.jce.provider.BouncyCastleProvider(), 1)
    }
}