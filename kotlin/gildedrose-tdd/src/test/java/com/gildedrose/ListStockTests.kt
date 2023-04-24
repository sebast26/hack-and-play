package com.gildedrose

import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Status
import org.http4k.core.Status.Companion.OK
import org.http4k.testing.ApprovalTest
import org.http4k.testing.Approver
import org.http4k.testing.assertApproved
import org.junit.jupiter.api.Assertions
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import java.time.Instant

@ExtendWith(ApprovalTest::class)
class ListStockTests {

    private val stockList = StockList(
        lastModified = Instant.parse("2023-03-25T12:00:00Z"),
        items = listOf(
            Item("banana", mar03.minusDays(1), 42),
            Item("kumquat", mar03.plusDays(1), 101),
            Item("undated", null, 50),
        )
    )

    @Test
    fun `list stock`(approver: Approver) {
        with(
            Fixture(stockList, now = Instant.parse("2023-03-09T12:00:00Z"))
        ) {
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
        }
    }

    @Test
    fun `list stock sees file updates`(approver: Approver) {
        with(
            Fixture(stockList, now = Instant.parse("2023-03-09T12:00:00Z"))
        ) {
            assertEquals(OK, routes(Request(Method.GET, "/")).status)

            save(StockList(Instant.now(), emptyList()))
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
        }
    }

    @Test
    fun `doesn't update when lastModified in today`(approver: Approver) {
        val sameDayAsLastModified = Instant.parse("2023-03-25T23:59:59Z")
        with(
            Fixture(stockList, now = sameDayAsLastModified)
        ) {
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
            assertEquals(stockList, load())
        }
    }

    @Test
    fun `does update when lastModified was yesterday`(approver: Approver) {
        val nextDayFromLastModified = Instant.parse("2023-03-26T00:00:01Z")
        with(
            Fixture(stockList, now = nextDayFromLastModified)
        ) {
            approver.assertApproved(routes(Request(Method.GET, "/")), Status.OK)
            Assertions.assertNotEquals(stockList, load())
        }
    }
}

