package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.annotation.StringRes
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentRegisterBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.wallet.WalletFragment
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class RegisterFragment : Fragment() {

    companion object {

        private const val PASSPORT = "PASSPORT"
        private const val CHECKSUM = "CHECKSUM"

        fun newInstance(passport: String, checksum: ByteArray) = RegisterFragment().apply {
            val bundle = Bundle()
            bundle.putString(PASSPORT, passport)
            bundle.putByteArray(CHECKSUM, checksum)
            arguments = bundle
        }
    }

    private val registerViewModel: RegisterViewModel by viewModels(ownerProducer = { requireActivity() })
    private var _binding: FragmentRegisterBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    private lateinit var passport: String
    private lateinit var checksum: ByteArray

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentRegisterBinding.inflate(inflater, container, false).apply {
            // Set the passport number field to the value given as argument
            requireArguments().getString(PASSPORT)!!.let {
                passport = it
                passportNumber.text = it
            }

            requireArguments().getByteArray(CHECKSUM)!!.let {
                checksum = it
                passportChecksum.text =
                    it.joinToString(separator = "") { b -> "%02x".format(b) }.substring(0, 16)
            }
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val usernameEditText = binding.registerUsername
        val pincodeEditText = binding.registerPincode
        val registerButton = binding.registerRegister
        val roleButton = binding.radioGroup

        // Set values from view model to entries

        registerViewModel.registerFormState.observe(
            viewLifecycleOwner,
            Observer { registerFormState ->
                if (registerFormState == null) {
                    return@Observer
                }

                registerFormState.usernameError?.let {
                    usernameEditText.error = getString(it)
                }

                registerFormState.pincodeError?.let {
                    pincodeEditText.error = getString(it)
                }

                registerButton.isEnabled = registerFormState.isDataValid
                registerButton.isVisible = registerFormState.isDataValid
            })

        registerViewModel.registerResult.observe(viewLifecycleOwner,
            Observer { registerResult ->
                registerResult ?: return@Observer
                if (registerResult.error != null)
                    onRegisterFailed(registerResult.error)
                else {
                    onRegisterSuccess()
                }
            })

        val afterTextChangedListener = object : TextWatcher {
            override fun beforeTextChanged(s: CharSequence, start: Int, count: Int, after: Int) {
                // ignore
            }

            override fun onTextChanged(s: CharSequence, start: Int, before: Int, count: Int) {
                // ignore
            }

            override fun afterTextChanged(s: Editable) {
                registerViewModel.registerDataChanged(
                    usernameEditText.text.toString(),
                    pincodeEditText.text.toString()
                )
            }
        }
        usernameEditText.addTextChangedListener(afterTextChangedListener)
        pincodeEditText.addTextChangedListener(afterTextChangedListener)

        registerButton.setOnClickListener {
            val role = when (roleButton.checkedRadioButtonId) {
                R.id.radioButtonBeneficiary -> Role.BENEFICIARY
                R.id.radioButtonMerchant -> Role.MERCHANT
                else -> throw Error("Unhandled role type")
            }

            registerViewModel.register(
                usernameEditText.text.toString(),
                pincodeEditText.text.toString(),
                passport,
                checksum,
                role
            )
        }

        // Set the default result with the current texts
        registerViewModel.registerDataChanged(
            usernameEditText.text.toString(),
            pincodeEditText.text.toString()
        )
    }

    private fun onRegisterFailed(@StringRes errorString: Int) {
        val appContext = context?.applicationContext ?: return
        Toast.makeText(appContext, errorString, Toast.LENGTH_LONG).show()
    }

    private fun onRegisterSuccess() {
        MainActivity.setCurrentFragment(parentFragmentManager, WalletFragment())
    }

    override fun onDestroyView() {
        //TODO: nfcReader?.stop()
        super.onDestroyView()
        _binding = null
    }
}
