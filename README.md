# AWS SQS Dumper

CLI script to dump messages from an AWS SQS queue into a directory with one file for each message.

# Installation

    go install github.com/mattrx/sqs-dumper

# Usage

    Usage:
      sqs-dumper [flags]

    Flags:
      -h, --help                       help for sqs-dumper
          --loop-count int32           number of loops for receive (default 1000)
      -o, --output string              output directory (default "./messages")
      -p, --profile string             aws profile
      -q, --queue string               queue url
          --visibility-timeout int32   visibility timeout in seconds (default 60)

You have to provide an AWS profile configured in `~/.aws/accounts` and the queue url in the form of `https://sqs.eu-central-1.amazonaws.com/{accountID}/{queueName}`.
