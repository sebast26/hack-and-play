package me.sgorecki.avro

import example.avro.User
import org.apache.avro.Schema
import org.apache.avro.file.DataFileReader
import org.apache.avro.file.DataFileWriter
import org.apache.avro.generic.GenericData
import org.apache.avro.generic.GenericDatumReader
import org.apache.avro.generic.GenericDatumWriter
import org.apache.avro.generic.GenericRecord
import org.apache.avro.specific.SpecificDatumReader
import org.apache.avro.specific.SpecificDatumWriter
import java.io.File

class SerdeWithCodeGeneration {
    fun run() {
        val (user1, user2, user3) = createUsers()
        serializeUsers(user1, user2, user3)
        deserializeAndPrint()
    }

    private fun createUsers(): Triple<User, User, User> {
        val user1: User = User()
        user1.name = "Sebastian"
        user1.favoriteColor = "blue"

        val user2: User = User("Ignacy", 7, "green")

        val user3: User = User.newBuilder()
            .setName("Jadwiga")
            .setFavoriteColor("red")
            .setFavoriteNumber(null)
            .build()

        return Triple(user1, user2, user3)
    }

    /**
     * We create a DatumWriter, which converts Java objects into an in-memory serialized format.
     * The SpecificDatumWriter class is used with generated classes and extracts the schema from the specified generated type.
     *
     * Next we create a DataFileWriter, which writes the serialized records, as well as the schema, to the file
     * specified in the dataFileWriter.create call. We write our users to the file via calls to the dataFileWriter.append method.
     * When we are done writing, we close the data file.
     */
    private fun serializeUsers(user1: User, user2: User, user3: User) {
        val userDatumWriter: SpecificDatumWriter<User> = SpecificDatumWriter(User::class.java)

        DataFileWriter(userDatumWriter).use {
            it.create(user1.schema, File("users.avro"))
            it.append(user1)
            it.append(user2)
            it.append(user3)
        }
    }

    private fun deserializeAndPrint() {
        val userDatumReader: SpecificDatumReader<User> = SpecificDatumReader(User::class.java)
        DataFileReader(File("users.avro"), userDatumReader).use {
            var user: User
            while (it.hasNext()) {
                user = it.next()
                println(user)
            }
        }
    }
}

/**
 * Data in Avro is always stored with its corresponding schema, meaning we can always read a serialized item regardless
 * of whether we know the schema ahead of time.
 *
 * This allows us to perform serialization and deserialization without code generation.
 */
class SerdeWithoutCodeGeneration {
    fun run() {
        // First, we use a Parser to read our schema definition and create a Schema object.
        val schema: Schema = Schema.Parser().parse(this::class.java.getResourceAsStream("/user.avsc"))

        val (user1, user2) = createUsers(schema)
        serializeUsers(schema, user1, user2)
        deserializeAndPrint(schema)
    }

    private fun createUsers(schema: Schema): Pair<GenericRecord, GenericRecord> {
        // Since we're not using code generation, we use GenericRecords to represent users.
        // GenericRecord uses the schema to verify that we only specify valid fields.
        // If we try to set a non-existent field (e.g., user1.put("favorite_animal", "cat")), we'll get an AvroRuntimeException
        // when we run the program.
        val user1: GenericRecord = GenericData.Record(schema)
        user1.put("name", "Sebastian")
        user1.put("favorite_number", 256)

        val user2: GenericRecord = GenericData.Record(schema)
        user2.put("name", "Ignacy")
        user2.put("favorite_number", 17)
        user2.put("favorite_color", "green")

        return Pair(user1, user2)
    }

    /**
     * Serializing and deserializing is almost identical to the example above which uses code generation.
     * The main difference is that we use generic instead of specific readers and writers.
     */
    private fun serializeUsers(schema: Schema, user1: GenericRecord, user2: GenericRecord) {
        val datumWriter: GenericDatumWriter<GenericRecord> = GenericDatumWriter<GenericRecord>(schema)
        DataFileWriter(datumWriter).use {
            it.create(schema, File("users_generic.avro"))
            it.append(user1)
            it.append(user2)
        }
    }

    private fun deserializeAndPrint(schema: Schema) {
        val datumReader: GenericDatumReader<GenericRecord> = GenericDatumReader<GenericRecord>(schema)
        DataFileReader(File("users_generic.avro"), datumReader).use {
            var user: GenericRecord
            while (it.hasNext()) {
                user = it.next()
                println(user)
            }
        }
    }
}

fun main() {
    SerdeWithCodeGeneration().run()
    println("-------------------")
    SerdeWithoutCodeGeneration().run()
}
