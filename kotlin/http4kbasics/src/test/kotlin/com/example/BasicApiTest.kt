package com.example

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status.Companion.OK
import org.junit.jupiter.api.Test

class BasicApiTest {

    @Test
    fun `say hello - in memory`() {
        val request = Request(Method.GET, "/hello/http4k")

        val response = helloApi(request)

        assertThat(response.status, equalTo(OK))
        assertThat(response.bodyString(), equalTo("""{"message":"hello http4k"}"""))
    }

}