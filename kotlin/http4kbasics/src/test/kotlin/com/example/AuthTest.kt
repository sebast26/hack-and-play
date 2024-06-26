package com.example

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status
import org.http4k.core.Status.Companion.UNAUTHORIZED
import org.http4k.core.then
import org.http4k.filter.ClientFilters
import org.junit.jupiter.api.Test

class AuthTest {
    @Test
    fun unauthorized() {
        val request = Request(Method.GET, "")
        val response = authenticatedApi(request)

        assertThat(response.status, equalTo(UNAUTHORIZED))
    }

    @Test
    fun ok() {
        val withToken = ClientFilters.BearerAuth("letmein")
            .then(authenticatedApi)

        val request = Request(Method.GET, "")
        val response = withToken(request)

        assertThat(response.status, equalTo(Status.OK))
        assertThat(greetingLens(response), equalTo(Greeting("hello user1")))
    }
}