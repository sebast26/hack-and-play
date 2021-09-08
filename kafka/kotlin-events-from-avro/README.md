# kotlin-events-from-avro

Sample project that should demonstrate how to generate classes from AVRO schema files using Gradle plugin.

## How to?

### Get AVRO schemas from schema registry

To fetch schema definitions from schema registry you can use Schema Registry API.

API Documentation: [Schema Registry API Reference](https://docs.confluent.io/platform/current/schema-registry/develop/api.html)

Example usages:

```
curl -f http://localhost:8081/subjects
curl -f http://localhost:8081/subjects/me.sgorecki.iam.UserDisabled/versions/1/schema
```

Then you have to put AVRO files in `src/main/avro` folder.

### How to generate Java files from AVRO schemas?

You can use `com.github.davidmc24.gradle.plugin.avro` Gradle plugin.

To generate Java files you can use: `gradle generateAvroJava`. This will generate Java classes inside build folder.