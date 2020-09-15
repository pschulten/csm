# CSM

Silly tool for debugging AWS IAM policies.

Refer to the [official docs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-SDK-Metrics.html) for details

## Installation

```bash
go get github.com/pschulten/csm
```

or download the binary from the [releases](https://github.com/pschulten/csm/releases/)

## Usage
Run the server in one shell:
```bash
csm
```

Enable CSM and do some AWS API calls in another shell:
```bash
export AWS_CSM_ENABLED=true
aws kms describe-key --key-id 5492cc29-0843-40e7-87e5-c677959bc7dd
aws s3 ls
aws s3 ls s3://ALLOWED_BUCKET
aws s3 ls s3://FORBIDDEN_BUCKET
```

See your report in the first shell. Stop the server with `C-c`. 
```console
$ csm
200 KMS:DescribeKey                                    kms.eu-central-1.amazonaws.com
200 S3:ListBuckets                                     s3.eu-central-1.amazonaws.com
200 S3:ListObjectsV2                                   ALLOWED_BUCKET.s3.eu-central-1.amazonaws.com
403 S3:ListObjectsV2                                   FORBIDDEN_BUCKET.s3.eu-central-1.amazonaws.com
^C
$ 
```

