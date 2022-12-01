package com.epfl.dedis.hbt.ui

import android.app.PendingIntent
import android.app.PendingIntent.FLAG_MUTABLE
import android.app.PendingIntent.FLAG_UPDATE_CURRENT
import android.content.Intent
import android.nfc.NfcAdapter
import android.os.Build
import android.os.Bundle
import android.os.Parcelable
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentManager
import androidx.lifecycle.Lifecycle
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.databinding.ActivityMainBinding
import com.epfl.dedis.hbt.ui.login.LoginFragment
import dagger.hilt.android.AndroidEntryPoint
import java.io.Serializable

@AndroidEntryPoint
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    private val nfcViewModel: NFCViewModel by viewModels()
    private var nfcAdapter: NfcAdapter? = null

    private val pendingIntent =
        PendingIntent.getActivity(
            this,
            0,
            Intent(this, this.javaClass).addFlags(Intent.FLAG_ACTIVITY_SINGLE_TOP),
            (if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) FLAG_MUTABLE else 0) or FLAG_UPDATE_CURRENT
        )

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        nfcAdapter = NfcAdapter.getDefaultAdapter(this)

        // Set default fragment (MainFragment)
        setCurrentFragment(supportFragmentManager, LoginFragment.newInstance(), false)

        nfcViewModel.listenToNFC.observe(this) {
            if (lifecycle.currentState != Lifecycle.State.RESUMED) return@observe

            if (it) {
                // TODO Show toast if NFC is not available
                nfcAdapter?.enableForegroundDispatch(
                    this, pendingIntent, null, null
                )
            } else {
                nfcAdapter?.disableForegroundDispatch(this)
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        if (NfcAdapter.ACTION_TAG_DISCOVERED == intent.action || NfcAdapter.ACTION_TECH_DISCOVERED == intent.action) {
            nfcViewModel.onNewIntent(intent)
        } else {
            super.onNewIntent(intent)
        }
    }

    companion object {

        /**
         * This function should be used to set the current fragment shown on the main activity.
         *
         * @param fm the fragment manager of the activity
         * @param fragment to launch
         * @param addToBackstack whether the transaction should be added to the backstack (and thus reversible)
         *                       default value : true
         */
        @JvmStatic
        fun setCurrentFragment(
            fm: FragmentManager,
            fragment: Fragment,
            addToBackstack: Boolean = true
        ) {
            val transaction = fm.beginTransaction()
                .replace(R.id.container, fragment)
            if (addToBackstack)
                transaction.addToBackStack(null)
            transaction.commit()
        }

        @Suppress("DEPRECATION")
        inline fun <reified T : Serializable> Bundle.getSafeSerializable(key: String): T? {
            return if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
                getSerializable(key, T::class.java)
            } else {
                getSerializable(key) as? T
            }
        }


        @Suppress("DEPRECATION")
        inline fun <reified T : Parcelable> Bundle.getSafeParcelable(key: String): T? {
            return if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
                getParcelable(key, T::class.java)
            } else {
                getParcelable(key) as? T
            }
        }
    }
}
