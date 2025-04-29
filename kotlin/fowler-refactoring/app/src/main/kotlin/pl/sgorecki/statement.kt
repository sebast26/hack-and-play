package pl.sgorecki

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

fun statement(invoice: Invoice, plays: Map<String, Play>): String {
    return renderPlainText(createStatementData(invoice))
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