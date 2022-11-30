package com.epfl.dedis.hbt.ui

import android.content.Intent
import android.util.Log
import androidx.lifecycle.*
import com.epfl.dedis.hbt.test.fragment.FragmentScenario.Companion.TAG
import dagger.hilt.android.lifecycle.HiltViewModel

@HiltViewModel
class NFCViewModel : ViewModel() {

    private var curCallBack: ((Intent) -> Unit)? = null

    private val _listenToNFC = MutableLiveData(false)
    val listenToNFC: LiveData<Boolean> = _listenToNFC

    fun onNewIntent(intent: Intent) {
        curCallBack?.invoke(intent)
    }

    fun setCallback(lifecycle: Lifecycle, callback: (Intent) -> Unit) {
        if (lifecycle.currentState == Lifecycle.State.DESTROYED)
            return

        if (curCallBack != null) {
            Log.e(TAG, "A callback is already defined, cannot set a new one yet")
            return
        }

        this.curCallBack = callback
        // Subscribe to lifecycle to remove callback when the observer gets destroyed
        lifecycle.addObserver(LifecycleEventObserver { _, event ->
            if (event.targetState == Lifecycle.State.DESTROYED) curCallBack = null
        })
    }

    fun listenToNFC() {
        _listenToNFC.value = true
    }

    fun completeNFC() {
        _listenToNFC.value = false
    }
}