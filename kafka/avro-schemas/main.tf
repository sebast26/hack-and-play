terraform {
  required_providers {
    kafka = {
      source  = "Mongey/kafka"
      version = "0.3.3"
    }

    schemaregistry = {
      source = "arkiaconsulting/confluent-schema-registry"
      version = "0.7.0"
    }
  }
}

provider "kafka" {
  bootstrap_servers = ["localhost:9092"]
  tls_enabled       = false
}

provider "schemaregistry" {
  schema_registry_url = "http://localhost:8081"
  username            = ""
  password            = ""
}

