# Development Setup Guide

## First-Time Setup

### 1. Clone Repository

```bash
git clone git@github.com:sfumato00/content-analyzer.git
cd content-analyzer
```

### 2. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` and update:

- `GEMINI_API_KEY`: Get your free API key from https://makersuite.google.com/app/apikey
- `JWT_SECRET`: Generate a random 32+ character string (e.g., `openssl rand -base64 32`)

### 3. Start Infrastructure

```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis

# Wait for services to be healthy (check status)
docker-compose ps
```

### 4. Run Migrations

The application will automatically run migrations in development mode when you start it.

Alternatively, you can run migrations manually:

```bash
cd backend
./scripts/migrate.sh up
```

### 5. Start Application

**Option A: Local Development (Go)**

```bash
cd backend
go run cmd/api/main.go
```

**Option B: Docker**

```bash
docker-compose up -d
```

### 6. Verify Setup

```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "healthy",
#   "uptime": "...",
#   "version": "1.0.0",
#   "components": {
#     "database": "connected",
#     "redis": "connected"
#   }
# }

# Test database connection
docker-compose exec postgres psql -U postgres -d content_analyzer -c "SELECT COUNT(*) FROM users;"
```

## Troubleshooting

### Database Connection Failed

- Ensure PostgreSQL is running: `docker-compose ps postgres`
- Check logs: `docker-compose logs postgres`
- Verify DATABASE_URL in .env matches: `postgresql://postgres:dev@localhost:5432/content_analyzer`

### Redis Connection Failed

- Ensure Redis is running: `docker-compose ps redis`
- Check logs: `docker-compose logs redis`
- Verify REDIS_URL in .env: `redis://localhost:6379`

### Migrations Failed

- Check migration files in `backend/migrations/`
- Verify database exists: `content_analyzer`
- Force to specific version if needed: `cd backend && ./scripts/migrate.sh force 1`
- Check migration status: `cd backend && ./scripts/migrate.sh version`

### Port Already in Use

If port 8080 is already in use:

```bash
# Find what's using the port
lsof -i :8080

# Kill the process or change PORT in .env
PORT=8081
```

## Development Workflow

### Making Code Changes

1. Edit code in `backend/` directory
2. Restart the application:
   - Local: Stop (Ctrl+C) and run `go run cmd/api/main.go` again
   - Docker: `docker-compose restart api`
3. Test changes at http://localhost:8080

### Creating Database Changes

1. Create migration:
   ```bash
   cd backend
   ./scripts/migrate.sh create add_new_column
   ```

2. Edit the generated `.up.sql` and `.down.sql` files in `migrations/`

3. Run migration:
   ```bash
   ./scripts/migrate.sh up
   ```

4. Test rollback:
   ```bash
   ./scripts/migrate.sh down
   ```

### Running Tests

```bash
cd backend
go test ./...
```

### Viewing Logs

```bash
# API logs
docker-compose logs -f api

# PostgreSQL logs
docker-compose logs -f postgres

# Redis logs
docker-compose logs -f redis

# All services
docker-compose logs -f
```

## Database Management

### Connecting to PostgreSQL

```bash
# Via Docker
docker-compose exec postgres psql -U postgres -d content_analyzer

# Via local psql client
psql postgresql://postgres:dev@localhost:5432/content_analyzer
```

### Common SQL Commands

```sql
-- List all tables
\dt

-- Describe table structure
\d users

-- View table data
SELECT * FROM users LIMIT 10;

-- Check migration status
SELECT * FROM schema_migrations;
```

### Connecting to Redis

```bash
# Via Docker
docker-compose exec redis redis-cli

# Common Redis commands
PING           # Test connection
KEYS *         # List all keys
GET key_name   # Get value
DEL key_name   # Delete key
FLUSHALL       # Clear all data (use with caution!)
```

## Useful Commands

### Docker

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# View service status
docker-compose ps

# Remove volumes (clears database data)
docker-compose down -v
```

### Go

```bash
# Install dependencies
go mod download

# Update dependencies
go get -u ./...

# Tidy dependencies
go mod tidy

# Build binary
go build -o bin/api cmd/api/main.go

# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Environment Variables

### Required Variables

- `GEMINI_API_KEY`: Google Gemini API key (get from https://makersuite.google.com/app/apikey)
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: Secret for JWT token signing (32+ characters)

### Optional Variables

- `PORT`: API server port (default: 8080)
- `ENV`: Environment mode - `development` or `production` (default: development)
- `ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins

## Next Steps

Once your setup is complete, you're ready to start development:

1. **Week 1-2**: Backend foundation is complete
2. **Week 3**: Implement authentication (register, login, JWT)
3. **Week 4**: Integrate Gemini AI for content analysis
4. **Week 5**: Build frontend with React and TypeScript

## Getting Help

- Check the main [README.md](README.md) for project overview
- Review [PRODUCT_PLAN.md](PRODUCT_PLAN.md) for detailed architecture
- Run verification script: `./scripts/verify-setup.sh`

## Clean Slate

If you want to start fresh:

```bash
# Stop services and remove volumes
docker-compose down -v

# Remove built binaries
rm -rf backend/bin

# Restart from step 3
docker-compose up -d postgres redis
```
