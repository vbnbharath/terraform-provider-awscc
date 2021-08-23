// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"strings"
	"text/template"

	cfschema "github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go"
	"github.com/hashicorp/terraform-provider-awscc/internal/naming"
	"github.com/hashicorp/terraform-provider-awscc/internal/provider/generators/resource/codegen"
	"github.com/mitchellh/cli"
)

var (
	cfTypeSchemaFile = flag.String("cfschema", "", "CloudFormation resource type schema file; required")
	packageName      = flag.String("package", "", "override package name for generated code")
	tfResourceType   = flag.String("resource", "", "Terraform resource type; required")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tmain.go [flags] -resource <TF-resource-type> -cfschema <CF-type-schema-file> <generated-schema-file> <generated-acctests-file>\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 || *tfResourceType == "" || *cfTypeSchemaFile == "" {
		flag.Usage()
		os.Exit(2)
	}

	destinationPackage := os.Getenv("GOPACKAGE")
	if *packageName != "" {
		destinationPackage = *packageName
	}

	schemaFilename := args[0]
	acctestsFilename := args[1]

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	generator := NewGenerator(ui, *tfResourceType, *cfTypeSchemaFile)

	if err := generator.Generate(destinationPackage, schemaFilename, acctestsFilename); err != nil {
		ui.Error(fmt.Sprintf("error generating Terraform %s resource: %s", *tfResourceType, err))
		os.Exit(1)
	}
}

type Generator struct {
	cfTypeSchemaFile string
	tfResourceType   string
	ui               cli.Ui
}

