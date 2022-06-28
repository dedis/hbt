package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.databinding.FragmentWalletPincodeBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class PincodeFragment : Fragment() {

    companion object {
        private const val USERNAME = "USERNAME"

        fun newInstance(username: String?) = ScanFragment().apply {
            val bundle = Bundle()
            bundle.putString(USERNAME, username)
            arguments = bundle
        }
    }

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletPincodeBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletPincodeBinding.inflate(inflater, container, false).apply {
            // Set the username field to the value given as argument (if present)
            arguments?.getString(USERNAME)?.let {
                walletName.text = it
                walletBalance.text = "0 HBT"
            }
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletPincodeOkButton
        val amount = binding.walletPincodeAmount

        walletViewModel.walletFormState.observe(
            viewLifecycleOwner,
            Observer { walletFormState ->
                if (walletFormState == null) {
                    return@Observer
                }
            })

        walletViewModel.walletResult.observe(viewLifecycleOwner,
            Observer { walletResult ->
                walletResult ?: return@Observer
                if (walletResult.error != null) {
//                    onRegisterFailed(walletResult.error)
                } else {
                    //                  onRegisterSuccess(usernameEditText.text.toString())
                }
            })

        okButton.setOnClickListener {
            walletViewModel.send(amount.text.toString().toFloat())
        }
    }
}
