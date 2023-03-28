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
}