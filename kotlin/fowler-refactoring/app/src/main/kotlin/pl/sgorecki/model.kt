package pl.sgorecki

import pl.sgorecki.PlayType.COMEDY
import pl.sgorecki.PlayType.TRAGEDY


enum class PlayType {
    TRAGEDY,
    COMEDY,
    OTHER
}

data class Play(val name: String, val type: PlayType)

data class Performance(val playId: String, val audience: Int)

data class Invoice(val customer: String, val performances: List<Performance>)

val plays = mapOf(
    "hamlet" to Play("Hamlet", TRAGEDY),
    "as-like" to Play("As You Like It", COMEDY),
    "othello" to Play("Othello", TRAGEDY)
)

val invoices = listOf(
    Invoice(
        customer = "BigCo",
        performances = listOf(
            Performance(playId = "hamlet", audience = 55),
            Performance(playId = "as-like", audience = 35),
            Performance(playId = "othello", audience = 40),
        ),
    )
)

