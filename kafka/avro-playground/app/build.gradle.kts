plugins {
    id("org.jetbrains.kotlin.jvm") version "1.5.31"
    id("com.github.davidmc24.gradle.plugin.avro") version "1.2.1"
    application
}

repositories {
    maven {
        name = "confluent"
        url = uri("https://packages.confluent.io/maven/")
    }
    mavenLocal()
    mavenCentral()
}

dependencies {
    implementation(platform("org.jetbrains.kotlin:kotlin-bom"))
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    implementation("org.apache.avro:avro:1.11.0")
    implementation("org.slf4j:slf4j-nop:1.7.36")

    testImplementation("org.jetbrains.kotlin:kotlin-test")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit")
}

kotlin {
    tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile>().configureEach {
        kotlinOptions {
            jvmTarget = "1.8"
        }
    }
}

application {
    // Define the main class for the application.
    mainClass.set("me.sgorecki.avro.AppKt")
}