func NewGenerator(ui cli.Ui, tfResourceType, cfTypeSchemaFile string) *Generator {
	return &Generator{
		cfTypeSchemaFile: cfTypeSchemaFile,
		tfResourceType:   tfResourceType,
		ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
}

// Generate generates the resource's type factory into the specified file.
func (g *Generator) Generate(packageName, schemaFilename, acctestsFilename string) error {
	g.infof("generating Terraform resource code for %[1]q from %[2]q into %[3]q and %[4]q", g.tfResourceType, g.cfTypeSchemaFile, schemaFilename, acctestsFilename)

	// Create target directories.
	for _, filename := range []string{schemaFilename, acctestsFilename} {
		dirname := path.Dir(filename)
		err := os.MkdirAll(dirname, 0755)

		if err != nil {
			return fmt.Errorf("creating target directory %s: %w", dirname, err)
		}
	}

	resource, err := NewResource(g.tfResourceType, g.cfTypeSchemaFile)

	if err != nil {
		return fmt.Errorf("reading CloudFormation resource schema for %s: %w", g.tfResourceType, err)
	}

	cfTypeName := *resource.CfResource.TypeName
	org, svc, res, err := naming.ParseCloudFormationTypeName(cfTypeName)

	if err != nil {
		return fmt.Errorf("incorrect format for CloudFormation Resource Provider Schema type name: %s", cfTypeName)
	}

	// e.g. "logGroupResourceType"
	factoryFunctionName := string(bytes.ToLower([]byte(res[:1]))) + res[1:] + "ResourceType"

	// e.g. "TestAccAWSLogsLogGroup"
	acceptanceTestFunctionPrefix := fmt.Sprintf("TestAcc%[1]s%[2]s%[3]s", org, svc, res)

	sb := strings.Builder{}
	codeEmitter := codegen.Emitter{
		CfResource: resource.CfResource,
		Ui:         g.ui,
		Writer:     &sb,
	}

	// Generate code for the CloudFormation root properties schema.
	attributeNameMap := make(map[string]string) // Terraform attribute name to CloudFormation property name.
	codeFeatures, err := codeEmitter.EmitRootPropertiesSchema(attributeNameMap)

	if err != nil {
		return fmt.Errorf("emitting schema code: %w", err)
	}

	rootPropertiesSchema := sb.String()
	sb.Reset()

	templateData := TemplateData{
		AcceptanceTestFunctionPrefix: acceptanceTestFunctionPrefix,
		AttributeNameMap:             attributeNameMap,
		CloudFormationTypeName:       cfTypeName,
		FactoryFunctionName:          factoryFunctionName,
		HasRequiredAttribute:         true,
		HasUpdateMethod:              true,
		PackageName:                  packageName,
		RootPropertiesSchema:         rootPropertiesSchema,
		SchemaVersion:                1,
		SyntheticIDAttribute:         true,
		TerraformTypeName:            resource.TfType,
	}

	if codeFeatures&codegen.HasRequiredRootProperty == 0 {
		templateData.HasRequiredAttribute = false
	}
	if codeFeatures&codegen.HasUpdatableProperty == 0 {
		templateData.HasUpdateMethod = false
	}
	if codeFeatures&codegen.UsesInternalTypes > 0 {
		templateData.ImportInternalTypes = true
	}
	if codeFeatures&codegen.UsesValidation > 0 {
		templateData.ImportValidate = true
	}
	if codeFeatures&codegen.HasIDRootProperty > 0 {
		templateData.SyntheticIDAttribute = false
	}

	if description := resource.CfResource.Description; description != nil {
		templateData.SchemaDescription = *description
	}

	for _, path := range resource.CfResource.WriteOnlyProperties {
		templateData.WriteOnlyPropertyPaths = append(templateData.WriteOnlyPropertyPaths, string(path))
	}

	if v, ok := resource.CfResource.Handlers[cfschema.HandlerTypeCreate]; ok {
		templateData.CreateTimeoutInMinutes = v.TimeoutInMinutes
	}
	if v, ok := resource.CfResource.Handlers[cfschema.HandlerTypeUpdate]; ok {
		templateData.UpdateTimeoutInMinutes = v.TimeoutInMinutes
	}
	if v, ok := resource.CfResource.Handlers[cfschema.HandlerTypeDelete]; ok {
		templateData.DeleteTimeoutInMinutes = v.TimeoutInMinutes
	}

	err = g.applyAndWriteTemplate(schemaFilename, resourceSchemaTemplateBody, &templateData)

	if err != nil {
		return err
	}

	err = g.applyAndWriteTemplate(acctestsFilename, acceptanceTestsTemplateBody, &templateData)

	if err != nil {
		return err
	}

	return nil
}

// applyAndWriteTemplate applies the template body to the specified data and writes it to file.
func (g *Generator) applyAndWriteTemplate(filename, templateBody string, templateData *TemplateData) error {
	tmpl, err := template.New("function").Parse(templateBody)

	if err != nil {
		return fmt.Errorf("parsing function template: %w", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		g.infof("%s", buffer.String())
		return fmt.Errorf("formatting generated source code: %w", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("creating file (%s): %w", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		return fmt.Errorf("writing to file (%s): %w", filename, err)
	}

	return nil
}

func (g *Generator) infof(format string, a ...interface{}) {
	g.ui.Info(fmt.Sprintf(format, a...))
}

type TemplateData struct {
	AcceptanceTestFunctionPrefix string
	AttributeNameMap             map[string]string
	CloudFormationTypeName       string
	CreateTimeoutInMinutes       int
	DeleteTimeoutInMinutes       int
	FactoryFunctionName          string
	HasRequiredAttribute         bool
	HasUpdateMethod              bool
	ImportInternalTypes          bool
	ImportValidate               bool
	PackageName                  string
	RootPropertiesSchema         string
	SchemaDescription            string
	SchemaVersion                int64
	SyntheticIDAttribute         bool
	TerraformTypeName            string
	UpdateTimeoutInMinutes       int
	WriteOnlyPropertyPaths       []string
}

// Terraform resource schema definition.
var resourceSchemaTemplateBody = `
// Code generated by generators/resource/main.go; DO NOT EDIT.

package {{ .PackageName }}

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tflog "github.com/hashicorp/terraform-plugin-log"
	. "github.com/hashicorp/terraform-provider-awscc/internal/generic"
	"github.com/hashicorp/terraform-provider-awscc/internal/registry"
	{{ if .ImportInternalTypes }}providertypes "github.com/hashicorp/terraform-provider-awscc/internal/types"{{- end }}
	{{ if .ImportValidate }}"github.com/hashicorp/terraform-provider-awscc/internal/validate"{{- end }}
)

func init() {
	registry.AddResourceTypeFactory("{{ .TerraformTypeName }}", {{ .FactoryFunctionName }})
}

// {{ .FactoryFunctionName }} returns the Terraform {{ .TerraformTypeName }} resource type.
// This Terraform resource type corresponds to the CloudFormation {{ .CloudFormationTypeName }} resource type.
func {{ .FactoryFunctionName }}(ctx context.Context) (tfsdk.ResourceType, error) {
	attributes := {{ .RootPropertiesSchema }}

{{ if .SyntheticIDAttribute }}
	attributes["id"] = tfsdk.Attribute{
		Description: "Uniquely identifies the resource.",
		Type:        types.StringType,
		Computed:    true,
	}
{{- end }}

	schema := tfsdk.Schema{
		Description: "{{ .SchemaDescription }}",
		Version:     {{ .SchemaVersion }},
		Attributes:  attributes,
	}

	var opts ResourceTypeOptions

	opts = opts.WithCloudFormationTypeName("{{ .CloudFormationTypeName }}").WithTerraformTypeName("{{ .TerraformTypeName }}")
	opts = opts.WithTerraformSchema(schema)
	opts = opts.WithSyntheticIDAttribute({{ .SyntheticIDAttribute }})
	opts = opts.WithAttributeNameMap(map[string]string{
{{- range $key, $value := .AttributeNameMap }}
		"{{ $key }}": "{{ $value }}",
{{- end }}
	})
{{ if not .HasUpdateMethod }}
	opts = opts.IsImmutableType(true)
{{- end }}
{{ if .WriteOnlyPropertyPaths }}
	opts = opts.WithWriteOnlyPropertyPaths([]string{
	{{- range .WriteOnlyPropertyPaths }}
		"{{ . }}",
	{{- end }}
	})
{{- end }}
	opts = opts.WithCreateTimeoutInMinutes({{ .CreateTimeoutInMinutes }}).WithDeleteTimeoutInMinutes({{ .DeleteTimeoutInMinutes }})
{{ if .HasUpdateMethod }}
	opts = opts.WithUpdateTimeoutInMinutes({{ .UpdateTimeoutInMinutes }})
{{- end }}

	resourceType, err := NewResourceType(ctx, opts...)

	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Generated schema", "tfTypeName", "{{ .TerraformTypeName }}", "schema", hclog.Fmt("%v", schema))

	return resourceType, nil
}
`

// Terraform acceptance tests.
var acceptanceTestsTemplateBody = `
// Code generated by generators/resource/main.go; DO NOT EDIT.

package {{ .PackageName }}_test

import (
	{{ if .HasRequiredAttribute }}"regexp"{{- end }}
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-awscc/internal/acctest"
)

{{ if .HasRequiredAttribute }}
func {{ .AcceptanceTestFunctionPrefix }}_basic(t *testing.T) {
	td := acctest.NewTestData(t, "{{ .CloudFormationTypeName }}", "{{ .TerraformTypeName }}", "test")

	td.ResourceTest(t, []resource.TestStep{
		{
			Config: td.EmptyConfig(),
			ExpectError: regexp.MustCompile("Missing required argument"),
		},
	})
}
{{- else }}
func {{ .AcceptanceTestFunctionPrefix }}_basic(t *testing.T) {
	td := acctest.NewTestData(t, "{{ .CloudFormationTypeName }}", "{{ .TerraformTypeName }}", "test")

	td.ResourceTest(t, []resource.TestStep{
		{
			Config: td.EmptyConfig(),
			Check: resource.ComposeTestCheckFunc(
				td.CheckExistsInAWS(),
			),
		},
	})
}

func {{ .AcceptanceTestFunctionPrefix }}_disappears(t *testing.T) {
	td := acctest.NewTestData(t, "{{ .CloudFormationTypeName }}", "{{ .TerraformTypeName }}", "test")

	td.ResourceTest(t, []resource.TestStep{
		{
			Config: td.EmptyConfig(),
			Check: resource.ComposeTestCheckFunc(
				td.CheckExistsInAWS(),
				td.DeleteResource(),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}
{{- end }}
`

type Resource struct {
	CfResource *cfschema.Resource
	TfType     string
}

func NewResource(resourceType, cfTypeSchemaFile string) (*Resource, error) {
	resourceSchema, err := cfschema.NewResourceJsonSchemaPath(cfTypeSchemaFile)

	if err != nil {
		return nil, fmt.Errorf("reading CloudFormation Resource Type Schema: %w", err)
	}

	resource, err := resourceSchema.Resource()

	if err != nil {
		return nil, fmt.Errorf("parsing CloudFormation Resource Type Schema: %w", err)
	}

	if err := resource.Expand(); err != nil {
		return nil, fmt.Errorf("expanding JSON Pointer references: %w", err)
	}

	return &Resource{
		CfResource: resource,
		TfType:     resourceType,
	}, nil
}
