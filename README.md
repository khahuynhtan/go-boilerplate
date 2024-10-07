# How to run the source code

```shell
docker compose up -d
```

- Run docker compose

```shell
docker compose up -d --build
```

## Open Swagger API

[http://localhost:8081](http://localhost:8081)

## Track logs

```
docker compose logs --follow app
docker compose logs --follow db
```

## Create migrations
Let's create table called `users`:
```shell
  docker exec -it $(docker ps -q -f name=app) sh -c "migrate create -ext sql -dir migrations -seq create_users_table"
```
If there were no errors, we should have two files available under `migrations` folder:
- 000001_create_users_table.down.sql
- 000001_create_users_table.up.sql

Note the `sql` extension that we provided.

In the `.up.sql` file let's create the table:
```sql
CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
);
```
And in the `.down.sql` let's delete it:
```sql
DROP TABLE IF EXISTS users;
```
By adding `IF EXISTS/IF NOT EXISTS` we are making migrations idempotent

## Run migrations
```shell
  docker exec -it $(docker ps -q -f name=app) sh -c "migrate -database ${POSTGRESQL_URL} -path migrations up"
```
Let's check if the table was created properly by running `psql example -c "\d users"` or using the `Database editor`
## Define a route
In the `routes/routing.go` 
```js

func (r Routing) RegisterRoutes() []types.Endpoint {
	var endpoints []types.Endpoint

	createAuthorSchema := types.ExtendSchema{
		Request: new(entities.CreateAuthorDto), // the api's input, it will be validated
		Responses: map[int]interface{}{         // defined the api's output
			200: map[string]interface{}{
				"items": map[string]interface{}{
					"ID":    1,
					"Name":  "Author 1",
					"Email": "String",
				},
			},
			500: map[string]interface{}{
				"error":       "Internal Server Error",
				"status_code": 500,
			}},
		Description: "Create an author",
		Tag:         "Authors",
		IsAuth:      utils.BoolPtr(false),  // whether need to Authorize
	}
	scripts.RegisterEndpoint("/authors", "POST", "Create an author", "createAuthorId", createAuthorSchema, r.auth.CreateAuthorHandler, &endpoints, middlewares.ValidateRequestMiddleware(new(entities.CreateAuthorDto)))
  return endpoints
}
```
- Note: `createAuthorSchema` will be used to generate the swagger docs

In the `entities/dtos`
```js
type CreateAuthorDto struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
}
```
All endpoints defined in the `RegisterRoutes` function will be reflected in the [http://localhost:8081](Swagger UI)