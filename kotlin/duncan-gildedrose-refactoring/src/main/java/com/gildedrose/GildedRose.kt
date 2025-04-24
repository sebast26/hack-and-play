package com.gildedrose

internal class GildedRose(var items: List<Item>) {
    fun updateQuality() {
        for (item in items) {
            item.update()
        }
    }
}
