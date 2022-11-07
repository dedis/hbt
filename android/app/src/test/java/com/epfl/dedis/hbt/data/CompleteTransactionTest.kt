package com.epfl.dedis.hbt.data

import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.model.CompleteTransaction
import com.epfl.dedis.hbt.di.JsonModule
import com.epfl.dedis.hbt.utility.json.JsonService
import com.epfl.dedis.hbt.utility.json.JsonServiceTest.Companion.jsonEq
import com.epfl.dedis.hbt.utility.json.JsonType
import org.hamcrest.MatcherAssert.assertThat
import org.junit.Test
import org.junit.runner.RunWith
import org.hamcrest.CoreMatchers.`is` as eq

@RunWith(AndroidJUnit4::class)
class CompleteTransactionTest {

    private val transaction = CompleteTransaction("source", "dest", 10.5F, 104320)
    private val jsonRepresentation =
        "{\"datetime\": 104320,\"source\": \"source\",\"destination\": \"dest\",\"amount\": 10.5}"

    private val jsonService = JsonService(JsonModule.provideObjectMapper())
        .apply { loadSchemas() }

    @Test
    fun serializationTest() {
        assertThat(
            jsonService.toJson(JsonType.COMPLETE_TRANSACTION, transaction),
            jsonEq(jsonRepresentation)
        )
    }

    @Test
    fun deserializationTest() {
        assertThat(
            jsonService.fromJson(
                JsonType.COMPLETE_TRANSACTION,
                jsonRepresentation,
                CompleteTransaction::class
            ),
            eq(transaction)
        )
    }
}