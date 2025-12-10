package pl.sgorecki.k8s

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.ConfigSpec
import org.http4k.config.Host
import org.http4k.config.Port
import pl.sgorecki.storage.databaseUrlFor

object Database : ConfigSpec() {
    val name by required<String>()

    val writerHost by required<String>()
    private val writerPort by optional(default = 5432)

    val readerHost by required<String>()
    private val readerPort by optional(default = 5432)

    val primaryUrl by required<String>()
    val readerUrl by required<String>()

    fun primaryUrlFor(config: Config) = databaseUrlFor(
        databaseHost = Host(config[writerHost]),
        databasePort = Port(config[writerPort]),
        databaseName = config[name]
    ).toString()

    fun readerUrlFor(config: Config) = databaseUrlFor(
        databaseHost = Host(config[readerHost]),
        databasePort = Port(config[readerPort]),
        databaseName = config[name]
    ).toString()
}