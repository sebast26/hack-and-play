package me.sgorecki.spacexlaunches

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform