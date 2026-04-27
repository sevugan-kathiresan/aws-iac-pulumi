package cloudfront

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sevugan-kathiresan/aws-iac-pulumi/helper-modules/tags"
)

// InputDistributionArgs holds the configuration for creating a CloudFront distribution.
type InputDistributionArgs struct {
	BucketID  pulumi.StringOutput
	BucketArn pulumi.StringOutput
	Tags      tags.Tags
}

// Distribution holds the outputs of the created CloudFront distribution.
type Distribution struct {
	pulumi.ResourceState
	DistributionID     pulumi.StringOutput
	DistributionDomain pulumi.StringOutput
}

// NewDistribution creates a CloudFront distribution with OAC for the given S3 bucket.
func NewDistribution(ctx *pulumi.Context, name string, args *InputDistributionArgs, opts ...pulumi.ResourceOption) (*Distribution, error) {

	// Create Origin Access Control (OAC)
	oac, err := cloudfront.NewOriginAccessControl(ctx, name+"-oac", &cloudfront.OriginAccessControlArgs{
		Name:                          pulumi.String(name + "-oac"), // Label for the OAC
		OriginAccessControlOriginType: pulumi.String("s3"),
		SigningBehavior:               pulumi.String("always"), // CloudFront will always sign the request
		SigningProtocol:               pulumi.String("sigv4"),  // AWS Signature Version 4
	}, opts...)

	if err != nil {
		return nil, err
	}

	// Create CloudFront distribution and link the OAC to the distribution
	distribution, err := cloudfront.NewDistribution(ctx, name, &cloudfront.DistributionArgs{
		Enabled:           pulumi.Bool(true),
		DefaultRootObject: pulumi.String("index.html"),
		PriceClass:        pulumi.String("PriceClass_200"),
		//Origin Definition
		Origins: cloudfront.DistributionOriginArray{
			&cloudfront.DistributionOriginArgs{
				OriginId: args.BucketID, // used to link cache behaviour to origin
				DomainName: args.BucketID.ApplyT(func(id string) string { // creating s3 domain name from bucket id -> this domain is a the S3 rest endpoint not website URL
					return id + ".s3.amazonaws.com"
				}).(pulumi.StringOutput),
				OriginAccessControlId: oac.ID(),
			},
		},
		DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
			TargetOriginId:       args.BucketID, // used to link cache behaviour to origin
			ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
			AllowedMethods: pulumi.StringArray{ // Since we are using this distribution to host a static website, we are only allowing read methods
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			CachedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{ // Defines what values are forwarded to the origin
				QueryString: pulumi.Bool(false),
				Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
					Forward: pulumi.String("none"),
				},
			},
		},
		Restrictions: &cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},
		ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
			CloudfrontDefaultCertificate: pulumi.Bool(true),
		},
		Tags: args.Tags.ToPulumiMap(),
	}, pulumi.Parent(oac))
	if err != nil {
		return nil, err
	}

	// Attach bucket policy to allow CloudFront OAC access
	_, err = s3.NewBucketPolicy(ctx, name+"-bucket-policy", &s3.BucketPolicyArgs{
		Bucket: args.BucketID,
		Policy: pulumi.All(args.BucketArn, distribution.Arn).ApplyT(func(vals []interface{}) string {
			// pulumi.All(...) combines multiple Outputs (BucketArn, DistributionArn) into a single Output that resolves when BOTH values are ready.
			// The resolved values are passed as a []interface{} (generic slice). Generic because this interface{} does not have any methods defined.
			// vals[0] corresponds to args.BucketArn, vals[1] to distribution.Arn
			// We use type assertions (.(string)) to convert from interface{} to string
			// because Pulumi returns a generic container interface{} even if we know the types.
			bucketArn := vals[0].(string)
			distributionArn := vals[1].(string)
			return `{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "AllowCloudFrontOAC",
					"Effect": "Allow",
					"Principal": {
						"Service": "cloudfront.amazonaws.com"
					},
					"Action": "s3:GetObject",
					"Resource": "` + bucketArn + `/*",
					"Condition": {
						"StringEquals": {
							"AWS:SourceArn": "` + distributionArn + `"
						}
					}
				}]
			}`
		}).(pulumi.StringOutput),
	})
	if err != nil {
		return nil, err
	}

	return &Distribution{
		DistributionID:     distribution.ID().ToStringOutput(),
		DistributionDomain: distribution.DomainName,
	}, nil

}
