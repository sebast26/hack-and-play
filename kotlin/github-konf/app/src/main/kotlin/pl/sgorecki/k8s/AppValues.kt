package pl.sgorecki.k8s

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.ConfigSpec
import io.fabric8.kubernetes.api.model.EnvVar
import org.http4k.config.Host
import org.http4k.config.Port
import pl.sgorecki.server.TcpPorts
import pl.sgorecki.storage.databaseUrlFor

object App : ConfigSpec(prefix = "") {
    val deploymentName by required<String>()
    val namespace by optional(default = "hack-and-play")

    val department by optional(default = "seba")

    val imageName by optional(default = "{{ .Values.imageName }}", description = "name of the container image")
    val entrypointCommand by optional(default = emptyList<String>())
    val entrypointArgs by optional(default = emptyList<String>())

    val projectName by required<String>()
    val teamName by optional(default = "gumisie")

    val tcpPorts by optional(default = TcpPorts())

    object Deployment : ConfigSpec() {
        object PodTemplate : ConfigSpec() {
            val serviceAccountName by optional(default = "")

            object Containers : ConfigSpec() {
                val environmentVariables by optional(default = emptyList<EnvVar>())

                object Resources : ConfigSpec() {
                    object Limits : ConfigSpec() {
                        val cpu by required<CpuResource>()
                        val memory by required<MemoryResource>()
                    }

                    object Requests : ConfigSpec() {
                        val cpu by required<CpuResource>()
                        val memory by required<MemoryResource>()
                    }
                }

                object ReadinessProbe : ConfigSpec() {
                    val path by optional(default = "/internal/readiness")
                    val timeoutSeconds by optional(default = 5)
                }

                object LivenessProbe : ConfigSpec() {
                    val path by optional(default = "/internal/liveness")
                    val timeoutSeconds by optional(default = 5)
                }
            }
        }
    }
}

object Ingress : ConfigSpec() {
    val securityGroups by optional<String?>(default = null)
    val certificateArn by required<String>()
    val hostname by required<String>()
}

object Database : ConfigSpec() {
    val name by required<String>()

    val writerHost by required<String>()
    private val writerPort by optional(default = 5432)

    val readerHost by required<String>()
    private val readerPort by optional(default = 5432)

    val primaryUrl by required<String>()
    val readerUrl by required<String>()

    fun primaryUrlFor(config: Config) = databaseUrlFor(
        databaseHost = Host(config[writerHost]),
        databasePort = Port(config[writerPort]),
        databaseName = config[name]
    ).toString()

    fun readerUrlFor(config: Config) = databaseUrlFor(
        databaseHost = Host(config[readerHost]),
        databasePort = Port(config[readerPort]),
        databaseName = config[name]
    ).toString()
}