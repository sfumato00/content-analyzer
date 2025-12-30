# AI-Powered Content Analyzer

A personal portfolio project demonstrating modern backend engineering skills through an AI-powered content analysis tool.

## Tech Stack

- **Backend**: Go 1.21+
- **Frontend**: TypeScript, React, Vite
- **Database**: PostgreSQL
- **Cache**: Redis
- **AI**: Google Gemini API (free tier)
- **Infrastructure**: Docker, Docker Compose

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Google Gemini API key (free from https://makersuite.google.com/app/apikey)

### Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd content-analyzer
   ```

2. **Set up environment variables**
   ```bash
   # Copy the example file
   cp .env.example .env

   # Edit .env and add your actual API keys
   # IMPORTANT: Add your real Gemini API key!
   nano .env
   ```

3. **Start the database and Redis**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Install Go dependencies**
   ```bash
   cd backend
   go mod download
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

   The application will automatically run database migrations in development mode.

   You should see:
   ```
   🚀 Content Analyzer API
   ========================
   Environment: development
   Port: 8080
   ...
   ✅ Server ready on port 8080
   ```

### Database Setup

The application automatically runs migrations on startup in development mode. You can also manage migrations manually:

```bash
cd backend

# Run all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Check current migration version
./scripts/migrate.sh version

# Create new migration
./scripts/migrate.sh create add_new_table
```

**Database Schema**: The initial migration creates three tables:
- `users` - User accounts with email and password
- `submissions` - User-submitted content for analysis
- `analyses` - AI analysis results from Gemini

### Running with Docker

You can run the entire stack in Docker:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up -d --build
```

### Verification

After starting the services, verify everything is working:

```bash
# Check API health
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "status": "healthy",
  "uptime": "1m30s",
  "version": "1.0.0",
  "components": {
    "database": "connected",
    "redis": "connected"
  }
}
```

**Verify database tables:**
```bash
docker-compose exec postgres psql -U postgres -d content_analyzer -c "\dt"
```

You should see: `users`, `submissions`, `analyses`

**Verify Redis:**
```bash
docker-compose exec redis redis-cli ping
# Should return: PONG
```

**Run automated verification:**
```bash
./scripts/verify-setup.sh
```

### Troubleshooting

**Database connection failed:**
- Ensure PostgreSQL is running: `docker-compose ps postgres`
- Check DATABASE_URL in .env
- View logs: `docker-compose logs postgres`

**Redis connection failed:**
- Ensure Redis is running: `docker-compose ps redis`
- Check REDIS_URL in .env
- View logs: `docker-compose logs redis`

**Port already in use:**
```bash
# Find what's using port 8080
lsof -i :8080

# Or change PORT in .env
PORT=8081
```

**For detailed setup instructions**, see [SETUP.md](./SETUP.md)

## API Endpoints

### Authentication (Public)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `POST /api/v1/auth/logout` - Logout (client-side token removal)

### User (Protected - Requires JWT)
- `GET /api/v1/me` - Get current user info
- `GET /api/v1/me/stats` - Get user statistics (coming soon)

### Submissions (Protected - Requires JWT)
- `POST /api/v1/submissions` - Submit content for analysis (coming soon)
- `GET /api/v1/submissions` - List user's submissions (coming soon)
- `GET /api/v1/submissions/:id` - Get submission details (coming soon)
- `GET /api/v1/submissions/:id/analysis` - Get AI analysis (coming soon)

### Health
- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check
- `GET /live` - Liveness check

**Example usage:**
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get current user (requires token from login)
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Project Structure

```
content-analyzer/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go           # Application entry point
│   ├── internal/
│   │   ├── config/               # Configuration management
│   │   ├── auth/                 # Authentication (JWT, middleware) ✅
│   │   ├── database/             # PostgreSQL setup ✅
│   │   ├── handlers/             # HTTP handlers ✅
│   │   ├── models/               # Data models ✅
│   │   ├── middleware/           # Security middleware ✅
│   │   ├── response/             # Response helpers ✅
│   │   ├── cache/                # Redis client ✅
│   │   └── services/             # Business logic (coming soon)
│   │       ├── ai/               # Gemini integration
│   │       └── queue/            # Background jobs
│   ├── migrations/               # SQL migrations ✅
│   ├── Dockerfile                # ✅
│   └── go.mod                    # ✅
├── frontend/                     # (planned for Week 4)
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── api/                  # API client
│   │   └── App.tsx
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml            # ✅
├── Makefile                      # ✅
├── .env.example                  # ✅
├── .gitignore                    # ✅
└── README.md                     # ✅
```

## Current Status

**Week 1-2: Backend Foundation** ✅ COMPLETE

- ✅ PostgreSQL schema and migrations
- ✅ User model with bcrypt password hashing
- ✅ JWT authentication (register, login, logout, /me)
- ✅ Protected routes with JWT middleware
- ✅ Comprehensive unit tests
- ✅ Makefile for streamlined development

**Next: Week 3 - AI Integration**

See [PRODUCT_PLAN.md](./PRODUCT_PLAN.md) for the complete development roadmap.

### Week 3: AI Integration
- [ ] Gemini API integration
- [ ] Background job queue
- [ ] Redis caching layer
- [ ] Rate limiting

### Week 4: Frontend & Deployment
- [ ] React app with TypeScript
- [ ] Authentication UI
- [ ] Content submission form
- [ ] Deploy to Fly.io or Railway

## Development Commands

The project includes a comprehensive Makefile for streamlined development:

```bash
# Show all available commands
make help

# Development setup
make dev-setup              # Install deps + start Docker services
make run                    # Run the backend server
make test                   # Run all tests
make test-coverage          # Run tests with coverage report

# Docker management
make docker-up              # Start all services (postgres, redis, api)
make docker-up-db           # Start only database services
make docker-down            # Stop all services
make docker-logs            # View logs from all services
make docker-rebuild         # Rebuild and restart

# Database
make db-shell               # Open PostgreSQL shell
make migrate-up             # Run pending migrations
make migrate-down           # Rollback last migration

# Code quality
make fmt                    # Format Go code
make lint                   # Run linter
make ci                     # Run CI checks (lint + test)

# Testing API endpoints
make test-register          # Test registration endpoint
make test-login             # Test login endpoint
make health                 # Check API health

# Utilities
make verify                 # Verify entire setup is working
make clean                  # Clean build artifacts
```

**Or use raw commands:**

```bash
# Start all services
docker-compose up -d

# Run Go application
cd backend && go run cmd/api/main.go

# Run tests
cd backend && go test ./...
```

## Environment Variables

See `.env.example` for all available configuration options.

**Required**:
- `GEMINI_API_KEY` - Get from https://makersuite.google.com/app/apikey
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - Random secret string (min 32 characters)

**Optional**:
- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development/production)
- `ALLOWED_ORIGINS` - CORS allowed origins

## Security Notes

- Never commit `.env` file to git
- Use strong JWT secrets (min 32 characters)
- In production, use platform secrets (Fly.io secrets, Railway env vars)
- API keys are masked in logs automatically

## Cost Estimate

**Development**: $0 (everything runs locally)

**Production**:
- Fly.io free tier: $0/month
- Gemini API free tier: $0/month (1500 requests/day)
- **Total: $0-5/month**

## License

MIT

## Author

Your Name - [GitHub](https://github.com/yourusername)

## Acknowledgments

Built as a portfolio project to demonstrate modern backend engineering skills for 2026 job market.
