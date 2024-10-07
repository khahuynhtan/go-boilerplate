package utils

import (
	"fmt"
	"myapp/types"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type MediaType struct {
	Schema interface{} `yaml:"schema"`
}

type OpenAPI struct {
	OpenAPI string              `yaml:"openapi"`
	Info    Info                `yaml:"info"`
	Paths   map[string]PathItem `yaml:"paths"`
}

type Info struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

type PathItem struct {
	Post   *Operation `yaml:"post,omitempty"`
	Get    *Operation `yaml:"get,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty"`
}

type ResponseSchema struct {
	Description string               `yaml:"description,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
}

type Operation struct {
	ID          string                    `yaml:"operationId"`
	Tags        []string                  `yaml:"tags,omitempty"`
	Summary     string                    `yaml:"summary,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	RequestBody *RequestBody              `yaml:"requestBody,omitempty"`
	Responses   map[string]ResponseSchema `yaml:"responses,omitempty"`
	Parameters  []map[string]interface{}  `yaml:"parameters,omitempty"`
}

type RequestBody struct {
	Content map[string]MediaType `yaml:"content,omitempty"`
}
type Property struct {
	Type        string               `yaml:"type,omitempty"`
	Default     interface{}          `yaml:"default,omitempty"`
	Properties  map[string]*Property `yaml:"properties,omitempty"`
	Items       *Property            `yaml:"items,omitempty"`
	Required    []string             `yaml:"required,omitempty"`
	Description *string              `yaml:"description,omitempty"`
	Example     interface{}          `yaml:"example,omitempty"`
}

func detectType(val interface{}) string {
	value := reflect.ValueOf(val)
	switch value.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Map:
		return "object"
	case reflect.Slice, reflect.Array:
		return "array"
	default:
		return "unknown"
	}
}
func generateSchema(data interface{}) *Property {
	property := &Property{}

	// Use reflect to detect the type of data
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		Logger(InfoLevel, "Slice or Array 2")
		property.Type = "object"
		property.Properties = make(map[string]*Property)
		for key, value := range data.(map[string]interface{}) {
			property.Properties[key] = generateSchema(value)
		}
	case reflect.Slice, reflect.Array:
		property.Type = "array"
		if value.Len() > 0 {
			Logger(InfoLevel, "Slice or Array")

			property.Items = generateSchema(value.Index(0).Interface())
		}
	case reflect.String:
		Logger(InfoLevel, "String")
		property.Type = "string"
		property.Default = data
	case reflect.Float64, reflect.Int, reflect.Int64:
		Logger(InfoLevel, "Number")
		property.Type = "number"
		property.Default = data
	case reflect.Bool:
		Logger(InfoLevel, "Boolean")
		property.Type = "boolean"
		property.Default = data
	default:
		Logger(InfoLevel, "Unknown")
	}
	return property
}

func GenerateOpenAPI(endpoints []types.Endpoint, title string, version string) (*OpenAPI, error) {
	openAPI := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:   title,
			Version: version,
		},
		Paths: make(map[string]PathItem),
	}

	for _, endpoint := range endpoints {
		pathItem, exists := openAPI.Paths[endpoint.Path]
		if !exists {
			pathItem = PathItem{}
		}

		operation := Operation{
			ID:          endpoint.OperationID,
			Tags:        []string{endpoint.ExtendSchema.Tag},
			Summary:     endpoint.Summary,
			Description: endpoint.ExtendSchema.Description,
			Responses:   make(map[string]ResponseSchema),
		}

		// Handle Parameters
		if endpoint.ExtendSchema.Parameters != nil {
			for _, param := range endpoint.ExtendSchema.Parameters {
				operation.Parameters = append(operation.Parameters, map[string]interface{}{
					"name":        param.Name,
					"in":          param.Location,
					"required":    param.Required,
					"description": param.Description,
					"schema": map[string]interface{}{
						"type":    param.Type,
						"default": param.Default,
					},
				})
			}
		}

		// Handle Request body
		if endpoint.ExtendSchema.Request != nil {
			requestSchema := map[string]Property{}
			// get value of request (dereference pointer)
			reqValue := reflect.ValueOf(endpoint.ExtendSchema.Request).Elem()
			for i := 0; i < reqValue.NumField(); i++ {
				field := reqValue.Type().Field(i)
				fieldVal := reqValue.Field(i).Interface()
				requestSchema[field.Name] = Property{
					Type: detectType(fieldVal),
				}
			}

			operation.RequestBody = &RequestBody{
				Content: map[string]MediaType{
					"application/json": {
						Schema: map[string]interface{}{
							"type":       "object",
							"properties": requestSchema,
						},
					},
				},
			}
		}

		// Handle Response body
		for statusCode, responseSchema := range endpoint.ExtendSchema.Responses {
			responseContent := map[string]*Property{}
			// detect type of responseSchema
			value := reflect.ValueOf(responseSchema)
			switch value.Kind() {
			case reflect.Map:
				for key, val := range value.Interface().(map[string]interface{}) {
					responseContent[key] = generateSchema(val)
				}

			case reflect.Slice, reflect.Array:
				responseContent["items"] = generateSchema(value.Index(0).Interface())

			default:
				responseContent["default"] = generateSchema(value.Index(0).Interface())
			}

			operation.Responses[fmt.Sprintf("%d", statusCode)] = ResponseSchema{
				Description: fmt.Sprintf("Response for status code %d", statusCode),
				Content: map[string]MediaType{
					"application/json": {
						Schema: map[string]interface{}{
							"type":       "object",
							"properties": responseContent,
						},
					},
				},
			}
		}

		switch endpoint.Method {
		case "POST":
			pathItem.Post = &operation
		case "GET":
			pathItem.Get = &operation
		case "PUT":
			pathItem.Put = &operation
		case "DELETE":
			pathItem.Delete = &operation
		case "PATCH":
			pathItem.Patch = &operation
		}

		openAPI.Paths[endpoint.Path] = pathItem
	}
	out, err := yaml.Marshal(openAPI)
	if err != nil {
		panic(err)
	}
	os.WriteFile("api.yaml", out, 0644)
	return openAPI, nil
}
