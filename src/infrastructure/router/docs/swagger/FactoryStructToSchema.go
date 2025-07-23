package swagger

import (
	"fmt"
	"reflect"
	"strings"
)

type FactoryStructToSchema struct {
	seen    map[reflect.Type]string
	schemas map[string]Schema
}

func NewFactoryStructToSchema() *FactoryStructToSchema {
	return &FactoryStructToSchema{
		seen:    make(map[reflect.Type]string),
		schemas: make(map[string]Schema),
	}
}

func (f *FactoryStructToSchema) Components() *Components {
	return &Components{
		Schemas: f.schemas,
	}
}

func (f *FactoryStructToSchema) MakeSchema(root any) (map[string]Schema, *Schema, error) {
	t := reflect.TypeOf(root)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	ref, isVector, err := f.collectSchema(t)
	if err != nil {
		return f.schemas, nil, err
	}

	if isVector {
		return f.schemas, &Schema{
			Items: &Schema{
				Ref: ref,
			},
		}, nil
	}

	return f.schemas, &Schema{
		Ref: ref,
	}, nil
}

func (f *FactoryStructToSchema) collectSchema(t reflect.Type) (string, bool, error) {
	isVector := f.isVector(t)

	t = f.deferencePointer(t)

	if t.Kind() != reflect.Struct {
		return "", isVector, nil
	}

	if ref, ok := f.seen[t]; ok {
		return ref, isVector, nil
	}

	properties := make(map[string]*Schema)
	required := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			continue
		}

		isJson, jsonTag, jsonOmitempty := f.isJsonField(field)
		if !isJson {
			continue
		}

		propRef, err := f.inferSchema(field.Type)
		if err != nil {
			return "", isVector, err
		}

		propRef.Description = field.Tag.Get("description")
		
		properties[jsonTag] = propRef

		if !jsonOmitempty && f.required(field) {
			required = append(required, jsonTag)
		}
	}

	name := t.Name()
	if name == "" {
		name = fmt.Sprintf("Anon%s", t.PkgPath())
	}

	ref := f.makeRefString(name)

	f.seen[t] = ref

	f.schemas[name] = Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}

	return ref, isVector, nil
}

func (f *FactoryStructToSchema) isJsonField(field reflect.StructField) (bool, string, bool) {
	jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
	if jsonTag == "-" {
		return false, "", false
	}

	if jsonTag == "" {
		jsonTag = field.Name
	}

	return true, jsonTag, strings.Contains(field.Tag.Get("json"), "omitempty")
}

func (f *FactoryStructToSchema) required(field reflect.StructField) bool {
	return field.Type.Kind() != reflect.Ptr &&
		field.Type.Kind() != reflect.Slice &&
		field.Type.Kind() != reflect.Map 
}

func (f *FactoryStructToSchema) inferSchema(fieldType reflect.Type) (*Schema, error) {
	switch fieldType.Kind() {
	case reflect.Ptr:
		return f.inferSchema(fieldType.Elem())
	case reflect.Struct:
		return f.inferStruct(fieldType)
	case reflect.Slice, reflect.Array:
		return f.inferArray(fieldType)
	case reflect.Map:
		return f.inferMap(fieldType)
	case reflect.String:
		return &Schema{Type: "string"}, nil
	case reflect.Bool:
		return &Schema{Type: "boolean"}, nil
	case reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer"}, nil
	case reflect.Float32, reflect.Float64:
		return &Schema{Type: "number"}, nil
	default:
		return &Schema{Type: "string"}, nil
	}
}

func (f *FactoryStructToSchema) inferStruct(fieldType reflect.Type) (*Schema, error) {
	ref, isVector, err := f.collectSchema(fieldType)
	if err != nil {
		return nil, err
	}

	if isVector {
		return &Schema{
			Items: &Schema{
				Ref: ref,
			},
		}, nil
	}

	return &Schema{Ref: ref}, nil
}

func (f *FactoryStructToSchema) inferArray(fieldType reflect.Type) (*Schema, error) {
	itemRef, err := f.inferSchema(fieldType.Elem())
	if err != nil {
		return nil, err
	}

	return &Schema{
		Type:  "array",
		Items: itemRef,
	}, nil
}

func (f *FactoryStructToSchema) inferMap(fieldType reflect.Type) (*Schema, error) {
	properties, err := f.inferSchema(fieldType.Elem())
	if err != nil {
		return nil, err
	}

	return &Schema{
		Type: "object",
		AdditionalProperties: properties,
	}, nil
}

func (f *FactoryStructToSchema) deferencePointer(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}
	return t
}

func (f *FactoryStructToSchema) isVector(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

func (f *FactoryStructToSchema) makeRefString(name string) string {
	return fmt.Sprintf("#/components/schemas/%s", name)
}
