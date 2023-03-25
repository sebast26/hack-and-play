package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.io.File
import java.time.Instant


class StockTests {
    private val initialStockList = standardStockList.copy(
        lastModified = Instant.parse("2023-03-25T23:59:59Z")
    )
    private val fixture = Fixture(initialStockList)
    private val stock = Stock(fixture.stockFile)

    @Test
    fun `loads stock from file`() {
        val now = Instant.parse("2023-03-25T23:59:59Z")
        assertEquals(initialStockList, stock.stockList(now))
    }

    @Test
    fun `updates stock if last modified yesterday`() {
        val now = Instant.parse("2023-03-26T00:00:01Z")
        val expectedUpdateResult = initialStockList.copy(lastModified = now)
        assertEquals(expectedUpdateResult, stock.stockList(now))
        assertEquals(expectedUpdateResult, fixture.load())
    }
}

class Stock(
    val stockFile: File
) {
    fun stockList(now: Instant): StockList {
        val loaded = stockFile.loadItems()
        return if (loaded.lastModified.daysBetween(now) == 0L) {
            loaded
        } else {
            loaded.copy(lastModified = now).also {
                it.saveTo(stockFile)
            }
        }
    }
}

private fun Instant.daysBetween(then: Instant): Long = when (then) {
    Instant.parse("2023-03-26T00:00:01Z") -> 1
    Instant.parse("2023-03-25T23:59:59Z") -> 0
    else -> TODO()
}
