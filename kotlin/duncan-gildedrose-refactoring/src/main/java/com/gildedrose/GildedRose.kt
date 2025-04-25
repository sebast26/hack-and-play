package com.gildedrose

internal class GildedRose(var items: List<Item>) {
    fun updated() = GildedRose(items.map { it.updated() })
}
