package pointofsale

import products.Product

interface PointOfSalePresenter {
    fun format(products: List<Product>): String
}

class PointOfSaleOutputPresenter() : PointOfSalePresenter {
    override fun format(products: List<Product>): String {
        return products.joinToString("\n") { "${it.uid} - ${it.name} | ${it.price}" }
    }

}