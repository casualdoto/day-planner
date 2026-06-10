# Day Planner

Local-first day planner for daily notes, tasks, tags, importance levels, and a calendar view.

Data is stored in a local PostgreSQL database. There are no accounts and no cloud backend.

## Stack

* Go 
* PostgreSQL
* Vue 3 + Vite
* TypeScript
* Docker Compose

## Run 

```powershell
docker compose up --build
```

Open:

```text
http://localhost:8080
```

Stop:

```powershell
docker compose down
```

Docker stores PostgreSQL data in the volume `day-planner_day_planner_postgres_data`.

Delete containers and data:

```powershell
docker compose down --volumes
```

## Tests

Backend:

```powershell
go test ./...
```

PostgreSQL integration tests:

```powershell
$env:DAY_PLANNER_TEST_DATABASE_URL = "postgres://day_planner:day_planner@localhost:5432/day_planner?sslmode=disable"
go test ./...
```

Frontend build/type check:

```powershell
cd frontend
npm run build
```

Docker:

```powershell
docker compose run --rm test
```