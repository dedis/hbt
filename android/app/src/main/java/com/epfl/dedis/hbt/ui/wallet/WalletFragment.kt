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
import com.epfl.dedis.hbt.ui.login.LoginFragment
import com.epfl.dedis.hbt.ui.wallet.TransactionState.*
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
            val role = walletViewModel.user?.role ?: Role.BENEFICIARY
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)

            walletButtonLogout.setOnClickListener {
                walletViewModel.logout()
                MainActivity.setCurrentFragment(parentFragmentManager, LoginFragment.newInstance())
            }
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val sendButton = binding.walletButtonSend
        val receiveButton = binding.walletButtonReceive
        val logoutButton = binding.walletButtonLogout

        receiveButton.setOnClickListener {
            MainActivity.setCurrentFragment(
                parentFragmentManager,
                RxAmountFragment()
            )
        }

        sendButton.setOnClickListener {
            walletViewModel.transitionTo(SenderRead)
        }

        logoutButton.setOnClickListener {
            walletViewModel.logout()
            MainActivity.setCurrentFragment(parentFragmentManager, LoginFragment.newInstance())
        }

        walletViewModel.transactionState.observe(viewLifecycleOwner) {
            when (it) {
                is ReceiverRead, is SenderRead ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ScanFragment()
                    )
                is ReceiverShow, is SenderShow ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ShowQrFragment()
                    )
                else -> {}
            }
        }
    }
}
