package s3

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sevugan-kathiresan/aws-iac-pulumi/helper-modules/tags"
)

// Struct to hold the input argument for the bucket
type InputBucketArgs struct {
	BucketName string
	Tags       tags.Tags // Imported from our helper-modules
}

// Struct to hold the output of the created S3 bucket
type Bucket struct {
	pulumi.ResourceState                     // Embedding the pulumi resourcesate type without the field name. When you embed a type without a field name, Go automatically promotes all the fields and methods of that type into your struct.
	BucketID             pulumi.StringOutput // String.Output -> An asynchronous asynchronous "promise" that may not be known until after a resource is created.
	BucketArn            pulumi.StringOutput
}

func NewBucket(ctx *pulumi.Context, name string, args *InputBucketArgs, opts ...pulumi.ResourceOption) (*Bucket, error) {

	// Create a new bucket
	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		Bucket: pulumi.String(args.BucketName),
		Tags:   args.Tags.ToPulumiMap(),
	}, opts...)

	if err != nil {
		return nil, err
	}

	// Explicitly blocking public access -> Best Practice in IaC
	_, err = s3.NewBucketPublicAccessBlock(ctx, name+"-public-access-block", &s3.BucketPublicAccessBlockArgs{
		Bucket:                bucket.ID(), // This is what attaches the Access Block to the bucket
		BlockPublicAcls:       pulumi.Bool(true),
		BlockPublicPolicy:     pulumi.Bool(true),
		IgnorePublicAcls:      pulumi.Bool(true),
		RestrictPublicBuckets: pulumi.Bool(true),
	}, pulumi.Parent(bucket))

	if err != nil {
		return nil, err
	}

	// Enable Versioning
	_, err = s3.NewBucketVersioning(ctx, name+"-versioning", &s3.BucketVersioningArgs{
		Bucket: bucket.ID(),
		VersioningConfiguration: &s3.BucketVersioningVersioningConfigurationArgs{
			Status: pulumi.String("Enabled"),
		},
	}, pulumi.Parent(bucket))

	if err != nil {
		return nil, err
	}

	return &Bucket{
		BucketID:  bucket.ID().ToStringOutput(), // Pulumi returns the bucket id as the type IDOutput so we need to explicityly convert the type to StringOutput
		BucketArn: bucket.Arn,
	}, nil

}

/*
*** Educational Comment ***

NOTE: In Pulumi, resource ID is accessed via .ID() (method), not a struct field like .Arn.

Reason:
- ID is part of Pulumi’s core resource lifecycle
- Returned by the Create API (primary identifier)
- Managed by Pulumi engine/state

In contrast:
- Arn is a normal AWS attribute
- Usually returned via Describe/Get APIs

Mental model:
- ID() = lifecycle identity
- Arn  = AWS property

*/
