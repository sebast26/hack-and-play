package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

data class StatementData(
    val customer: String,
    val performances: List<EnrichedPerformance>,
    val totalAmount: Int,
    val totalVolumeCredits: Int,
)

data class EnrichedPerformance(
    val performance: Performance,
    val play: Play,
    val amount: Int,
    val volumeCredits: Int
)

fun playFor(performance: Performance) = plays[performance.playId]!!

fun usd(number: Int) = NumberFormat.getCurrencyInstance(Locale.US).format(number / 100)

fun amountFor(perf: Performance): Int {
    var result = 0
    when (playFor(perf).type) {
        PlayType.TRAGEDY -> {
            result = 40000
            if (perf.audience > 30) {
                result += 1000 * (perf.audience - 30)
            }
        }

        COMEDY -> {
            result = 30000
            if (perf.audience > 20) {
                result += 10000 + 500 * (perf.audience - 20)
            }
            result += 300 * perf.audience
        }

        else -> error("Unknown play type: ${playFor(perf).type}")
    }
    return result
}

fun volumeCreditsFor(performance: Performance): Int {
    var result = 0
    result += max(performance.audience - 30, 0)
    if (COMEDY == playFor(performance).type) result += floor(performance.audience.toDouble() / 5).toInt()
    return result
}

fun totalAmount(performances: List<EnrichedPerformance>): Int {
    var result = 0
    for (perf in performances) {
        result += perf.amount
    }
    return result
}

fun totalVolumeCredits(performances: List<EnrichedPerformance>): Int {
    var result = 0
    for (perf in performances) {
        result += perf.volumeCredits
    }
    return result
}

fun enrichPerformance(performance: Performance) = EnrichedPerformance(
    performance = performance,
    play = playFor(performance),
    amount = amountFor(performance),
    volumeCredits = volumeCreditsFor(performance),
)

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    val performances = invoice.performances.map { it -> enrichPerformance(it) }
    val statementData = StatementData(
        customer = invoice.customer,
        performances = performances,
        totalAmount = totalAmount(performances),
        totalVolumeCredits = totalVolumeCredits(performances),
    )

    return renderPlainText(statementData)
}

fun renderPlainText(data: StatementData): String {
    var result = "Statement for ${data.customer}\n"
    for (perf in data.performances) {
        result += "    ${perf.play.name}: ${usd(perf.amount)} (${perf.performance.audience} seats)\n"
    }
    result += "Amount owed is ${usd(data.totalAmount)}\n"
    result += "You earned ${data.totalVolumeCredits} credits\n"
    return result
}