# AI-Powered Content Analyzer
## Personal Side Project Plan

**Last Updated**: 2025-12-28
**Status**: Development Phase - Week 1-2 Infrastructure Complete
**Developer**: Solo side project
**Target Stack**: Go 1.24, TypeScript, Docker, PostgreSQL, Redis, Gemini AI
**Estimated Cost**: $0-15/month

---

## Executive Summary

A personal portfolio project demonstrating modern backend engineering skills through an AI-powered content analysis tool. Users can submit text content (articles, social media posts, customer reviews) and receive AI-driven insights including sentiment analysis, key topics, and automated summaries.

**Why This Project**: Showcases the most in-demand 2026 backend skills (Go, TypeScript, AI integration, modern infrastructure) while being realistic for a solo developer to build and maintain as a side project.

---

## Product Vision

**What**: A web application where users can paste or upload text content and get AI-powered analysis.

**Core Features**:
- Text submission via web UI or API
- AI sentiment analysis using Gemini (free tier)
- Keyword/topic extraction
- Content summarization
- Historical dashboard showing past analyses
- Simple user authentication

**Value for Portfolio**:
- Demonstrates Go backend development
- Shows AI/LLM integration experience
- Modern TypeScript frontend skills
- Infrastructure as code
- Production deployment experience

---

## Simplified Architecture

### Architecture Pattern
**Modular Monolith** - Single Go application with clear internal boundaries, can be split into microservices later if needed.

### System Components

#### 1. **Backend API** (Go)
**Technology**: Go 1.21+, Chi router (lightweight), structured logging

**Responsibilities**:
- RESTful API endpoints
- Request validation
- Authentication (JWT)
- Database operations
- Background job processing (AI requests)
- Rate limiting

**Why Go**:
- Performance and concurrency for job processing
- Most in-demand backend language for 2026
- Easy deployment (single binary)

**Key Packages**:
- `chi` - HTTP router
- `pgx` - PostgreSQL driver
- `golang-jwt` - JWT auth
- `slog` - Structured logging
- `testify` - Testing

#### 2. **Frontend** (TypeScript + React)
**Technology**: TypeScript, React, Vite, TailwindCSS

**Responsibilities**:
- Content submission form
- Display AI analysis results
- User authentication UI
- Dashboard for historical analyses
- Responsive design

**Why TypeScript**:
- In-demand skill for 2026
- Type safety reduces bugs
- Better developer experience

**Note**: Keep frontend simple - focus is on backend skills.

#### 3. **AI Service Layer** (Go → Gemini)
**Technology**: Google Gemini API (free tier), HTTP client

**Responsibilities**:
- Queue AI analysis requests
- Call Gemini API for:
  - Sentiment analysis
  - Topic extraction
  - Text summarization
- Handle rate limiting (60 requests/min free tier)
- Retry logic for failed requests
- Cache results to avoid duplicate API calls

**Why Gemini**:
- FREE tier (15 requests/min, 1500 requests/day)
- Good quality results
- Simple REST API

---

## Data Architecture

### Primary Data Store: PostgreSQL

**Why PostgreSQL**:
- Free (self-hosted or Supabase free tier)
- ACID compliance
- JSON support for flexible AI response storage
- Full-text search capabilities

**Schema Design**:

```sql
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Content submissions
CREATE TABLE submissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  content TEXT NOT NULL,
  status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
  created_at TIMESTAMP DEFAULT NOW()
);

-- AI analysis results
CREATE TABLE analyses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  submission_id UUID REFERENCES submissions(id),
  sentiment VARCHAR(50), -- positive, neutral, negative
  sentiment_score FLOAT,
  topics JSONB, -- Array of extracted topics
  summary TEXT,
  raw_response JSONB, -- Full Gemini response
  processing_time_ms INT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_analyses_submission_id ON analyses(submission_id);
```

**Optimizations**:
- Indexes on foreign keys and status columns
- JSONB for flexible AI response storage
- Timestamps for analytics

### Caching Layer: Redis

