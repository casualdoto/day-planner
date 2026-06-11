# syntax=docker/dockerfile:1

FROM node:22-bookworm AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build


FROM golang:1.23-bookworm AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/day-planner ./cmd/app


FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN useradd --create-home --shell /usr/sbin/nologin appuser \
    && chown -R appuser:appuser /app

COPY --from=go-builder /out/day-planner /app/day-planner
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

ENV DAY_PLANNER_ADDR=0.0.0.0:8080

EXPOSE 8080

USER appuser

CMD ["/app/day-planner"]
