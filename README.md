# Store Backend Services using Go

## Overview

This repository contains a Golang backend Store and Product service that provides various functionalities related to store and products management. 
The service is built using the Go programming language and follows best practices for building scalable and maintainable backend services.

## Features

The Store service includes the following features:

- Store Management: Allows users to create store, update store, show store and show product list in the store including filter, pagination and searching.

- Product Management : Allows users to create product, update product, show product, delete product and show all product list including filter, pagination and searching.

## Project Structure

- cmd/ # Main application entry point
- internal/ # Internal application packages
- - internal/adapter/ # Adapters for database and external services
- - - internal/adapter/repository/ # Adapters for database
- - internal/business/ # Business layer
- - - internal/business/domain/ # Business domain entities and data structures
- - - internal/business/port/ # Business layer
- - - internal/business/service/ # Business logic layer
- - internal/handler/ # Service implementations
- - - internal/handler/http/ # Http Service implementations
- migration/ # Database migration files
- pkg/ # Shared packages and utilities
- .env # Environment variables
- .env.example # Environment variables example
- .gitignore # Git ignore rules
- go.mod # Go module file
- go.sum # Go module checksum file
- Makefile # Makefile for build and run
- README.md # Readme file

## Installation and Setup

1. Clone the repository:
    `git clone git@github.com:ijlik/store-app.git`
2. Install the required dependencies:
   `go mod download`
3. Set up the environment variables:
   `cp .env.example .env`
4. Set up the database config in .env file and run the migrations:
   `make migrate-up`
5. Run the application:
   `make local-run`
6. Run the linter:
   `make lint`

## Testing
To run the unit tests, use the following command:
    `make test`

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvement, please feel free to open an issue or submit a pull request.