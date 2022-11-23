![gorm-cli](./gorm-cli.png)

Online doc: https://gorm-cli.bugfree.app

# gorm-cli
The CLI tool for [gorm ORM](https://gorm.io/). Such as migration, and seed.

Currently, [gorm-cli](https://github.com/iKala/gorm-cli) supports two dialects: `mysql` and `postgres`.

## How gorm-cli works
[gorm-cli](https://github.com/iKala/gorm-cli) maintains a table - `gorm_meta` to record the migration execution history. Support developers to keep the DB migration in code, and easy to roll back.

## Installation
```shell
$ go install github.com/iKala/gorm-cli@latest
$ gorm-cli
Option of migration is needed - (db:prebuild / db:init / db:migrate / db:rollback / db:create_migration)
```

### macOS
Let's say you got `command not found`.
```shell
# Add the env to your cli (whatever fish, zsh, bash...) then try again
export PATH=$PATH:$(go env GOPATH)/bin
```

# How to use [gorm-cli](https://github.com/iKala/gorm-cli)

## Initialization
First, prepare your connection settings.

[gorm-cli](https://github.com/iKala/gorm-cli) will load the gorm drivers automatically to support various dialects such as `mysql`, `postgres`...like `gorm` do, but you need to install the drive yourself.

```shell
# Get the mysql driver when you use the mysql dialect.
# When the module was recorded in your go.mod, you don't have to do this for far.
go get gorm.io/driver/mysql
```

Generate your connection file.
```shell
$ gorm-cli db:init
```

[gorm-cli](https://github.com/iKala/gorm-cli) supports two way to configure your DB connection.

### YAML
Place the yaml config at your project root where you run `gorm-cli` command and name it `.gorm-cli.yaml`.

```yaml
db:
  host: localhost
  port: 3306
  dialects: mysql
  user: root
  password: password
  dbname: hentai
  charset: utf8mb4

migration:
  path: ./migration
```

```yaml
db:
  host: localhost
  port: 5432
  dialects: postgres
  user: root
  password: password
  dbname: hentai
  sslmode: disable
  timezone: UTC

migration:
  path: ./migration
```

### ENV

| Name | Description | Default Value |
|------|-------------|---------------|
| DB_HOST | DB Host | |
| DB_DIALECTS | mysql, postgres...etc, check the doc [here](https://gorm.io/docs/write_driver.html#Write-new-driver) | |
| DB_USER | DB user | |
| DB_PASSWORD | DB password | |
| DB_DBNAME | Target DB name | |
| DB_CHARSET | Connection charset (Available when dialect is `mysql`) | |
| DB_SSLMODE | Connection charset (Available when dialect is `postgres`) | `disable` |
| DB_TIMEZONE | Connection charset (Available when dialect is `postgres`) | `UTC` |
| MIGRATION | The folder you place the migrations. | `./migration` |

## Build your migration

```shell
# Replace `create_user` to YOUR_MIGRATION_TITLE
$ db:create_migration create_user
Migration created. ./migration/20221120125320_create_user.go

# The empty migration will be created in the migration folder you configured.
```

Since we're building migrations on top of gorm, you should check the gorm's doc to know what you need. [Toturials - Migration](https://gorm.io/docs/migration.html)

Here are some examples

### Create table via model declearation.
```go
package main

import (
  "github.com/foo/bar/model"
  "gorm.io/gorm"
)

type migration string

// Up - Changes for the migration.
func (m migration) Up(db *gorm.DB) error {
  return db.AutoMigrate(&model.User{})
}

// Down - Rollback changes for the migration.
func (m migration) Down(db *gorm.DB) error {
  return db.Migrator().DropTable(&model.User{})
}

var Migration migration
```

### Prepare some seed data.
```go
package main

import (
  "github.com/foo/bar/model"
  "gorm.io/gorm"
)

type migration string

// Up - Changes for the migration.
func (m migration) Up(db *gorm.DB) error {
  return db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Exec("INSERT INTO tag_type_weights (type, weight) VALUES (?, ?)", model.TagTypeCharacter, 1.6).Error; err != nil {
      return err
    }
    if err := tx.Exec("INSERT INTO tag_type_weights (type, weight) VALUES (?, ?)", model.TagTypeCopyright, 1.4).Error; err != nil {
      return err
    }
    if err := tx.Exec("INSERT INTO tag_type_weights (type, weight) VALUES (?, ?)", model.TagTypeArtist, 1.3).Error; err != nil {
      return err
    }
    if err := tx.Exec("INSERT INTO tag_type_weights (type, weight) VALUES (?, ?)", model.TagTypeGeneral, 1.1).Error; err != nil {
      return err
    }

    return nil
  })
}

// Down - Rollback changes for the migration.
func (m migration) Down(db *gorm.DB) error {
  return db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Exec("DELETE FROM tag_type_weights WHERE type = ? AND weight = ?", model.TagTypeCharacter, 1.6).Error; err != nil {
      return err
    }
    if err := tx.Exec("DELETE FROM tag_type_weights WHERE type = ? AND weight = ?", model.TagTypeCopyright, 1.4).Error; err != nil {
      return err
    }
    if err := tx.Exec("DELETE FROM tag_type_weights WHERE type = ? AND weight = ?", model.TagTypeArtist, 1.3).Error; err != nil {
      return err
    }
    if err := tx.Exec("DELETE FROM tag_type_weights WHERE type = ? AND weight = ?", model.TagTypeGeneral, 1.1).Error; err != nil {
      return err
    }

    return nil
  })
}

var Migration migration
```

## Execute migrations
```shell
# `db:migrate` command will execute the `Up` function in your migrations.
$ gorm-cli db:migrate
Migrated. 1 20221120125320_create_user.so

# `db:rollback` or `db:rollback [STEPS]` executes the `Down` function to revert the migration.
$ gorm-cli db:rollback
✔ Are you sure to rollback the migration with all steps? (Yes/No):
```

[gorm-cli](https://github.com/iKala/gorm-cli) will compile the migrations into `.so` files and cache them at `.plugins` folder. Once the `.so` file built, [gorm-cli](https://github.com/iKala/gorm-cli) will not change it unless the `-f` flag is applied.

```shell
# The `-f` flag to force rebuild the `.so` file, it's useful when you're testing and need to retry.
$ gorm-cli db:migrate -f
$ gorm-cli db:rollback -f
```

## Prebuild migrations
As mentioned above, [gorm-cli](https://github.com/iKala/gorm-cli) will compile the migrations into `.so` automatically. However, when `go` is not installed in the production environment, you need to pre-build the migrations to be portable.

```shell
$ gorm-cli db:prebuild
Migrated. 1 20221120125320_create_user.so
connection.so created

# Check generated files
$ ls -l ./migration/.plugins
# Move files to where ever you want.
$ cp -r ./migration/.plugins [ANYWHERE]
```

## Dockerfile
You can also use docker to run the migration.

Here is the example.
```dockerfile
##
## Build
##
FROM golang:alpine AS build

WORKDIR /

RUN apk update && apk add musl-dev gcc

COPY ./model ./model
COPY ./go.mod ./
COPY ./go.sum ./

ARG DB_HOST
ARG DB_DIALECTS
ARG DB_USER
ARG DB_PASSWORD
ARG DB_DBNAME
ARG DB_CHARSET

ENV MIGRATION_PATH=/model/migration

RUN go install github.com/iKala/gorm-cli@latest && go mod download
RUN gorm-cli db:prebuild && gorm-cli db:init

##
## Deploy
##
FROM alpine

COPY --from=build /model/migration/ /migration/
COPY --from=build /go/bin/gorm-cli /go/bin/gorm-cli

ENV MIGRATION_PATH=/migration

CMD ["/go/bin/gorm-cli", "db:migrate"]

```

```shell
$ docker run --rm gorm-cli-migration
```

---

## Thanks for all of you guys liking this project
It's very welcome to give your PRs to make this little tool much easier to use. 🙂

