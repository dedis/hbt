package com.epfl.dedis.hbt.utility

import android.content.Context
import android.nfc.NfcAdapter
import android.nfc.NfcAdapter.getDefaultAdapter
import android.widget.Toast

class NfcReader(activity: Context?) {
    private val ctx: Context? = activity
    private lateinit var nfcAdapter: NfcAdapter

    fun start() {
        // if NFC is not supported, toast it !
        try {
            nfcAdapter = getDefaultAdapter(ctx)
            if (nfcAdapter == null) {
                TODO("Throw proper exception")
            }
        } catch (e: Exception) {
            Toast.makeText(
                ctx,
                "Nfc is not supported on this device, received exception $e",
                Toast.LENGTH_SHORT
            ).show()
            stop()
        }

        // if NFC is not enabled on the device, toast it !
        if (!nfcAdapter?.isEnabled!!) {
            Toast.makeText(
                ctx,
                "NFC disabled on this device. Turn on to proceed",
                Toast.LENGTH_SHORT
            ).show()
        }
    }

    fun stop() {
        TODO("Not yet implemented")
    }

}