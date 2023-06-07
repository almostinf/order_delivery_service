## Order Delivery Service

### Quickstart
To run a project in `docker` containers, use the command:
```shell
make compose-up
```

To run unit tests:
```shell
make test
```

To run integration tests:
```shell
make integration-test-docker-up
```

### Dependencies
**For the tests:**
- **Unit tests:**
  - [testify](https://github.com/stretchr/testify)
  - [gomock](https://github.com/golang/mock)
- **Integration tests:**
  - [go-hit](https://github.com/Eun/go-hit)

**Basic web framework:** [Gin](https://github.com/gin-gonic/gin)

**Database operation:**
- [pgx](https://github.com/jackc/pgx)
- [goose](https://github.com/pressly/goose)

**For working with configuration files:** [cleanenv](https://github.com/ilyakaznacheev/cleanenv)

**Logger:** [zerolog](https://github.com/rs/zerolog)

**Linters**:
- [Smart Imports](https://github.com/pav5000/smartimports)
- [Golang-ci lint](https://golangci-lint.run/)

### Licence
**Licensed under:**
- MIT license ([LICENSE-MIT](https://github.com/seanmonstar/httparse/blob/master/LICENSE-MIT) or [https://opensource.org/licenses/MIT](https://opensource.org/licenses/MIT))