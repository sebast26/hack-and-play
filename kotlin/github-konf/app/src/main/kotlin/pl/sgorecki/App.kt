package pl.sgorecki

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.ConfigSpec
import com.uchuhimo.konf.source.yaml

object ServerSpec : ConfigSpec() {
    val host by optional("0.0.0.0")
    val tcpPort by required<Int>()
}

fun main() {

    val config = Config {
        addSpec(ServerSpec)
    }.from.yaml.string("""
       server:
         host: 0.0.0.0
         tcp_port: 8080
    """.trimIndent())

    println("Server ${config[ServerSpec.host]} running on port ${config[ServerSpec.tcpPort]}")

}