package openapi_generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const defaultNamePrefix = "nexustempmodule/apis"

type definition struct {
	schema       *extensionsv1.JSONSchemaProps
	dependencies []string
	resolved     bool
}

func (d *definition) fixEnumsTypes() {
	fixEnumsTypesInSchema(d.schema)
}

// fixEnumsTypesInSchema checks if schema desription has openAPIEnum annotation
// in it, and replaces it's type from `integer` to `string`. This is needed to
// properly handle enums from proto. See README for more info.
func fixEnumsTypesInSchema(schema *extensionsv1.JSONSchemaProps) {
	if strings.Contains(schema.Description, openAPIEnum) {
		schema.Description = strings.Trim(
			strings.ReplaceAll(schema.Description, openAPIEnum, ""), " ")
		if schema.Type == "integer" {
			schema.Type = "string"
			schema.Format = ""
			schema.Default = nil
		}
		if schema.Type == "array" {
			if schema.Items != nil {
				if schema.Items.Schema != nil {
					if schema.Items.Schema.Type == "integer" {
						schema.Items.Schema.Type = "string"
						schema.Items.Schema.Format = ""
						schema.Items.Schema.Default = nil
					}
					fixEnumsTypesInSchema(schema.Items.Schema)
				}
				if len(schema.Items.JSONSchemas) != 0 {
					toReplace := make([]extensionsv1.JSONSchemaProps, len(schema.Items.JSONSchemas))
					for i, s := range schema.Items.JSONSchemas {
						if s.Type == "integer" {
							s.Type = "string"
							s.Format = ""
							s.Default = nil
						}
						fixEnumsTypesInSchema(&s)
						toReplace[i] = s
					}
					schema.Items.JSONSchemas = toReplace
				}
			}
		}
	}
	for name := range schema.Properties {
		prop := schema.Properties[name]
		fixEnumsTypesInSchema(&prop)
		schema.Properties[name] = prop
	}
	additionalProperties := schema.AdditionalProperties
	if additionalProperties != nil {
		if additionalProperties.Schema != nil {
			fixEnumsTypesInSchema(additionalProperties.Schema)
		}
	}
	additionalItems := schema.AdditionalItems
	if additionalItems != nil {
		if additionalItems.Schema != nil {
			fixEnumsTypesInSchema(additionalItems.Schema)
		}
	}
}

// NewGenerator returns a new Generator instance and an error.
// rawDefinitions is a value returned by GetOpenAPIDefinitions functions created
// by openapi-gen from k8s.io/kube-openapi.
func NewGenerator(rawDefinitions map[string]common.OpenAPIDefinition) (*Generator, error) {
	for k, v := range defaultSchemas() {
		rawDefinitions[strings.ToLower(k)] = v
	}
	definitions := make(map[string]*definition, len(rawDefinitions))
	for name, rawDef := range rawDefinitions {
		definition, err := convertDefinition(rawDef)
		if err != nil {
			return nil, fmt.Errorf("converting schema %q: %v", name, err)
		}
		definition.fixEnumsTypes()
		definitions[strings.ToLower(name)] = definition
	}
	return &Generator{
		definitions:        definitions,
		missingDefinitions: map[string]struct{}{},
		namePrefix:         defaultNamePrefix,
	}, nil
}

func convertDefinition(input common.OpenAPIDefinition) (*definition, error) {
	serialized, err := json.Marshal(input.Schema.SchemaProps)
	if err != nil {
		return nil, fmt.Errorf("serializing schema: %v", err)
	}
	var schemaProps extensionsv1.JSONSchemaProps
	err = json.Unmarshal(serialized, &schemaProps)
	if err != nil {
		return nil, fmt.Errorf("deserializing schema: %v", err)
	}
	dependencies := make([]string, len(input.Dependencies))
	for i, dep := range input.Dependencies {
		dependencies[i] = strings.ToLower(dep)
	}
	return &definition{
		schema:       &schemaProps,
		dependencies: dependencies,
	}, nil
}

type Generator struct {
	definitions        map[string]*definition
	missingDefinitions map[string]struct{}
	namePrefix         string
}

// MissingDefinitions is a simple getter for g.missingDefinitions.
func (g *Generator) MissingDefinitions() map[string]struct{} {
	return g.missingDefinitions
}

// SetNamePrefix sets the namePrefix used in name creation.
// This should not be needed in day to day usage, useful for testing.
func (g *Generator) SetNamePrefix(namePrefix string) {
	g.namePrefix = namePrefix
}

