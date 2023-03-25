package com.gildedrose

import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNotEquals
import org.junit.jupiter.api.Disabled
import org.junit.jupiter.api.Test
import java.time.Instant
import java.time.LocalDate

class UpdateStockTest {

    private val stockList = standardStockList.copy(lastModified = Instant.parse("2023-03-25T12:00:00Z"))

    @Test
    fun `doesn't update when lastModified in today`() {
        val sameDayAsLastModified = LocalDate.parse("2023-03-25")
        with(Fixture(standardStockList, now = sameDayAsLastModified)) {
            assertEquals(Status.OK, routes(Request(Method.GET, "/")).status)
            assertEquals(stockList, load())
        }
    }

    @Disabled("WIP")
    @Test
    fun `does update when lastModified was yesterday`() {
        val nextDayFromLastModified = LocalDate.parse("2023-03-26")
        with(Fixture(standardStockList, now = nextDayFromLastModified)) {
            assertEquals(Status.OK, routes(Request(Method.GET, "/")).status)
            assertNotEquals(stockList, load())
        }
    }
}