**Why Redis**:
- Free (self-hosted or Upstash free tier - 10K requests/day)
- Fast response caching
- Session storage
- Rate limiting counters

**Use Cases**:
```
- Cache AI results: cache:analysis:{content_hash} → analysis result (TTL: 7 days)
- Rate limiting: ratelimit:{user_id}:{endpoint} → counter (TTL: 1 minute)
- Session storage: session:{token} → user data (TTL: 24 hours)
```

---

## Infrastructure (Minimal Cost)

### Local Development
**Docker Compose** - Run entire stack locally

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: content_analyzer
      POSTGRES_PASSWORD: dev
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  api:
    build: ./backend
    environment:
      DATABASE_URL: postgresql://postgres:dev@postgres/content_analyzer
      REDIS_URL: redis://redis:6379
      GEMINI_API_KEY: ${GEMINI_API_KEY}
    ports:
      - "8080:8080"

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
```

### Production Deployment (Free/Low Cost Options)

**Option 1: Fly.io (FREE tier)**
- 3 shared-cpu VMs (256MB RAM each) - FREE
- Deploy backend + PostgreSQL + Redis
- Estimated cost: **$0/month** (within free tier)

**Option 2: Railway.dev (FREE tier)**
- $5 free credits/month
- PostgreSQL, Redis, and app hosting
- Estimated cost: **$0-5/month**

**Option 3: Self-hosted VPS**
- Hetzner Cloud (smallest VPS): €4.5/month (~$5)
- Oracle Cloud Free Tier: 2 VMs FREE forever
- Estimated cost: **$0-5/month**

**Frontend Hosting**:
- **Vercel** or **Netlify** - FREE tier (perfect for React apps)
- Cloudflare Pages - FREE tier

**Total Infrastructure Cost: $0-10/month**

---

## AI Integration (FREE Tier)

### Google Gemini API

**Free Tier Limits** (as of 2025):
- 15 requests per minute
- 1,500 requests per day
- 1 million tokens per month

**For Personal Project**: More than sufficient!

**API Integration**:

```go
// Example Go code for Gemini API call
type GeminiRequest struct {
    Contents []Content `json:"contents"`
}

type Content struct {
    Parts []Part `json:"parts"`
}

type Part struct {
    Text string `json:"text"`
}

func analyzeContent(text string) (*Analysis, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("analysis:%s", hashContent(text))
    if cached, err := redis.Get(ctx, cacheKey).Result(); err == nil {
        return parseAnalysis(cached), nil
    }

    // Build prompt for Gemini
    prompt := fmt.Sprintf(`Analyze the following text and provide:
1. Sentiment (positive/neutral/negative) with confidence score
2. Top 5 key topics or themes
3. A 2-3 sentence summary

Text: %s

Respond in JSON format.`, text)

    // Call Gemini API
    resp, err := http.Post(
        "https://generativelanguage.googleapis.com/v1/models/gemini-pro:generateContent?key="+apiKey,
        "application/json",
        buildGeminiRequest(prompt),
    )

    // Parse and cache result
    analysis := parseGeminiResponse(resp)
    redis.Set(ctx, cacheKey, analysis, 7*24*time.Hour)

    return analysis, nil
}
```

**Rate Limiting Strategy**:
- Queue requests if hitting rate limit
- Process queue with worker (15 requests/min max)
- Show "processing" status to users
- Email when analysis is complete (for async processing)

**Cost**: **$0/month** (free tier)

---

## API Design

### RESTful API (Simple & Clean)

```
Authentication:
POST   /api/v1/auth/register       - Create account
POST   /api/v1/auth/login          - Login (returns JWT)
POST   /api/v1/auth/logout         - Logout

Content Analysis:
POST   /api/v1/submissions         - Submit text for analysis
GET    /api/v1/submissions         - List user's submissions
GET    /api/v1/submissions/:id     - Get submission details
GET    /api/v1/submissions/:id/analysis - Get AI analysis result

