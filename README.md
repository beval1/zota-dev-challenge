# Zota Dev Challenge

"*Your task is to integrate Zota as a payment provider. You will need to follow our public
documentation here. Your solution should implement two flows: deposit (non-credit card) and
status check. We love good test coverage so at least one unit test is mandatory. You can
assume there is a service that will call your implementation directly. The way you go about
structuring is up to your imagination and mastery.*"

### Requirements
* Tested on Go 1.22.0
* Internally using `uber/fx` for dependency injection, `go/chi` for routing, `uber/zap` for logging

### Running the project
* You will need .env file in the root directory, refer to `config.go`
* Run `go run cmd/main.go` to start the server.

### Running with docker
* Building the image `docker build -t zota-challenge .`
* Running the image `docker run -p 8080:8080 zota-challenge`

### Tests
* Run `go test ./...` to run all tests.
* Run `go test -cover ./...` to run all tests with coverage.
* The mocked dependencies were generated automatically using `mockgen`

### Api Documentation
* The API documentation was generated automatically using `swaggo/swag` 
  * To generate again - `swag init --dir cmd,internal`
* You can find the OpenAPI specification in the `docs` folder.
* You can also run swagger-ui locally with the specification file.
  * Run the program
  * Go to `http://localhost:8080/swagger/index.html` and paste this into the search bar `http://localhost:8080/swagger/doc.json`

### Folder Structure - inspired by DDD
* `cmd`: Contains the main application code.
* `internal`: Contains the internal packages.
    * `deposit`: Contains the deposit flow.
    * `status`: Contains the status flow.
    * `config`: Contains the configuration for the application.
* `docs`: Contains the OpenAPI specification.



