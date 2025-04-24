package com.gildedrose

open class Item(
    val name: String,
    var sellIn: Int = 0,
    var quality: Int = 0
) {
    override fun toString() = "$name, $sellIn, $quality"
}

open class BaseItem(name: String, sellIn: Int, quality: Int) : Item(name, sellIn, quality) {
    fun update() {
        sellIn = sellIn - aging()
        quality = saturation(degradation(sellIn, quality))
    }

    protected open fun aging(): Int = 1

    protected open fun degradation(sellIn: Int, quality: Int) = when {
        sellIn < 0 -> quality - 2
        else -> quality - 1
    }

    protected open fun saturation(quality: Int) = when {
        quality < 0 -> 0
        quality > 50 -> 50
        else -> quality
    }
}

class Brie(name: String, sellIn: Int, quality: Int) : BaseItem(name, sellIn, quality) {
    override fun degradation(sellIn: Int, quality: Int) = when {
        sellIn < 0 -> quality + 2
        else -> quality + 1
    }
}

class Pass(name: String, sellIn: Int, quality: Int) : BaseItem(name, sellIn, quality) {
    override fun degradation(sellIn: Int, quality: Int) = when {
        sellIn < 0 -> 0
        sellIn < 5 -> quality + 3
        sellIn < 10 -> quality + 2
        else -> quality + 1

    }
}

class Sulfuras(name: String, sellIn: Int, quality: Int) : BaseItem(name, sellIn, quality) {
    override fun aging() = 0
    override fun degradation(sellIn: Int, quality: Int) = quality
    override fun saturation(quality: Int) = quality
}
