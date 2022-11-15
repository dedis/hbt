package com.epfl.dedis.hbt.ui.wallet

import androidx.camera.core.ImageAnalysis
import androidx.camera.mlkit.vision.MlKitAnalyzer
import androidx.core.util.Consumer
import com.google.mlkit.vision.barcode.BarcodeScanner
import java.util.concurrent.Executor
import javax.inject.Inject

class ImageAnalyzerProvider @Inject constructor() {

    fun provide(
        detector: BarcodeScanner,
        targetCoordinateSystem: Int,
        executor: Executor,
        consumer: Consumer<String>
    ): ImageAnalysis.Analyzer =
        MlKitAnalyzer(
            listOf(detector),
            targetCoordinateSystem,
            executor
        ) { result: MlKitAnalyzer.Result? ->
            val barcodeResults = result?.getValue(detector)
            // Test result value
            if (barcodeResults != null && barcodeResults.size != 0 && barcodeResults.first() != null) {
                consumer.accept(barcodeResults[0].rawValue ?: "")
            }
        }
}