User:
GET    /api/v1/me                  - Get current user info
GET    /api/v1/me/stats            - Get usage stats
```

**Authentication**: JWT tokens in Authorization header
```
Authorization: Bearer <jwt_token>
```

**Rate Limiting**: 100 requests/hour per user (free tier)

---

## Development Workflow

### Tech Stack Summary

**Backend**:
- Go 1.21+
- Chi router
- PostgreSQL (pgx driver)
- Redis client

**Frontend**:
- TypeScript
- React 18
- Vite
- TailwindCSS
- React Query (data fetching)

**Infrastructure**:
- Docker & Docker Compose
- GitHub Actions (CI/CD)
- Fly.io or Railway (hosting)

### Project Structure

```
content-analyzer/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── auth/          # JWT, middleware
│   │   ├── database/      # PostgreSQL setup
│   │   ├── handlers/      # HTTP handlers
│   │   ├── models/        # Data models
│   │   ├── services/      # Business logic
│   │   │   ├── ai/        # Gemini integration
│   │   │   └── queue/     # Background jobs
│   │   └── cache/         # Redis client
│   ├── migrations/        # SQL migrations
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── api/           # API client
│   │   └── App.tsx
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── .github/
│   └── workflows/
│       └── ci.yml
└── README.md
```

### CI/CD Pipeline (GitHub Actions - FREE)

```yaml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: |
          cd backend
          go test -v ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v3

  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Fly.io
        uses: superfly/flyctl-actions@v1
        with:
          args: "deploy"
