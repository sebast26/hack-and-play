package com.example

import org.http4k.client.JavaHttpClient
import org.http4k.core.Method.GET
import org.http4k.core.Request
import org.http4k.server.Http4kServer
import org.http4k.server.Undertow
import org.http4k.server.asServer

fun main() {
    val app: Http4kServer = Pokemon4k(RealPokeCoApi()).asServer(Undertow(0)).start()
    val http = JavaHttpClient()
    http(Request(GET, "http://localhost:${app.port()}/b"))
}
