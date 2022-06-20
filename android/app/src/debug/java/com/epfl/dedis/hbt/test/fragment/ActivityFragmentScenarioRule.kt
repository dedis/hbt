package com.epfl.dedis.hbt.test.fragment

import android.os.Bundle
import androidx.annotation.IdRes
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentFactory
import org.junit.rules.ExternalResource
import java.util.function.Supplier

/**
 * A test rule that creates a Fragment for each test
 *
 * @param <F> the fragment type
</F> */
open class ActivityFragmentScenarioRule<A : AppCompatActivity, F : Fragment>(private val scenarioSupplier: Supplier<FragmentScenario<A, F>>) :
    ExternalResource() {
    private var scenario: FragmentScenario<A, F>? = null

    @Throws(Throwable::class)
    override fun before() {
        scenario = scenarioSupplier.get()
    }

    override fun after() {
        if (scenario != null) {
            scenario!!.close()
            scenario = null
        }
    }

    fun getScenario(): FragmentScenario<A, F> {
        return scenario!!
    }

    companion object {
        fun <A : AppCompatActivity, F : Fragment> launchIn(
            activityClass: Class<A>,
            @IdRes contentId: Int,
            fragmentClass: Class<F>,
            supplier: Supplier<F>,
            activityArgs: Bundle? = Bundle.EMPTY,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): ActivityFragmentScenarioRule<A, F> {
            return ActivityFragmentScenarioRule {
                FragmentScenario.launchIn(
                    activityClass,
                    contentId,
                    fragmentClass,
                    supplier,
                    activityArgs,
                    fragmentArgs
                )
            }
        }

        fun <A : AppCompatActivity, F : Fragment> launchIn(
            activityClass: Class<A>,
            @IdRes contentId: Int,
            fragmentClass: Class<F>,
            activityArgs: Bundle? = Bundle.EMPTY,
            fragmentArgs: Bundle? = Bundle.EMPTY,
            factory: FragmentFactory? = null
        ): ActivityFragmentScenarioRule<A, F> {
            return ActivityFragmentScenarioRule {
                FragmentScenario.launchIn(
                    activityClass,
                    contentId,
                    fragmentClass,
                    factory,
                    activityArgs,
                    fragmentArgs
                )
            }
        }
    }
}