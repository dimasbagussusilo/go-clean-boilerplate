
# Go Clean Boilerplate

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A RESTful API built with Go following Clean Architecture principles. This project demonstrates how to build maintainable, testable, and scalable Go applications by strictly adhering to Clean Architecture principles.

## 🏗️ Clean Architecture Overview

This project implements Clean Architecture to separate concerns and create a system that is:

1. **Independent of Frameworks:** The core business logic doesn't depend on external frameworks or libraries.
2. **Testable:** Business rules can be tested without UI, database, web server, or any external elements.
3. **Independent of UI:** The UI can change without changing the rest of the system.
4. **Independent of Database:** Your business rules aren't bound to a specific database implementation.
5. **Independent of External Services:** Business rules don't know anything about interfaces to the outside world.

### Layers in this Project

The application is organized into distinct layers with clear responsibilities:

* **Domain Layer (`domain/`):** The enterprise business rules
  * `entity/`: Core business objects with no dependencies
  * `repository/`: Interfaces defining data access contracts

* **Use Case Layer (`usecase/`):** Application-specific business rules
  * Orchestrates data flow between domain entities and external layers
  * Depends on domain layer, but not on outer layers

* **Infrastructure Layer (`infrastructure/`):** Implementation details
  * `repository/`: Concrete implementations of repository interfaces
    * `memory/`: In-memory data store implementation

* **Delivery Layer (`delivery/`):** How the outside world interacts with the application
  * `http/`: HTTP-specific delivery mechanisms
    * `middleware/`: HTTP middleware components

* **Configuration (`config/`):** Application configuration management

## 📁 Project Structure

```
.
├── config/             # Application configuration
├── domain/             # Enterprise business rules
│   ├── entity/         # Business objects
│   └── repository/     # Repository interfaces
├── infrastructure/     # Implementation details
│   └── repository/     # Repository implementations
│       └── memory/     # In-memory data store
├── usecase/            # Application business rules
├── delivery/           # External interfaces
│   └── http/           # HTTP delivery
│       └── middleware/ # HTTP middleware
├── k8s/                # Kubernetes manifests
│   ├── configmap.yml   # ConfigMap and Secret
│   ├── deployment.yml  # Deployment configuration
│   └── service.yml     # Service and Ingress
├── .env                # Local environment variables
├── .env.example        # Example environment variables
├── .gitignore          # Git ignore patterns
├── docker-compose.yml  # Docker Compose configuration
├── Dockerfile          # Docker image definition
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── main.go             # Application entry point
```

## 🚀 Getting Started

### Prerequisites

* Go 1.24 or later ([Download Go](https://golang.org/dl/))
* Git

### Installation & Running

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/go-clean-boilerplate.git
   cd go-clean-boilerplate
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

   The API server will start on port 8080 (configurable via environment variables).

4. **Build for production:**
   ```bash
   go build -o ./bin/server
   ```

## 🐳 Docker Deployment

### Prerequisites

* Docker ([Get Docker](https://docs.docker.com/get-docker/))
* Docker Compose ([Get Docker Compose](https://docs.docker.com/compose/install/))

### Running with Docker Compose

1. **Build and run:**
   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the Go application
   - Start a PostgreSQL database
   - Connect the application to the database
   - Expose the API on port 8080

2. **Run in detached mode:**
   ```bash
   docker-compose up -d
   ```

3. **Stop the containers:**
   ```bash
   docker-compose down
   ```

### Building Docker Image Manually

```bash
docker build -t go-clean-boilerplate .
```

## ☸️ Kubernetes Deployment

For production environments, this project includes Kubernetes manifests in the `k8s/` directory.

### Prerequisites

* Kubernetes cluster (local like Minikube or remote)
* kubectl configured to communicate with your cluster

### Deployment Steps

1. **Apply all resources:**
   ```bash
   kubectl apply -f k8s/
   ```

   Or apply individual components:

   ```bash
   kubectl apply -f k8s/configmap.yml
   kubectl apply -f k8s/deployment.yml
   kubectl apply -f k8s/service.yml
   ```

### Customizing Kubernetes Deployment

* Update the image name in `deployment.yml` to point to your Docker registry
* Modify the host in `service.yml` ingress rules to match your domain
* Adjust resource limits in `deployment.yml` based on your application needs

## ⚙️ Configuration

The application can be configured using environment variables:

| Variable               | Description                       | Default Value    |
|:-----------------------|:----------------------------------|:-----------------|
| `PORT`                 | HTTP server port                  | `8080`           |
| `SERVER_READ_TIMEOUT`  | Request read timeout              | `10` (seconds)   |
| `SERVER_WRITE_TIMEOUT` | Response write timeout            | `10` (seconds)   |
| `SERVER_IDLE_TIMEOUT`  | Idle connection timeout           | `120` (seconds)  |
| `DB_DRIVER`            | Database driver                   | `memory`         |
| `DB_HOST`              | Database host                     | `localhost`      |
| `DB_PORT`              | Database port                     | `5432`           |
| `DB_USER`              | Database username                 | `postgres`       |
| `DB_PASSWORD`          | Database password                 | `postgres`       |
| `DB_NAME`              | Database name                     | `go_clean_boilerplate`   |
| `DB_SSL_MODE`          | Database SSL mode                 | `disable`        |
| `LOG_LEVEL`            | Logging level                     | `info`           |

**Note:** To use PostgreSQL instead of the default in-memory database:
1. Implement the repository interfaces for PostgreSQL
2. Set `DB_DRIVER=postgres` and configure the other database variables

## 🔌 API Endpoints

| Method   | Path           | Description                      |
|:---------|:---------------|:---------------------------------|
| `GET`    | `/health`      | Health check                     |
| `GET`    | `/users`       | List all users                   |
| `POST`   | `/users`       | Create a new user                |
| `GET`    | `/users/{id}`  | Get user by ID                   |
| `PUT`    | `/users/{id}`  | Update user by ID                |
| `DELETE` | `/users/{id}`  | Delete user by ID                |

**Example Request Body for POST /users:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

## 🧪 Testing

Run tests using the standard Go tool:

```bash
go test ./...
```

Clean Architecture facilitates unit testing of domain logic and use cases without needing a running server or database.

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
#   g o - c l e a n - b o i l e r p l a t e  
 