package pointofsale

import products.ProductManager
import receipts.ReceiptDisplay

class PointOfSale(
    private val productsManager: ProductManager,
    private val receiptDisplay: ReceiptDisplay,
    private val screenOutput: ScreenOutput,
    private val presenter: PointOfSalePresenter,
) {
    private var transaction = Transaction(emptyList())

    fun showAvailableItems() {
        presenter.format(productsManager.getAllProducts())
            .also { screenOutput.print(it) }
    }

    fun addItemToTransaction(id: String) {
        val product = productsManager.getProductById(id)
            ?: throw IllegalArgumentException("No product matching id $id")
        transaction = transaction.copy(products = transaction.products + product)

        showTransaction()
    }

    private fun showTransaction() {
        for (product in transaction.products) {
            screenOutput.print("${product.name} | ${product.price}")
        }
        screenOutput.print("Total ${transaction.total}")
    }

    fun printReceipt() {
        receiptDisplay.renderReceipt(transaction.products)
    }
}
