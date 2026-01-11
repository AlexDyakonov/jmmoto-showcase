# Telegram Mini App Template

A full-stack template for building Telegram Mini Apps with Go backend and React frontend.

[Читать на русском](README_RU.md)

## 🚀 Features

- **Backend**: Go with Huma REST API framework
- **Frontend**: React 19 + Vite + TypeScript + Tailwind CSS
- **Database**: PostgreSQL with migrations
- **Telegram Integration**: Bot SDK and Mini App support
- **API Documentation**: Auto-generated OpenAPI/Swagger docs
- **Docker**: Docker Compose for easy development setup

## 📋 Prerequisites

- Go 1.24+ 
- Node.js 18+ and npm
- Docker and Docker Compose
- PostgreSQL 16+ (or use Docker Compose)

## 🛠️ Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd tma-template
```

### 2. Backend Setup

Navigate to the server directory:

```bash
cd server
```

Create a `.env` file in the `server` directory:

```env
# Server Configuration
HTTP_PORT=8000
HTTP_HOST=0.0.0.0
DEBUG=true

# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=cats

# Telegram Configuration
TG_BOT_TOKEN=your_bot_token_here
WEBAPP_NAME=your_webapp_name

# Storage Configuration (optional)
STORAGE_IMAGES_PATH=images

# S3 Configuration (optional)
S3_ACCESS_KEY_ID=
S3_SECRET_KEY=
S3_REGION=
S3_BUCKET=
S3_ENDPOINT_URL=
S3_ROOT_DIRECTORY=

# Logging
LOG_HANDLER=tint
```

Install Go dependencies:

```bash
go mod download
```

### 3. Frontend Setup

Navigate to the client directory:

```bash
cd ../client
```

Install dependencies:

```bash
npm install
```

Update `public/config.js` with your configuration:

```javascript
window.api = {
  BOT_USERNAME: "your_bot_username",
  API_URL: "http://localhost:8000",
};
```

## 🏃 Running the Application

### Option 1: Using Docker Compose (Recommended for Database)

Start the database and Adminer:

```bash
cd server
make db-compose-up
```

Run database migrations:

```bash
make db-migrate-up
```

Start the backend server:

```bash
make server-run
```

In another terminal, start the frontend:

```bash
cd client
npm run dev
```

### Option 2: Full Docker Compose

Start all services:

```bash
cd server
make compose-up
```

### Option 3: Manual Setup

1. Start PostgreSQL database
2. Run migrations: `make db-migrate-up`
3. Start backend: `make server-run`
4. Start frontend: `cd client && npm run dev`

## 📚 API Documentation

Once the server is running, access the API documentation:

- **OpenAPI YAML**: http://localhost:8000/api/v1/openapi.yaml
- **OpenAPI JSON**: http://localhost:8000/api/v1/openapi.json
- **Interactive Docs**: http://localhost:8000/api/v1/docs

## 🗄️ Database Migrations

### Run migrations up

```bash
cd server
make db-migrate-up
```

### Rollback migrations

```bash
make db-migrate-down
```

## 🧪 Development

### Backend Commands

```bash
cd server

# Run server
make server-run

# Run Telegram bot
make bot-run

# Lint code
make lint

# Fix linting issues
make lint-fix

# Generate Swagger documentation
make swag
```

### Frontend Commands

```bash
cd client

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint

# Fix linting issues
npm run lint-fix
```

## 🔧 Configuration

### Environment Variables

See the `.env` example above for all available configuration options.

### Frontend Configuration

Edit `client/public/config.js` to configure:
- Bot username
- API URL

## 🐳 Docker

### Build and run with Docker Compose

```bash
cd server
docker compose up --build
```

### Stop services

```bash
docker compose down
```