resource "schemaregistry_schema" "me_sgorecki_iam_UserCreated" {
  subject = "me.sgorecki.iam.UserCreated"
  schema = jsonencode({
    "type" : "record",
    "name" : "UserCreated",
    "namespace" : "me.sgorecki.iam",
    "doc" : "When user was created",
    "fields" : [
      {
        name : "userUuid",
        type : {
          type: "string",
          logicalType : "uuid"
        }
        doc : "User UUID"
      },
      {
        name : "email",
        type : "string",
        doc : "email address provided during user registration"
      },
      {
        name : "createAt",
        type : {
          type : "long",
          logicalType : "timestamp-mills"
        },
        doc : "Timestamp when user was created"
      }
    ]
  })
}

resource "schemaregistry_schema" "me_sgorecki_iam_UserDisabled" {
  subject = "me.sgorecki.iam.UserDisabled"
  schema = jsonencode({
    "type" : "record",
    "name" : "UserDisabled",
    "namespace" : "me.sgorecki.iam",
    "doc" : "When user was disabled",
    "fields" : [
      {
        name : "userUuid",
        type : {
          type : "string",
          logicalType : "uuid"
        },
        doc : "User UUID"
      },
      {
        name : "disabledAt",
        type : {
          type : "long",
          logicalType : "timestamp-mills"
        },
        doc : "Timestamp when user was disabled"
      }
    ]
  })
}

resource "schemaregistry_schema" "me_sgorecki_iam_User" {
  subject = "${kafka_topic.me_sgorecki_iam_user.name}-value"

  schema = jsonencode([
    "me.sgorecki.iam.UserCreated",
    "me.sgorecki.iam.UserDisabled"
  ])

  reference {
    name    = schemaregistry_schema.me_sgorecki_iam_UserCreated.subject
    subject = schemaregistry_schema.me_sgorecki_iam_UserCreated.subject
    version = schemaregistry_schema.me_sgorecki_iam_UserCreated.version
  }

  reference {
    name    = schemaregistry_schema.me_sgorecki_iam_UserDisabled.subject
    subject = schemaregistry_schema.me_sgorecki_iam_UserDisabled.subject
    version = schemaregistry_schema.me_sgorecki_iam_UserDisabled.version
  }
}