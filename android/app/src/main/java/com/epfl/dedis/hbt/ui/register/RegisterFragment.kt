package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.databinding.FragmentRegisterBinding
import com.epfl.dedis.hbt.utility.NfcReader
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class RegisterFragment : Fragment() {

    companion object {
        private const val USERNAME = "USERNAME"
        private const val PINCODE = "PINCODE"

        fun newInstance(username: String?, pincode: String?) = RegisterFragment().apply {
            val bundle = Bundle()
            bundle.putString(USERNAME, username)
            bundle.putString(PINCODE, pincode)
            arguments = bundle
        }
    }

    private val viewModel: RegisterViewModel by viewModels()
    private var nfcReader: NfcReader? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        nfcReader = NfcReader(this.activity).also {
            //TODO: it.start()
        }

        return FragmentRegisterBinding.inflate(inflater, container, false).apply {
            // Set the username field to the value given as argument (if present)
            arguments?.getString(USERNAME)?.let {
                username.setText(it)
            }
            // Same for pincode
            arguments?.getString(PINCODE)?.let {
                pincode.setText(it)
            }
        }.root
    }

    override fun onDestroyView() {
        //TODO: nfcReader?.stop()
        super.onDestroyView()
    }
}