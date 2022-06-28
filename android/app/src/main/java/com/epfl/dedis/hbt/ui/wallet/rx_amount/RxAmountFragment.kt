package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.databinding.FragmentWalletRxAmountBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class RxAmountFragment : Fragment() {

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletRxAmountBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletRxAmountBinding.inflate(inflater, container, false).apply {
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val receiveButton = binding.walletButtonReceive
        val amount = binding.walletRxAmount

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

        receiveButton.setOnClickListener {
            walletViewModel.send(amount.text.toString().toFloat())
        }
    }
}
