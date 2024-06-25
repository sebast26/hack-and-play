package com.example

import org.http4k.core.Filter
import org.http4k.core.HttpHandler
import org.http4k.core.Request

fun Debug() = Filter { next: HttpHandler ->
    { req: Request ->
        next(req).also { response -> println(response) }
    }
}