package com.epfl.dedis.hbt.service.passport.ncf

import android.util.Log
import com.epfl.dedis.hbt.service.passport.mrz.BACData
import net.sf.scuba.smartcards.CardServiceException
import org.jmrtd.BACKey
import org.jmrtd.PACEKeySpec
import org.jmrtd.PassportService
import org.jmrtd.lds.*
import org.jmrtd.lds.icao.COMFile
import org.jmrtd.lds.icao.DG11File
import org.jmrtd.lds.icao.DG1File
import java.io.IOException
import java.security.GeneralSecurityException


/**
 * Creates a document by reading it from a service.
 *
 * Strongly inspired by https://github.com/jllarraz/AndroidPassportReader
 *
 * @param service the service to read from
 * @param bacData the BAC entries
 *
 * @throws GeneralSecurityException if certain security primitives are not supported
 */
class PassportNFC @Throws(GeneralSecurityException::class)
constructor(service: PassportService, bacData: BACData) {

    var sodFile: SODFile? = null
        private set
    var dg1File: DG1File? = null
        private set
    var dg11File: DG11File? = null
        private set

    init {
        val hasSAC: Boolean
        var isSACSucceeded = false
        try {
            service.open()
            /* Find out whether this MRTD supports SAC. */
            hasSAC =
                try {
                    Log.i(TAG, "Inspecting card access file")
                    val cardAccessFile =
                        CardAccessFile(
                            service.getInputStream(
                                PassportService.EF_CARD_ACCESS,
                                PassportService.DEFAULT_MAX_BLOCKSIZE
                            )
                        )
                    // Supports SAC if any of the security info is a PACEInfo
                    cardAccessFile.securityInfos.any { it is PACEInfo }
                } catch (e: Exception) {
                    /* NOTE: No card access file, continue to test for BAC. */
                    Log.i(TAG, "DEBUG: failed to get card access file: " + e.message)
                    e.printStackTrace()
                    false
                }

            if (hasSAC) {
                isSACSucceeded = try {
                    service.doPACE(bacData)
                    true
                } catch (e: Exception) {
                    Log.i(TAG, "PACE failed, falling back to BAC", e)
                    false
                }

            }
            service.sendSelectApplet(isSACSucceeded)
        } catch (cse: CardServiceException) {
            throw cse
        } catch (e: Exception) {
            e.printStackTrace()
            throw CardServiceException("Cannot open document. " + e.message)
        }

        /* Find out whether this MRTD supports BAC. */
        val hasBAC = try {
            /* Attempt to read EF.COM before BAC. */
            COMFile(
                service.getInputStream(
                    PassportService.EF_COM,
                    PassportService.DEFAULT_MAX_BLOCKSIZE
                )
            )
            false
        } catch (e: Exception) {
            Log.i(TAG, "Attempt to read EF.COM before BAC failed with: " + e.message)
            true
        }

        /* If we have to do BAC, try to do BAC. */
        if (hasBAC && !(hasSAC && isSACSucceeded)) {
            val bacKey = bacData.bacKey
            val triedBACEntries = ArrayList<BACKey>()
            triedBACEntries.add(bacKey)
            try {
                service.doBAC(bacKey)
            } catch (e: Exception) {
                Log.i(TAG, "Failed to do BAC", e)
            }
        }

        try {
            sodFile = service.getSodFile()
            dg1File = service.getDG1File()
            dg11File = service.getDG11File()
        } catch (ioe: IOException) {
            ioe.printStackTrace()
            Log.w(TAG, "Could not read file")
        }
    }

    @Throws(IOException::class, CardServiceException::class, GeneralSecurityException::class)
    private fun PassportService.doPACE(bacData: BACData) {
        val paceKeySpec = PACEKeySpec.createMRZKey(bacData.bacKey)

        getInputStream(PassportService.EF_CARD_ACCESS, PassportService.DEFAULT_MAX_BLOCKSIZE).use {
            val cardAccessFile = CardAccessFile(it)
            val securityInfos = cardAccessFile.securityInfos
            val securityInfo = securityInfos.iterator().next()
            val paceInfos = ArrayList<PACEInfo>()
            if (securityInfo is PACEInfo) {
                paceInfos.add(securityInfo)
            }

            if (paceInfos.size > 0) {
                val paceInfo = paceInfos.iterator().next()
                doPACE(
                    paceKeySpec,
                    paceInfo.objectIdentifier,
                    PACEInfo.toParameterSpec(paceInfo.parameterId),
                    null
                )
            }
        }
    }

    @Throws(CardServiceException::class, IOException::class)
    private fun PassportService.getSodFile(): SODFile =
        getFile(PassportService.EF_SOD)

    @Throws(CardServiceException::class, IOException::class)
    private fun PassportService.getDG1File(): DG1File =
        getFile(PassportService.EF_DG1)

    @Throws(CardServiceException::class, IOException::class)
    private fun PassportService.getDG11File(): DG11File =
        getFile(PassportService.EF_DG11)

    @Throws(CardServiceException::class, IOException::class)
    private fun <T : AbstractTaggedLDSFile> PassportService.getFile(id: Short) =
        getInputStream(id, PassportService.DEFAULT_MAX_BLOCKSIZE).use {
            @Suppress("INACCESSIBLE_TYPE", "UNCHECKED_CAST")
            LDSFileUtil.getLDSFile(id, it) as T
        }

    companion object {
        private val TAG = PassportNFC::class.java.simpleName
    }
}