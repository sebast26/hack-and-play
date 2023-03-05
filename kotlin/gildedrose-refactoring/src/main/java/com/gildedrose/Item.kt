package com.gildedrose

class Item(name: String?, sellIn: Int, quality: Int) {
    var name: String? = null
    var sellIn = 0
    var quality = 0

    init {
        this.name = name
        this.sellIn = sellIn
        this.quality = quality
    }

    override fun toString(): String {
        return name + ", " + sellIn + ", " + quality
    }
}