#!/bin/sh
go build -o ./builds/ ./cmd/...
# sudo docker run --name reg-rdb -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=xcpc_team_reg -v ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql -d postgres
# sudo docker run --name reg-redis -p 6379:6379 -d redis
