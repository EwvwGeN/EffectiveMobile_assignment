# EffectiveMobile_assignment

Simple api to work with catalog of cars

# Table of contents

- [Startup](#startup)
    - [Configuration](#configuration)
    - [Preparing environment variables](#preparing-environment-variables)
    - [Direct startup](#direct-startup)
    - [Docker startup](#docker-startup)
- [Http handlers description](#ttp_handlers_description)

## Startup

This section describes the configuration of the application and how to run it

### Configuration

The config has the following fields:

```yaml
log_level: debug
http:
  port: 9099
  host: 0.0.0.0
postgres:
  db_con_format: postgres
  db_host: postgres
  db_port: 5432
  db_user: user
  db_pass: pwd
  db_name: test-db
  db_tbl_car: car_table
car_info_getter: http://localhost:8080/info
```

- `log_level` - level reports the minimum record level that will be logged.
- `http` - settings for http server.
- `postgres` - setting for connection and name of tabbles that will be used.
- `data_collect_time` - interval for auto collecting data (products and categories) from source.
- `car_info_getter` - the link of source from which data will be collected.

Also, the following path `storage/migrations/init.sql` contains a migration for creating a database.

### Preparing environment variables

You can use script to convert .yaml to .env file

`go run config_to_env.go <path_to_config>`

### Direct startup

You can use build command to get bin file :</br>
`CGO_ENABLED=0 GOOS=linux go build -o <output path> <path to main.go>`

Than you need to start file via command:</br>
`file -config=<path_to_config>` - if you want to use config file.
Or just run file without flag to use env variables.

Also you can run service by using `go run`: `go run ./cmd/server/main.go -config=./configs/config.yaml` where you also can use config file or env variables.

### Docker startup

You can start only service by launching Dockerfile or start service with the database by launchig docker-compose file: `docker-compose up`

## Http handlers description

You can see all http handlers by visiting the swagger documentation via link:

`($service_host):($service_port)/api/swagger/index.html`