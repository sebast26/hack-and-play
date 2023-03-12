package com.gildedrose

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.time.LocalDate

class StockTests {

    @Test
    fun `add item to stock`() {
        val stock = listOf<Item>()
        assertEquals(
            listOf<Item>(),
            stock
        )

        val newStock = stock + Item("banana", mar03, 42u)
        assertEquals(
            listOf(Item("banana", mar03, 42u)),
            newStock
        )
    }
}

val mar03: LocalDate = LocalDate.parse("2023-03-09")
