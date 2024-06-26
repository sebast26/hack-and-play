package com.example

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status.Companion.OK
import org.junit.jupiter.api.Test

class FilterTest {

    @Test
    fun `test filter`() {
        val request = Request(Method.GET, "/hello/http4k")
        val response = filteredApi(request)

        assertThat(response.status, equalTo(OK))
        assertThat(response.header("Foo"), equalTo("Bar"))
        assertThat(greetingLens(response), equalTo(Greeting("hello http4k")))
    }
}