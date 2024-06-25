package com.example

import org.http4k.core.HttpHandler
import org.http4k.core.then
import org.http4k.routing.routes

fun Pokemon4k(pokecoHttp: HttpHandler) = Debug().then(
    routes(FindPokemonWithPrefix(pokecoHttp))
)