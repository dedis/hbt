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
import com.epfl.dedis.hbt.databinding.FragmentWalletShowqrBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class ShowQrFragment : Fragment() {

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletShowqrBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletShowqrBinding.inflate(inflater, container, false).apply {
            walletName.text = walletViewModel.user?.name.toString()
            val role = walletViewModel.user?.role ?: Role.BENEFICIARY
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletButtonOk

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
        }
    }
}
