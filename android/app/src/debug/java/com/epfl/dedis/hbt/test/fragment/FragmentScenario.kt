package com.epfl.dedis.hbt.test.fragment

import android.content.ComponentName
import android.content.Intent
import android.os.Bundle
import androidx.annotation.IdRes
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentFactory
import androidx.lifecycle.Lifecycle
import androidx.test.core.app.ActivityScenario
import androidx.test.core.app.ActivityScenario.ActivityAction
import androidx.test.core.app.ApplicationProvider
import java.util.*
import java.util.function.Consumer
import java.util.function.Supplier

/**
 * This class allows easy testing of fragments. It is greatly inspired by Android's fragment scenario and
 * https://homanhuang.medium.com/to-make-fragment-test-under-hilt-installed-65ff2d5e5eb6#1736
 *
 * It allows Hilt injection and custom activity holding the fragment
 *
 * @param A Activity class
 * @param F Fragment class
 */
class FragmentScenario<A : AppCompatActivity, F : Fragment> private constructor(
    private val activityScenario: ActivityScenario<A>, private val fragmentClass: Class<F>
) {
    /**
     * Recreate the scenario
     *
     * @return the scenario
     */
    fun recreate(): FragmentScenario<A, F> {
        activityScenario.recreate()
        return this
    }

    /**
     * Execute on action on the activity the scenario is running on
     *
     * @param action to execute on the activity
     * @return the scenario
     */
    fun onActivity(action: ActivityAction<A>): FragmentScenario<A, F> {
        activityScenario.onActivity(action)
        return this
    }

    /**
     * Execute an action on the fragment the scenario is running on
     *
     * @param action to execute on the fragment
     * @return the scenario
     */
    fun onFragment(action: Consumer<F>): FragmentScenario<A, F> {
        activityScenario.onActivity { activity: A ->
            action.accept(
                Objects.requireNonNull(
                    fragmentClass.cast(
                        activity.supportFragmentManager.findFragmentByTag(TAG)
                    )
                )
            )
        }
        return this
    }

    /**
     * Advance the fragment to a new state
     *
     * @param newState to move onto
     * @return the scenario
     */
    fun moveToState(newState: Lifecycle.State): FragmentScenario<A, F> {
        if (newState == Lifecycle.State.DESTROYED) {
            activityScenario.onActivity { activity: A ->
                val fragment = activity.supportFragmentManager.findFragmentByTag(TAG)

                if (fragment != null) {
                    activity.supportFragmentManager.beginTransaction().remove(fragment).commitNow()
                }
            }
        } else {
            activityScenario.onActivity { activity: A ->
                val fragment = activity.supportFragmentManager.findFragmentByTag(
                    TAG
                ) ?: throw IllegalStateException("fragment is already destroyed")

                activity
                    .supportFragmentManager
                    .beginTransaction()
                    .setMaxLifecycle(fragment, newState)
                    .commitNow()
            }
        }
        return this
    }

    /** Close the scenario. This should be called at the end of any test */
    fun close() {
        activityScenario.close()
    }

    private class SimpleFragmentFactory(private val supplier: Supplier<out Fragment>) :
        FragmentFactory() {
        override fun instantiate(classLoader: ClassLoader, className: String): Fragment {
            return supplier.get()
        }
    }

    companion object {

        /**
         * Launch a new FragmentScenario
         *
         * @param fragmentClass fragment to launch the scenario on
         * @param F Fragment type
         * @return the launched scenario
         */
        fun <F : Fragment> launch(
            fragmentClass: Class<F>,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenario<EmptyHiltActivity, F> {
            return launchIn(
                EmptyHiltActivity::class.java,
                android.R.id.content,
                fragmentClass,
                fragmentArgs = fragmentArgs
            )
        }

        /**
         * Launch a new FragmentScenario
         *
         * @param fragmentClass fragment to launch the scenario on
         * @param fragmentSupplier supplier that creates the fragment object
         * @param F Fragment type
         * @return the launched scenario
         */
        fun <F : Fragment> launch(
            fragmentClass: Class<F>,
            fragmentSupplier: Supplier<F>,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenario<EmptyHiltActivity, F> {
            return launchIn(
                EmptyHiltActivity::class.java,
                android.R.id.content,
                fragmentClass,
                fragmentSupplier,
                fragmentArgs = fragmentArgs
            )
        }

        /**
         * Launch a new FragmentScenario
         *
         * @param activityClass activity to launch the scenario on
         * @param activityArgs arguments of the activity
         * @param contentId id of the placeholder where the fragment will be put
         * @param fragmentClass fragment to launch the scenario on
         * @param fragmentSupplier supplier that creates the fragment object
         * @param fragmentArgs arguments of the fragment
         * @param A Activity type
         * @param F Fragment type
         * @return the launched scenario
         */
        fun <A : AppCompatActivity, F : Fragment> launchIn(
            activityClass: Class<A>,
            @IdRes contentId: Int,
            fragmentClass: Class<F>,
            fragmentSupplier: Supplier<F>,
            activityArgs: Bundle? = Bundle.EMPTY,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenario<A, F> {
            return launchIn(
                activityClass,
                contentId,
                fragmentClass,
                factory(fragmentSupplier),
                activityArgs,
                fragmentArgs
            )
        }

        private fun <F : Fragment> factory(fragmentSupplier: Supplier<F>): FragmentFactory {
            return SimpleFragmentFactory(fragmentSupplier)
        }

        const val TAG = "FRAGMENT"

        /**
         * Launch a new FragmentScenario with following arguments :
         *
         * @param activityClass activity to launch the fragment on
         * @param contentId id of the placeholder where the fragment will be put
         * @param fragmentClass fragment to launch
         * @param factory that produces the fragment object. If null, the android default will be used.
         * @param activityArgs arguments of the activity
         * @param fragmentArgs arguments of the fragment
         *
         * @param A Activity type
         * @param F Fragment type
         * @return the launched FragmentScenario
         */
        fun <A : AppCompatActivity, F : Fragment> launchIn(
            activityClass: Class<A>,
            @IdRes contentId: Int,
            fragmentClass: Class<F>,
            factory: FragmentFactory? = null,
            activityArgs: Bundle? = Bundle.EMPTY,
            fragmentArgs: Bundle? = Bundle.EMPTY
        ): FragmentScenario<A, F> {
            val mainActivityIntent = Intent.makeMainActivity(
                ComponentName(ApplicationProvider.getApplicationContext(), activityClass)
            )
            mainActivityIntent.putExtras(activityArgs!!)
            val scenario = ActivityScenario.launch<A>(mainActivityIntent)
            val fragmentScenario = FragmentScenario(scenario, fragmentClass)
            scenario.onActivity { activity: A ->
                if (factory != null) {
                    activity.supportFragmentManager.fragmentFactory = factory
                }
                val fragment = activity
                    .supportFragmentManager
                    .fragmentFactory
                    .instantiate(
                        fragmentClass.classLoader!!,
                        fragmentClass.name
                    )
                fragment.arguments = fragmentArgs
                activity
                    .supportFragmentManager
                    .beginTransaction()
                    .replace(contentId, fragment, TAG)
                    .setMaxLifecycle(fragment, Lifecycle.State.RESUMED)
                    .commitNow()
            }
            return fragmentScenario
        }
    }
}