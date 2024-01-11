package com.epfl.dedis.hbt.service.http

import io.reactivex.rxjava3.core.Single
import io.reactivex.rxjava3.subjects.SingleSubject
import org.chromium.net.CronetEngine
import org.chromium.net.CronetException
import org.chromium.net.UrlRequest
import org.chromium.net.UrlResponseInfo
import java.nio.ByteBuffer
import java.util.concurrent.Executors
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class HttpService @Inject constructor(private val engine: CronetEngine) {

    private val executor = Executors.newSingleThreadExecutor()

    fun get(url: String): Single<String> {
        val single = SingleSubject.create<String>()

        engine.newUrlRequestBuilder(url, Callback(single), executor)
            .setHttpMethod("GET")
            .build()
            .start()

        return single
    }

    private class Callback(private val single: SingleSubject<String>) : UrlRequest.Callback() {
        override fun onRedirectReceived(
            request: UrlRequest?,
            info: UrlResponseInfo?,
            newLocationUrl: String?
        ) {
            request?.followRedirect()
        }

        override fun onResponseStarted(request: UrlRequest?, info: UrlResponseInfo?) {
            request?.read(ByteBuffer.allocateDirect(BYTE_BUFFER_CAPACITY_BYTES))
        }

        override fun onReadCompleted(
            request: UrlRequest?,
            info: UrlResponseInfo?,
            byteBuffer: ByteBuffer?
        ) {
            byteBuffer?.clear()
            request?.read(byteBuffer)
        }

        override fun onSucceeded(request: UrlRequest?, info: UrlResponseInfo?) {
            single.onSuccess(info?.httpStatusCode.toString())
        }

        override fun onFailed(
            request: UrlRequest?,
            info: UrlResponseInfo?,
            error: CronetException?
        ) {
            single.onError(error ?: Exception("Unknown error"))
        }

        companion object {
            private const val BYTE_BUFFER_CAPACITY_BYTES = 100 * 1024
        }
    }
}