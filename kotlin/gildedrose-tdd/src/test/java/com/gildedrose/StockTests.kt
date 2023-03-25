package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.time.Instant
import java.time.ZoneId
import java.util.concurrent.Callable
import java.util.concurrent.CyclicBarrier
import java.util.concurrent.Executors


class StockTests {
    private val initialStockList = StockList(
        lastModified = Instant.parse("2023-03-24T23:59:59Z"),
        items = listOf(
            Item("banana", mar03.minusDays(1), 42u),
            Item("kumquat", mar03.plusDays(1), 101u)
        )
    )
    private val fixture = Fixture(initialStockList, now = initialStockList.lastModified)
    private val stock = Stock(
        stockFile = fixture.stockFile,
        zoneID = ZoneId.of("Europe/London"),
        update = ::updateItems
    )

    @Test
    fun `loads stock from file`() {
        val now = Instant.parse("2023-03-24T23:59:59Z")
        assertEquals(initialStockList, stock.stockList(now))
    }

    @Test
    fun `updates stock if last modified yesterday`() {
        val now = Instant.parse("2023-03-25T00:00:01Z")
        val expectedUpdateResult = StockList(
            lastModified = now,
            items = listOf(
                Item("banana", mar03.minusDays(1), 41u),
                Item("kumquat", mar03.plusDays(1), 100u)
            )
        )
        assertEquals(expectedUpdateResult, stock.stockList(now))
        assertEquals(expectedUpdateResult, fixture.load())
    }

    @Test
    fun `updates stock by two days if last modified the day before yesterday`() {
        val now = Instant.parse("2023-03-26T00:00:01Z")
        val expectedUpdateResult = StockList(
            lastModified = now,
            items = listOf(
                Item("banana", mar03.minusDays(1), 40u),
                Item("kumquat", mar03.plusDays(1), 99u)
            )
        )
        assertEquals(expectedUpdateResult, stock.stockList(now))
        assertEquals(expectedUpdateResult, fixture.load())
    }

    @Test
    fun `does not update stock if modified tomorrow`() {
        val now = Instant.parse("2023-03-23T00:00:01Z")
        assertEquals(initialStockList, stock.stockList(now))
        assertEquals(initialStockList, fixture.load())
    }

    @Test
    fun `parallel execution`() {
        val count = 8
        val executor = Executors.newFixedThreadPool(count)
        val barrier = CyclicBarrier(count)
        val futures = executor.invokeAll(
            (1..count).map {
                Callable() {
                    barrier.await()
                    `updates stock if last modified yesterday`()
                }
            }
        )
        futures.forEach { it.get() }
    }

    @Test
    fun `sanity check`() {
        for (i in 1..10) {
            `parallel execution`()
        }
    }
}

private fun updateItems(items: List<Item>, days: Int) = items.map { it.copy(quality = it.quality - days.toUInt()) }
