package com.gildedrose

import org.http4k.core.Method
import org.http4k.core.Status
import org.http4k.testing.ApprovalTest
import org.http4k.testing.Approver
import org.http4k.testing.assertApproved
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import java.time.LocalDate

@ExtendWith(ApprovalTest::class)
class ListStockTests {
    private val now = LocalDate.parse("2023-03-09")

    @Test
    fun `list stock`(approver: Approver) {
        val stock = listOf(
            Item("banana", now.minusDays(1), 42u),
            Item("kumquat", now.plusDays(1), 101u)
        )
        val server = Server(stock) { now }
        val client = server.routes
        val response = client(org.http4k.core.Request(Method.GET, "/"))
        approver.assertApproved(response, Status.OK)
    }
}
