package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.EditorInfo
import android.widget.Toast
import androidx.annotation.StringRes
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.databinding.FragmentRegisterBinding
import com.epfl.dedis.hbt.utility.NfcReader
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class RegisterFragment : Fragment() {

    companion object {
        private const val USERNAME = "USERNAME"
        private const val PINCODE = "PINCODE"

        fun newInstance(username: String?, pincode: String?) = RegisterFragment().apply {
            val bundle = Bundle()
            bundle.putString(USERNAME, username)
            bundle.putString(PINCODE, pincode)
            arguments = bundle
        }
    }

    private var nfcReader: NfcReader? = null
    private val registerViewModel: RegisterViewModel by viewModels()
    private var _binding: FragmentRegisterBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentRegisterBinding.inflate(inflater, container, false).apply {
            // Set the username field to the value given as argument (if present)
            arguments?.getString(USERNAME)?.let {
                username.setText(it)
            }
            // Same for pincode
            arguments?.getString(PINCODE)?.let {
                pincode.setText(it)
            }
        }

        nfcReader = NfcReader(this.activity).also {
            //TODO: it.start()
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val usernameEditText = binding.username
        val pincodeEditText = binding.pincode
        val passportEditText = binding.passport
        val registerButton = binding.register

        registerViewModel.registerFormState.observe(viewLifecycleOwner,
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

                registerFormState.passportError?.let {
                    passportEditText.error = getString(it)
                }

                registerFormState.isDataValid?.let {
                    registerButton.isEnabled = true
                }

            })

        registerViewModel.registerResult.observe(viewLifecycleOwner,
            Observer { registerResult ->
                registerResult ?: return@Observer
                registerResult.error?.let {
                    showRegisterFailed(it)
                }
                registerResult.success?.let {
                    showRegisterSuccess()
                    // TODO: display Wallet tab and fragment
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
                    pincodeEditText.text.toString(),
                    passportEditText.text.toString()
                )
            }
        }
        usernameEditText.addTextChangedListener(afterTextChangedListener)
        pincodeEditText.addTextChangedListener(afterTextChangedListener)
        pincodeEditText.setOnEditorActionListener { _, actionId, _ ->
            if (actionId == EditorInfo.IME_ACTION_DONE) {
                registerViewModel.register(
                    usernameEditText.text.toString(),
                    pincodeEditText.text.toString(),
                    passportEditText.text.toString()
                )
            }
            false
        }

        registerButton.setOnClickListener {
            //TODO: switch to login or to wallet fragment
/*
            val rf = RegisterFragment.newInstance(
                usernameEditText.text.toString(),
                pincodeEditText.text.toString()
            )

            MainActivity.setCurrentFragment(activity?.supportFragmentManager ?: parentFragmentManager, rf )
*/
        }
    }

    private fun showRegisterSuccess() {
        val appContext = context?.applicationContext ?: return
        Toast.makeText(appContext, "Registered !", Toast.LENGTH_LONG).show()
    }

    private fun showRegisterFailed(@StringRes errorString: Int) {
        val appContext = context?.applicationContext ?: return
        Toast.makeText(appContext, errorString, Toast.LENGTH_LONG).show()
    }

    override fun onDestroyView() {
        //TODO: nfcReader?.stop()
        super.onDestroyView()
        _binding = null
    }
}
