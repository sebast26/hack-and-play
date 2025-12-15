package pl.sgorecki

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.ConfigSpec
import com.uchuhimo.konf.Feature.FAIL_ON_UNKNOWN_PATH
import pl.sgorecki.k8s.App
import pl.sgorecki.k8s.AppEnvironment
import pl.sgorecki.k8s.Database
import pl.sgorecki.k8s.Ingress

object HackAndPlay : ConfigSpec(prefix = "") {

}

fun Config.initialise() {
    addSpec(App)
    addSpec(Database)
    addSpec(Ingress)
    enable(FAIL_ON_UNKNOWN_PATH)
}

fun AppEnvironment.applySharedSettings(config: Config) {
    config[App.department] = "home"
    config[App.projectName] = "hack-and-play"
    config[App.deploymentName] = "${config[App.department]}-${config[App.projectName]}"

    config[Database.name] = "hackandplay"
    config[Database.primaryUrl] = Database.primaryUrlFor(config)
    config[Database.readerUrl] = Database.readerUrlFor(config)

}