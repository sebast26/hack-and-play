package com.gildedrose

import org.http4k.core.Method
import org.http4k.core.Response
import org.http4k.core.Status
import org.http4k.routing.bind
import org.http4k.routing.routes
import org.http4k.server.Undertow
import org.http4k.server.asServer
import org.http4k.template.HandlebarsTemplates
import org.http4k.template.ViewModel
import org.intellij.lang.annotations.Language
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit

val dateFormat: DateTimeFormatter = DateTimeFormatter.ofPattern("d LLLL yyyy")

class Server(
    stock: List<Item>,
    clock: () -> LocalDate = LocalDate::now
) {
    val routes = routes(
        "/" bind Method.GET to { _ ->
            val now = clock()
            Response(Status.OK).body(handlebars(
                StockListViewModel(
                    now = dateFormat.format(now),
                    items = stock.map { it.toMap(now) }
                )
            ))
        }
    )

    private val handlebars = HandlebarsTemplates().HotReload("src/main/java")

    fun start() {
        val http4kServer = routes.asServer(Undertow(8080))
        http4kServer.start()
    }
}

private data class StockListViewModel(
    val now: String,
    val items: List<Map<String, String>>
) : ViewModel

private fun Item.toMap(now: LocalDate): Map<String, String> = mapOf(
    "name" to name,
    "sellByDate" to dateFormat.format(sellByDate),
    "sellByDays" to this.daysUntilSellBy(now).toString(),
    "quality" to this.quality.toString()
)

@Language("HTML")
private val templateSource = """
    <html lang="en">
    <body>
    <table>
    <h1>{{this.now}}</h1>
    <tr>
        <th>Name</th>
        <th>Sell By Date</th>
        <th>Sell By Days</th>
        <th>Quality</th>
    </tr>
    {{#each this.items}}<tr>
        <td>{{this.name}}</td>
        <td>{{this.sellByDate}}</td>
        <td>{{this.sellByDays}}</td>
        <td>{{this.quality}}</td>
    </tr>
    {{/each}}
    </table>
    </body>
    </html>
""".trimIndent()

private fun Item.daysUntilSellBy(now: LocalDate): Long =
    ChronoUnit.DAYS.between(now, sellByDate)