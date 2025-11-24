# 🔧 Backend Setup & Troubleshooting

## ✅ Error Fixed!

### Issues Resolved:
1. **404 Error on `/api/standings`** - Server restarted successfully
2. **"main redeclared" lint error** - `seed.go` moved to `cmd/` folder

---

## 🚀 Running the Backend

### Start the main server:
```bash
cd Backend
go run main.go
```

Server will run on: `http://localhost:8080`

---

## 👤 Create Admin User

If you need to create or reset the admin user:

```bash
cd Backend
go run cmd/seed.go
```

This will:
- Create admin user with credentials: `admin` / `admin123`
- Or reset password if user already exists

---

## 📡 API Endpoints

### Authentication
- `POST /api/register` - Create admin (dev only)
- `POST /api/login` - Login

### Tournaments
- `GET /api/tournaments` - List all tournaments
- `POST /api/tournaments` - Create tournament
- `PUT /api/tournaments/:id` - Update tournament

### Teams
- `GET /api/teams?tournament_id=X` - List teams
- `POST /api/teams` - Create team

### Rounds
- `GET /api/rounds?tournament_id=X` - List rounds
- `POST /api/rounds` - Create round

### Matches
- `GET /api/matches?round_id=X` - List matches
- `POST /api/matches` - Create match

### Ballots
- `POST /api/ballots` - Submit scores

### Standings
- `GET /api/standings?tournament_id=X` - Get team standings

### Articles
- `GET /api/articles` - List articles
- `POST /api/articles` - Create article

---

## 🐛 Common Issues

### Port 8080 already in use:
```bash
# Find process using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID with actual number)
taskkill /PID <PID> /F

# Restart server
go run main.go
```

### Database connection error:
Check `models/setup.go` and ensure PostgreSQL is running:
- Host: localhost
- Port: 5433
- User: admin
- Password: password123
- Database: eds_upi

---

## 📝 Notes

- `seed.go` is now in `cmd/` folder to avoid conflicts with `main.go`
- CORS is enabled for all origins (development only)
- Auto-migration runs on server start
