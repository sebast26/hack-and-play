package com.gildedrose

import org.junit.jupiter.api.Assertions
import org.junit.jupiter.api.Test
import java.time.LocalDate

class PrintingTest {

    private val now = LocalDate.parse("2023-03-09")

    @Test
    fun `print empty stock list`() {
        val stock = listOf<Item>()
        val expected = listOf("9 March 2023")

        Assertions.assertEquals(expected, stock.printout(now))
    }

    @Test
    fun `print non empty stock list`() {
        val stock = listOf(
            Item("banana", now.minusDays(1), 42u),
            Item("kumquat", now.plusDays(1), 101u)
        )
        val expected = listOf(
            "9 March 2023",
            "banana, -1, 42",
            "kumquat, 1, 101",
        )

        Assertions.assertEquals(expected, stock.printout(now))
    }
}