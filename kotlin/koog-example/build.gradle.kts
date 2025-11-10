plugins {
    kotlin("jvm") version "2.2.20"
    alias(libs.plugins.kotlin.serialization)
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(libs.koog.agents)
    implementation(libs.koog.tools)
    implementation(libs.koog.executor.openai.client)
    testImplementation(kotlin("test"))
}

tasks.test {
    useJUnitPlatform()
}
kotlin {
    jvmToolchain(21)
}