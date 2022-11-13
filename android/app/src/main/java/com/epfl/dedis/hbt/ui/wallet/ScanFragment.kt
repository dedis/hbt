package com.epfl.dedis.hbt.ui.wallet

import android.Manifest
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.OnBackPressedCallback
import androidx.activity.result.contract.ActivityResultContracts
import androidx.camera.mlkit.vision.MlKitAnalyzer
import androidx.camera.view.CameraController.COORDINATE_SYSTEM_VIEW_REFERENCED
import androidx.camera.view.LifecycleCameraController
import androidx.camera.view.PreviewView
import androidx.core.content.ContextCompat
import androidx.core.content.PermissionChecker.PERMISSION_GRANTED
import androidx.core.content.PermissionChecker.checkSelfPermission
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Observer
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.CompleteTransaction
import com.epfl.dedis.hbt.data.model.PendingTransaction
import com.epfl.dedis.hbt.databinding.FragmentWalletScanBinding
import com.epfl.dedis.hbt.ui.wallet.TransactionState.*
import com.epfl.dedis.hbt.utility.json.JsonService
import com.epfl.dedis.hbt.utility.json.JsonType
import com.google.mlkit.vision.barcode.BarcodeScanner
import com.google.mlkit.vision.barcode.BarcodeScannerOptions
import com.google.mlkit.vision.barcode.BarcodeScanning
import com.google.mlkit.vision.barcode.common.Barcode
import dagger.hilt.android.AndroidEntryPoint
import java.util.concurrent.ExecutorService
import javax.inject.Inject


@AndroidEntryPoint
class ScanFragment : Fragment() {

    companion object {
        private val TAG: String = ScanFragment::class.java.simpleName
    }

    private val walletViewModel: WalletViewModel by viewModels(ownerProducer = { requireActivity() })
    private var _binding: FragmentWalletScanBinding? = null

    @Inject
    lateinit var jsonService: JsonService

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    private lateinit var cameraExecutor: ExecutorService
    private lateinit var barcodeScanner: BarcodeScanner

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

        _binding = FragmentWalletScanBinding.inflate(inflater, container, false).apply {
            walletName.text = walletViewModel.user.name
            val role = walletViewModel.user.role
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)

            requestPermissionButton.setOnClickListener {
                // Launch the permission request on click
                requestPermissionLauncher.launch(
                    Manifest.permission.CAMERA
                )
            }
        }

        // Override back button such that it cancels current transaction
        requireActivity().onBackPressedDispatcher.addCallback(
            viewLifecycleOwner,
            object : OnBackPressedCallback(true) {
                override fun handleOnBackPressed() {
                    walletViewModel.transitionTo(None)
                }
            }
        )

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        walletViewModel.walletFormState.observe(
            viewLifecycleOwner,
            Observer { walletFormState ->
                if (walletFormState == null) {
                    return@Observer
                }
            })

        walletViewModel.walletResult.observe(viewLifecycleOwner,
            Observer { walletResult ->
                walletResult ?: return@Observer
                if (walletResult.error != null) {
//                  onRegisterFailed(walletResult.error)
                } else {
//                  onRegisterSuccess(usernameEditText.text.toString())
                }
            })
    }

    override fun onResume() {
        super.onResume()
        applyPermissionToView()
    }

    override fun onDestroy() {
        super.onDestroy()
        cameraExecutor.shutdown()
        barcodeScanner.close()
    }

    private fun startCamera() {
        val cameraController = LifecycleCameraController(requireContext())
        val previewView: PreviewView = binding.previewView

        val options = BarcodeScannerOptions.Builder()
            .setBarcodeFormats(Barcode.FORMAT_QR_CODE)
            .build()
        barcodeScanner = BarcodeScanning.getClient(options)

        cameraController.setImageAnalysisAnalyzer(
            ContextCompat.getMainExecutor(requireActivity()),
            MlKitAnalyzer(
                listOf(barcodeScanner),
                COORDINATE_SYSTEM_VIEW_REFERENCED,
                ContextCompat.getMainExecutor(requireActivity())
            ) { result: MlKitAnalyzer.Result? ->
                val barcodeResults = result?.getValue(barcodeScanner)
                // Test result value
                if ((barcodeResults == null) ||
                    (barcodeResults.size == 0) ||
                    (barcodeResults.first() == null)
                ) {
                    return@MlKitAnalyzer
                }

                onResult(barcodeResults[0])
            }
        )

        cameraController.bindToLifecycle(this)
        previewView.controller = cameraController
    }

    private fun onResult(barcode: Barcode) {
        when (val state = walletViewModel.transactionState.value) {
            is ReceiverRead -> {
                val trx = jsonService.fromJson(
                    JsonType.COMPLETE_TRANSACTION,
                    barcode.rawValue ?: "",
                    CompleteTransaction::class
                )
                walletViewModel.receive(trx)
            }
            SenderRead -> {
                val trx = jsonService.fromJson(
                    JsonType.PENDING_TRANSACTION,
                    barcode.rawValue ?: "",
                    PendingTransaction::class
                )
                walletViewModel.transitionTo(SenderShow(trx.withSource(walletViewModel.user.name)))
            }
            else -> {
                Log.e(TAG, "Unhandled state in the ShowQrFragment : $state")
                Toast.makeText(context, "Invalid transaction state", Toast.LENGTH_SHORT).show()
            }
        }
    }

    private fun applyPermissionToView() {
        // Depending on the current permission state,
        // show the preview or the request permission views
        if (checkPermission()) {
            binding.requestPermission.visibility = View.GONE
            binding.previewView.visibility = View.VISIBLE
            startCamera()
        } else {
            binding.requestPermission.visibility = View.VISIBLE
            binding.previewView.visibility = View.GONE
        }
    }

    private fun checkPermission() =
        checkSelfPermission(requireContext(), Manifest.permission.CAMERA) == PERMISSION_GRANTED
}
