package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    fun amountFor(perf: Performance, play: Play): Int {
        var result = 0
        when (play.type) {
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

            else -> error("Unknown play type: ${play.type}")
        }
        return result
    }

    var totalAmount = 0
    var volumeCredits = 0
    var result = "Statement for ${invoice.customer}\n"
    val format: (Int) -> String = { num -> NumberFormat.getCurrencyInstance(Locale.US).format(num) }

    for (perf in invoice.performances) {
        val play = plays[perf.playId]
        val thisAmount = amountFor(perf, play!!)

        volumeCredits += max(perf.audience - 30, 0)
        // add extra credits for every ten comedy attendees
        if (COMEDY == play?.type) volumeCredits += floor(perf.audience.toDouble() / 5).toInt()

        // print line for this order
        result += "    ${play?.name}: ${format(thisAmount / 100)} (${perf.audience} seats)\n"
        totalAmount += thisAmount
    }

    result += "Amount owed is ${format(totalAmount / 100)}\n"
    result += "You earned $volumeCredits credits\n"
    return result
}