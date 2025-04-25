package com.gildedrose

data class Item(
    val name: String,
    val sellIn: Int = 0,
    val quality: Int = 0,
    private val aging: () -> Int = Aging.standard,
    private val degradation: (Int, Int) -> Int = Degradation.standard,
    private val saturation: (Int) -> Int = Saturation.standard
) {
    override fun toString() = "$name, $sellIn, $quality"

    fun updated(): Item {
        val sellIn = sellIn - aging()
        return this.copy(
            sellIn = sellIn,
            quality = saturation(this.quality - degradation(sellIn, this.quality))
        )
    }
}
