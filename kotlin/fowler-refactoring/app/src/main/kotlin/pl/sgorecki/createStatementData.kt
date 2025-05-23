package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import pl.sgorecki.PlayType.TRAGEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

fun playFor(performance: Performance) = plays[performance.playId]!!

fun usd(number: Int) = NumberFormat.getCurrencyInstance(Locale.US).format(number / 100)

fun totalAmount(performances: List<EnrichedPerformance>) = performances.sumOf { it.amount }

fun totalVolumeCredits(performances: List<EnrichedPerformance>) = performances.sumOf { it.volumeCredits }

abstract class PerformanceCalculator(
    private val performance: Performance,
    val play: Play
) {
    abstract fun amount(): Int

    open fun volumeCredits(): Int {
        return max(performance.audience - 30, 0)
    }
}

class TragedyCalculator(val performance: Performance, play: Play) : PerformanceCalculator(performance, play) {
    override fun amount(): Int {
        var result = 40000
        if (performance.audience > 30) {
            result += 1000 * (performance.audience - 30)
        }
        return result
    }
}

class ComedyCalculator(val performance: Performance, play: Play) : PerformanceCalculator(performance, play) {
    override fun amount(): Int {
        var result = 30000
        if (performance.audience > 20) {
            result += 10000 + 500 * (performance.audience - 20)
        }
        result += 300 * performance.audience
        return result
    }

    override fun volumeCredits(): Int {
        return super.volumeCredits() + floor(performance.audience.toDouble() / 5).toInt()
    }
}

fun enrichPerformance(performance: Performance): EnrichedPerformance {
    val calculator = createPerformanceCalculator(performance, playFor(performance))
    return EnrichedPerformance(
        performance = performance,
        play = calculator.play,
        amount = calculator.amount(),
        volumeCredits = calculator.volumeCredits(),
    )
}

private fun createPerformanceCalculator(performance: Performance, play: Play) = when (play.type) {
    TRAGEDY -> TragedyCalculator(performance, play)
    COMEDY -> ComedyCalculator(performance, play)
    else -> error("Unknown play type: ${play.type}")
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