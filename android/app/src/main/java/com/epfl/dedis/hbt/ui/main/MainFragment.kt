package com.epfl.dedis.hbt.ui.main

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentManager
import androidx.lifecycle.Lifecycle
import androidx.viewpager2.adapter.FragmentStateAdapter
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.databinding.FragmentMainBinding
import com.epfl.dedis.hbt.ui.login.LoginFragment
import com.epfl.dedis.hbt.ui.wallet.WalletDefaultFragment
import com.google.android.material.tabs.TabLayoutMediator
import dagger.hilt.android.AndroidEntryPoint

// Object argument to instantiate the Fragment
private const val POSITION = "position"

// Array of the tab titles
private val TAB_TITLES = arrayOf(
    R.string.tab_user_name,
    R.string.tab_wallet
)

/**
 * A simple [Fragment] subclass.
 * Use the [MainFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
@AndroidEntryPoint
class MainFragment : Fragment() {

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        return FragmentMainBinding.inflate(inflater, container, false).apply {
            // Create the pager adapter
            mainPager.adapter = SectionsPagerAdapter(childFragmentManager, lifecycle)
            // Link tab layout with the pager
            TabLayoutMediator(mainTabLayout, mainPager) { tab, pos ->
                tab.text = getString(TAB_TITLES[pos])
            }.attach()

            // Set opened page
            arguments?.let {
                mainPager.currentItem = it.getInt(POSITION)
            }
        }.root
    }

    class SectionsPagerAdapter(fm: FragmentManager, lifecycle: Lifecycle) :
        FragmentStateAdapter(fm, lifecycle) {

        override fun getItemCount(): Int = TAB_TITLES.size

        override fun createFragment(position: Int): Fragment {
            return when (position) {
                0 -> LoginFragment.newInstance()
                1 -> WalletDefaultFragment.newInstance()
                else -> TODO("Tab number $position is not setup.")
            }
        }
    }

    companion object {

        /**
         * Use this factory method to create a new instance of
         * this fragment using the provided parameters.
         *
         * @param position The tab index to open
         * @return A new instance of fragment MainFragment.
         */
        @JvmStatic
        fun newInstance(position: Int = 0) =
            MainFragment().apply {
                arguments = Bundle().apply {
                    putInt(POSITION, position)
                }
            }
    }
}