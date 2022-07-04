package com.epfl.dedis.hbt.ui.wallet

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentWalletScanBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class ScanFragment : Fragment() {

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletScanBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletScanBinding.inflate(inflater, container, false).apply {
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
    }
}
