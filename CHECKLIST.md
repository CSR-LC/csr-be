# CSR Backend - Quick Setup Checklist

Use this checklist to ensure everything is set up correctly.

## ✅ Pre-Setup Checklist

- [ ] Go 1.25+ installed (`go version`)
- [ ] Docker installed and running (`docker info`)
- [ ] Docker Compose installed (`docker-compose --version`)
- [ ] Make installed (`make --version`)
- [ ] Git installed (`git --version`)

## ✅ Repository Setup

- [ ] Clone repository: `git clone https://github.com/CSR-LC/csr-be.git`
- [ ] Enter directory: `cd csr-be`
- [ ] Check go.mod exists
- [ ] Update Go version in go.mod to `go 1.25`
- [ ] Run: `go mod tidy`

## ✅ Replace Updated Files

From the files I created, replace these in your project:

- [ ] Replace `Makefile` with `Makefile.new`

  ```bash
  cp Makefile Makefile.backup
  cp Makefile.new Makefile
  ```

- [ ] Replace `README.md` with the new version

  ```bash
  cp README.md README.md.backup
  cp README.md.new README.md
  ```

- [ ] Add `SETUP_GUIDE.md` to your project (new file)

## ✅ Tool Installation

- [ ] Run: `make setup`
- [ ] Verify swagger: `swagger version`
- [ ] Verify ent: `ent version`
- [ ] Verify mockery: `mockery --version`

## ✅ Code Generation

- [ ] Run: `make generate`
- [ ] Verify `internal/generated/swagger/` has files
- [ ] Verify `internal/generated/ent/` has files
- [ ] Verify `internal/generated/mocks/` has files
- [ ] No errors during generation

## ✅ Configuration

- [ ] `config.json` exists
- [ ] Database host is set to `"localhost"` for local development
- [ ] Database port is `5432`
- [ ] Database user is `"csr"`
- [ ] Database name is `"csr"`

## ✅ Database Setup

- [ ] Run: `make db`
- [ ] Wait 10 seconds
- [ ] Verify container running: `docker ps | grep db-local`
- [ ] Test connection: `docker exec -it db-local psql -U csr -d csr -c "SELECT 1;"`

## ✅ Application Start

- [ ] Run: `make run`
- [ ] No errors in console
- [ ] Server starts on port 8080
- [ ] Can access http://127.0.0.1:8080/api
- [ ] Can access http://127.0.0.1:8080/api/docs

## ✅ Verification Tests

- [ ] Test API: `curl http://127.0.0.1:8080/api`
- [ ] Test endpoint: `curl -X POST http://127.0.0.1:8080/api/v1/users/ -v`
- [ ] Response is valid JSON with user ID
- [ ] Swagger UI loads properly in browser

## ✅ Optional: Docker Compose Setup

- [ ] Update config.json: database host to `"postgres"`
- [ ] Run: `make rebuild_project`
- [ ] Verify both containers running: `docker ps`
- [ ] Test API: `curl http://localhost:8080/api`
- [ ] Stop: `make stop_project`

## ✅ Optional: Run Tests

- [ ] Run: `make test`
- [ ] All tests pass
- [ ] Run: `make coverage`
- [ ] Coverage report displays

## Common Issues Checklist

If something doesn't work, check these:

- [ ] GOPATH/bin in PATH: `echo $PATH | grep "$(go env GOPATH)/bin"`
- [ ] Port 8080 available: `lsof -i :8080`
- [ ] Docker daemon running: `docker info`
- [ ] Database container healthy: `docker ps` (STATUS should show "healthy")
- [ ] No conflicting PostgreSQL on 5432: `lsof -i :5432`
- [ ] All go.mod dependencies downloaded: `go mod download`

## File Changes Summary

### Files to Update:

1. **Makefile** - Added `rebuild_project`, improved organization, added `help` command
2. **README.md** - Complete rewrite with better structure, correct instructions, Go 1.25
3. **go.mod** - Update `go` directive to `1.25`

### New Files to Add:

1. **SETUP_GUIDE.md** - Detailed step-by-step setup instructions
2. **CHECKLIST.md** - This file (optional, for tracking progress)

## Success Criteria

✅ **All green checkboxes above**

✅ **Application runs without errors**

✅ **Can access Swagger UI**

✅ **API endpoints respond correctly**

✅ **Database connection works**

---

**When everything is checked:** You're ready to develop! 🚀

**If stuck:** Refer to SETUP_GUIDE.md for detailed troubleshooting steps.
