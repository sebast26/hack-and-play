package pointofsale

import money.USD
import org.junit.jupiter.api.Test
import products.NetworkProductsManager
import products.Product
import products.ProductManager
import receipts.ReceiptDisplay

class PointOfSaleTest {
    private fun pointOfSale(
        productsManager: ProductManager = NetworkProductsManager(),
        receiptDisplay: ReceiptDisplay = ReceiptDisplay(),
        output: ScreenOutput = TerminalOutput(),
        presenter: PointOfSalePresenter = PointOfSaleOutputPresenter(),
    ): PointOfSale = PointOfSale(
        productsManager = productsManager,
        receiptDisplay = receiptDisplay,
        screenOutput = output,
        presenter = presenter,
    )

    @Test
    fun `addItemToTransaction should print transaction`() {
        val screenOutput = screenOutputMock()
        val productManager = ProductManagerStub(
            getProductByIdStub = { id -> Product("", "", USD(0)) }
        )
        val pointOfSale = pointOfSale(
            productsManager = productManager,
            output = screenOutput
        )

        pointOfSale.addItemToTransaction("")

        assert(screenOutput.wasPrintCalled)
    }

    @Test
    fun `showAvailableItems prints a line for each item containing uid, name, and price`() {
        val screenOutput = screenOutputMock()
        val pointOfSale = pointOfSale(
            output = screenOutput
        )

        pointOfSale.showAvailableItems()

        assert(screenOutput.wasPrintCalled)
        assert(screenOutput.printCallContents == "some format")

    }

    interface ScreenOutputMock : ScreenOutput {
        var wasPrintCalled: Boolean
        var printCallContents: String?
    }

    private fun screenOutputMock(): ScreenOutputMock {
        return object : ScreenOutputMock {
            override var wasPrintCalled: Boolean = false
            override var printCallContents: String? = null

            override fun print(output: String) {
                wasPrintCalled = true
                printCallContents = output
            }
        }
    }

    class ProductManagerStub(
        val getAllProductsStub: () -> List<Product> = { emptyList() },
        val getProductByIdStub: (id: String) -> Product? = { null }
    ) : ProductManager {
        override fun getAllProducts(): List<Product> {
            return getAllProductsStub()
        }

        override fun getProductById(id: String): Product? {
            return getProductByIdStub(id)
        }

    }

    class PresenterStub(
        formatStub: (List<Product>) -> String = { "" }
    ) : PointOfSalePresenter by PointOfSaleOutputPresenter()
}
