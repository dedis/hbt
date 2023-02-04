package com.epfl.dedis.hbt.ui.login

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
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.databinding.FragmentLoginBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.register.ScanPassportFragment
import com.epfl.dedis.hbt.ui.wallet.WalletFragment
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class LoginFragment : Fragment() {

    companion object {
        fun newInstance() = LoginFragment()
    }

    private val loginViewModel: LoginViewModel by viewModels()
    private var _binding: FragmentLoginBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentLoginBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val usernameEditText = binding.loginUsername
        val pincodeEditText = binding.loginPincode
        val loginButton = binding.loginSignin
        val registerButton = binding.loginRegister

        loginViewModel.loginFormState.observe(
            viewLifecycleOwner,
            Observer { loginFormState ->
                if (loginFormState == null) {
                    return@Observer
                }

                loginButton.isEnabled =
                    loginFormState.isDataValid && loginFormState.isUserRegistered

                loginFormState.usernameError?.let {
                    usernameEditText.error = getString(it)
                }
                loginFormState.pincodeError?.let {
                    pincodeEditText.error = getString(it)
                }
            })

        loginViewModel.loginResult.observe(viewLifecycleOwner,
            Observer { loginResult ->
                loginResult ?: return@Observer
                loginResult.error?.let {
                    showLoginFailed(it)
                }
                loginResult.success?.let {
                    updateUiWithUser(it)
                    // Open Wallet fragment
                    MainActivity.setCurrentFragment(parentFragmentManager, WalletFragment())
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
                loginViewModel.loginDataChanged(
                    usernameEditText.text.toString(),
                    pincodeEditText.text.toString()
                )
            }
        }
        usernameEditText.addTextChangedListener(afterTextChangedListener)
        pincodeEditText.addTextChangedListener(afterTextChangedListener)
        pincodeEditText.setOnEditorActionListener { _, actionId, _ ->
            if (actionId == EditorInfo.IME_ACTION_DONE) {
                loginViewModel.login(
                    usernameEditText.text.toString(),
                    pincodeEditText.text.toString()
                )
            }
            false
        }

        loginButton.setOnClickListener {
            loginViewModel.login(
                usernameEditText.text.toString(),
                pincodeEditText.text.toString()
            )
        }

        registerButton.setOnClickListener {
            MainActivity.setCurrentFragment(parentFragmentManager, ScanPassportFragment())
        }
    }

    private fun updateUiWithUser(model: LoggedInUserView) {
        val welcome = getString(R.string.welcome) + model.displayName
        val appContext = context?.applicationContext ?: return
        Toast.makeText(appContext, welcome, Toast.LENGTH_LONG).show()
    }

    private fun showLoginFailed(@StringRes errorString: Int) {
        val appContext = context?.applicationContext ?: return
        Toast.makeText(appContext, errorString, Toast.LENGTH_LONG).show()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}