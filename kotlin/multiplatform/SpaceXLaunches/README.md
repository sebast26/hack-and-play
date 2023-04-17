# SpaceX Launches project

This is tutorial project based on [Create a multiplatform app using Ktor and SQLDelight – tutorial]
(https://kotlinlang.org/docs/multiplatform-mobile-ktor-sqldelight.html).

## Description

This tutorial demonstrates how to use Android Studio to create a mobile application for iOS and 
Android using Kotlin Multiplatform Mobile with Ktor and SQLDelight.

The application will include a module with shared code for both the iOS and Android platforms. 
The business logic and data access layers will be implemented only once in the shared module, while 
the UI of both applications will be native.

The output will be an app that retrieves data over the internet from the public SpaceX API, saves it 
in a local database, and displays a list of SpaceX rocket launches together with the launch date, 
results, and a detailed description of the launch.