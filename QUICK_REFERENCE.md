# CSR Backend - Quick Reference Card

## 🚀 First Time Setup

```bash
# 1. Clone and enter
git clone https://github.com/CSR-LC/csr-be.git
cd csr-be

# 2. Update Go version in go.mod to: go 1.25
# Then run:
go mod tidy

# 3. Run automated setup (recommended)
chmod +x setup.sh
./setup.sh

# 4. Start the app
make run
```

## 📍 Access Points

- API Base: `http://127.0.0.1:8080/api`
- Swagger UI: `http://127.0.0.1:8080/api/docs`

## 🛠️ Common Commands

### Daily Development

```bash
make db          # Start database
make run         # Run the application
make test        # Run tests
make generate    # Regenerate code (after schema/API changes)
```

### Building

```bash
make build       # Build binary
make clean       # Clean generated files
```

### Docker Compose

```bash
make build_project    # Build containers
make rebuild_project  # Rebuild from scratch
make start_project    # Start all services
make stop_project     # Stop all services
make restart_project  # Quick restart
```

### Code Generation

```bash
make generate/swagger  # Only Swagger
make generate/ent      # Only Ent ORM
make generate/mocks    # Only test mocks
make generate          # All of the above
```

### Testing

```bash
make test              # Unit tests
make int-test          # Integration tests
make coverage          # Coverage report
make lint              # Run linter
```

### Help

```bash
make help              # Show all commands
```

## 🔧 Configuration Quick Tips

### Local Development

In `config.json`:

```json
{
  "database": {
    "host": "localhost"
  }
}
```

### Docker Compose

In `config.json`:

```json
{
  "database": {
    "host": "postgres"
  }
}
```

## ⚡ Quick Troubleshooting

### Database Connection Failed

```bash
docker ps                    # Check if db-local is running
make stop_project           # Stop everything
make db                     # Start database
# Wait 10 seconds
make run                    # Start app
```

### Port 8080 Already in Use

```bash
lsof -i :8080              # Find what's using it
kill -9 <PID>              # Kill the process
# Or change port in config.json
```

### Swagger/Ent Generation Fails

```bash
make setup                 # Reinstall tools
make clean                # Clean old files
make generate             # Regenerate
```

### Go Module Issues

```bash
go clean -modcache
go mod download
go mod tidy
```

## 📋 Pre-Flight Checklist

Before you start coding:

- [ ] Database running: `docker ps | grep db-local`
- [ ] Config correct for your mode (local/docker)
- [ ] Code generated: `make generate`
- [ ] App starts: `make run`
- [ ] API responds: `curl http://127.0.0.1:8080/api`

## 🧪 Testing Your Changes

```bash
# After making code changes:
make test                   # Run tests
make lint                   # Check code quality

# After changing swagger.yaml:
make generate/swagger
make run

# After changing DB schema:
make generate/ent
make run
```

## 📂 Important Files

- `config.json` - Application configuration
- `swagger.yaml` - API specification
- `Makefile` - Build commands
- `internal/ent/schema/` - Database schemas
- `cmd/swagger/` - Application entry point

## 🎯 Go Version

**Required**: Go 1.25+

Check your version:

```bash
go version
```

## 📚 Documentation Files

- `README.md` - Complete overview and reference
- `SETUP_GUIDE.md` - Step-by-step setup instructions
- `CHECKLIST.md` - Setup verification checklist
- `IMPROVEMENTS_SUMMARY.md` - What was changed and why
- `QUICK_REFERENCE.md` - This file

## 💡 Pro Tips

1. **Use make help** - See all available commands
2. **Keep database running** - Don't stop it between coding sessions
3. **Regenerate after schema changes** - Run `make generate`
4. **Check logs** - `docker logs db-local` for database issues
5. **Use Swagger UI** - Test your API changes interactively

## 🔄 Typical Development Workflow

```bash
# Morning startup
make db                    # Start database
make run                   # Start app in terminal 1

# In another terminal (terminal 2)
# Make your code changes...

# After changing API (swagger.yaml)
make generate/swagger
# Restart app in terminal 1 (Ctrl+C, then make run)

# After changing DB schema
make generate/ent
# Restart app in terminal 1 (Ctrl+C, then make run)

# Before committing
make test
make lint

# End of day
# Ctrl+C in terminal 1 (stops app)
# Keep database running for next day
```

## 🆘 Getting Help

1. Check error message carefully
2. Look in SETUP_GUIDE.md troubleshooting
3. Verify configuration in config.json
4. Check Docker containers: `docker ps`
5. Check database logs: `docker logs db-local`
6. Regenerate code: `make generate`
7. Clean and rebuild: `make clean && make generate`

---

**Print this card and keep it handy!** 📌

For detailed information, see:

- Full docs: **README.md**
- Setup help: **SETUP_GUIDE.md**
- Changes made: **IMPROVEMENTS_SUMMARY.md**
