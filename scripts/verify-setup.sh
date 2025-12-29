#!/bin/bash
# Setup verification script

set -e

echo "=================================="
echo "Content Analyzer Setup Verification"
echo "=================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Test 1: Check Docker
echo "1. Checking Docker..."
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker is installed"
else
    echo -e "${RED}✗${NC} Docker is not installed"
    exit 1
fi

# Test 2: Check Docker Compose
echo "2. Checking Docker Compose..."
if command -v docker-compose &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker Compose is installed"
else
    echo -e "${RED}✗${NC} Docker Compose is not installed"
    exit 1
fi

# Test 3: Check .env file
echo "3. Checking .env file..."
if [ -f .env ]; then
    echo -e "${GREEN}✓${NC} .env file exists"

    # Check critical variables
    if grep -q "GEMINI_API_KEY=" .env && ! grep -q "GEMINI_API_KEY=API_KEY_HERE" .env; then
        echo -e "${GREEN}✓${NC} GEMINI_API_KEY is configured"
    else
        echo -e "${RED}✗${NC} GEMINI_API_KEY needs to be set in .env"
    fi

    if grep -q "JWT_SECRET=" .env && ! grep -q "change_this_to_a_random_secret" .env; then
        echo -e "${GREEN}✓${NC} JWT_SECRET is configured"
    else
        echo -e "${RED}✗${NC} JWT_SECRET needs to be set in .env"
    fi
else
    echo -e "${RED}✗${NC} .env file not found. Copy from .env.example"
    exit 1
fi

# Test 4: Check Docker services
echo "4. Checking Docker services..."
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓${NC} Docker services are running"
else
    echo -e "${RED}✗${NC} Docker services are not running. Run: docker-compose up -d"
fi

# Test 5: Check PostgreSQL
echo "5. Checking PostgreSQL..."
if docker-compose exec -T postgres pg_isready -U postgres &> /dev/null; then
    echo -e "${GREEN}✓${NC} PostgreSQL is ready"
else
    echo -e "${RED}✗${NC} PostgreSQL is not ready"
fi

# Test 6: Check Redis
echo "6. Checking Redis..."
if docker-compose exec -T redis redis-cli ping &> /dev/null; then
    echo -e "${GREEN}✓${NC} Redis is ready"
else
    echo -e "${RED}✗${NC} Redis is not ready"
fi

# Test 7: Check API health
echo "7. Checking API health..."
if curl -s http://localhost:8080/health | grep -q "healthy"; then
    echo -e "${GREEN}✓${NC} API is healthy"
else
    echo -e "${RED}✗${NC} API is not responding"
fi

# Test 8: Check database tables
echo "8. Checking database schema..."
if docker-compose exec -T postgres psql -U postgres -d content_analyzer -c "\dt" 2>/dev/null | grep -q "users"; then
    echo -e "${GREEN}✓${NC} Database migrations have run"
else
    echo -e "${RED}✗${NC} Database migrations need to run"
fi

echo ""
echo "=================================="
echo "Verification Complete!"
echo "=================================="
