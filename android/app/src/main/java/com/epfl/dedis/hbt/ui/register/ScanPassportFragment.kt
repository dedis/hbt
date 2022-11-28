package com.epfl.dedis.hbt.ui.register

import android.Manifest
import android.os.Bundle
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
import com.epfl.dedis.hbt.databinding.FragmentPassportScanBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.wallet.ImageAnalyzerProvider
import com.epfl.dedis.hbt.utility.json.JsonService
import com.google.mlkit.vision.text.TextRecognition
import com.google.mlkit.vision.text.TextRecognizer
import com.google.mlkit.vision.text.latin.TextRecognizerOptions
import dagger.hilt.android.AndroidEntryPoint
import java.util.regex.Pattern
import javax.inject.Inject

@AndroidEntryPoint
class ScanPassportFragment : Fragment() {

    companion object {
        private val LINE_1_PATTERN = Pattern.compile("(P[A-Z0-9<]{1})([A-Z]{3})([A-Z0-9<]{39})")
        private val LINE_2_PATTERN =
            Pattern.compile("([A-Z0-9<]{9})([0-9]{1})([A-Z]{3})([0-9]{6})([0-9]{1})([M|F|X|<]{1})([0-9]{6})([0-9]{1})([A-Z0-9<]{14})([0-9<]{1})([0-9]{1})")
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
                val text = raw.replace(" ", "").replace("«", "<<")
                val matcher1 = LINE_1_PATTERN.matcher(text)
                val matcher2 = LINE_2_PATTERN.matcher(text)

                // https://en.wikipedia.org/wiki/Machine-readable_passport
                if (matcher1.find() && matcher2.find()) {
                    val names = matcher1.group()
                        .substring(5)
                        .split("<")
                        .filter { s -> s.isNotEmpty() }
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        RegisterFragment.newInstance(
                            names.joinToString(" "),
                            "",
                            matcher2.group().substring(0, 9).replace("<", "")
                        )
                    )
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