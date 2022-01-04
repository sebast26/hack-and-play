buildscript {
    repositories {
        maven {
            url = uri("https://plugins.gradle.org/m2/")
        }
        maven {
            url = uri("https://packages.confluent.io/maven/")
        }
        maven {
            url = uri("https://jitpack.io")
        }
    }
}

plugins {
    kotlin("jvm") version "1.6.0"
    id("com.github.imflog.kafka-schema-registry-gradle-plugin") version "1.6.0"
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(kotlin("stdlib"))
}

schemaRegistry {
    url.set("")
    credentials {
        username.set("")
        password.set("")
    }
    download {
        subject("avroSubject", "src/main/avro")
//        subjectPattern("avro.*", )
    }
}