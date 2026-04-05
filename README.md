# Entain BE Technical Test

This project extends the original technical test with the requested backend features across the `racing`, `sports`, and `api` services.

## Features Implemented

### Racing

- Added optional `visible` filtering to `ListRaces`
- Added ordering by `advertised_start_time`
- Added optional sort support through the request filter
- Added derived `status` field:
  - `OPEN`
  - `CLOSED`
- Added `GetRace` RPC to fetch a single race by ID

### Sports

- Added a new standalone `sports` service
- Added `ListEvents` for sports event listing
- Added optional `visible` filtering
- Added ordering by `advertised_start_time`
- Added derived `status` field:
  - `OPEN`
  - `CLOSED`

### API Gateway

- Wired the API gateway to expose both:
  - `racing`
  - `sports`

## Project Structure

```text
entain/
├─ api/
├─ racing/
├─ sports/
├─ README.md
```

## Notes

- `racing` and `sports` are implemented as separate services
- Both services include test coverage for the main feature behaviour
