package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.utility.NfcReader

class RegisterFragment : Fragment() {

    companion object {
        fun newInstance() = RegisterFragment()
    }

    private lateinit var viewModel: RegisterViewModel
    private var nfcReader: NfcReader? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        NfcReader(this.activity).also { nfcReader = it }
        nfcReader?.start()

        return inflater.inflate(R.layout.fragment_register, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        viewModel = ViewModelProvider(this)[RegisterViewModel::class.java]

        // TODO: Use the ViewModel
    }

    override fun onDestroyView() {
        nfcReader?.stop()
        super.onDestroyView()
    }
}