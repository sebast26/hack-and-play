package pl.sgorecki.storage

import org.http4k.config.Host
import org.http4k.config.Port
import org.http4k.core.Uri

fun databaseUrlFor(
    databaseHost: Host,
    databasePort: Port,
    databaseName: String,
) = Uri.of("jdbc:postgresql://${databaseHost.value}:${databasePort.value}/$databaseName")