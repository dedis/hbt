package com.epfl.dedis.hbt.ui.wallet

import androidx.test.espresso.action.ViewActions.click
import androidx.test.espresso.action.ViewActions.replaceText
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.withId
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import androidx.test.rule.GrantPermissionRule
import com.epfl.dedis.hbt.data.UserRepository
import com.epfl.dedis.hbt.data.model.CompleteTransaction
import com.epfl.dedis.hbt.data.model.User
import com.epfl.dedis.hbt.data.model.Wallet
import com.epfl.dedis.hbt.test.fragment.FragmentScenarioRule
import com.epfl.dedis.hbt.test.ui.page.MainActivityPage.currentFragment
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmount
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmountFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.RxAmountFragmentPage.rxAmountOk
import com.epfl.dedis.hbt.test.ui.page.wallet.ScanFragmentPage.scanFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.ShowFragmentPage.showFragmentId
import com.epfl.dedis.hbt.test.ui.page.wallet.ShowFragmentPage.showOk
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.receive
import com.epfl.dedis.hbt.test.ui.page.wallet.WalletFragmentPage.walletFragmentId
import com.epfl.dedis.hbt.utility.json.JsonService
import com.epfl.dedis.hbt.utility.json.JsonType
import com.google.android.gms.tasks.Tasks
import com.google.android.odml.image.MlImage
import com.google.mlkit.vision.barcode.BarcodeScanner
import com.google.mlkit.vision.barcode.BarcodeScanning
import com.google.mlkit.vision.barcode.common.Barcode
import dagger.hilt.android.testing.BindValue
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import org.junit.After
import org.junit.Ignore
import org.junit.Rule
import org.junit.Test
import org.junit.rules.ExternalResource
import org.junit.runner.RunWith
import org.mockito.MockedStatic
import org.mockito.Mockito
import org.mockito.kotlin.*
import javax.inject.Inject

@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class TransactionScenarioTest {


    @BindValue
    lateinit var userRepo: UserRepository

    @Inject
    lateinit var jsonService: JsonService
    private lateinit var barcodeScanningMock: MockedStatic<BarcodeScanning>

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

    private val user = User("Jon Smith", 12345, "passport")
    private val wallet = Wallet()

    // A null value will be seen as no result
    private var qrCodeContent: String? = null

    fun setup() {
        hiltRule.inject()

        qrCodeContent = null

        // Mock the barcode scanner such that it returns the qrCodeContent value
        val fakeBarcode = mock<Barcode> {
            on { rawValue } doAnswer { qrCodeContent }
        }

        val fakeScanner = mock<BarcodeScanner> {
            on { process(isA<MlImage>()) } doAnswer {
                Tasks.forResult(
                    if (qrCodeContent == null) emptyList() else listOf(fakeBarcode)
                )
            }
        }

        barcodeScanningMock = Mockito.mockStatic(BarcodeScanning::class.java)
        barcodeScanningMock.`when`<BarcodeScanner> { BarcodeScanning.getClient(any()) }
            .thenReturn(fakeScanner)

        // Mock user repo
        userRepo = mock {
            on { loggedInUser } doAnswer { user }

            on { isLoggedIn } doAnswer { true }

            on { wallet } doAnswer { wallet }

            on { isRegistered(any()) } doReturn true
        }
    }

    @After
    fun teardown() {
        barcodeScanningMock.close()
    }

    @Test
    @Ignore("Not working yet")
    fun receiverScenario() {
        qrCodeContent = jsonService.toJson(
            JsonType.CompleteTransactionType, CompleteTransaction(
                "ben",
                user.name,
                115.5F,
                33917321
            )
        )

        receive().perform(click())

        currentFragment().check(matches(withId(rxAmountFragmentId())))

        rxAmount().perform(replaceText("115.5"))
        rxAmountOk().perform(click())

        currentFragment().check(matches(withId(showFragmentId())))

        showOk().perform(click())

        currentFragment().check(matches(withId(scanFragmentId())))

        InstrumentationRegistry.getInstrumentation().waitForIdleSync()

        currentFragment().check(matches(withId(walletFragmentId())))
    }
}