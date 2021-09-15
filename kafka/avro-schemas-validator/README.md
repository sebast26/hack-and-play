If you run

java -jar avro-tools-1.10.2.jar compile schema invalid_avro out

it will output WARNings (invalid use of logicalType)

If you run

java -jar avro-tools-1.10.2.jar compile schema valid_avro out

there will be no warnings