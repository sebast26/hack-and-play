package com.gildedrose

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.databind.SerializationFeature
import com.fasterxml.jackson.datatype.jdk8.Jdk8Module
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule
import com.fasterxml.jackson.module.paramnames.ParameterNamesModule
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.time.Instant

class AnalyticsTests {
    @Test
    fun `outputs json of the events`() {
        val logged = mutableListOf<String>()
        val now = Instant.parse("2023-04-23T13:31:53.298304Z")
        val analytics = LoggingAnalytics(
            logger = logged::add,
            clock = { now }
        )

        assertEquals(0, logged.size)
        analytics(TestEvent("banana"))
        assertEquals(
            listOf("""{"timestamp":"2023-04-23T13:31:53.298304Z","event":{"value":"banana","eventName":"TestEvent"}}"""),
            logged
        )
    }
}

data class TestEvent(val value: String) : AnalyticsEvent


