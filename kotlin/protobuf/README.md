# protobuf in Kotlin


Link to [tutorial](https://protobuf.dev/getting-started/kotlintutorial/).

Compiling the protobuf files:

```bash
mkdir -p app/build/gen/kotlin
mkdir -p app/build/gen/java
protoc --proto_path=app/src --java_out=app/build/gen/java --kotlin_out=app/build/gen/kotlin app/src/addressbook.proto
```

