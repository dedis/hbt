package com.epfl.dedis.hbt.ui.wallet

import android.graphics.Bitmap
import android.graphics.Matrix
import android.media.Image
import androidx.camera.core.ImageAnalysis
import androidx.camera.mlkit.vision.MlKitAnalyzer
import com.google.android.gms.tasks.Task
import com.google.mlkit.vision.interfaces.Detector
import java.nio.ByteBuffer
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
            produceRotationsDetector(detector),
            targetCoordinateSystem,
            executor
        ) {
            consumer(it?.getValue(detector))
        }

    private fun <T> produceRotationsDetector(detector: Detector<T>): List<Detector<T>> =
        listOf(
            detector,
            RotatedDetector(detector, 90),
            RotatedDetector(detector, -90),
            RotatedDetector(detector, 180)
        )

    /** A simple delegation pattern to feed the underlying detector a rotated image */
    private class RotatedDetector<T>(val detector: Detector<T>, val offset: Int) :
        Detector<T> by detector {

        override fun process(bitmap: Bitmap, rotation: Int): Task<T> =
            detector.process(bitmap, rotation + offset)

        override fun process(image: Image, rotation: Int): Task<T> =
            detector.process(image, rotation + offset)

        override fun process(image: Image, rotation: Int, coordinatesMatrix: Matrix): Task<T> =
            detector.process(image, rotation + offset, coordinatesMatrix)

        override fun process(
            byteBuffer: ByteBuffer,
            width: Int,
            height: Int,
            rotation: Int,
            format: Int
        ): Task<T> {
            return detector.process(byteBuffer, width, height, rotation + offset, format)
        }
    }
}