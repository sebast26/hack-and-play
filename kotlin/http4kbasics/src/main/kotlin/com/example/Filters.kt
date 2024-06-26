package com.example

import org.http4k.core.Filter
import org.http4k.core.HttpHandler
import org.http4k.core.Method
import org.http4k.core.Method.GET
import org.http4k.core.Request
import org.http4k.core.then
import org.http4k.filter.ResponseFilters

val printTransaction = Filter { next ->
    {
        println("Request:\n$it")
        val response = next(it)
        println("Response:\n$response")
        response
    }
}

val setFooHeader = ResponseFilters
    .SetHeader("Foo", "Bar")

val filteredApi: HttpHandler = printTransaction
    .then(setFooHeader)
    .then(helloApi)

fun main() {
    val req = Request(GET, "/hello/http4k")
    val res = filteredApi(req)
}