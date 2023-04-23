package com.gildedrose

import HttpEvent
import UncaughtExceptionEvent
import org.http4k.core.Method.GET
import org.http4k.core.Request
import org.http4k.core.Response
import org.http4k.core.Status.Companion.INTERNAL_SERVER_ERROR
import org.http4k.core.Status.Companion.OK
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.time.Instant

class ApplicationEventsTests {

    private val fixture = Fixture(
        StockList(
            lastModified = Instant.parse("2023-04-23T12:00:00Z"),
            items = emptyList()
        ), now = Instant.now()
    )

    @Test
    fun `uncaught exceptions raise an event`() {
        with(fixture) {
            assertEquals(0, events.size)
            val response = routes(Request(GET, "/errors"))
            assertEquals(
                Response(INTERNAL_SERVER_ERROR).body("Something went wrong, sorry"),
                response

            )
            assertEquals(UncaughtExceptionEvent::class, events[0]::class)
            assertEquals(HttpEvent::class, events[1]::class)
        }
    }

    @Test
    fun `every request raises an event`() {
        with(fixture) {
            assertEquals(0, events.size)
            val response = routes(Request(GET, "/"))
            assertEquals(OK, response.status)
            assertEquals(1, events.size)
            assertEquals(HttpEvent::class, events.only()::class)
        }
    }
}

private fun <E> Collection<E>.only(): E =
    when {
        this.size != 1 -> error("Expected one item, got $this")
        else -> this.first()
    }
