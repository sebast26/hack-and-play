package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

data class StatementData(
    val customer: String,
    val performances: List<EnrichedPerformance>
)

data class EnrichedPerformance(val performance: Performance, val play: Play)

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    fun playFor(performance: Performance): Play = plays[performance.playId]!!

    fun enrichPerformance(performance: Performance) = EnrichedPerformance(performance, playFor(performance))

    val statementData = StatementData(
        customer = invoice.customer,
        performances = invoice.performances.map { it -> enrichPerformance(it) }
    )

    return renderPlainText(statementData, plays)
}

fun renderPlainText(data: StatementData, plays: Map<String, Play>): String {
    fun playFor(performance: Performance): Play = plays[performance.playId]!!

    fun amountFor(perf: EnrichedPerformance): Int {
        var result = 0
        when (playFor(perf.performance).type) {
            PlayType.TRAGEDY -> {
                result = 40000
                if (perf.performance.audience > 30) {
                    result += 1000 * (perf.performance.audience - 30)
                }
            }

            COMEDY -> {
                result = 30000
                if (perf.performance.audience > 20) {
                    result += 10000 + 500 * (perf.performance.audience - 20)
                }
                result += 300 * perf.performance.audience
            }

            else -> error("Unknown play type: ${playFor(perf.performance).type}")
        }
        return result
    }

    fun volumeCreditsFor(performance: EnrichedPerformance): Int {
        var result = 0
        result += max(performance.performance.audience - 30, 0)
        if (COMEDY == playFor(performance.performance).type) result += floor(performance.performance.audience.toDouble() / 5).toInt()
        return result
    }

    fun usd(number: Int): String {
        return NumberFormat.getCurrencyInstance(Locale.US).format(number / 100)
    }

    fun totalVolumeCredits(): Int {
        var result = 0
        for (perf in data.performances) {
            result += volumeCreditsFor(perf)
        }
        return result
    }

    fun totalAmount(): Int {
        var result = 0
        for (perf in data.performances) {
            result += amountFor(perf)
        }
        return result
    }

    var result = "Statement for ${data.customer}\n"
    for (perf in data.performances) {
        result += "    ${perf.play.name}: ${usd(amountFor(perf))} (${perf.performance.audience} seats)\n"
    }
    result += "Amount owed is ${usd(totalAmount())}\n"
    result += "You earned ${totalVolumeCredits()} credits\n"
    return result
}