```

**Cost**: **$0/month** (GitHub Actions free for public repos)

---

## Monitoring & Logging (FREE Tools)

### Application Monitoring

**Option 1: Built-in Logging**
- Structured logging with Go's `slog`
- Log to stdout → Fly.io dashboard (free)
- Store logs in PostgreSQL for simple analytics

**Option 2: Free Monitoring Tools**
- **Better Stack** (formerly Logtail): Free tier - 1GB logs/month
- **Grafana Cloud**: Free tier - 10K metrics, 50GB logs
- **Sentry** (errors): Free tier - 5K errors/month

**Metrics to Track**:
```
- Total submissions count
- AI analysis success rate
- Average processing time
- API response times (p50, p95, p99)
- Cache hit rate
- Active users count
```

### Simple Health Check Endpoint

```go
GET /health
Response:
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "gemini_api": "available"
}
```

**Cost**: **$0/month** (using free tiers)

---

## Security (Free Best Practices)

### Authentication & Authorization
- **JWT** for stateless auth
- **bcrypt** for password hashing (cost factor: 12)
- HTTP-only cookies for web clients
- Rate limiting per user

### Data Security
- **Environment variables** for secrets (never commit)
- **HTTPS only** in production (Fly.io provides free SSL)
- **Input validation** and sanitization
- **Prepared statements** (prevent SQL injection)

### Security Headers
```go
// Middleware for security headers
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=63072000")
        next.ServeHTTP(w, r)
    })
}
```

**No Paid Tools Required**: Follow OWASP best practices

---

## Development Roadmap (Realistic for Solo Developer)

### Phase 1: MVP (Weeks 1-4) - Part-time work

**Week 1-2: Backend Foundation**
- [x] Set up Go project structure
- [x] PostgreSQL schema and migrations
- [x] Docker Compose setup (with Dockerfile, health checks)
- [x] Database connection layer (pgx with connection pooling)
- [x] Redis cache layer (go-redis client)
- [x] Database migrations tooling (golang-migrate)
- [x] User model with password hashing (bcrypt)
- [x] JWT authentication (token generation, validation, middleware)
- [x] Auth endpoints (register, login, logout, /me)
- [x] Protected routes with JWT middleware
- [x] Unit tests for auth and models
- [x] Makefile for development workflow
- [ ] Basic CRUD API endpoints (routes stubbed, implementation pending)

**Week 3: AI Integration**
- [ ] Gemini API integration
- [ ] Background job queue for AI requests
- [x] Redis caching layer (infrastructure ready)
- [ ] Rate limiting implementation
- [ ] Error handling and retries

**Week 4: Frontend & Deployment**
- [ ] React app with TypeScript
- [ ] Authentication UI (login/register)
- [ ] Content submission form
- [ ] Results display page
- [ ] Deploy to Fly.io or Railway
- [ ] Set up CI/CD pipeline

**Success Criteria**:
- Can submit text and get AI analysis
- User authentication works
- Deployed and accessible online
- Costs $0-5/month

### Phase 2: Polish & Portfolio Prep (Weeks 5-6) - Optional

**Enhancements**:
- [ ] Dashboard with analytics charts
- [ ] Export results as PDF/JSON
- [ ] Improve UI/UX with better design
- [ ] Add API documentation (Swagger/OpenAPI)
- [ ] Write comprehensive README
- [ ] Add integration tests
- [ ] Performance optimization

**Portfolio Materials**:
- [ ] Architecture diagram
- [ ] Demo video or screenshots
- [ ] Blog post about technical decisions
- [ ] Link to live demo

### Phase 3: Advanced Features (Future) - If time permits

**Ideas for Expansion**:
- [ ] Batch processing (upload CSV with multiple texts)
- [ ] Comparison view (compare two texts)
- [ ] Public API for other developers
- [ ] Webhook integrations
- [ ] Real-time updates with WebSockets
- [ ] Switch to microservices (split AI service)

**Note**: Only if you want to continue developing after getting a job!

---

## Cost Breakdown (Monthly)

### Infrastructure

**Option A: Fly.io (Recommended)**
```
PostgreSQL:        $0 (256MB free tier)
Redis:             $0 (Upstash free tier - 10K requests/day)
App hosting:       $0 (3 shared VMs free)
Domain (optional): $12/year (~$1/month)
Total:             $0-1/month
```

**Option B: Railway**
```
PostgreSQL:        Included in $5 credits
Redis:             Included in $5 credits
App hosting:       Included in $5 credits
Total:             $0-5/month
```

**Option C: Oracle Cloud Free Tier (Most Complex)**
```
VM (ARM):          $0 (always free)
Block storage:     $0 (200GB free)
Total:             $0/month
```

### Services

```
Gemini API:        $0 (free tier - 1500 requests/day)
GitHub Actions:    $0 (free for public repos)
Frontend hosting:  $0 (Vercel/Netlify free tier)
Monitoring:        $0 (Better Stack or built-in logging)
Domain:            $1/month (optional, use *.fly.dev subdomain for free)
```

### Total Monthly Cost: $0-5

**One-time Setup**:
- 0-2 hours to set up accounts
- $0 cost

---

## Success Metrics (For Portfolio & Learning)

### Technical Skills Demonstrated

**Backend Engineering**:
- [x] Go production application
- [x] RESTful API design
- [x] PostgreSQL database design
- [x] Redis caching strategies
- [x] Background job processing
- [x] AI/LLM API integration

**Infrastructure & DevOps**:
- [x] Docker containerization
- [x] CI/CD pipeline setup
- [x] Cloud deployment (Fly.io/Railway)
- [x] Monitoring and logging
- [x] Production-ready security practices

**Frontend**:
- [x] TypeScript application
- [x] Modern React development
- [x] API integration

### Portfolio Value

**For Job Interviews**:
- Live demo to show recruiters
- Source code on GitHub (shows code quality)
- Architecture diagram (shows systems thinking)
- Cost-conscious decisions (shows business awareness)
- Production deployment (shows full-stack capability)

**Talking Points**:
- "Built with Go to learn the most in-demand backend language for 2026"
- "Integrated Google's Gemini AI for sentiment analysis"
- "Deployed on Fly.io with CI/CD, costs less than $5/month"
- "Demonstrates microservices-ready architecture (modular monolith)"
- "Production-ready: auth, caching, rate limiting, monitoring"

---

## Technology Decision Log

| Decision | Options Considered | Chosen | Rationale |
|----------|-------------------|--------|-----------|
| **Backend Language** | Go, Rust, TypeScript | **Go** | Most in-demand for 2026, great for learning, easy deployment |
| **Frontend** | TypeScript, plain JS | **TypeScript** | In-demand skill, catches bugs early |
| **Database** | PostgreSQL, MySQL, SQLite | **PostgreSQL** | Production-ready, free tier available, JSON support |
| **Caching** | Redis, Memcached, in-memory | **Redis** | Industry standard, free tier, simple |
| **AI Provider** | OpenAI, Anthropic, Gemini, local models | **Gemini** | FREE tier with generous limits, good quality |
| **Hosting** | AWS, GCP, Fly.io, Railway, Heroku | **Fly.io** | Best free tier, easy deployment, PostgreSQL included |
| **Architecture** | Microservices, Monolith | **Modular Monolith** | Realistic for solo dev, can split later |
| **Message Queue** | Kafka, RabbitMQ, SQS, Redis | **In-memory (Go channels)** | Simple, free, sufficient for side project |

---

## Risk Assessment & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Free tier limits exceeded** | Medium | Low | Monitor usage, add rate limiting, upgrade if needed ($5-10/month) |
| **Gemini API changes/pricing** | Medium | Low | Abstract AI service, easy to swap providers |
| **Time commitment** | High | Medium | Keep scope minimal, MVP in 4 weeks, rest is optional |
| **Scope creep** | Medium | High | Stick to roadmap, resist adding features until MVP done |
| **Learning curve (Go)** | Medium | Medium | Go is beginner-friendly, excellent docs, start simple |

---

## Next Steps (Start Now!)

### Week 1 Actions

**Day 1-2: Setup** ✅ COMPLETE
1. ✅ Create GitHub repository
2. ✅ Initialize Go project: `go mod init github.com/sfumato00/content-analyzer`
3. ✅ Set up Docker Compose with PostgreSQL + Redis
4. ✅ Create Dockerfile with multi-stage build
5. ✅ Implement database migrations (users, submissions, analyses tables)
6. ✅ Set up database connection layer with pgx
7. ✅ Set up Redis cache layer
8. ✅ Get Gemini API key (user action): https://makersuite.google.com/app/apikey
9. ✅ Create basic project structure

**Day 3-5: Core Backend** ✅ COMPLETE
1. ✅ Implement PostgreSQL connection and migrations (completed in Day 1-2)
2. ✅ Create User model and auth endpoints (register/login/logout)
3. ✅ Implement JWT middleware
4. ✅ Write tests for auth logic
5. ✅ Create Makefile for development workflow

**Day 6-7: First Feature**
1. Create submissions endpoint (POST /api/v1/submissions)
2. Implement basic Gemini API call
3. Store results in database
4. Test end-to-end flow

### Useful Resources (All Free)

**Learning Go**:
- [Go by Example](https://gobyexample.com/)
- [Official Go Tour](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)

**Gemini API**:
- [Gemini API Quickstart](https://ai.google.dev/tutorials/rest_quickstart)
- [API Reference](https://ai.google.dev/api/rest)

**Deployment**:
- [Fly.io Go Guide](https://fly.io/docs/languages-and-frameworks/golang/)
- [Railway Deployment Docs](https://docs.railway.app/)

**Portfolio Building**:
- [How to present side projects](https://www.youtube.com/results?search_query=present+side+projects+interview)

---

## Appendix: Scaling Path (Future)

If this project takes off or you want to add it later:

### From Modular Monolith → Microservices

**Step 1**: Extract AI service
- Move `internal/services/ai` to separate Go application
- Communicate via REST or gRPC
- Deploy as separate container

**Step 2**: Add real message queue
- Introduce Redis Streams or RabbitMQ (cheap)
- Decouple API from AI processing
- Better handling of rate limits

**Step 3**: Add Kubernetes
- Local: k3s or kind
- Cloud: GKE Autopilot (cheaper than EKS)
- Learn K8s without overspending

**Total Cost After Scaling**: Still ~$10-30/month if done right

---

**Document Version**: 2.0 - Personal Side Project Edition
**Maintainer**: You!
**Next Review**: After completing MVP
