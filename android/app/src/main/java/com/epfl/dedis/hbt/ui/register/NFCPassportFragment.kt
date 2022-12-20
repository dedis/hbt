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
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.databinding.FragmentNfcPassportBinding
import com.epfl.dedis.hbt.service.passport.mrz.BACData
import com.epfl.dedis.hbt.service.passport.ncf.NFCReader
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.MainActivity.Companion.getSafeSerializable
import com.epfl.dedis.hbt.ui.NFCViewModel
import kotlinx.coroutines.launch

private const val BAC_DATA = "bac_data"

/**
 * A simple [Fragment] subclass.
 * Use the [NFCPassportFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
class NFCPassportFragment : Fragment() {

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
                        MainActivity.setCurrentFragment(
                            parentFragmentManager,
                            RegisterFragment.newInstance(
                                passport.mrzInfo.number,
                                passport.dg11File!!.personalNumber
                            )
                        )
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

        return FragmentNfcPassportBinding.inflate(inflater, container, false).root
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
            NFCPassportFragment().apply {
                arguments = Bundle().apply {
                    putSerializable(BAC_DATA, bacData)
                }
            }
    }
}