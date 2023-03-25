package com.gildedrose

import routesFor
import java.io.File
import java.nio.file.Files
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId

class Fixture(
    initialStockList: StockList,
    val now: Instant,
    val stockFile: File = Files.createTempFile("stock", ".tsv").toFile()
) {
    init {
        save(initialStockList)
    }

    val routes = routesFor(
        stockFile = stockFile,
        clock = { now }
    )

    fun save(stockList: StockList) {
        stockList.saveTo(stockFile)
    }

    fun load(): StockList {
        return stockFile.loadItems()
    }
}