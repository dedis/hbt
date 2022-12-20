package com.epfl.dedis.hbt.ui.register

import android.Manifest
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.activity.result.contract.ActivityResultContracts
import androidx.camera.view.CameraController
import androidx.camera.view.LifecycleCameraController
import androidx.camera.view.PreviewView
import androidx.core.content.ContextCompat
import androidx.core.content.PermissionChecker
import androidx.fragment.app.Fragment
import com.epfl.dedis.hbt.data.Result.Error
import com.epfl.dedis.hbt.data.Result.Success
import com.epfl.dedis.hbt.databinding.FragmentPassportScanBinding
import com.epfl.dedis.hbt.service.passport.mrz.MRZExtractor
import com.epfl.dedis.hbt.service.passport.mrz.ValidationException
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.wallet.ImageAnalyzerProvider
import com.epfl.dedis.hbt.utility.json.JsonService
import com.google.mlkit.vision.text.TextRecognition
import com.google.mlkit.vision.text.TextRecognizer
import com.google.mlkit.vision.text.latin.TextRecognizerOptions
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class ScanPassportFragment : Fragment() {

    companion object {
        private val TAG = ScanPassportFragment::class.java.simpleName
    }

    private var _binding: FragmentPassportScanBinding? = null

    @Inject
    lateinit var jsonService: JsonService

    @Inject
    lateinit var imageAnalyzerProvider: ImageAnalyzerProvider

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    private var scanner: TextRecognizer? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Create a request permission launcher which will ask for permission when launched
        val requestPermissionLauncher = registerForActivityResult(
            ActivityResultContracts.RequestPermission(),
            requireActivity().activityResultRegistry
        ) {
            // This is the callback of the permission granter
            applyPermissionToView()
        }

        _binding = FragmentPassportScanBinding.inflate(inflater, container, false).apply {
            requestPermissionButton.setOnClickListener {
                // Launch the permission request on click
                requestPermissionLauncher.launch(
                    Manifest.permission.CAMERA
                )
            }

            manualInput.setOnClickListener {
                MainActivity.setCurrentFragment(
                    parentFragmentManager,
                    PassportDataFragment()
                )
            }
        }

        return binding.root
    }

    override fun onResume() {
        super.onResume()
        applyPermissionToView()
    }

    override fun onDestroy() {
        super.onDestroy()

        scanner?.close()
    }

    private fun startCamera() {
        val cameraController = LifecycleCameraController(requireContext())
        val previewView: PreviewView = binding.previewView

        scanner = TextRecognition.getClient(TextRecognizerOptions.DEFAULT_OPTIONS)

        cameraController.setImageAnalysisAnalyzer(
            ContextCompat.getMainExecutor(requireActivity()),
            imageAnalyzerProvider.provide(
                scanner!!,
                CameraController.COORDINATE_SYSTEM_VIEW_REFERENCED,
                ContextCompat.getMainExecutor(requireActivity())
            ) {
                val raw = it?.text ?: return@provide
                // The vision algorithm sometimes adds spaces and mistakes '<<' for '«'
                val text = raw.replace(" ", "").replace("«", "<<")
                when (val result = MRZExtractor.match(text)) {
                    is Success ->
                        MainActivity.setCurrentFragment(
                            parentFragmentManager,
                            NFCPassportFragment.newInstance(result.data)
                        )
                    is Error -> when (result.exception) {
                        is ValidationException ->
                            Log.i(TAG, result.exception.message ?: "")
                        else ->
                            Log.d(TAG, result.exception.message ?: "")
                    }
                }
            }
        )

        cameraController.bindToLifecycle(this)
        previewView.controller = cameraController
    }

    private fun applyPermissionToView() {
        // Depending on the current permission state,
        // show the preview or the request permission views
        if (isPermissionGranted()) {
            binding.requestPermission.visibility = View.GONE
            binding.previewView.visibility = View.VISIBLE
            startCamera()
        } else {
            binding.requestPermission.visibility = View.VISIBLE
            binding.previewView.visibility = View.GONE
        }
    }

    private fun isPermissionGranted() =
        PermissionChecker.checkSelfPermission(
            requireContext(),
            Manifest.permission.CAMERA
        ) == PermissionChecker.PERMISSION_GRANTED
}