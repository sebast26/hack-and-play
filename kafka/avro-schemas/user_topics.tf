resource "kafka_topic" "me_sgorecki_iam_user" {
  name               = "me.sgorecki.iam.User"
  replication_factor = 1
  partitions         = 1
  config = {
    "cleanup.policy" = "delete"
  }
}