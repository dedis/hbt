package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.widget.addTextChangedListener
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.transaction.PendingTransaction
import com.epfl.dedis.hbt.data.transaction.TransactionState
import com.epfl.dedis.hbt.data.transaction.TransactionState.ReceiverShow
import com.epfl.dedis.hbt.data.transaction.TransactionStateManager
import com.epfl.dedis.hbt.databinding.FragmentWalletRxAmountBinding
import com.epfl.dedis.hbt.ui.MainActivity
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class RxAmountFragment : Fragment() {

    @Inject
    lateinit var trxStateManager: TransactionStateManager

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
            trxStateManager.startReceivingTransaction(
                PendingTransaction(
                    walletViewModel.user.passport,
                    amount.text.toString().toFloat(),
                    datetime
                )
            )
        }

        trxStateManager.currentState.observe(viewLifecycleOwner) {
            when (it) {
                is TransactionState.ReceiverRead, is TransactionState.SenderRead ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ScanQrFragment()
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
