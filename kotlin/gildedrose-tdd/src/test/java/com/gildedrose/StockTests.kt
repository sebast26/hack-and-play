package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.io.File
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId


class StockTests {
    private val initialStockList = standardStockList.copy(
        lastModified = Instant.parse("2023-03-24T23:59:59Z")
    )
    private val fixture = Fixture(initialStockList)
    private val stock = Stock(fixture.stockFile, ZoneId.of("Europe/London"))

    @Test
    fun `loads stock from file`() {
        val now = Instant.parse("2023-03-24T23:59:59Z")
        assertEquals(initialStockList, stock.stockList(now))
    }

    @Test
    fun `updates stock if last modified yesterday`() {
        val now = Instant.parse("2023-03-25T00:00:01Z")
        val expectedUpdateResult = initialStockList.copy(lastModified = now)
        assertEquals(expectedUpdateResult, stock.stockList(now))
        assertEquals(expectedUpdateResult, fixture.load())
    }
}

class Stock(
    private val stockFile: File,
    private val zoneID: ZoneId
) {
    fun stockList(now: Instant): StockList {
        val loaded = stockFile.loadItems()
        return if (loaded.lastModified.daysTo(now, zoneID) == 0L) {
            loaded
        } else {
            loaded.copy(lastModified = now).also {
                it.saveTo(stockFile)
            }
        }
    }
}

internal fun Instant.daysTo(that: Instant, zone: ZoneId): Long =
    LocalDate.ofInstant(that, zone).toEpochDay() - LocalDate.ofInstant(this, zone).toEpochDay()
