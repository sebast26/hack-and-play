package com.gildedrose

import org.http4k.core.HttpHandler
import org.http4k.core.Response
import org.http4k.core.Status
import org.http4k.template.HandlebarsTemplates
import org.http4k.template.ViewModel
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit

private val dateFormat: DateTimeFormatter = DateTimeFormatter.ofPattern("d LLLL yyyy")
private val handlebars = HandlebarsTemplates().HotReload("src/main/java")

fun listHandler(
    clock: () -> Instant,
    zoneId: ZoneId,
    listing: (Instant) -> StockList
): HttpHandler = { _ ->
    val now = clock()
    val today = LocalDate.ofInstant(now, zoneId)
    val stockList = listing(now)
    Response(Status.OK).body(handlebars(
        StockListViewModel(
            now = dateFormat.format(today),
            items = stockList.map { it.toMap(today) }
        )
    ))
}

private data class StockListViewModel(
    val now: String,
    val items: List<Map<String, String>>
) : ViewModel

private fun Item.toMap(now: LocalDate): Map<String, String> = mapOf(
    "name" to name,
    "sellByDate" to this.sellByDateString(),
    "sellByDays" to this.daysUntilSellBy(now).toString(),
    "quality" to this.quality.toString()
)

private fun Item.sellByDateString(): String = when (this) {
    is DatedItem -> dateFormat.format(sellByDate)
    else -> ""
}

private fun Item.daysUntilSellBy(now: LocalDate): Long =
    when (this) {
        is DatedItem -> ChronoUnit.DAYS.between(now, sellByDate)
        else -> 0
    }