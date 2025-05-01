# Simple Wallet Management REST API

## Objective
The objective of this project is to build a simple RESTful backend service in GoLang to manage user wallets and handle basic transactions between them. This service will allow users to perform wallet management operations such as creating users, checking balances, transferring funds, and viewing transaction history.


## Tech Stack
- **Language**: GoLang
- **Web Framework**: Gin
- **Database**: PostgreSQL (using GORM ORM)
- **Docker**: Dockerized for easy setup and deployment
- **Environment Variables**: Used for database configuration

## Setup Instructions

### Prerequisites
Before you begin, make sure you have the following installed:
- GoLang (v1.16 or above)
- Docker
- PostgreSQL (if not using Dockerized DB)
- Git

### 1. Clone the repository
Clone this repository to your local machine:
```bash
git clone https://github.com/tegveer-singh123/wallet-api.git
cd wallet-api

### 2. Set up the  env variables 

APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=your-db-username
DB_PASSWORD=your-db-password
DB_NAME=wallet_api


### 3. Install Dependencies 

go mod tidy 

### 4. Run the Project locally

go run cmd/main.go

### 5. Run using Docker

docker-compose up --build

API Docs Link - https://documenter.getpostman.com/view/39144503/2sB2j4eWRy

