package com.epfl.dedis.hbt.ui

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentManager
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.databinding.ActivityMainBinding
import com.epfl.dedis.hbt.ui.main.MainFragment
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // Set default fragment (MainFragment)
        setCurrentFragment(supportFragmentManager, MainFragment.newInstance(), false)
    }

    companion object {

        /**
         * This function should be used to set the current fragment shown on the main activity.
         *
         * @param fm the fragment manager of the activity
         * @param fragment to launch
         * @param addToBackstack whether the transaction should be added to the backstack (and thus reversible)
         *                       default value : true
         */
        @JvmStatic
        fun setCurrentFragment(
            fm: FragmentManager,
            fragment: Fragment,
            addToBackstack: Boolean = true
        ) {
            val transaction = fm.beginTransaction()
                .replace(R.id.container, fragment)
            if (addToBackstack)
                transaction.addToBackStack(null)
            transaction.commit()
        }
    }
}