// AddDefinition is a helper for adding default definitions.
// This should not be needed in day to day usage, useful for testing.
func (g *Generator) AddDefinition(name string, rawDef common.OpenAPIDefinition) error {
	def, err := convertDefinition(rawDef)
	if err != nil {
		return err
	}
	g.definitions[name] = def
	return nil
}

// ResolveRefs calls resolveRefsForPackage for each package which has not been
// resolved yet
func (g *Generator) ResolveRefs() error {
	for pkg, schema := range g.definitions {
		if schema.resolved {
			continue
		}
		err := g.resolveRefsForPackage(pkg)
		if err != nil {
			return fmt.Errorf("resolving refs for package %q: %v", pkg, err)
		}
	}
	return nil
}

func (g *Generator) resolveRefsForPackage(pkg string) error {
	pkgSchema := g.getDefinition(pkg)
	if pkgSchema.resolved {
		return nil
	}

	// We need to have all the subtrees resolved not to miss something
	for _, dep := range pkgSchema.dependencies {
		err := g.resolveRefsForPackage(dep)
		if err != nil {
			return fmt.Errorf("resolving dependency %v for schema %v: %v", dep, pkg, err)
		}
	}

	fmt.Printf("Resolving refs for %v\n", pkg)
	// We are forcing camelCase in all field names for consistency
	for property, propSchema := range pkgSchema.schema.Properties {
		if strings.Contains(property, "_") {
			pkgSchema.schema.Properties[convertToCamelCase(property)] = propSchema
			delete(pkgSchema.schema.Properties, property)
		}
		toReplace := make([]string, len(pkgSchema.schema.Required))
		for i, required := range pkgSchema.schema.Required {
			toReplace[i] = convertToCamelCase(required)
		}
		pkgSchema.schema.Required = toReplace
	}

	g.resolveRefsInProperty(pkgSchema.schema)
	g.resolveRefsInProperties(pkgSchema.schema)
	g.resolveRefsInAdditionalProperties(pkgSchema.schema)

	pkgSchema.resolved = true
	return nil
}

func (g *Generator) getDefinition(pkg string) *definition {
	pkg = strings.ToLower(pkg)
	_, ok := g.definitions[pkg]
	if !ok {
		// g.missingDefinitions[pkg] = struct{}{}
		fmt.Printf("Missing schema for %v, using default\n", pkg)
		g.definitions[pkg] = g.definitions["default"]
	}
	def := g.definitions[pkg]
	g.addKubernetesExtensionsFlags(def.schema)
	return def
}

// addKubernetesExtensionsFlags adds the following flags:
//   - x-kubernetes-preserve-unknown-fields - required properties of type object
//     with no properties specified
//   - x-kubernetes-int-or-string - required for properties which use anyOf with
//     integer or string values
func (g *Generator) addKubernetesExtensionsFlags(schema *extensionsv1.JSONSchemaProps) {
	if len(schema.Properties) == 0 && schema.AdditionalProperties == nil && schema.Type == "object" {
		t := true
		schema.XPreserveUnknownFields = &t
	}
	if len(schema.AnyOf) > 0 {
		schema.XIntOrString = true
	}
	for name, prop := range schema.Properties {
		g.addKubernetesExtensionsFlags(&prop)
		schema.Properties[name] = prop
	}
	if schema.Items != nil {
		if schema.Items.Schema != nil {
			g.addKubernetesExtensionsFlags(schema.Items.Schema)
		}
		for arrName, arrSchema := range schema.Items.JSONSchemas {
			g.addKubernetesExtensionsFlags(&arrSchema)
			schema.Items.JSONSchemas[arrName] = arrSchema
		}
	}
}

// resolveRefsInProperties checks if any of the properties' schemas are refs
// and replaces them with proper schema
func (g *Generator) resolveRefsInProperties(schema *extensionsv1.JSONSchemaProps) {
	toInline := map[string]extensionsv1.JSONSchemaProps{}
	toDelete := map[string]struct{}{}
	for property, propSchema := range schema.Properties {
		if propSchema.Ref != nil {
			refSchema := g.getDefinition(*propSchema.Ref).schema
			if _, ok := refSchema.Properties[inline]; ok {
				toDelete[inline] = struct{}{}
				for k, v := range refSchema.Properties {
					toInline[k] = v
				}
				toDelete[property] = struct{}{}
			} else {
				schema.Properties[property] = *refSchema
			}
		}
		g.resolveRefsInProperty(&propSchema)
		g.resolveRefsInProperties(&propSchema)
		g.resolveRefsInAdditionalProperties(&propSchema)
	}
	for k, v := range toInline {
		schema.Properties[k] = v
	}
	for v := range toDelete {
		delete(schema.Properties, v)
	}
	required := []string{}
	for _, v := range schema.Required {
		if _, ok := toDelete[v]; ok {
			continue
		}
		required = append(required, v)
	}
	schema.Required = required
}

