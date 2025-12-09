package pl.sgorecki

import com.uchuhimo.konf.ConfigSpec

object ServerSpec : ConfigSpec() {
    val host by optional("0.0.0.0")
    val tcpPort by required<Int>()
}

fun main() {
    println(ServerSpec.prefix)
}