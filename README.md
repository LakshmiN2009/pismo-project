# Pismo Project

## Project Structure
```
pismo-project/
├── src/
│   ├── cmd/
│   │   └── main.go
│   ├── db/
│   │   └── db.go
│   ├── handler/
│   │   └── api_handler.go
│   │   └── api_handler_test.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── model/
│   │   └── account/
│   │       └── account.go
│   ├── service/
│   │   └── accountService/
│   │       └── accountService.go
│   ├── docker/
│   │   └── Dockerfile
│   │   └── docker-compose.yml
│   ├── go.mod
│   ├── go.sum

```

## Configurations
- **Go Version:** v1.23
- **Docker Version:** v27.3.1
- **Docker Compose Version:** v2.29.1
- **Postgres Version:** v14.13

### To Run the Project

1. **Initialize the Docker Daemon:**
   Ensure that the Docker daemon is running. You can start it with:
   ```bash
   sudo systemctl start docker
   ```

2. **Navigate to the Docker Directory:**
   Change to the directory containing the Docker configuration:
   ```bash
   cd src/docker
   ```

3. **Build and Start the Services:**
   Use Docker Compose to build and start the services:
   ```bash
   docker-compose up --build
   ``` 

This will set up the project and its dependencies, making the API available for use.
   ```

## API Details

The service provides the following APIs:

1. **Create Account**  
   **POST /accounts**  
   Adds a new account.

2. **Get Account Details**  
   **GET /accounts/:account_id**  
   Retrieves details for a specified account.

3. **Create Transaction**  
   **POST /transactions**  
   Adds a new transaction for an account based on the operation type.

### CURL Examples

1. **Create an Account**
   ```bash
   curl -X POST http://localhost:8080/accounts \
   -H 'Content-Type: application/json' \
   -d '{"document_number": "12345678900"}'
   ```

2. **Get Account Details**
   ```bash
   curl -X GET http://localhost:8080/accounts/1
   ```

3. **Create a Transaction**
   ```bash
   curl -X POST http://localhost:8080/transactions \
   -H 'Content-Type: application/json' \
   -d '{"account_id": 1, "operation_type_id": 4, "amount": 123.45}'
   ```

## Testing

Unit tests have been added for the API endpoints in `api_handler_test.go`. After starting the services with Docker Compose, run the test cases to validate both positive and negative scenarios.

---

This README provides instructions on running the service and details for using the API endpoints, including example curl commands. Test coverage ensures the API handles various scenarios correctly.
