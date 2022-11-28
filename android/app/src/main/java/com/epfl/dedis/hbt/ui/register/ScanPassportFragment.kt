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
        // https://en.wikipedia.org/wiki/Machine-readable_passport
        /**
         * Group 1 : Country code
         * Group 2 : Holder's name
         */
        private val LINE_1_PATTERN = Pattern.compile("P[A-Z<]([A-Z<]{3})([A-Z<]{39})")

        /**
         * Group 1 : Passport number
         * Group 2 : Passport number's checksum
         * Group 3 : Nationality
         * Group 4 : Date of birth (YYMMDD)
         * Group 5 : Date of birth checksum
         * Group 6 : Sex (M, F or < for male, female or unspecified)
         * Group 7 : Expiration date of passport (YYMMDD)
         * Group 8 : Expiration date's checksum
         * Group 9 : Personal number (may be used by the issuing country as it desires)
         * Group 10 : Personal number's checksum (may be < if all characters are <)
         * Group 11 : Checksum on Passport number, Date of birth, Expiration date and there checksums
         */
        private val LINE_2_PATTERN =
            Pattern.compile("([A-Z\\d<]{9})(\\d)([A-Z]{3})(\\d{6})(\\d)([A-B])(\\d{6})(\\d)([A-Z\\d<]{14})([\\d<])(\\d)")
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
                // The vision algorithm sometimes adds spaces and mistakes '<<' for '«'
                val text = raw.replace(" ", "").replace("«", "<<")
                val matcher1 = LINE_1_PATTERN.matcher(text)
                val matcher2 = LINE_2_PATTERN.matcher(text)

                if (matcher1.find() && matcher2.find()) {
                    val names = (matcher1.group(2) ?: "")
                        .split("<")
                        .filter { s -> s.isNotEmpty() }
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        RegisterFragment.newInstance(
                            names.joinToString(" "),
                            "",
                            matcher2.group(1)?.replace("<", "") ?: ""
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