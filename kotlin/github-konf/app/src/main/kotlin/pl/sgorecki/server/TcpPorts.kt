package pl.sgorecki.server

import dev.forkhandles.values.IntValue
import dev.forkhandles.values.IntValueFactory
import dev.forkhandles.values.minValue
import pl.sgorecki.server.TcpPortName.app
import pl.sgorecki.server.TcpPortName.ops

@Suppress("EnumEntryName")
enum class TcpPortName { app, ops }

class TcpPort private constructor(value: Int) : IntValue(value) {
    companion object : IntValueFactory<TcpPort>(::TcpPort, (-1).minValue) {
        val DEFAULT = TcpPort.of(80)
    }
}

data class TcpPorts(
    private val config: MutableMap<TcpPortName, TcpPort> = mutableMapOf(app to TcpPort.DEFAULT)
) : Map<TcpPortName, TcpPort> by config {
    fun add(name: TcpPortName, value: TcpPort) = config.put(name, value)
    fun forApp() = config.getValue(app)
    fun forOps() = config[ops] ?: config.getValue(app)
    fun forProbes() = (config[ops]?.let { ops } ?: app).name
}