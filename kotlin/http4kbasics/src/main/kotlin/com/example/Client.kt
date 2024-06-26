package com.example

import org.http4k.client.JavaHttpClient
import org.http4k.core.HttpHandler
import org.http4k.core.Method
import org.http4k.core.Request

fun main() {
    val internet: HttpHandler = JavaHttpClient()

    val request = Request(Method.GET, "https://httpbin.org/uuid")

    val response = internet(request)

    println(response.bodyString())
}