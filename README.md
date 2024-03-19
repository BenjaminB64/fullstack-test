# Fullstack test
Coding challenge designed specifically for backend developers skilled in React / Typescript and Go

# Challenge Description

The objective is to create a small engine of job processors designed to automate manual actions. For this exercise, we wish to automate the checking of the weather as well as the status of the opening of the Chaban-Delmas Bridge.

There is a service "jobservice" to which you can send a job to be processed. The job can be of two types: "weather" or "bridge".

There is a service "jobprocessor" that will process the jobs asynchronously. The jobprocessor will call the "weather" and "bridge" services to get the information and then send the result to the "jobservice".

# REST API
See the swagger file (back/jobservice/docs/swagger.json) for the REST API documentation.

# Setup .env
If you want to use docker-compose, just copy the template.env file to .env
    
```bash
cp ./back/jobservice/template.env ./back/jobservice/.env
```

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

# Generate protobuf go files
You can generate the protobuf go files using the following command:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

protoc --proto_path=./back/common/protobuf --go_out=./back/common/protobuf/jobs-proto --go_opt=paths=source_relative --go-grpc_out=./back/common/protobuf/jobs-proto --go-grpc_opt=paths=source_relative jobs.proto
```
