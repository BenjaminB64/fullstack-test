# Fullstack test
Coding challenge designed specifically for backend developers skilled in React / Typescript and Go

# Challenge Description

The objective is to create a small engine of job processors designed to automate manual actions. For this exercise, we wish to automate the checking of the weather as well as the status of the opening of the Chaban-Delmas Bridge.

There is a service "jobservice" to which you can send a job to be processed. The job can be of two types: "weather" or "bridge".

There is a service "jobprocessor" that will process the jobs asynchronously. The jobprocessor will call the "weather" and "bridge" services to get the information and then send the result to the "jobservice".

# REST API
See the swagger file (back/jobservice/docs/swagger.json) for the REST API documentation.

# Run the services
You can run the services using the following commands:

```bash
docker compose up
```
and then you can access the front-end at http://localhost:3000

The database is automatically created and populated with the necessary schema.

# Generate swagger documentation
You can generate the swagger documentation using the following command:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
cd back/jobservice
swag init -g infrastructure/http_server/http.go -ot json,yaml
```

# Generate typescript client from swagger
You can generate the typescript client using the following command:

```bash
npm install @openapitools/openapi-generator-cli -g
openapi-generator-cli generate -g typescript-fetch -i ./back/jobservice/docs/swagger.json -o ./front/src/apiClient --additional-properties=
```