# Gin Market

Welcome to Gin Market! 
This is a lightweight e-commerce API built using the Gin web framework for Go. 
With Gin Market, you can easily manage your inventory and handle user authentication.

## Features

- **Item Management:** Perform CRUD operations for managing items, including details like name, price, and description.
- **User Authentication:** Enable users to sign up and log in securely.

## Installation

### Prerequisites

- Go 1.15 or later installed on your system.
- A compatible database system (such as SQLite) set up and ready to go.

### Getting Started

1. Clone this repository to your local machine:

    ```bash
    git clone https://github.com/namashin/gin-market.git
    ```

2. Navigate into the project directory and install dependencies:

    ```bash
    cd gin-market
    go mod tidy
    ```

3. Start the server:

    ```bash
    go run main.go
    ```

    Your server should now be up and running at [http://localhost:8080](http://localhost:8080).

## Usage

### API Endpoints

**Items:**

- `GET /items`: Retrieve all items.
- `GET /items/mine`: Retrieve items belonging to the authenticated user.
- `GET /items/:id`: Retrieve a specific item by its ID.
- `POST /items`: Create a new item.
- `PUT /items/:id`: Update an existing item.
- `DELETE /items/:id`: Delete an item.

**Authentication:**

- `POST /auth/signup`: Sign up a new user.
- `POST /auth/login`: Log in an existing user.

### Configuration

You can configure the application's database settings in the `config.yaml` file.

## Contributing

We welcome contributions! If you have any suggestions, ideas, or bug fixes, please feel free to open a pull request or issue on GitHub.
