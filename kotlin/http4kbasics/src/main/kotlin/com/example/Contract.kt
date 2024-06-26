package com.example

import org.http4k.contract.ContractRoute
import org.http4k.contract.contract
import org.http4k.contract.div
import org.http4k.contract.meta
import org.http4k.contract.openapi.ApiInfo
import org.http4k.contract.openapi.v3.OpenApi3
import org.http4k.contract.security.BearerAuthSecurity
import org.http4k.contract.ui.swaggerUiLite
import org.http4k.core.Method.GET
import org.http4k.core.Request
import org.http4k.core.Response
import org.http4k.core.Status.Companion.OK
import org.http4k.core.then
import org.http4k.core.with
import org.http4k.filter.ServerFilters
import org.http4k.routing.routes
import org.http4k.server.Undertow
import org.http4k.server.asServer

val sayHelloSpec = ("/hello" / nameLens) meta {
    operationId = "sayHello"
    summary = "Say hello"

    returning(OK, greetingLens to Greeting("hello sample"))
} bindContract GET

val sayHello: ContractRoute = sayHelloSpec to { name ->
    { request: Request ->
        val userId = userIdLens(request)
        val greeting = Greeting("hello $userId/$name")
        Response(OK)
            .with(greetingLens of greeting)
    }
}

val api = contract {
    routes += sayHello
    security = BearerAuthSecurity(userIdLens, authLookup)

    renderer = OpenApi3(
        apiInfo = ApiInfo("Hello Api", "1.0.0")
    )
    descriptionPath = "openapi.json"
}

val ui = swaggerUiLite {
    url = "openapi.json"
}

val helloContractApi = ServerFilters
    .InitialiseRequestContext(requestContexts)
    .then(routes(api, ui))

fun main() {
    helloContractApi.asServer(Undertow(8002)).start()
}