# Plenum

Plenum is a Vercel-ready, serverless Golang backend for AeroGo pre-processing and telemetry workflows.

## Implemented API Endpoints

- `POST /api/foam/generate-dict`: Generates OpenFOAM dictionaries (`blockMeshDict`, `controlDict`, `fvSchemes`).
- `POST /api/auto/drivaer-bounds`: Computes DrivAer domain dimensions and wall distances.
- `GET /api/chem/arrhenius`: Computes Arrhenius rate constants over a temperature range.
- `POST /api/telemetry/update`: Ingests job telemetry updates.
- `GET /api/telemetry/update?jobId=<id>`: Reads latest cached telemetry for a job.

## Project Structure

```text
api/
  foam/generate.go
  auto/drivaer-bounds.go
  chem/arrhenius.go
  telemetry/update.go
pkg/
  foam/generator.go
  aerodynamics/calc.go
  chem/arrhenius.go
  telemetry/store.go
vercel.json
go.mod
```

## Notes

- Current telemetry cache uses in-memory storage for local/serverless compatibility.
- For production durability, replace `pkg/telemetry` store with Upstash Redis-backed implementation.
