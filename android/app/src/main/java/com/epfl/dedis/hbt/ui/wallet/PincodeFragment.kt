package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentWalletPincodeBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class PincodeFragment : Fragment() {

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
            walletName.text = walletViewModel.user?.name.toString()
            when (walletViewModel.user?.role) {
                Role.BENEFICIARY -> walletRole.text = getString(R.string.role_beneficiary)
                Role.MERCHANT -> walletRole.text = getString(R.string.role_merchant)
                else -> walletRole.text = getString(R.string.role_beneficiary)
            }
            walletBalance.text =
                walletViewModel.wallet?.balance.toString() + getString(R.string.hbt_currency)
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletPincodeOkButton
        val amount = binding.walletPincodeAmount

        okButton.setOnClickListener {
            walletViewModel.send(amount.text.toString().toFloat())
        }
    }
}
