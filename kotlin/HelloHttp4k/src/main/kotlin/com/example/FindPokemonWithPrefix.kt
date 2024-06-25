package com.example

import org.http4k.core.*
import org.http4k.format.Moshi.auto
import org.http4k.lens.Path
import org.http4k.routing.RoutingHttpHandler
import org.http4k.routing.bind

fun FindPokemonWithPrefix(pokecoHttp: HttpHandler): RoutingHttpHandler =
    "/{prefix}" bind Method.GET to { req: Request ->
        val client = PokemonClient(pokecoHttp)
        val results = client.list()
        val prefixLens = Path.of("prefix")
        val body = Body.auto<PokemonList>().toLens()
        val prefix: String = prefixLens(req)

        val pokemonList = PokemonList(results
            .map(Pokemon::name)
            .filter { it.startsWith(prefix) }
        )
        Response(Status.OK).with(body of pokemonList)
    }

data class PokemonList(val pokemon: List<String>)