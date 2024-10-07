package routes

import (
	authControllers "myapp/api/author/controllers"
	"myapp/api/author/entities"
	"myapp/middlewares"
	"myapp/scripts"
	"myapp/types"
)

type Routing struct {
	auth authControllers.AuthController
}

func (r Routing) RegisterRoutes() []types.Endpoint {
	var endpoints []types.Endpoint

	createAuthorSchema := types.ExtendSchema{
		Request: new(entities.CreateAuthorDto),
		Responses: map[int]interface{}{
			200: map[string]interface{}{
				"items": map[string]interface{}{
					"ID":    1,
					"Name":  "Author 1",
					"Email": "String",
					"Users": []map[string]interface{}{
						{
							"ID":   1,
							"Name": "User 1",
							"Profile": map[string]interface{}{
								"ID":   1,
								"Name": "Profile 1",
							},
						},
					},
				},
			},
			500: map[string]interface{}{
				"error":       "Internal Server Error",
				"status_code": 500,
			}},
		Description: "Create author",
		Tag:         "Authors",
	}
	scripts.RegisterEndpoint("/authors", "POST", "create", "create", createAuthorSchema, r.auth.CreateAuthorHandler, &endpoints, middlewares.ValidateRequestMiddleware(new(entities.CreateAuthorDto)))

	getListAuthorSchema := types.ExtendSchema{
		Parameters: []types.Parameter{
			{
				Name:        "limit",
				Type:        "integer",
				Location:    "query",
				Required:    false,
				Description: "Limit",
				Default:     10,
			},
		},
		Request: nil,
		Responses: map[int]interface{}{
			200: []map[string]interface{}{
				{
					"ID":   1,
					"Name": "Author 1",
					"Users": []map[string]interface{}{
						{
							"ID":   1,
							"Name": "User 1",
							"Profile": map[string]interface{}{
								"ID":   1,
								"Name": "Profile 1",
							},
						},
					},
					"author": map[string]interface{}{
						"ID":   1,
						"Name": "Author 1",
						"book": map[string]interface{}{
							"ID":   1,
							"Name": "Book 1",
							"category": map[string]interface{}{
								"ID":   1,
								"Name": "Category 1",
							},
						},
					},
				},
			},
			500: map[string]interface{}{
				"error":       "Internal Server Error",
				"status_code": 500,
			},
		},
		Description: "Get list author",
		Tag:         "Authors",
	}
	scripts.RegisterEndpoint("/authors", "GET", "get", "get", getListAuthorSchema, r.auth.GetListAuthorsHandler, &endpoints)
	return endpoints
}
