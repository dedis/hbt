package com.epfl.dedis.hbt.ui

import android.content.Intent
import android.util.Log
import androidx.lifecycle.*
import androidx.lifecycle.Lifecycle.State.CREATED
import androidx.lifecycle.Lifecycle.State.RESUMED
import com.epfl.dedis.hbt.test.fragment.FragmentScenario.Companion.TAG
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class NFCViewModel @Inject constructor() : ViewModel() {

    private var curCallBack: ((Intent) -> Unit)? = null

    private val _listenToNFC = MutableLiveData(false)
    val listenToNFC: LiveData<Boolean> = _listenToNFC

    fun onNewIntent(intent: Intent) {
        curCallBack?.invoke(intent)
    }

    fun setCallback(lifecycle: Lifecycle, callback: (Intent) -> Unit) {
        if (curCallBack != null) {
            Log.e(TAG, "A callback is already defined, cannot set a new one yet")
            return
        }

        // Subscribe to lifecycle set and the remove callback when the observer gets destroyed
        lifecycle.addObserver(LifecycleEventObserver { owner, event ->
            if (event.targetState == RESUMED) {
                curCallBack = callback
                _listenToNFC.value = true
            }
            // If we move out of resume state, remove callback
            if (owner.lifecycle.currentState == RESUMED && event.targetState == CREATED) {
                _listenToNFC.value = false
                curCallBack = null
            }
        })
    }
}