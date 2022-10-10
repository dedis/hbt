package com.epfl.dedis.hbt.test.fragment

import android.os.Bundle
import androidx.fragment.app.Fragment
import com.epfl.dedis.hbt.ui.MainActivity
import java.util.function.Supplier

/**
 * An [ActivityFragmentScenarioRule] where the activity doesn't matter (And is thus an empty activity)
 *
 * @param F Fragment type
 */
class FragmentScenarioRule<F : Fragment>(
    fragmentScenarioSupplier: Supplier<FragmentScenario<MainActivity, F>>
) : ActivityFragmentScenarioRule<MainActivity, F>(fragmentScenarioSupplier) {
    companion object {
        fun <F : Fragment> launch(
            fragmentClass: Class<F>,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenarioRule<F> {
            return FragmentScenarioRule {
                FragmentScenario.launch(
                    fragmentClass,
                    fragmentArgs
                )
            }
        }

        fun <F : Fragment> launch(
            fragmentClass: Class<F>,
            supplier: Supplier<F>,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenarioRule<F> {
            return FragmentScenarioRule {
                FragmentScenario.launch(
                    fragmentClass,
                    supplier,
                    fragmentArgs
                )
            }
        }
    }
}