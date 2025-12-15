package pl.sgorecki.k8s

import com.uchuhimo.konf.Config
import com.uchuhimo.konf.Item
import com.uchuhimo.konf.NoSuchItemException
import org.http4k.connect.amazon.core.model.AwsAccount
import org.http4k.connect.amazon.core.model.Region
import org.http4k.connect.amazon.core.model.Region.Companion.EU_WEST_1

@Suppress("EnumEntryName")
enum class AwsEnvironment(
    val shortName: String,
    val account: AwsAccount,
    val region: Region,
) {
    staging(
        shortName = "stg",
        account = AwsAccount.of("12345678901111"),
        region = EU_WEST_1
    ),
    production(
        shortName = "prd",
        account = AwsAccount.of("77777777777111"),
        region = EU_WEST_1
    )
}

abstract class AppEnvironment(
    val name: AwsEnvironment,
    initialise: Config.() -> Unit,
) {
    abstract fun configure(config: Config)

    val config: Config by lazy {
        val config = Config(initialise)
        try {
            configure(config)
        } catch (e: NoSuchItemException) {
            throw Exception("You supplied config ${e.name}, but it was not expected. Do you need to add Spec to Config.initialize() ?")
        }
        config.validateRequired()
    }

    val awsAccount: AwsAccount = name.account
    val awsRegion: Region = name.region

    val deploymentName by lazy { config[App.deploymentName] }
    val namespace by lazy { config[App.namespace] }

    operator fun <T> get(item: Item<T>): T = config[item]
}