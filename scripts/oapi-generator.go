package scripts

import (
	"fmt"
	"myapp/types"
	"myapp/utils"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type SecurityScheme struct {
	Type         string `yaml:"type"`
	Scheme       string `yaml:"scheme,omitempty"`
	BearerFormat string `yaml:"bearerFormat,omitempty"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes"`
}
type Server struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description,omitempty"`
}
type MediaType struct {
	Schema interface{} `yaml:"schema"`
}

type OpenAPI struct {
	OpenAPI    string              `yaml:"openapi"`
	Info       Info                `yaml:"info"`
	Paths      map[string]PathItem `yaml:"paths"`
	Components Components          `yaml:"components"`
	Servers    []Server            `yaml:"servers,omitempty"`
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
	Security    []map[string][]string     `yaml:"security,omitempty"`
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

func getDefaultValue(val interface{}) string {
	defaultValue := detectType(val)
	if val != "" {
		defaultValue = fmt.Sprintf("%v", val)
	}

	return defaultValue
}

func generateSchema(data interface{}) *Property {
	property := &Property{}

	// Use reflect to detect the type of data
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		property.Type = "object"
		property.Properties = make(map[string]*Property)
		for key, value := range data.(map[string]interface{}) {
			property.Properties[key] = generateSchema(value)
		}
	case reflect.Slice, reflect.Array:
		property.Type = "array"
		if value.Len() > 0 {
			property.Items = generateSchema(value.Index(0).Interface())
		}
	case reflect.Struct:
		property.Type = "object"
		property.Properties = make(map[string]*Property)
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			fieldVal := value.Field(i).Interface()
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				jsonTag = strings.ToLower(field.Name)
			}
			property.Properties[jsonTag] = generateSchema(fieldVal)
		}
	case reflect.String:
		property.Type = "string"

		property.Default = getDefaultValue(data)
	case reflect.Float64, reflect.Int, reflect.Int64:
		property.Type = "number"
		property.Default = getDefaultValue(data)
	case reflect.Bool:
		property.Type = "boolean"
		property.Default = getDefaultValue(data)
	default:
		utils.Logger(utils.InfoLevel, "Unknown")
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
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"bearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080",
				Description: "Local server",
			},
		},
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
		operation.Security = nil

		if endpoint.ExtendSchema.IsAuth == nil || *endpoint.ExtendSchema.IsAuth {
			operation.Security = []map[string][]string{
				{
					"bearerAuth": {},
				},
			}
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
			reqValue := reflect.ValueOf(endpoint.ExtendSchema.Request).Elem()

			for i := 0; i < reqValue.NumField(); i++ {
				field := reqValue.Type().Field(i)
				fieldVal := reqValue.Field(i).Interface()

				jsonTag := field.Tag.Get("json")
				if jsonTag == "" || jsonTag == "-" {
					jsonTag = strings.ToLower(field.Name)
				}

				fieldType := reflect.TypeOf(fieldVal)
				if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
					// get element type of array
					elemType := fieldType.Elem()
					requestSchema[jsonTag] = Property{
						Type:  "array",
						Items: generateSchema(reflect.New(elemType).Elem().Interface()),
					}
				} else {
					requestSchema[jsonTag] = Property{
						Type: detectType(fieldVal),
					}
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
		for statusCode, responseValue := range endpoint.ExtendSchema.Responses {
			responseContent := map[string]*Property{}
			value := reflect.ValueOf(responseValue)
			switch value.Kind() {
			case reflect.Map:
				for key, val := range value.Interface().(map[string]interface{}) {
					responseContent[key] = generateSchema(val)
				}

			case reflect.Slice, reflect.Array:
				operation.Responses[fmt.Sprintf("%d", statusCode)] = ResponseSchema{
					Description: fmt.Sprintf("Response for status code %d", statusCode),
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type":  "array",
								"items": generateSchema(value.Index(0).Interface()),
							},
						},
					},
				}
				continue

			case reflect.Float64, reflect.Int, reflect.Int64, reflect.String, reflect.Bool:
				// only return the type of the response, no need inside the object
				operation.Responses[fmt.Sprintf("%d", statusCode)] = ResponseSchema{
					Description: fmt.Sprintf("Response for status code %d", statusCode),
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type":    detectType(responseValue),
								"default": responseValue,
							},
						},
					},
				}
				continue

			default:
				responseContent["default"] = generateSchema(value.Interface())
			}

			operation.Responses[fmt.Sprintf("%d", statusCode)] = ResponseSchema{
				Description: fmt.Sprintf("Response for status code %d", statusCode),
				Content: map[string]MediaType{
					"application/json": {
						Schema: map[string]interface{}{
							"type":       detectType(responseValue),
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
