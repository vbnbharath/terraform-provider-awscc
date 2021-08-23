// Code generated by generators/resource/main.go; DO NOT EDIT.

package iotsitewise

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
	registry.AddResourceTypeFactory("awscc_iotsitewise_access_policy", accessPolicyResourceType)
}

// accessPolicyResourceType returns the Terraform awscc_iotsitewise_access_policy resource type.
// This Terraform resource type corresponds to the CloudFormation AWS::IoTSiteWise::AccessPolicy resource type.
func accessPolicyResourceType(ctx context.Context) (tfsdk.ResourceType, error) {
	attributes := map[string]tfsdk.Attribute{
		"access_policy_arn": {
			// Property: AccessPolicyArn
			// CloudFormation resource type schema:
			// {
			//   "description": "The ARN of the access policy.",
			//   "type": "string"
			// }
			Description: "The ARN of the access policy.",
			Type:        types.StringType,
			Computed:    true,
		},
		"access_policy_id": {
			// Property: AccessPolicyId
			// CloudFormation resource type schema:
			// {
			//   "description": "The ID of the access policy.",
			//   "type": "string"
			// }
			Description: "The ID of the access policy.",
			Type:        types.StringType,
			Computed:    true,
		},
		"access_policy_identity": {
			// Property: AccessPolicyIdentity
			// CloudFormation resource type schema:
			// {
			//   "additionalProperties": false,
			//   "description": "The identity for this access policy. Choose either an SSO user or group or an IAM user or role.",
			//   "properties": {
			//     "IamRole": {
			//       "additionalProperties": false,
			//       "description": "Contains information for an IAM role identity in an access policy.",
			//       "properties": {
			//         "arn": {
			//           "description": "The ARN of the IAM role.",
			//           "type": "string"
			//         }
			//       },
			//       "type": "object"
			//     },
			//     "IamUser": {
			//       "additionalProperties": false,
			//       "description": "Contains information for an IAM user identity in an access policy.",
			//       "properties": {
			//         "arn": {
			//           "description": "The ARN of the IAM user.",
			//           "type": "string"
			//         }
			//       },
			//       "type": "object"
			//     },
			//     "User": {
			//       "additionalProperties": false,
			//       "description": "Contains information for a user identity in an access policy.",
			//       "properties": {
			//         "id": {
			//           "description": "The AWS SSO ID of the user.",
			//           "type": "string"
			//         }
			//       },
			//       "type": "object"
			//     }
			//   },
			//   "type": "object"
			// }
			Description: "The identity for this access policy. Choose either an SSO user or group or an IAM user or role.",
			Attributes: tfsdk.SingleNestedAttributes(
				map[string]tfsdk.Attribute{
					"iam_role": {
						// Property: IamRole
						Description: "Contains information for an IAM role identity in an access policy.",
						Attributes: tfsdk.SingleNestedAttributes(
							map[string]tfsdk.Attribute{
								"arn": {
									// Property: arn
									Description: "The ARN of the IAM role.",
									Type:        types.StringType,
									Optional:    true,
								},
							},
						),
						Optional: true,
					},
					"iam_user": {
						// Property: IamUser
						Description: "Contains information for an IAM user identity in an access policy.",
						Attributes: tfsdk.SingleNestedAttributes(
							map[string]tfsdk.Attribute{
								"arn": {
									// Property: arn
									Description: "The ARN of the IAM user.",
									Type:        types.StringType,
									Optional:    true,
								},
							},
						),
						Optional: true,
					},
					"user": {
						// Property: User
						Description: "Contains information for a user identity in an access policy.",
						Attributes: tfsdk.SingleNestedAttributes(
							map[string]tfsdk.Attribute{
								"id": {
									// Property: id
									Description: "The AWS SSO ID of the user.",
									Type:        types.StringType,
									Optional:    true,
								},
							},
						),
						Optional: true,
					},
				},
			),
			Required: true,
		},
		"access_policy_permission": {
			// Property: AccessPolicyPermission
			// CloudFormation resource type schema:
			// {
			//   "description": "The permission level for this access policy. Valid values are ADMINISTRATOR or VIEWER.",
			//   "type": "string"
			// }
			Description: "The permission level for this access policy. Valid values are ADMINISTRATOR or VIEWER.",
			Type:        types.StringType,
			Required:    true,
		},
		"access_policy_resource": {
			// Property: AccessPolicyResource
			// CloudFormation resource type schema:
			// {
			//   "additionalProperties": false,
			//   "description": "The AWS IoT SiteWise Monitor resource for this access policy. Choose either portal or project but not both.",
			//   "properties": {
			//     "Portal": {
			//       "additionalProperties": false,
			//       "description": "A portal resource.",
			//       "properties": {
			//         "id": {
			//           "description": "The ID of the portal.",
			//           "type": "string"
			//         }
			//       },
			//       "type": "object"
			//     },
			//     "Project": {
			//       "additionalProperties": false,
			//       "description": "A project resource.",
			//       "properties": {
			//         "id": {
			//           "description": "The ID of the project.",
			//           "type": "string"
			//         }
			//       },
			//       "type": "object"
			//     }
			//   },
			//   "type": "object"
			// }
			Description: "The AWS IoT SiteWise Monitor resource for this access policy. Choose either portal or project but not both.",
			Attributes: tfsdk.SingleNestedAttributes(
				map[string]tfsdk.Attribute{
					"portal": {
						// Property: Portal
						Description: "A portal resource.",
						Attributes: tfsdk.SingleNestedAttributes(
							map[string]tfsdk.Attribute{
								"id": {
									// Property: id
									Description: "The ID of the portal.",
									Type:        types.StringType,
									Optional:    true,
								},
							},
						),
						Optional: true,
					},
					"project": {
						// Property: Project
						Description: "A project resource.",
						Attributes: tfsdk.SingleNestedAttributes(
							map[string]tfsdk.Attribute{
								"id": {
									// Property: id
									Description: "The ID of the project.",
									Type:        types.StringType,
									Optional:    true,
								},
							},
						),
						Optional: true,
					},
				},
			),
			Required: true,
		},
	}

	attributes["id"] = tfsdk.Attribute{
		Description: "Uniquely identifies the resource.",
		Type:        types.StringType,
		Computed:    true,
	}

	schema := tfsdk.Schema{
		Description: "Resource schema for AWS::IoTSiteWise::AccessPolicy",
		Version:     1,
		Attributes:  attributes,
	}

	var opts ResourceTypeOptions

	opts = opts.WithCloudFormationTypeName("AWS::IoTSiteWise::AccessPolicy").WithTerraformTypeName("awscc_iotsitewise_access_policy")
	opts = opts.WithTerraformSchema(schema)
	opts = opts.WithSyntheticIDAttribute(true)
	opts = opts.WithAttributeNameMap(map[string]string{
		"access_policy_arn":        "AccessPolicyArn",
		"access_policy_id":         "AccessPolicyId",
		"access_policy_identity":   "AccessPolicyIdentity",
		"access_policy_permission": "AccessPolicyPermission",
		"access_policy_resource":   "AccessPolicyResource",
		"arn":                      "arn",
		"iam_role":                 "IamRole",
		"iam_user":                 "IamUser",
		"id":                       "id",
		"portal":                   "Portal",
		"project":                  "Project",
		"user":                     "User",
	})

	opts = opts.WithCreateTimeoutInMinutes(0).WithDeleteTimeoutInMinutes(0)

	opts = opts.WithUpdateTimeoutInMinutes(0)

	resourceType, err := NewResourceType(ctx, opts...)

	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Generated schema", "tfTypeName", "awscc_iotsitewise_access_policy", "schema", hclog.Fmt("%v", schema))

	return resourceType, nil
}