// resolveRefsInAdditionalProperties checks if any of the additionalProperties'
// schemas are refs and replaces them with proper schema
func (g *Generator) resolveRefsInAdditionalProperties(schema *extensionsv1.JSONSchemaProps) {
	additionalProperties := schema.AdditionalProperties
	if additionalProperties != nil {
		if additionalProperties.Schema != nil {
			if additionalProperties.Schema.Ref != nil {
				additionalProperties.Schema = g.getDefinition(*additionalProperties.Schema.Ref).schema
			}
			g.resolveRefsInProperty(additionalProperties.Schema)
			g.resolveRefsInProperties(additionalProperties.Schema)
			g.resolveRefsInAdditionalProperties(additionalProperties.Schema)
		}
	}

	additionalItems := schema.AdditionalItems
	if additionalItems != nil {
		if additionalItems.Schema != nil {
			if additionalItems.Schema.Ref != nil {
				additionalItems.Schema = g.getDefinition(*additionalItems.Schema.Ref).schema
			}
			g.resolveRefsInProperty(additionalItems.Schema)
			g.resolveRefsInProperties(additionalItems.Schema)
			g.resolveRefsInAdditionalProperties(additionalItems.Schema)
		}
	}
}

// resolveRefsInProperty checks if property.Items schemas are refs and replaces them
// with proper schema
func (g *Generator) resolveRefsInProperty(propSchema *extensionsv1.JSONSchemaProps) {
	if propSchema.Items != nil {
		// in Items we can have either single schema defined in Schema, or a list of
		// them in JSONSchemas
		if propSchema.Items.Schema != nil {
			if propSchema.Items.Schema.Ref != nil {
				*propSchema.Items.Schema = *g.getDefinition(*propSchema.Items.Schema.Ref).schema
			}
		}
		newSchemas := make([]extensionsv1.JSONSchemaProps, len(propSchema.Items.JSONSchemas))
		for i, schema := range propSchema.Items.JSONSchemas {
			if schema.Ref == nil {
				newSchemas[i] = schema
				continue
			}
			newSchemas[i] = *g.getDefinition(*schema.Ref).schema
		}
		propSchema.Items.JSONSchemas = newSchemas
	}
}

func splitCRDs(content []byte) []string {
	return strings.Split(string(content), "---")
}

func CheckBackwardCompatibility(inCompatibleCRDs []*bytes.Buffer, crd extensionsv1.CustomResourceDefinition, oldCRDContent []byte) ([]*bytes.Buffer, error) {
	oldCRDParts := splitCRDs(oldCRDContent)
	for _, oldPart := range oldCRDParts {
		if oldPart == "" {
			continue
		}
		oldCRD := &extensionsv1.CustomResourceDefinition{}
		err := yaml.Unmarshal([]byte(oldPart), oldCRD)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling existing crd: %v", err)
		}

		newCRDPart, err := yaml.Marshal(&crd)
		if err != nil {
			return nil, fmt.Errorf("error marshaling new crd: %v", err)
		}

		if oldCRD.Name != crd.Name {
			continue
		}

		// When there is a backward incompatibility, we fail the build if we don't force an upgrade.
		isInCompatible, message, err := nexus_compare.CompareFiles([]byte(oldPart), newCRDPart)
		if err != nil {
			log.Errorf("Error occurred while checking CRD's %q backward compatibility: %v", crd.Name, err)
		}
		if isInCompatible {
			log.Warnf("CRD %q is incompatible with the previous version", crd.Name)
			inCompatibleCRDs = append(inCompatibleCRDs, message)
		}
	}
	return inCompatibleCRDs, nil
}

