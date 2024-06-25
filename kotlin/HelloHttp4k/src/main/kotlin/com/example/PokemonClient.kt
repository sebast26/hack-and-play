package com.example

import org.http4k.core.*
import org.http4k.core.Method.*
import org.http4k.format.Moshi.auto

data class Pokemon(val name: String)

data class Result(val results: List<Pokemon>)

class PokemonClient(private val httpHandler: HttpHandler) {
    fun list(): List<Pokemon> {
        val body = Body.auto<Result>().toLens()
        val response = httpHandler(Request(GET, "/api/v2/pokemon").query("limit", "100"))
        return body(response).results
    }
}