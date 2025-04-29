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
        var result = 0
        result += max(performance.audience - 30, 0)
        if (COMEDY == playFor(performance).type) result += floor(performance.audience.toDouble() / 5).toInt()
        return result
    }

    fun usd(number: Int): String {
        return NumberFormat.getCurrencyInstance(Locale.US).format(number / 100)
    }

    fun totalVolumeCredits(): Int {
        var volumeCredits = 0
        for (perf in invoice.performances) {
            volumeCredits += volumeCreditsFor(perf)
        }
        return volumeCredits
    }

    fun appleSouce(): Int {
        var totalAmount = 0
        for (perf in invoice.performances) {
            totalAmount += amountFor(perf)
        }
        return totalAmount
    }

    var result = "Statement for ${invoice.customer}\n"
    for (perf in invoice.performances) {
        // print line for this order
        result += "    ${playFor(perf).name}: ${usd(amountFor(perf))} (${perf.audience} seats)\n"
    }
    result += "Amount owed is ${usd(appleSouce())}\n"
    result += "You earned ${totalVolumeCredits()} credits\n"
    return result
}