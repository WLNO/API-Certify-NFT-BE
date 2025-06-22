# Certify NFT - Backend API

This is the backend API for the Certify NFT platform. It handles event management, user/vendor registration, authentication, and certificate issuance logic. The API is built with Go and containerized with Docker.

## Table of Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Running Locally](#installation--running-locally)
- [API Documentation](#api-documentation)
- [Stack](#stack)
- [Project Structure](#project-structure)

## Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation & Running Locally

1.  **Clone the repository:**
    ```sh
    git clone <your-repository-url>
    cd API-Certify-NFT-BE
    ```

2.  **Create your environment file:**
    Copy the example environment file `.env.example` to a new file named `.env`.
    ```sh
    cp .env.example .env
    ```
    Then, open `.env` and fill in the required values, especially your database password.

3.  **Run with Docker Compose:**
    This command will build the Go application image and start the service.
    ```sh
    docker-compose up --build -d
    ```

4.  **The API is now running!**
    You can access it at `http://localhost:4002`.

## API Documentation
For detailed information about all available endpoints, request/response formats, and examples, please see the [API Documentation](./apiDocs.md).

## Stack
Here are the technologies, tools, and frameworks used in this project:

-   **Language:** **Go (Golang)**
-   **Web Framework:** **Echo (v4)** - A high-performance, minimalist Go web framework.
-   **Database:** **PostgreSQL** - A powerful, open-source object-relational database system.
-   **Containerization:**
    -   **Docker** - For creating a consistent and isolated environment for the application.
    -   **Docker Compose** - For defining and running the multi-container Docker application.
-   **Go Packages:**
    -   `github.com/lib/pq`: A pure Go Postgres driver for the `database/sql` package.
    -   `github.com/joho/godotenv`: A library to load environment variables from a `.env` file.

## Project Structure

```
.
├── .git/
├── .github/          # GitHub Actions workflows
├── .gitignore
├── apiDocs.md        # Detailed API documentation
├── database_structure.md
├── docker-compose.yml # Defines the services, networks, and volumes
├── Dockerfile        # Defines the steps to build the Go application image
├── go.mod            # Manages dependencies for the Go project
├── go.sum
├── main.go           # Main application entry point, routing, and handlers
├── README.md         # This file
└── schema.sql        # Database schema definitions
```
