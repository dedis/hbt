package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentWalletBinding
import com.epfl.dedis.hbt.ui.MainActivity
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class WalletFragment : Fragment() {

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletBinding.inflate(inflater, container, false).apply {
            walletName.text = walletViewModel.user?.name.toString()
            when (walletViewModel.user?.role) {
                Role.BENEFICIARY -> walletRole.text = getString(R.string.role_beneficiary)
                Role.MERCHANT -> walletRole.text = getString(R.string.role_merchant)
                else -> walletRole.text = getString(R.string.role_beneficiary)
            }
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val sendButton = binding.walletButtonSend
        val receiveButton = binding.walletButtonReceive

        receiveButton.setOnClickListener {
            MainActivity.setCurrentFragment(
                parentFragmentManager,
                RxAmountFragment()
            )
        }

        sendButton.setOnClickListener {
            MainActivity.setCurrentFragment(
                parentFragmentManager,
                ScanFragment()
            )
        }
    }
}
