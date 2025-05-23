# Go Database Application

## Overview
This project is a simple Go application that demonstrates best practices for manipulating database tables. It includes a user management system with functionalities for creating and retrieving users.

## Project Structure
```
go-db-app
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── db
│   │   └── db.go       # Database connection logic
│   ├── models
│   │   └── user.go      # User model definition
│   ├── repository
│   │   └── user_repository.go # User repository for database operations
│   └── service
│       └── user_service.go    # Business logic for user operations
├── pkg
│   └── utils
│       └── logger.go    # Logging utility
├── go.mod                # Module definition
├── go.sum                # Module dependency checksums
└── README.md             # Project documentation
```

## Setup Instructions
1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-db-app
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Set up the database:**
   Ensure you have a database set up and update the connection details in `internal/db/db.go`.

4. **Run the application:**
   ```
   go run cmd/main.go
   ```

## Usage
- **Register a User:**
  Send a POST request to `/register` with user details (name, email, password).

- **Get User by ID:**
  Send a GET request to `/user/{id}` to retrieve user information.

- **Get All Users:**
  Send a GET request to `/users` to retrieve a list of all users.

## Logging
The application uses a simple logging utility located in `pkg/utils/logger.go`. You can use the `Info`, `Error`, and `Debug` functions to log messages at different levels.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License.