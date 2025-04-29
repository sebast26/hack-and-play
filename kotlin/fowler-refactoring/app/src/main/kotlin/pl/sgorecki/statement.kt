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

fun htmlStatement(invoice: Invoice, plays: Map<String, Play>): String {
    return renderHtml(createStatementData(invoice))
}

fun renderHtml(data: StatementData): String {
    var result = "<h1>Statement for ${data.customer}</h1>\n"
    result += "<table>\n"
    result += "<tr><th>play</th><th>seats</th><th>cost</th></tr>\n"
    for (perf in data.performances) {
        result += "    <tr><td>${perf.play.name}</td><td>${perf.performance.audience}</td>"
        result += "<td>${usd(perf.amount)}</td></tr>\n"
    }
    result += "</table>\n"
    result += "<p>Amount owed is <em>${usd(data.totalAmount)}</em></p>\n"
    result += "<p>You earned <em>${data.totalVolumeCredits}</em> credits</p>\n"
    return result
}