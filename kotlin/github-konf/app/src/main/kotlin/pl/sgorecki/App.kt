package pl.sgorecki

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.ConfigSpec
import pl.sgorecki.k8s.Database

object ServerSpec : ConfigSpec() {
    val host by optional("0.0.0.0")
    val tcpPort by required<Int>()
}

fun main() {
// 1. Simple example
//    val config = Config {
//        addSpec(ServerSpec)
//    }.from.yaml.string("""
//       server:
//         host: 0.0.0.0
//         tcp_port: 8080
//    """.trimIndent())
//
//    println("Server ${config[ServerSpec.host]} running on port ${config[ServerSpec.tcpPort]}")

// 2. Configuration of the database
    val config = Config {
        addSpec(Database)
    }

    // for Production
    config[Database.writerHost] = "aurora.prod.rds.amazonaws.com"
    config[Database.readerHost] = "aurora-ro.prod.rds.amazonaws.com"
    config[Database.name] = "github-konf"
    config[Database.primaryUrl] = Database.primaryUrlFor(config)
    config[Database.readerUrl] = Database.readerUrlFor(config)

    println("Database config: $config")


// 3. Full App configuration

    println("Staging config: ${Staging.config}")
}