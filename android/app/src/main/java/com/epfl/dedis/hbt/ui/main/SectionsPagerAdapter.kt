package com.epfl.dedis.hbt.ui.main

import android.content.Context
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentManager
import androidx.fragment.app.FragmentPagerAdapter
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.ui.login.LoginFragment
import com.epfl.dedis.hbt.ui.wallet.WalletDefaultFragment

private val TAB_TITLES = arrayOf(
    R.string.tab_text_1,
    R.string.tab_text_2
)

/**
 * A [FragmentPagerAdapter] that returns a fragment corresponding to
 * one of the sections/tabs/pages.
 */
class SectionsPagerAdapter(private val context: Context, fm: FragmentManager) :
    FragmentPagerAdapter(fm) {

    override fun getItem(position: Int): Fragment {
        // getItem is called to instantiate the fragment for the given page.
        var fragment = when (position) {
            0 -> {
                LoginFragment.newInstance()
            }
            1 -> {
                WalletDefaultFragment.newInstance()
            }
            else -> throw Exception("Only 2 tabs supported for now !!!")
        }
        return fragment
    }

    override fun getPageTitle(position: Int): CharSequence? {
        return context.resources.getString(TAB_TITLES[position])
    }

    override fun getCount(): Int {
        return TAB_TITLES.count()
    }
}