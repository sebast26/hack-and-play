package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

fun playFor(performance: Performance) = plays[performance.playId]!!

fun usd(number: Int) = NumberFormat.getCurrencyInstance(Locale.US).format(number / 100)

fun volumeCreditsFor(performance: Performance): Int {
    var result = 0
    result += max(performance.audience - 30, 0)
    if (COMEDY == playFor(performance).type) result += floor(performance.audience.toDouble() / 5).toInt()
    return result
}

fun totalAmount(performances: List<EnrichedPerformance>) = performances.sumOf { it.amount }

fun totalVolumeCredits(performances: List<EnrichedPerformance>) = performances.sumOf { it.volumeCredits }

class PerformanceCalculator(
    private val performance: Performance,
    val play: Play
) {
    fun amount(): Int {
        var result = 0
        when (play.type) {
            PlayType.TRAGEDY -> {
                result = 40000
                if (performance.audience > 30) {
                    result += 1000 * (performance.audience - 30)
                }
            }

            COMEDY -> {
                result = 30000
                if (performance.audience > 20) {
                    result += 10000 + 500 * (performance.audience - 20)
                }
                result += 300 * performance.audience
            }

            else -> error("Unknown play type: ${play.type}")
        }
        return result
    }
}

fun enrichPerformance(performance: Performance): EnrichedPerformance {
    val calculator = PerformanceCalculator(performance, playFor(performance))
    return EnrichedPerformance(
        performance = performance,
        play = calculator.play,
        amount = calculator.amount(),
        volumeCredits = volumeCreditsFor(performance),
    )
}

fun createStatementData(invoice: Invoice): StatementData {
    val performances = invoice.performances.map { it -> enrichPerformance(it) }
    return StatementData(
        customer = invoice.customer,
        performances = performances,
        totalAmount = totalAmount(performances),
        totalVolumeCredits = totalVolumeCredits(performances),
    )
}