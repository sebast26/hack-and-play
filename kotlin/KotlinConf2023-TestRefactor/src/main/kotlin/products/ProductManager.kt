package products

interface ProductManager {
    fun getAllProducts(): List<Product>

    fun getProductById(id: String): Product?
}