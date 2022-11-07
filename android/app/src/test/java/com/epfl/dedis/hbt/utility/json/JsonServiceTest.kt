package com.epfl.dedis.hbt.utility.json

import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.model.Transaction
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.registerKotlinModule
import com.networknt.schema.JsonSchemaException
import org.hamcrest.MatcherAssert.assertThat
import org.junit.Assert.assertThrows
import org.junit.Test
import org.junit.runner.RunWith
import org.hamcrest.CoreMatchers.`is` as eq

@RunWith(AndroidJUnit4::class)
class JsonServiceTest {

    private val mapper = ObjectMapper().registerKotlinModule()

    companion object {
        private val validTransaction = Transaction("Source", "Dest", 12.5f, 234)
        private const val validTransactionJson =
            "{\"datetime\": 234,\"source\": \"Source\",\"destination\": \"Dest\",\"amount\": 12.5}"
    }

    @Test
    fun validTransactionIsDeserialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        val transaction: Transaction =
            service.fromJson(
                JsonType.COMPLETE_TRANSACTION,
                validTransactionJson,
                Transaction::class.java
            )
        assertThat(transaction, eq(validTransaction))
    }

    @Test
    fun validTransactionIsSerialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        val json = service.toJson(JsonType.COMPLETE_TRANSACTION, validTransaction)
        assertThat(mapper.readTree(json), eq(mapper.readTree(validTransactionJson)))
    }

    @Test
    fun invalidTransactionIsNotDeserialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        assertThrows(JsonSchemaException::class.java) {
            service.fromJson(
                JsonType.COMPLETE_TRANSACTION,
                "{\"datetime\": -6,\"source\": \"Source\",\"destination\": \"Dest\",\"amount\": 12.5}",
                Transaction::class.java
            )
        }
    }
}