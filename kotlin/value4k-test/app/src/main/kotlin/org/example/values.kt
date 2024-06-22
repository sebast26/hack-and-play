package org.example

import dev.forkhandles.values.*

class Email private constructor(override val value: String) : StringValue(value) {
    companion object : StringValueFactory<Email>(::Email, "^[^@]+@[^@]+\\.[^@]+\$".regex)
}

class PostalCode private constructor(override val value: String) : StringValue(value) {
    companion object : StringValueFactory<PostalCode>(::PostalCode, "\\d{2}-\\d{3}".regex)
}

class OrderId private constructor(override val value: Int) : IntValue(value) {
    companion object : IntValueFactory<OrderId>(::OrderId, 1.minValue)
}

fun test() {
    val email = Email.of("example@example.com")
    val postalCode = PostalCode.of("88-310")
}