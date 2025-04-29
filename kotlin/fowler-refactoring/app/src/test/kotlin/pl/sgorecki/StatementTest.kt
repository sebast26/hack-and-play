package pl.sgorecki

import org.junit.jupiter.api.Test

class StatementTest {
    @Test
    fun appHasAGreeting() {
        println(statement(invoices[0], plays))
    }
}

private val expected = """
    Statement for BigCo
        Hamlet: $650.00 (55 seats)
        As You Like It: $580.00 (35 seats)
        Othello: $500.00 (40 seats)
    Amount owed is $1,730.00
    You earned 47 credits
""".trimIndent()