package pl.sgorecki

import com.uchuhimo.konf.Config
import pl.sgorecki.k8s.App.Deployment.PodTemplate.Containers.Resources.Limits
import pl.sgorecki.k8s.App.Deployment.PodTemplate.Containers.Resources.Requests
import pl.sgorecki.k8s.AppEnvironment
import pl.sgorecki.k8s.AwsEnvironment.staging
import pl.sgorecki.k8s.CpuResource
import pl.sgorecki.k8s.Database
import pl.sgorecki.k8s.Ingress
import pl.sgorecki.k8s.MemoryResource

object Staging : AppEnvironment(staging, Config::initialise) {
    override fun configure(config: Config) {
        config[Database.writerHost] = "home-hackandplay.cluster.aws"
        config[Database.readerHost] = "home-hackandplay-ro.cluster.aws"

        config[Limits.cpu] = CpuResource.of("1")
        config[Limits.memory] = MemoryResource.of("600Mi")

        config[Requests.cpu] = CpuResource.of("100m")
        config[Requests.memory] = config[Limits.memory]

        config[Ingress.certificateArn] = "abc:arn:com"
        config[Ingress.hostname] = "some.example.com"

        applySharedSettings(config)
    }
}