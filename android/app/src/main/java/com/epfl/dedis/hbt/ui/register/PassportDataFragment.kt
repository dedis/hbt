package com.epfl.dedis.hbt.ui.register

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import com.epfl.dedis.hbt.databinding.FragmentPassportDataBinding
import com.epfl.dedis.hbt.service.passport.mrz.MRZInfo
import com.epfl.dedis.hbt.ui.MainActivity

/**
 * A simple [Fragment] subclass.
 */
class PassportDataFragment : Fragment() {

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        return FragmentPassportDataBinding.inflate(inflater, container, false).apply {
            automaticScan.setOnClickListener {
                MainActivity.setCurrentFragment(
                    parentFragmentManager,
                    ScanPassportFragment()
                )
            }

            validate.setOnClickListener {
                MainActivity.setCurrentFragment(
                    parentFragmentManager,
                    NFCPassportFragment.newInstance(
                        MRZInfo(
                            passportNumber.text.toString(),
                            editBirthDate.text.toString(),
                            editExpirationDate.text.toString()
                        )
                    )
                )
            }
        }.root
    }
}