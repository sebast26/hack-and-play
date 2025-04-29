package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import java.text.NumberFormat
import java.util.*
import kotlin.math.floor
import kotlin.math.max

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    var totalAmount = 0
    var volumeCredits = 0
    var result = "Statement for ${invoice.customer}\n"
    val format: (Int) -> String = { num -> NumberFormat.getCurrencyInstance(Locale.US).format(num) }

    for (perf in invoice.performances) {
        val play = plays[perf.playId]
        var thisAmount = 0

        when (play?.type) {
            PlayType.TRAGEDY -> {
                thisAmount = 40000
                if (perf.audience > 30) {
                    thisAmount += 1000 * (perf.audience - 30)
                }
            }

            COMEDY -> {
                thisAmount = 30000
                if (perf.audience > 20) {
                    thisAmount += 10000 + 500 * (perf.audience - 20)
                }
                thisAmount += 300 * perf.audience
            }

            else -> error("Unknown play type: ${play?.type}")
        }

        volumeCredits += max(perf.audience - 30, 0)
        // add extra credits for every ten comedy attendees
        if (COMEDY == play.type) volumeCredits += floor(perf.audience.toDouble() / 5).toInt()

        // print line for this order
        result += "\t${play.name}: ${format(thisAmount / 100)} (${perf.audience} seats)\n"
        totalAmount += thisAmount
    }

    result += "Amount owed is ${format(totalAmount / 100)}\n"
    result += "You earned $volumeCredits credits\n"
    return result
}