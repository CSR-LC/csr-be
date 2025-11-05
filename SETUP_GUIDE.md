# Step-by-Step Guide: Running CSR Backend Locally

This guide will help you get the CSR backend running on your local machine for development.

## Prerequisites Check

Before starting, verify you have:

```bash
# Check Go version (should be 1.25+)
go version

# Check Docker
docker --version
docker-compose --version

# Check Make
make --version

# Check Git
git --version
```

## Step 1: Clone and Enter the Repository

```bash
git clone https://github.com/CSR-LC/csr-be.git
cd csr-be
```

## Step 2: Verify Go Version

Check your go version. It should specify Go 1.25 or later:

```bash
go version
```

## Step 3: Install Required Tools

```bash
make setup
```

This installs:

- `swagger` - For API code generation
- `ent` - For ORM code generation
- `mockery` - For test mock generation

**Verify installation:**

```bash
swagger version
ent version
mockery --version
```

## Step 4: Generate Required Code

```bash
make generate
```

This command will:

1. Generate Swagger server and client code
2. Generate Ent ORM code
3. Generate test mocks

**What to expect:**

- New files in `internal/generated/swagger/`
- New files in `internal/generated/ent/`
- New files in `internal/generated/mocks/`

**Common issues:**

- If you get "command not found" errors, re-run `make setup`
- If generation fails, check that `swagger.yaml` and schema files exist

## Step 5: Configure Database Connection

Check your `config.json` file. For local development, ensure the database host is set to `localhost`:

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "csr",
    "database": "csr",
    "password": ""
  },
  "server": {
    "port": 8080
  }
}
```

**Note:** If this file doesn't exist, create it based on `config.local.example.json` or the format above.

## Step 6: Start PostgreSQL Database

```bash
make db
```

This starts a PostgreSQL container. Wait about 10 seconds for it to be ready.

**Verify database is running:**

```bash
docker ps
```

You should see a container named `db-local` running on port 5432.

**Test database connection:**

```bash
docker exec -it db-local psql -U csr -d csr -c "SELECT version();"
```

## Step 7: Run the Application

```bash
make run
```

**What to expect:**

```
INFO Starting server on :8080
INFO Swagger UI available at /api/docs
```

The application should now be running!

## Step 8: Verify the Application

### Check the API is responding:

```bash
curl http://127.0.0.1:8080/api
```

### Open Swagger UI in your browser:

```
http://127.0.0.1:8080/api/docs
```

## Step 9: Running Tests (Optional)

```bash
# Run unit tests
make test

# View coverage
make coverage
```

## Troubleshooting

### Issue: "Port 8080 already in use"

**Solution:**

```bash
# Find what's using the port
lsof -i :8080

# Kill the process or change the port in config.json
```

### Issue: "Cannot connect to database"

**Solution:**

```bash
# Check if PostgreSQL container is running
docker ps

# If not running, start it
make db

# Check logs
docker logs db-local

# Verify connection
docker exec -it db-local psql -U csr -d csr -c "SELECT 1;"
```

### Issue: "swagger command not found"

**Solution:**

```bash
# Reinstall tools
make setup

# Verify installation
which swagger
swagger version

# Check GOPATH/bin is in PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" || echo "Add $(go env GOPATH)/bin to PATH"
```

### Issue: "go.mod errors" or dependency issues

**Solution:**

```bash
# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy up
go mod tidy
```

### Issue: Database migrations needed

**Solution:**

```bash
# If you need to reset the database
docker-compose down -v
make db

# Wait for database to be ready, then run
make run
```

## Development Workflow

### Making changes:

1. **Edit code** in your preferred editor
2. **Regenerate if needed**: `make generate` (only if you changed swagger.yaml or schemas)
3. **Stop the running app**: Ctrl+C
4. **Restart**: `make run`

### Adding new API endpoints:

1. Edit `swagger.yaml`
2. Run `make generate/swagger`
3. Implement handlers
4. Restart: `make run`

### Modifying database schema:

1. Edit schema files in `internal/ent/schema/`
2. Run `make generate/ent`
3. Handle any migrations
4. Restart: `make run`

## Next Steps

Now that you have the application running:

1. **Explore the API** using Swagger UI at http://127.0.0.1:8080/api/docs
2. **Read the codebase** starting from `cmd/swagger/main.go`
3. **Run tests** with `make test`
4. **Try Docker Compose** setup (see README.md)

## Quick Reference

```bash
# Start development
make db          # Start database
make generate    # Generate code
make run         # Run app

# During development
make test        # Run tests
make lint        # Check code quality

# Stop everything
Ctrl+C           # Stop app
make stop_project # Stop Docker containers
```

## Getting Help

- Check the main [README.md](README.md) for more details
- Review the [Makefile](Makefile) for all available commands
- Run `make help` to see command descriptions
- Check logs: `docker logs db-local` for database issues

---

**Success!** You should now have the CSR backend running locally at http://127.0.0.1:8080/api
