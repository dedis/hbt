package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.widget.addTextChangedListener
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.databinding.FragmentWalletRxAmountBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.wallet.TransactionState.ReceiverShow
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class RxAmountFragment : Fragment() {

    private val walletViewModel: WalletViewModel by viewModels(ownerProducer = { requireActivity() })
    private var _binding: FragmentWalletRxAmountBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletRxAmountBinding.inflate(inflater, container, false).apply {
            walletName.text = walletViewModel.user.name
            val role = walletViewModel.user.role
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val amount = binding.walletRxAmount
        val okButton = binding.rxAmountOk

        amount.addTextChangedListener {
            okButton.isEnabled = it?.toString()?.toFloatOrNull() != null
        }

        okButton.setOnClickListener {
            val datetime = System.currentTimeMillis()
            walletViewModel.transitionTo(ReceiverShow(amount, datetime))
        }

        walletViewModel.transactionState.observe(viewLifecycleOwner) {
            when (it) {
                is TransactionState.ReceiverRead, is TransactionState.SenderRead ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ScanFragment()
                    )
                is ReceiverShow, is TransactionState.SenderShow ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ShowQrFragment()
                    )
                else -> {}
            }
        }
    }
}
