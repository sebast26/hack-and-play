package pl.sgorecki.k8s

import dev.forkhandles.values.NonBlankStringValueFactory
import dev.forkhandles.values.StringValue
import java.math.BigInteger

class CpuResource private constructor(private val raw: String) : StringValue(raw), Comparable<CpuResource> {
    companion object : NonBlankStringValueFactory<CpuResource>(::CpuResource)

    val bytes = when {
        raw.endsWith("m") -> raw.removeSuffix("m").toInt()
        raw.toDoubleOrNull() != null -> (raw.toDouble() * 1000).toInt()
        else -> error("Value not supported yet: $this")
    }

    override fun compareTo(other: CpuResource) = bytes.compareTo(other.bytes)
}

class MemoryResource private constructor(private val raw: String) : StringValue(raw), Comparable<MemoryResource> {
    companion object : NonBlankStringValueFactory<MemoryResource>(::MemoryResource)

    private val exponentRegex = "(\\d+)e(\\d+)".toRegex()

    val bytes = when {
        raw.endsWith("Gi") -> raw("Gi") * powerOf1024(3)
        raw.endsWith("G") -> raw("G") * powerOf1000(3)
        raw.endsWith("Mi") -> raw("Mi") * powerOf1024(2)
        raw.endsWith("M") -> raw("M") * powerOf1000(2)
        raw.endsWith("Ki") -> raw("Ki") * powerOf1024(1)
        raw.endsWith("k") -> raw("k") * powerOf1000(1)
        raw.endsWith("m") -> raw("m").div(1_000.toBigInteger())
        raw.matches(exponentRegex) -> {
            val (quantity, exponent) = exponentRegex.find(raw)!!.destructured
            quantity.toBigInteger() * 10.toBigInteger().pow(exponent.toInt())
        }

        raw.toIntOrNull() != null -> raw.toBigInteger()
        else -> error("Value not supported yet: $this")
    }

    override fun compareTo(other: MemoryResource) = bytes.compareTo(other.bytes)

    private operator fun String.invoke(suffix: String) = removeSuffix(suffix).toBigInteger()
}

private fun powerOf1000(factor: Int): BigInteger = 10.toBigInteger().pow(3 * factor)
private fun powerOf1024(factor: Int): BigInteger = 2.toBigInteger().pow(10 * factor)