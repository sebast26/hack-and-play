package com.gildedrose

import java.time.Instant
import java.time.LocalDate

val mar03: LocalDate = LocalDate.parse("2023-03-09")
val standardStockList = StockList(
    Instant.now(),
    listOf(
        Item("banana", mar03.minusDays(1), 42u),
        Item("kumquat", mar03.plusDays(1), 101u)
    )
)