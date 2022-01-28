# DONE Bucket - destination for events from Stream -> Firehose
resource "aws_s3_bucket" "firehose_bucket" {
  bucket = "fh-stream-bucket"
  acl    = "private"
}

# DONE Allow Firehose to assume the role
data "aws_iam_policy_document" "firehose_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"

      identifiers = [
        "firehose.amazonaws.com",
      ]
    }

    actions = [
      "sts:AssumeRole",
    ]
  }
}

# DONE Role for Firehose (S3 + Kinesis)
resource "aws_iam_role" "firehose_stream_s3_role" {
  name = "firehose-stream-s3-role"
  assume_role_policy = data.aws_iam_policy_document.firehose_assume_role.json
}

# DONE Data for Firehose role (write to bucket, read from Kinesis)
data "aws_iam_policy_document" "firehose_write" {
  statement {
    resources = [
      aws_s3_bucket.firehose_bucket.arn,
      "${aws_s3_bucket.firehose_bucket.arn}/*",
    ]

    effect = "Allow"

    actions = [
      "s3:AbortMultipartUpload",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:ListBucketMultipartUploads",
      "s3:PutObject",
    ]
  }

  statement {
    resources = [
      aws_kinesis_stream.fh_stream.arn
    ]

    effect = "Allow"

    actions = [
      "kinesis:DescribeStream",
      "kinesis:GetShardIterator",
      "kinesis:GetRecords",
      "kinesis:ListShards",
    ]
  }
}

# Policy with data
DONE resource "aws_iam_policy" "policy_firehose_role" {
  name = "policy_firehose_role"
  policy = data.aws_iam_policy_document.firehose_write.json
}

# DONE Attach policy to the role
resource "aws_iam_role_policy_attachment" "policy_firehose_role" {
  role = aws_iam_role.firehose_stream_s3_role.name
  policy_arn = aws_iam_policy.policy_firehose_role.arn
}

# Firehose delivery stream
resource "aws_kinesis_firehose_delivery_stream" "fh-delivery-stream" {
  name        = "fh-delivery-stream"
  destination = "extended_s3"

  kinesis_source_configuration {
    kinesis_stream_arn = aws_kinesis_stream.fh_stream.arn
    role_arn = aws_iam_role.firehose_stream_s3_role.arn
  }

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_stream_s3_role.arn
    bucket_arn = aws_s3_bucket.firehose_bucket.arn

    buffer_size = 1
    buffer_interval = 60

    processing_configuration {
      enabled = "false"
    }
  }
}