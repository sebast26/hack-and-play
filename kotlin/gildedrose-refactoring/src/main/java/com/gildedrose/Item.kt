package com.gildedrose

open class Item(
    val name: String, var sellIn: Int, var quality: Int
) {

    override fun toString(): String = "$name, $sellIn, $quality"
}

class BaseItem(
    name: String,
    sellIn: Int,
    quality: Int,
    private val aging: () -> Int = { 1 },
    private val degradation: (Int, Int) -> Int = { sellIn: Int, quality: Int ->
        when {
            sellIn < 0 -> 2
            else -> 1
        }
    },
    private val saturation: (Int) -> Int = { quality: Int ->
        when {
            quality < 0 -> 0
            quality > 50 -> 50
            else -> quality
        }
    }
) : Item(name, sellIn, quality) {
    fun update() {
        sellIn = sellIn - aging()
        quality = saturation(quality - degradation(sellIn, quality))
    }
}

fun Sulfuras(name: String, sellIn: Int, quality: Int) = BaseItem(
    name,
    sellIn,
    quality,
    aging = { 0 },
    degradation = fun(sellIn: Int, quality: Int) = 0,
    saturation = fun(quality: Int) = quality
)

fun Brie(name: String, sellIn: Int, quality: Int) = BaseItem(
    name,
    sellIn,
    quality,
    degradation = fun(sellIn: Int, quality: Int) = when {
        sellIn < 0 -> -2
        else -> -1
    }
)

fun Pass(name: String, sellIn: Int, quality: Int) = BaseItem(
    name,
    sellIn,
    quality,
    degradation = fun(sellIn: Int, quality: Int): Int = when {
        sellIn < 0 -> quality
        sellIn < 5 -> -3
        sellIn < 10 -> -2
        else -> -1
    }
)
