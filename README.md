# Cloud Trail logs to AWS Privileges converter

Ever wondered how much privileges you really need when using AWS. This tool helps
you build a minimal permissions file for AWS.

## Setup

1. Setup a cloudtrail log according to this [documentation](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-a-trail-using-the-console-first-time.html)
2. Export the cloudtrail logs in json
3. Run `./cloudtrailToAWSPrivileges -file cloudtrail.json > new-policy.json`
4. Run `aws iam create-policy --policy-name <fancy-policy-name> --policy-document file://new-policy.json`
