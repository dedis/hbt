package com.epfl.dedis.hbt.ui.wallet

import androidx.camera.core.ImageAnalysis
import androidx.camera.mlkit.vision.MlKitAnalyzer
import com.google.mlkit.vision.interfaces.Detector
import java.util.concurrent.Executor
import javax.inject.Inject

class ImageAnalyzerProvider @Inject constructor() {

    fun <T> provide(
        detector: Detector<T>,
        targetCoordinateSystem: Int,
        executor: Executor,
        consumer: (T?) -> Unit
    ): ImageAnalysis.Analyzer =
        MlKitAnalyzer(
            listOf(detector),
            targetCoordinateSystem,
            executor
        ) {
            consumer(it?.getValue(detector))
        }
}