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
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.PendingTransaction
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentWalletShowqrBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.epfl.dedis.hbt.ui.wallet.TransactionState.*
import com.epfl.dedis.hbt.utility.json.JsonService
import com.epfl.dedis.hbt.utility.json.JsonType
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
            walletName.text = walletViewModel.user?.name.toString()
            val role = walletViewModel.user?.role ?: Role.BENEFICIARY
            walletRole.text = getString(role.roleName)
            walletBalance.text =
                getString(R.string.hbt_currency, walletViewModel.wallet?.balance ?: 0.0f)
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletButtonOk

        okButton.setOnClickListener {
            when (val state = walletViewModel.transactionState.value) {
                is ReceiverShow -> walletViewModel.transitionTo(
                    ReceiverRead(
                        PendingTransaction(
                            walletViewModel.user!!.passport,
                            state.amount,
                            state.datetime
                        )
                    )
                )
                is SenderShow ->
                    // TODO Say it is complete ?
                    walletViewModel.transitionTo(None)
                else -> {
                    Log.e(TAG, "Unhandled state in the ShowQrFragment : $state")
                    Toast.makeText(context, "Invalid transaction state", Toast.LENGTH_SHORT).show()
                }
            }
        }

        walletViewModel.transactionState.observe(viewLifecycleOwner) {
            when (it) {
                is SenderRead, is ReceiverRead ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        ScanFragment()
                    )
                None ->
                    MainActivity.setCurrentFragment(
                        parentFragmentManager,
                        WalletFragment()
                    )
                is ReceiverShow -> {
                    generateQrCode(
                        JsonType.PENDING_TRANSACTION,
                        PendingTransaction(
                            walletViewModel.user!!.passport,
                            it.amount,
                            it.datetime
                        )
                    )
                }
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
