package com.example

import org.http4k.core.Body
import org.http4k.core.HttpHandler
import org.http4k.core.Method
import org.http4k.core.Response
import org.http4k.core.Status.Companion.OK
import org.http4k.core.with
import org.http4k.format.Jackson.auto
import org.http4k.routing.bind
import org.http4k.routing.path
import org.http4k.routing.routes
import org.http4k.server.Undertow
import org.http4k.server.asServer

data class Greeting(val message: String)

val helloApi: HttpHandler = routes(
    "/hello/{name}" bind Method.GET to { req ->
        val name = req.path("name")
        val greeting = Greeting("hello $name")
        val lens = Body.auto<Greeting>().toLens()
        Response(OK).with(lens of greeting)
    }
)

fun main() {
    helloApi.asServer(Undertow(8001)).start()
}


