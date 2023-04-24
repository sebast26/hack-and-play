package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class UpdatingTests {

    @Test fun `item decrease in quality one per day`() {
        assertEquals(
            Item("banana", mar03, 41),
            Item("banana", mar03, 42).updatedBy(days = 1, on = mar03)
        )
        assertEquals(
            Item("banana", mar03, 42),
            Item("banana", mar03, 42).updatedBy(days = 0, on = mar03)
        )
        assertEquals(
            Item("banana", mar03, 40),
            Item("banana", mar03, 42).updatedBy(days = 2, on = mar03)
        )
    }

    @Test fun `quality doesn't become negative`() {
        assertEquals(
            Item("banana", mar03, 0),
            Item("banana", mar03, 0).updatedBy(days = 1, on = mar03)
        )
        assertEquals(
            Item("banana", mar03, 0),
            Item("banana", mar03, 1).updatedBy(days = 2, on = mar03)
        )
    }

    @Test fun `item decrease in quality two per day after sell by date`() {
        assertEquals(
            Item("banana", mar03, 40),
            Item("banana", mar03, 42).updatedBy(days = 1, on = mar03.plusDays(1))
        )
        assertEquals(
            Item("banana", mar03, 42),
            Item("banana", mar03, 42).updatedBy(days = 0, on = mar03.plusDays(1))
        )
        assertEquals(
            Item("banana", mar03, 38),
            Item("banana", mar03, 42).updatedBy(days = 2, on = mar03.plusDays(2))
        )
        assertEquals(
            Item("banana", mar03, 39),
            Item("banana", mar03, 42).updatedBy(days = 2, on = mar03.plusDays(1))
        )
    }

    @Test fun `items with no sellBy don't change quality`() {
        assertEquals(
            Item("banana", null, 42),
            Item("banana", null, 42).updatedBy(days = 1, on = mar03)
        )
    }

    @Test fun `Aged Brie increases in quality by one every day until its sell by date`() {
        assertEquals(
            Item("Aged Brie", mar03, 43),
            Item("Aged Brie", mar03, 42).updatedBy(days = 1, on = mar03)
        )
    }

    @Test fun `Aged Brie increases in quality by two every day after its sell by date`() {
        assertEquals(
            Item("Aged Brie", mar03, 44),
            Item("Aged Brie", mar03, 42).updatedBy(days = 1, on = mar03.plusDays(1))
        )
    }

    @Test fun `Aged Brie doesn't get better than 50`() {
        assertEquals(
            Item("Aged Brie", mar03, 50),
            Item("Aged Brie", mar03, 50).updatedBy(days = 1, on = mar03)
        )
        assertEquals(
            Item("Aged Brie", mar03, 50),
            Item("Aged Brie", mar03, 49).updatedBy(days = 1, on = mar03.plusDays(1))
        )
    }
}