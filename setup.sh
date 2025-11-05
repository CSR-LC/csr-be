#!/bin/bash

# CSR Backend - Quick Setup Script
# This script automates the initial setup process

set -e  # Exit on error

echo "🚀 CSR Backend Setup Script"
echo "=============================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check prerequisites
echo "Step 1: Checking prerequisites..."
echo "-----------------------------------"

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go is installed: $GO_VERSION"
else
    print_error "Go is not installed. Please install Go 1.25+ from https://golang.org/dl/"
    exit 1
fi

# Check Docker
if command -v docker &> /dev/null; then
    print_success "Docker is installed"
else
    print_error "Docker is not installed. Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker Compose
if command -v docker-compose &> /dev/null; then
    print_success "Docker Compose is installed"
else
    print_error "Docker Compose is not installed"
    exit 1
fi

# Check Make
if command -v make &> /dev/null; then
    print_success "Make is installed"
else
    print_error "Make is not installed"
    exit 1
fi

echo ""

# Update go.mod if needed
echo "Step 2: Checking Go version in go.mod..."
echo "------------------------------------------"

if [ -f "go.mod" ]; then
    CURRENT_GO_VERSION=$(grep "^go " go.mod | awk '{print $2}')
    print_info "Current Go version in go.mod: $CURRENT_GO_VERSION"
    
    if [ "$CURRENT_GO_VERSION" != "1.25" ]; then
        print_info "Updating go.mod to use Go 1.25..."
        sed -i.bak 's/^go .*/go 1.25/' go.mod
        print_success "Updated go.mod to Go 1.25"
        
        print_info "Running go mod tidy..."
        go mod tidy
        print_success "Go modules updated"
    else
        print_success "go.mod already uses Go 1.25"
    fi
else
    print_error "go.mod not found. Are you in the project root directory?"
    exit 1
fi

echo ""

# Install tools
echo "Step 3: Installing required tools..."
echo "--------------------------------------"

print_info "This may take a few minutes..."
if make setup; then
    print_success "Tools installed successfully"
else
    print_error "Failed to install tools"
    exit 1
fi

echo ""

# Generate code
echo "Step 4: Generating code..."
echo "----------------------------"

print_info "Generating Swagger, Ent, and Mock code..."
if make generate; then
    print_success "Code generated successfully"
else
    print_error "Failed to generate code"
    exit 1
fi

echo ""

# Check config.json
echo "Step 5: Checking configuration..."
echo "-----------------------------------"

if [ -f "config.json" ]; then
    print_success "config.json found"
    
    # Check if host is set to localhost
    if grep -q '"host": "localhost"' config.json; then
        print_success "Database host is configured for local development"
    else
        print_info "Database host is not set to localhost"
        print_info "For local development, it should be 'localhost'"
        print_info "For Docker Compose, it should be 'postgres'"
    fi
else
    print_error "config.json not found"
    print_info "Please create config.json based on config.example.json or the documentation"
    exit 1
fi

echo ""

# Start database
echo "Step 6: Starting PostgreSQL database..."
echo "-----------------------------------------"

print_info "Starting database container..."
if make db; then
    print_success "Database container started"
    
    print_info "Waiting for database to be ready (10 seconds)..."
    sleep 10
    
    # Check if database is healthy
    if docker ps | grep -q "db-local.*healthy"; then
        print_success "Database is healthy and ready"
    else
        print_info "Database is starting up..."
        print_info "You may need to wait a bit longer before running the app"
    fi
else
    print_error "Failed to start database"
    exit 1
fi

echo ""
echo "=============================="
echo " Setup Complete!"
echo "=============================="
echo ""
echo "Next steps:"
echo "1. Start the application:"
echo -e "   ${GREEN}make run${NC}"
echo ""
echo "2. Access the API:"
echo -e "   - API: ${GREEN}http://127.0.0.1:8080/api${NC}"
echo -e "   - Swagger: ${GREEN}http://127.0.0.1:8080/api/docs${NC}"
echo ""
echo "3. Test the API:"
echo -e "   ${GREEN}curl -X POST http://127.0.0.1:8080/api/v1/users/ -v${NC}"
echo ""
echo "4. Run tests (optional):"
echo -e "   ${GREEN}make test${NC}"
echo ""
echo "For more information, see:"
echo "- README.md - Overview and reference"
echo "- SETUP_GUIDE.md - Detailed step-by-step guide"
echo "- CHECKLIST.md - Setup verification checklist"
echo ""
echo -e "Run ${GREEN}make help${NC} to see all available commands"
echo ""
