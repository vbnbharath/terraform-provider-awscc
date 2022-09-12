// Code generated by generators/resource/main.go; DO NOT EDIT.

package gamelift

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/hashicorp/terraform-provider-awscc/internal/generic"
	"github.com/hashicorp/terraform-provider-awscc/internal/registry"
	"github.com/hashicorp/terraform-provider-awscc/internal/validate"
)

func init() {
	registry.AddResourceTypeFactory("awscc_gamelift_alias", aliasResourceType)
}

// aliasResourceType returns the Terraform awscc_gamelift_alias resource type.
// This Terraform resource type corresponds to the CloudFormation AWS::GameLift::Alias resource type.
func aliasResourceType(ctx context.Context) (provider.ResourceType, error) {
	attributes := map[string]tfsdk.Attribute{
		"alias_id": {
			// Property: AliasId
			// CloudFormation resource type schema:
			// {
			//   "description": "Unique alias ID",
			//   "type": "string"
			// }
			Description: "Unique alias ID",
			Type:        types.StringType,
			Computed:    true,
			PlanModifiers: []tfsdk.AttributePlanModifier{
				resource.UseStateForUnknown(),
			},
		},
		"description": {
			// Property: Description
			// CloudFormation resource type schema:
			// {
			//   "description": "A human-readable description of the alias.",
			//   "maxLength": 1024,
			//   "minLength": 1,
			//   "type": "string"
			// }
			Description: "A human-readable description of the alias.",
			Type:        types.StringType,
			Optional:    true,
			Computed:    true,
			Validators: []tfsdk.AttributeValidator{
				validate.StringLenBetween(1, 1024),
			},
			PlanModifiers: []tfsdk.AttributePlanModifier{
				resource.UseStateForUnknown(),
			},
		},
		"name": {
			// Property: Name
			// CloudFormation resource type schema:
			// {
			//   "description": "A descriptive label that is associated with an alias. Alias names do not need to be unique.",
			//   "maxLength": 1024,
			//   "minLength": 1,
			//   "pattern": ".*\\S.*",
			//   "type": "string"
			// }
			Description: "A descriptive label that is associated with an alias. Alias names do not need to be unique.",
			Type:        types.StringType,
			Required:    true,
			Validators: []tfsdk.AttributeValidator{
				validate.StringLenBetween(1, 1024),
				validate.StringMatch(regexp.MustCompile(".*\\S.*"), ""),
			},
		},
		"routing_strategy": {
			// Property: RoutingStrategy
			// CloudFormation resource type schema:
			// {
			//   "additionalProperties": false,
			//   "anyOf": [
			//     {
			//       "required": [
			//         "FleetId"
			//       ]
			//     },
			//     {
			//       "required": [
			//         "Message"
			//       ]
			//     }
			//   ],
			//   "description": "A routing configuration that specifies where traffic is directed for this alias, such as to a fleet or to a message.",
			//   "properties": {
			//     "FleetId": {
			//       "description": "A unique identifier for a fleet that the alias points to. If you specify SIMPLE for the Type property, you must specify this property.",
			//       "pattern": "^fleet-\\S+",
			//       "type": "string"
			//     },
			//     "Message": {
			//       "description": "The message text to be used with a terminal routing strategy. If you specify TERMINAL for the Type property, you must specify this property.",
			//       "type": "string"
			//     },
			//     "Type": {
			//       "description": "Simple routing strategy. The alias resolves to one specific fleet. Use this type when routing to active fleets.",
			//       "enum": [
			//         "SIMPLE",
			//         "TERMINAL"
			//       ],
			//       "type": "string"
			//     }
			//   },
			//   "required": [
			//     "Type"
			//   ],
			//   "type": "object"
			// }
			Description: "A routing configuration that specifies where traffic is directed for this alias, such as to a fleet or to a message.",
			Attributes: tfsdk.SingleNestedAttributes(
				map[string]tfsdk.Attribute{
					"fleet_id": {
						// Property: FleetId
						Description: "A unique identifier for a fleet that the alias points to. If you specify SIMPLE for the Type property, you must specify this property.",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						Validators: []tfsdk.AttributeValidator{
							validate.StringMatch(regexp.MustCompile("^fleet-\\S+"), ""),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.UseStateForUnknown(),
						},
					},
					"message": {
						// Property: Message
						Description: "The message text to be used with a terminal routing strategy. If you specify TERMINAL for the Type property, you must specify this property.",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.UseStateForUnknown(),
						},
					},
					"type": {
						// Property: Type
						Description: "Simple routing strategy. The alias resolves to one specific fleet. Use this type when routing to active fleets.",
						Type:        types.StringType,
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							validate.StringInSlice([]string{
								"SIMPLE",
								"TERMINAL",
							}),
						},
					},
				},
			),
			Required: true,
			Validators: []tfsdk.AttributeValidator{
				validate.RequiredAttributes(
					validate.AnyOfRequired(
						validate.Required(
							"fleet_id",
						),
						validate.Required(
							"message",
						),
					),
				),
			},
		},
	}

	attributes["id"] = tfsdk.Attribute{
		Description: "Uniquely identifies the resource.",
		Type:        types.StringType,
		Computed:    true,
		PlanModifiers: []tfsdk.AttributePlanModifier{
			resource.UseStateForUnknown(),
		},
	}

	schema := tfsdk.Schema{
		Description: "The AWS::GameLift::Alias resource creates an alias for an Amazon GameLift (GameLift) fleet destination.",
		Version:     1,
		Attributes:  attributes,
	}

	var opts ResourceTypeOptions

	opts = opts.WithCloudFormationTypeName("AWS::GameLift::Alias").WithTerraformTypeName("awscc_gamelift_alias")
	opts = opts.WithTerraformSchema(schema)
	opts = opts.WithSyntheticIDAttribute(true)
	opts = opts.WithAttributeNameMap(map[string]string{
		"alias_id":         "AliasId",
		"description":      "Description",
		"fleet_id":         "FleetId",
		"message":          "Message",
		"name":             "Name",
		"routing_strategy": "RoutingStrategy",
		"type":             "Type",
	})

	opts = opts.WithCreateTimeoutInMinutes(0).WithDeleteTimeoutInMinutes(0)

	opts = opts.WithUpdateTimeoutInMinutes(0)

	resourceType, err := NewResourceType(ctx, opts...)

	if err != nil {
		return nil, err
	}

	return resourceType, nil
}
