# Backend_Master_Class

## Learn everything about backend web development: Golang, Postgres, Redis, Gin, gRPC, Docker, Kubernetes, AWS, CI/CD by TECH SCHOOL

## Tools

[Docker]('')
[SQLC]('https://sqlc.dev/')
[Migrate]('https://github.com/golang-migrate/migrate')
[Testify]('github.com/stretchr/testify')

# Docker

---

Run new image Docker

```shell
# docker run --name <container_name> -p <host_port>:<container_port> -e <environment> -e <environment> -d <image>:<tag>
docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine
```

How exec command in Docker Container

```shell
# docker exec -it <continer_nameORcontainer_id> <command>
docker exec -it postgres12 createdb --username=root --owner=root simple_bank
```

# Migrate

---

Start migrate

```shell
# migrate create -ext <extension> -dir <location-path> -seq <name>
migrate create -ext sql -dir db/migration -seq init_schema

```

Migrate UP

```shell
# migrate -path <location-path> -database <string_connection_db> --verbose <up_or_down>
migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose up
```

Migrate DOWN

```shell
# migrate -path <location-path> -database <string_connection_db> --verbose <up_or_down>
migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose down

```
