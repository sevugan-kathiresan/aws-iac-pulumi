package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/sevugan-kathiresan/aws-iac-pulumi/helper-modules/tags"
	infraCF "github.com/sevugan-kathiresan/aws-iac-pulumi/infra-modules/cloudfront"
	infraS3 "github.com/sevugan-kathiresan/aws-iac-pulumi/infra-modules/s3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Load the config from our stack (we only have dev.yaml)
		cfg := config.New(ctx, "")

		// Read the config values
		bucketName := cfg.Require("bucketName")
		creator := cfg.Require("creator")
		team := cfg.Require("team")
		service := cfg.Require("service")
		env := cfg.Require("env")

		// Creating the tags struct using out tags helper package (helper-modules/tags)
		t := tags.NewTags(creator, team, service, env)

		//Create the S3 bucket using our S3 package from infra-modules
		bucket, err := infraS3.NewBucket(ctx, bucketName, &infraS3.InputBucketArgs{
			BucketName: bucketName,
			Tags:       t,
		})

		if err != nil {
			return err
		}

		// Create CloudFront distribution
		distribution, err := infraCF.NewDistribution(ctx, "static-website-cf", &infraCF.InputDistributionArgs{
			BucketID:  bucket.BucketID,
			BucketArn: bucket.BucketArn,
			Tags:      t,
		})
		if err != nil {
			return err
		}

		// Export the CloudFront URL
		ctx.Export("websiteURL", distribution.DistributionDomain.ApplyT(func(domain string) string {
			return "https://" + domain
		}).(pulumi.StringOutput))

		return nil
	})
}
