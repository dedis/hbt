package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.coroutineScope
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.document.Portrait
import com.epfl.dedis.hbt.databinding.FragmentPassportNfcBinding
import com.epfl.dedis.hbt.service.passport.Passport
import com.epfl.dedis.hbt.service.passport.mrz.BACData
import com.epfl.dedis.hbt.service.passport.ncf.NFCReader
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.MainActivity.Companion.getSafeSerializable
import com.epfl.dedis.hbt.ui.NFCViewModel
import kotlinx.coroutines.launch
import java.io.FileNotFoundException

private const val USE_PERSONAL_DATA = true

private const val BAC_DATA = "bac_data"

/**
 * A simple [Fragment] subclass.
 * Use the [PassportNfcFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
class PassportNfcFragment : Fragment() {

    private val nfcViewModel: NFCViewModel by viewModels(ownerProducer = { requireActivity() })

    private lateinit var bacData: BACData

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        arguments?.let {
            bacData = it.getSafeSerializable(BAC_DATA)!!
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        nfcViewModel.setCallback(lifecycle) { intent ->
            lifecycle.coroutineScope.launch {
                when (val result = NFCReader.readPassport(intent, bacData)) {
                    is Result.Success -> {
                        val passport = result.data
                        val personalData = extractPersonalData(passport)

                        if (personalData == null) {
                            Toast.makeText(
                                requireContext(),
                                getString(R.string.error_invalid_passport),
                                Toast.LENGTH_LONG
                            ).show()
                        } else {
                            MainActivity.setCurrentFragment(
                                parentFragmentManager,
                                RegisterFragment.newInstance(
                                    passport.mrzInfo.number,
                                    personalData,
                                    passport.portrait
                                )
                            )
                        }

                    }
                    is Result.Error -> {
                        Log.e("NFC-PASSPORT", "Error in communication", result.exception)
                        Toast.makeText(
                            requireActivity(),
                            result.exception.message,
                            Toast.LENGTH_LONG
                        ).show()
                    }
                }
            }
        }

        return FragmentPassportNfcBinding.inflate(inflater, container, false).apply {
            skipNfc.setOnClickListener {
                MainActivity.setCurrentFragment(
                    parentFragmentManager,
                    RegisterFragment.newInstance(
                        "10AZ000001",
                        "some checksum",
                        getMockPortrait()
                    )
                )
            }

        }.root
    }

    fun getMockPortrait(): Portrait {
        val stream = PassportNfcFragment::class.java.getResourceAsStream("/utopia.png")
            ?: throw FileNotFoundException()

        stream.use {
            return Portrait("image/png", it.readBytes())
        }
    }

    private fun extractPersonalData(passport: Passport): String? {
        return if (USE_PERSONAL_DATA) {
            passport.dg11File?.personalNumber
        } else {
            // generate a dummy personal number based on the hash of the name and surname
            (passport.mrzInfo.name.hashCode() + passport.mrzInfo.surname.hashCode())
                .toString()
                .padStart(14, 'X')
                .substring(0 until 14)
        }
    }

    companion object {

        /**
         * Use this factory method to create a new instance of
         * this fragment using the provided parameters.
         *
         * @param bacData items needed to compute the BAC key
         * @return A new instance of fragment NCFPassportFragment.
         */
        @JvmStatic
        fun newInstance(bacData: BACData) =
            PassportNfcFragment().apply {
                arguments = Bundle().apply {
                    putSerializable(BAC_DATA, bacData)
                }
            }
    }
}