package com.epfl.dedis.hbt.utility.json

import androidx.test.ext.junit.runners.AndroidJUnit4
import com.epfl.dedis.hbt.data.model.CompleteTransaction
import com.epfl.dedis.hbt.di.JsonModule
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.registerKotlinModule
import com.networknt.schema.JsonSchemaException
import org.hamcrest.Description
import org.hamcrest.Matcher
import org.hamcrest.MatcherAssert.assertThat
import org.hamcrest.TypeSafeMatcher
import org.junit.Assert.assertThrows
import org.junit.Test
import org.junit.runner.RunWith
import org.hamcrest.CoreMatchers.`is` as eq

@RunWith(AndroidJUnit4::class)
class JsonServiceTest {

    private val mapper = ObjectMapper().registerKotlinModule()

    companion object {
        private val validTransaction = CompleteTransaction("Ben", "Marc", 12.5f, 234)
        private const val validTransactionJson =
            "{\"datetime\": 234,\"source\": \"Ben\",\"destination\": \"Marc\",\"amount\": 12.5}"

        fun jsonEq(
            jsonString: String,
            objectMapper: ObjectMapper = JsonModule.provideObjectMapper()
        ): Matcher<String> {
            return object : TypeSafeMatcher<String>(String::class.java) {
                override fun describeTo(description: Description) {
                    description.appendText("json equivalent to '$jsonString'")
                }

                override fun matchesSafely(item: String): Boolean {
                    val expected = objectMapper.readTree(jsonString)
                    val got = objectMapper.readTree(item)

                    return expected == got
                }
            }
        }
    }

    @Test
    fun validTransactionIsDeserialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        val transaction: CompleteTransaction =
            service.fromJson(
                JsonType.COMPLETE_TRANSACTION,
                validTransactionJson
            )
        assertThat(transaction, eq(validTransaction))
    }

    @Test
    fun validTransactionIsSerialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        assertThat(
            service.toJson(JsonType.COMPLETE_TRANSACTION, validTransaction),
            jsonEq(validTransactionJson)
        )
    }

    @Test
    fun invalidTransactionIsNotDeserialized() {
        val service = JsonService(mapper)
        service.loadSchemas()

        assertThrows(JsonSchemaException::class.java) {
            service.fromJson(
                JsonType.COMPLETE_TRANSACTION,
                "{\"datetime\": -6,\"source\": \"Source\",\"destination\": \"Dest\",\"amount\": 12.5}"
            )
        }
    }
}