package com.epfl.dedis.hbt.ui.wallet

import android.graphics.Bitmap
import android.graphics.Color
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.epfl.dedis.hbt.R
import com.epfl.dedis.hbt.data.model.Role
import com.epfl.dedis.hbt.databinding.FragmentWalletShowqrBinding
import com.epfl.dedis.hbt.ui.MainActivity
import com.google.zxing.BarcodeFormat
import com.google.zxing.qrcode.QRCodeWriter
import dagger.hilt.android.AndroidEntryPoint
import java.lang.Integer.max


@AndroidEntryPoint
class ShowQrFragment : Fragment() {

    companion object {
        private const val AMOUNT = "AMOUNT"

        fun newInstance(amount: Float) = ShowQrFragment().apply {
            val bundle = Bundle()
            bundle.putFloat(AMOUNT, amount)
            arguments = bundle
        }
    }

    private val walletViewModel: WalletViewModel by viewModels()
    private var _binding: FragmentWalletShowqrBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    private var transferAmount = 0F

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

            arguments?.getFloat(AMOUNT)?.let {
                transferAmount = it
            }
        }

        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        val okButton = binding.walletButtonOk

        generateQrCode(view)

        okButton.setOnClickListener {
            walletViewModel.receive(transferAmount)
            //TODO: move on to ScanFragment instead
            MainActivity.setCurrentFragment(parentFragmentManager, WalletFragment())
        }
    }

    private fun generateQrCode(view: View) {
        val imageView: ImageView = view.findViewById(R.id.walletQrImage) as ImageView
        val size = max(imageView.layoutParams.width, imageView.layoutParams.height)
        val qrCodeContent = getString(
            R.string.hbt_rx_transaction,
            walletViewModel.wallet?.pk.toString(),
            transferAmount
        )
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