func (g *Generator) UpdateYAMLs(yamlsPath, oldYamlsPath string, force bool) error {
	var inCompatibleCRDs []*bytes.Buffer
	if err := filepath.Walk(yamlsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking files: %v", err)
		}

		if info.IsDir() {
			fmt.Printf("Skipping dir %q\n", path)
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file %q: %v", path, err)
		}

		parts := splitCRDs(content)
		crds := make([]extensionsv1.CustomResourceDefinition, len(parts))
		for _, part := range parts {
			if part == "" {
				continue
			}
			var crd extensionsv1.CustomResourceDefinition
			err := yaml.Unmarshal([]byte(part), &crd)
			if err != nil {
				return fmt.Errorf("unmarshalling: %v", err)
			}
			// TODO we assume here that we have only one version, which is correct
			// at the time of writing this, but may not be in the future.
			if len(crd.Spec.Versions) != 1 {
				return fmt.Errorf("crd %v has no, or more than 1 versions", crd.Name)
			}
			// TODO another hack. I don't know how to prevent Status field being
			// generated in YAML file, so we add StoredVersions by hand
			crd.Status.StoredVersions = []string{crd.Spec.Versions[0].Name}
			err = g.addCustomResourceValidation(&crd)
			if err != nil {
				return err
			}

			/*
				yamlsPath - indicates the directory path of new crd yamls
						Ex: the path will be `_generated/crds`

				path - indicates the file path of node
				        Ex: For the node `Config` the path will be `_generated/crds/config_config.yaml`

				oldYamlsPath - indicates the directory path of the existing crd yamls
						Ex: the path will be `tmp/old-files/old-crds`
			*/
			oldFilePath := oldYamlsPath + strings.TrimPrefix(path, yamlsPath)
			oldCRDContent, err := os.ReadFile(oldFilePath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("error reading the existing crd file on path %q: %v", oldFilePath, err)
				}
				log.Debugf("Can't find the existing crd file on path %s", oldFilePath)
			}

			if oldCRDContent != nil {
				if inCompatibleCRDs, err = CheckBackwardCompatibility(inCompatibleCRDs, crd, oldCRDContent); err != nil {
					return err
				}
			}
			crds = append(crds, crd)
		}

		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("opening file %q for write: %v", path, err)
		}
		defer f.Close()
		for _, crd := range crds {
			// TODO IDK why this happens. It this out for now.
			if crd.Name == "" {
				continue
			}
			fmt.Printf("Writing schema %q to %q\n", crd.Name, path)
			_, err := f.Write([]byte("---\n"))
			if err != nil {
				return fmt.Errorf("writing separator: %v", err)
			}
			serialized, err := yaml.Marshal(crd)
			if err != nil {
				return fmt.Errorf("serializing crd %q: %v", crd.Name, err)
			}
			_, err = f.Write(serialized)
			if err != nil {
				return fmt.Errorf("writing crd %q: %v", crd.Name, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if inCompatibleCRDs != nil {
		if !force {
			// If the CRD are incompatible with the previous version, this will fail the build.
			return fmt.Errorf("datamodel upgrade failed due to backward incompatible changes:\n %v", inCompatibleCRDs)
		}
		log.Warnf("Upgrading the data model that is incompatible with the previous version: %v", inCompatibleCRDs)
	}
	return nil
}

func (g *Generator) addCustomResourceValidation(crd *extensionsv1.CustomResourceDefinition) error {
	name := g.createName(crd.Spec.Group, crd.Spec.Versions[0].Name, crd.Spec.Names.Kind)
	schemaProps := g.getDefinition(name).schema
	toReplace := []string{}
	for _, propName := range schemaProps.Required {
		// We want only "spec" and "data" fields required in schemas. "status" and all
		// others should remain optional.
		if propName == "spec" || propName == "data" {
			toReplace = append(toReplace, propName)
		}
	}
	schemaProps.Required = toReplace
	// We need to ensure that crd.spec.versions[*].schema.openapiV3Schema["metadata"]
	// does not contain any fields except for `type: object`. This is required by K8s
	// API and we added x-kubernetes-preserve-unknown-fields: true` here in `getDefinition`,
	// so we are fixing our wrongly added flag.
	if len(schemaProps.Properties) == 0 {
		schemaProps.Properties = map[string]extensionsv1.JSONSchemaProps{}
	}
	schemaProps.Properties["metadata"] = extensionsv1.JSONSchemaProps{
		Type: "object",
	}

	crd.Spec.PreserveUnknownFields = false
	crd.Spec.Versions[0].Schema = &extensionsv1.CustomResourceValidation{OpenAPIV3Schema: schemaProps}
	return nil
}

func (g *Generator) createName(group, apiVersion, name string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s/%s.%s", g.namePrefix, group, apiVersion, name))
}

func convertToCamelCase(input string) string {
	if !strings.Contains(input, "_") {
		return input
	}
	parts := strings.Split(input, "_")
	camelCaseName := parts[0]
	for _, p := range parts[1:] {
		camelCaseName += strings.ToUpper(string(p[0]))
		camelCaseName += p[1:]
	}
	return camelCaseName
}
