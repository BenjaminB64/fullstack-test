FROM golang:1.22
WORKDIR /app
COPY . .
RUN go work vendor

CMD ["go", "run", "main.go"]