package com.example

import org.http4k.core.Body
import org.http4k.core.Method
import org.http4k.core.Response
import org.http4k.core.Status
import org.http4k.core.with
import org.http4k.format.Jackson.auto
import org.http4k.lens.Path
import org.http4k.routing.bind
import org.http4k.routing.routes

val nameLens = Path.of("name")
val greetingLens = Body.auto<Greeting>().toLens()

val helloRoute: String = "/hello/$nameLens"

val helloApiWithLenses = routes(
    helloRoute bind Method.GET to { request ->
        val name = nameLens(request)
        val greeting = Greeting("hello $name")
        Response(Status.OK).with(greetingLens of greeting)
    }
)