package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import com.epfl.dedis.hbt.R

class WalletDefaultFragment : Fragment() {

    companion object {
        fun newInstance() = WalletDefaultFragment()
    }

    private lateinit var viewModel: WalletDefaultViewModel

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_wallet_default, container, false)
    }

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)
        viewModel = ViewModelProvider(this).get(WalletDefaultViewModel::class.java)
        // TODO: Use the ViewModel
    }

}