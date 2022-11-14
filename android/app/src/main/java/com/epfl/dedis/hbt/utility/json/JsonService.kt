package com.epfl.dedis.hbt.utility.json

import android.util.Log
import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import com.networknt.schema.*
import java.net.URI
import javax.inject.Inject
import javax.inject.Singleton
import kotlin.reflect.KClass

@Singleton
class JsonService @Inject constructor(private val objectMapper: ObjectMapper) {

    companion object {
        private val TAG = JsonService::class.simpleName
    }

    private var schemas: Map<JsonType<*>, JsonSchema>? = null

    fun loadSchemas() {
        Log.i(TAG, "Loading json schemas")
        val config = SchemaValidatorsConfig().apply {
            isHandleNullableField = false
        }

        val factory = JsonSchemaFactory.getInstance(SpecVersion.VersionFlag.V7)
        // Create the schema map by associating each schema to its generated validator
        schemas = JsonType.types().associateWith {
            factory.getSchema(URI.create("resource:/" + it.schemaPath), config)
                // Preload the schema now such that it isn't done later
                .apply { preloadJsonSchema() }
        }

        Log.i(TAG, "Schemas loaded successfully")
    }

    fun <T : Any> fromJson(jsonType: JsonType<T>, json: String): T =
        fromJson(jsonType, json, jsonType.type)

    fun <T1 : Any, T2 : T1> fromJson(jsonType: JsonType<T1>, json: String, type: KClass<T2>): T2 {
        val node = objectMapper.readTree(json)
        validate(jsonType, node)
        return objectMapper.treeToValue(node, type.java)
    }

    fun <T : Any> toJson(jsonType: JsonType<out T>, obj: T): String {
        val node: JsonNode = objectMapper.valueToTree(obj)
        validate(jsonType, node)
        return objectMapper.writeValueAsString(node)
    }

    private fun validate(jsonType: JsonType<*>, node: JsonNode) {
        val errors = getSchemas()[jsonType]!!.validate(node)

        if (errors.isNotEmpty()) {
            throw JsonSchemaException("ValidationMessage errors : $errors")
        }
    }

    private fun getSchemas(): Map<JsonType<*>, JsonSchema> {
        if (schemas == null) loadSchemas()
        return schemas!!
    }

}