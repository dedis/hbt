package com.epfl.dedis.hbt.ui.wallet

import android.graphics.Bitmap
import android.graphics.Color
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.Toast
import androidx.activity.OnBackPressedCallback
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.transaction.TransactionState.*
import com.epfl.dedis.hbt.data.transaction.TransactionStateManager
import com.epfl.dedis.hbt.databinding.FragmentWalletShowqrBinding
import com.epfl.dedis.hbt.service.json.JsonService
import com.epfl.dedis.hbt.service.json.JsonType
import com.epfl.dedis.hbt.ui.MainActivity
import com.google.zxing.BarcodeFormat
import com.google.zxing.qrcode.QRCodeWriter
import dagger.hilt.android.AndroidEntryPoint
import java.lang.Integer.max
import javax.inject.Inject


@AndroidEntryPoint
class ShowQrFragment : Fragment() {

    companion object {
        private val TAG: String = ShowQrFragment::class.java.simpleName
    }

    @Inject
    lateinit var jsonService: JsonService

    @Inject
    lateinit var trxStateManager: TransactionStateManager

    private val walletViewModel: WalletViewModel by viewModels(ownerProducer = { requireActivity() })
    private var _binding: FragmentWalletShowqrBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWalletShowqrBinding.inflate(inflater, container, false).apply {
            walletName.text = walletViewModel.user.name
            val role = walletViewModel.user.role
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)
        }

        // Override back button such that it cancels current transaction
        requireActivity().onBackPressedDispatcher.addCallback(
            viewLifecycleOwner,
            object : OnBackPressedCallback(true) {
                override fun handleOnBackPressed() {
                    trxStateManager.cancelTransaction()
                }
            }
        )

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletButtonOk

        okButton.setOnClickListener {
            when (val state = trxStateManager.currentState.value) {
                is ReceiverShow -> trxStateManager.readCompleteTransaction(state.transaction)
                is SenderShow -> trxStateManager.completeSending(state.transaction)
                else -> {
                    Log.e(TAG, "Unhandled state in the ShowQrFragment : $state")
                    Toast.makeText(context, "Invalid transaction state", Toast.LENGTH_SHORT).show()
                }
            }
        }

        trxStateManager.currentState.observe(viewLifecycleOwner) {
            when (it) {
                is SenderRead, is ReceiverRead ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ScanQrFragment()
                    )
                None ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        WalletFragment()
                    )
                is ReceiverShow -> generateQrCode(JsonType.PENDING_TRANSACTION, it.transaction)
                is SenderShow -> generateQrCode(JsonType.COMPLETE_TRANSACTION, it.transaction)
            }
        }
    }

    private fun generateQrCode(type: JsonType, content: Any) {
        val imageView: ImageView = binding.walletQrImage
        val size = max(imageView.layoutParams.width, imageView.layoutParams.height)

        val qrCodeContent = jsonService.toJson(type, content)

        val bits = QRCodeWriter().encode(qrCodeContent, BarcodeFormat.QR_CODE, size, size)
        val bitmap = Bitmap.createBitmap(size, size, Bitmap.Config.RGB_565).also {
            for (x in 0 until size) {
                for (y in 0 until size) {
                    it.setPixel(x, y, if (bits[x, y]) Color.BLACK else Color.WHITE)
                }
            }
        }

        imageView.setImageBitmap(bitmap)
    }
}
