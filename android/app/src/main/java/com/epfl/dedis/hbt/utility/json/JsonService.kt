package com.epfl.dedis.hbt.utility.json

import android.util.Log
import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import com.networknt.schema.*
import java.net.URI
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class JsonService @Inject constructor(private val objectMapper: ObjectMapper) {

    companion object {
        private val TAG = JsonService::class.simpleName
    }

    private lateinit var schemas: Map<JsonType, JsonSchema>

    fun loadSchemas() {
        Log.i(TAG, "Loading json schemas")
        val config = SchemaValidatorsConfig().apply {
            isHandleNullableField = false
        }

        val factory = JsonSchemaFactory.getInstance(SpecVersion.VersionFlag.V7)
        // Create the schema map by associating each schema to its generated validator
        schemas = JsonType.values().associateWith {
            factory.getSchema(URI.create("resource:/" + it.schemaPath), config)
                // Preload the schema now such that it isn't done later
                .apply { preloadJsonSchema() }
        }

        Log.i(TAG, "Schemas loaded successfully")
    }

    fun <T> fromJson(jsonType: JsonType, json: String, type: Class<T>): T {
        val node = objectMapper.readTree(json)
        validate(jsonType, node)
        return objectMapper.treeToValue(node, type)
    }

    fun <T> toJson(jsonType: JsonType, obj: T): String {
        val node: JsonNode = objectMapper.valueToTree(obj)
        validate(jsonType, node)
        return objectMapper.writeValueAsString(node)
    }

    private fun validate(jsonType: JsonType, node: JsonNode) {
        val errors = schemas[jsonType]!!.validate(node)

        if (errors.isNotEmpty()) {
            throw JsonSchemaException("ValidationMessage errors : $errors")
        }
    }

}