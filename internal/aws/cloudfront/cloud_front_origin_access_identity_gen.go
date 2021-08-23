// Code generated by generators/resource/main.go; DO NOT EDIT.

package cloudfront

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tflog "github.com/hashicorp/terraform-plugin-log"
	. "github.com/hashicorp/terraform-provider-awscc/internal/generic"
	"github.com/hashicorp/terraform-provider-awscc/internal/registry"
)

func init() {
	registry.AddResourceTypeFactory("awscc_cloudfront_cloud_front_origin_access_identity", cloudFrontOriginAccessIdentityResourceType)
}

// cloudFrontOriginAccessIdentityResourceType returns the Terraform awscc_cloudfront_cloud_front_origin_access_identity resource type.
// This Terraform resource type corresponds to the CloudFormation AWS::CloudFront::CloudFrontOriginAccessIdentity resource type.
func cloudFrontOriginAccessIdentityResourceType(ctx context.Context) (tfsdk.ResourceType, error) {
	attributes := map[string]tfsdk.Attribute{
		"cloud_front_origin_access_identity_config": {
			// Property: CloudFrontOriginAccessIdentityConfig
			// CloudFormation resource type schema:
			// {
			//   "additionalProperties": false,
			//   "properties": {
			//     "Comment": {
			//       "type": "string"
			//     }
			//   },
			//   "required": [
			//     "Comment"
			//   ],
			//   "type": "object"
			// }
			Attributes: tfsdk.SingleNestedAttributes(
				map[string]tfsdk.Attribute{
					"comment": {
						// Property: Comment
						Type:     types.StringType,
						Required: true,
					},
				},
			),
			Required: true,
		},
		"id": {
			// Property: Id
			// CloudFormation resource type schema:
			// {
			//   "type": "string"
			// }
			Type:     types.StringType,
			Computed: true,
		},
		"s3_canonical_user_id": {
			// Property: S3CanonicalUserId
			// CloudFormation resource type schema:
			// {
			//   "type": "string"
			// }
			Type:     types.StringType,
			Computed: true,
		},
	}

	schema := tfsdk.Schema{
		Description: "Resource Type definition for AWS::CloudFront::CloudFrontOriginAccessIdentity",
		Version:     1,
		Attributes:  attributes,
	}

	var opts ResourceTypeOptions

	opts = opts.WithCloudFormationTypeName("AWS::CloudFront::CloudFrontOriginAccessIdentity").WithTerraformTypeName("awscc_cloudfront_cloud_front_origin_access_identity")
	opts = opts.WithTerraformSchema(schema)
	opts = opts.WithSyntheticIDAttribute(false)
	opts = opts.WithAttributeNameMap(map[string]string{
		"cloud_front_origin_access_identity_config": "CloudFrontOriginAccessIdentityConfig",
		"comment":              "Comment",
		"id":                   "Id",
		"s3_canonical_user_id": "S3CanonicalUserId",
	})

	opts = opts.WithCreateTimeoutInMinutes(0).WithDeleteTimeoutInMinutes(0)

	opts = opts.WithUpdateTimeoutInMinutes(0)

	resourceType, err := NewResourceType(ctx, opts...)

	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Generated schema", "tfTypeName", "awscc_cloudfront_cloud_front_origin_access_identity", "schema", hclog.Fmt("%v", schema))

	return resourceType, nil
}
