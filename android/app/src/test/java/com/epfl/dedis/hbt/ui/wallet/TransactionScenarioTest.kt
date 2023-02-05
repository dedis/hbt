package com.epfl.dedis.hbt.ui.wallet

import androidx.camera.core.ImageAnalysis.Analyzer
import androidx.test.espresso.action.ViewActions.click
import androidx.test.espresso.action.ViewActions.replaceText
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.*
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import androidx.test.rule.GrantPermissionRule
import com.epfl.dedis.hbt.data.transaction.CompleteTransaction
import com.epfl.dedis.hbt.data.transaction.PendingTransaction
import com.epfl.dedis.hbt.data.transaction.TransactionState
import com.epfl.dedis.hbt.data.transaction.TransactionStateManager
import com.epfl.dedis.hbt.data.user.User
import com.epfl.dedis.hbt.data.user.UserRepository
import com.epfl.dedis.hbt.data.user.Wallet
import com.epfl.dedis.hbt.service.json.JsonService
import com.epfl.dedis.hbt.service.json.JsonType
import com.epfl.dedis.hbt.test.fragment.FragmentScenarioRule
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.currentFragment
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmount
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmountFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmountOk
import com.epfl.dedis.hbt.test.ui.page.wallet.ScanQrFragmentPage.scanQrFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.ShowQrFragmentPage.showOk
import com.epfl.dedis.hbt.test.ui.page.wallet.ShowQrFragmentPage.showQrFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.receive
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.send
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.walletFragmentId
import com.google.mlkit.common.sdkinternal.MlKitContext
import com.google.mlkit.vision.barcode.BarcodeScanner
import com.google.mlkit.vision.barcode.common.Barcode
import dagger.hilt.android.testing.BindValue
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.junit.Rule
import org.junit.Test
import org.junit.rules.ExternalResource
import org.junit.runner.RunWith
import org.mockito.kotlin.*
import javax.inject.Inject

@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class TransactionScenarioTest {

    @BindValue
    lateinit var userRepo: UserRepository

    @BindValue
    lateinit var fakeImageAnalyzerProvider: ImageAnalyzerProvider
    private lateinit var resultConsumer: (List<Barcode>?) -> Unit

    @Inject
    lateinit var jsonService: JsonService

    @Inject
    lateinit var transactionStateManager: TransactionStateManager

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val setupRule = object : ExternalResource() {
        override fun before() = setup()
    }

    @get:Rule(order = 2)
    val fragmentRule = FragmentScenarioRule.launch(WalletFragment::class.java)

    @get:Rule(order = 3)
    val permissionRule: GrantPermissionRule =
        GrantPermissionRule.grant(android.Manifest.permission.CAMERA)

    private val user = User("Jon Smith", 12345, "Jon's passport")
    private val wallet = Wallet()

    fun setup() {
        hiltRule.inject()

        // We need to manually initialize MLKit's context
        MlKitContext.initializeIfNeeded(InstrumentationRegistry.getInstrumentation().context)

        // Create a fake image analyzer whose sole purpose is to retrieve to result consumer
        // of the qrcode scanning pipeline
        fakeImageAnalyzerProvider = mock {
            on { provide(isA<BarcodeScanner>(), any(), any(), any()) } doAnswer {
                resultConsumer = it.getArgument(3) as (List<Barcode>?) -> Unit
                it.callRealMethod() as Analyzer
            }
        }

        // Mock user repo
        userRepo = mock {
            on { loggedInUser } doAnswer { user }

            on { isLoggedIn } doAnswer { true }

            on { wallet } doAnswer { wallet }

            on { isRegistered(any()) } doReturn true
        }
    }

    @Test
    fun receiverScenario() {
        // Start the receiver process
        receive().perform(click())

        currentFragment().check(matches(withId(rxAmountFragmentId())))

        // Set the amount to 115.5 and press ok
        rxAmountOk().check(matches(isNotEnabled()))
        rxAmount().perform(replaceText("115.5"))
        rxAmountOk().check(matches(isEnabled())).perform(click())

        currentFragment().check(matches(withId(showQrFragmentId())))
        // TODO It is probably a good idea to actually check the shown qrcode value
        val dateTime = getTrxDatetime()

        // Act as if the sender scanned the QRCode and click on Ok
        showOk().perform(click())

        currentFragment().check(matches(withId(scanQrFragmentId())))


        // Provide a fake qrcode result that is a valid complete transaction
        resultConsumer(
            createBarcode(
                jsonService.toJson(
                    JsonType.COMPLETE_TRANSACTION, CompleteTransaction(
                        "Ben's passport",
                        user.passport,
                        115.5F,
                        dateTime
                    )
                )
            )
        )

        // the transaction is complete, we should be back to the wallet fragment
        currentFragment().check(matches(withId(walletFragmentId())))
    }

    private fun getTrxDatetime(): Long =
        when (val state = transactionStateManager.currentState.value) {
            is TransactionState.ReceiverShow -> state.transaction.datetime
            else -> throw IllegalStateException("Wrong transaction state")
        }


    @Test
    fun senderScenario() {
        // Start the sender transaction process
        send().perform(click())

        currentFragment().check(matches(withId(scanQrFragmentId())))

        // Provide a fake qrcode result that is a valid pending transaction
        resultConsumer(
            createBarcode(
                jsonService.toJson(
                    JsonType.PENDING_TRANSACTION, PendingTransaction(
                        "Ben's passport",
                        115.5F,
                        33917321
                    )
                )
            )
        )

        currentFragment().check(matches(withId(showQrFragmentId())))

        // After the pending transaction is scanned, the receiver scan the complete transaction
        // Then the sender presses Ok
        showOk().perform(click())

        // the transaction is complete, we should be back to the wallet fragment
        currentFragment().check(matches(withId(walletFragmentId())))
    }

    private fun createBarcode(value: String): List<Barcode> =
        listOf(
            mock {
                on { rawValue } doReturn value
            }
        )
}