# Development stage
FROM golang:1.23.4 AS development

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy Go source code
COPY cmd/server/*.go ./
COPY . ./ 

RUN GOOS=linux CGO_ENABLED=0 go build -o /mock-ses ./cmd/server

# Production stage
FROM alpine:latest

ARG DB_NAME
ARG DB_HOST
ARG DB_PASSWORD
ARG DB_PORT
ARG BUILD

# Set environment variables
ENV CRUD_DB_NAME=$DB_NAME
ENV CRUD_DB_HOST=$DB_HOST
ENV CRUD_DB_PASSWORD=$DB_PASSWORD
ENV CRUD_DB_PORT=$DB_PORT
ENV CRUD_BUILD=$BUILD

# Copy the Go binary and other required files from the development stage
COPY --from=development /mock-ses /

EXPOSE 8080

CMD ["/mock-ses"]