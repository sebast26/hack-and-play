package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class UpdatingTests {
    val items = listOf(Item("banana", mar03, 42))

    @Test fun `item decrease in quality one per day`() {
        assertEquals(
            listOf(Item("banana", mar03, 41)),
            updateItems(items, days = 1)
        )
        assertEquals(
            listOf(Item("banana", mar03, 42)),
            updateItems(items, days = 0)
        )
        assertEquals(
            listOf(Item("banana", mar03, 40)),
            updateItems(items, days = 2)
        )
    }

    @Test fun `quality doesn't become negative`() {
        assertEquals(
            listOf(Item("banana", mar03, 0)),
            updateItems(listOf(Item("banana", mar03, 0)), days = 1)
        )
        assertEquals(
            listOf(Item("banana", mar03, 0)),
            updateItems(listOf(Item("banana", mar03, 1)), days = 2)
        )
    }
}