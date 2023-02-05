package com.epfl.dedis.hbt

import android.app.Application
import com.epfl.dedis.hbt.service.json.JsonService
import dagger.hilt.android.HiltAndroidApp
import javax.inject.Inject

@HiltAndroidApp
class HBTApplication : Application() {

    @Inject
    lateinit var jsonService: JsonService

    override fun onCreate() {
        super.onCreate()
        // Load the schema at the start of the application
        jsonService.loadSchemas()
    }
}