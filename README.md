# Fullstack test
Coding challenge designed specifically for backend developers skilled in React / Typescript and Go

# Challenge Description

The objective is to create a small engine of job processors designed to automate manual actions. For this exercise, we wish to automate the checking of the weather as well as the status of the opening of the Chaban-Delmas Bridge.

There is a service "jobservice" to which you can send a job to be processed. The job can be of two types: "weather" or "bridge".

There is a service "jobprocessor" that will process the jobs asynchronously. The jobprocessor will call the "weather" and "bridge" services to get the information and then send the result to the "jobservice".

# REST API
See the swagger file

# Run the services
You can run the services using the following commands:

```bash
docker compose up
```

A database will be created and populated with some data.
