package com.gildedrose

import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status
import org.http4k.core.Status.Companion.OK
import org.http4k.testing.ApprovalTest
import org.http4k.testing.Approver
import org.http4k.testing.assertApproved
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import java.time.Instant

@ExtendWith(ApprovalTest::class)
class ListStockTests {

    @Test
    fun `list stock`(approver: Approver) {
        with(Fixture(standardStockList, now = Instant.parse("2023-03-09T12:00:00Z"))) {
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
        }
    }

    @Test
    fun `list stock sees file updates`(approver: Approver) {
        with(Fixture(standardStockList, now = Instant.parse("2023-03-09T12:00:00Z"))) {
            assertEquals(OK, routes(Request(Method.GET, "/")).status)

            save(StockList(Instant.now(), emptyList()))
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
        }
    }
}

