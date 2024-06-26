package com.example

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status
import org.http4k.core.with
import org.junit.jupiter.api.Test

class LensesTest {
    @Test
    fun `lenses test`() {
        val request = Request(Method.GET, helloRoute)
            .with(nameLens of "Http4k")
        val response = helloApiWithLenses(request)

        assertThat(response.status, equalTo(Status.OK))
        assertThat(greetingLens(response), equalTo(Greeting("hello Http4k")))
    }
}