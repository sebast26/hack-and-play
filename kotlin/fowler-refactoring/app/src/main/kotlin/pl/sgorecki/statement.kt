package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    fun playFor(performance: Performance): Play = plays[performance.playId]!!

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
        var volumeCredits = 0
        volumeCredits += max(performance.audience - 30, 0)
        if (COMEDY == playFor(performance).type) volumeCredits += floor(performance.audience.toDouble() / 5).toInt()
        return volumeCredits
    }

    var totalAmount = 0
    var volumeCredits = 0
    var result = "Statement for ${invoice.customer}\n"
    val format: (Int) -> String = { num -> NumberFormat.getCurrencyInstance(Locale.US).format(num) }

    for (perf in invoice.performances) {
        volumeCredits += volumeCreditsFor(perf)

        // print line for this order
        result += "    ${playFor(perf).name}: ${format(amountFor(perf) / 100)} (${perf.audience} seats)\n"
        totalAmount += amountFor(perf)
    }

    result += "Amount owed is ${format(totalAmount / 100)}\n"
    result += "You earned $volumeCredits credits\n"
    return result
}