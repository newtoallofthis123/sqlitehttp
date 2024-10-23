# Sqlite Http

An experimental project that aims to explore the possibility of safely
exposing a sqlite database over the http protocol.

The aim is to make an actual usable http server that can be self hosted and
make it easy to interact with a sqlite database.

To now, only local databases ie file based sqlite3 databases are supported.
However, the plan is to support the use of sqlite3 over a GET based network protocol as well.

## Functionaility Implemented

- [x] Safe SELECT queries
- [x] Safe DDL queries
- [x] Local SQLite3 database support
- [ ] HTTP routes for queries
- [ ] Better Code Coverage
- [ ] Better Code Organization

## Usage

```bash
go install github.com/newtoallofthis123/sqlitehttp
sqlitehttp -db <path to sqlite3 database>
```

## License

The code is licensed under the MIT license. See the [LICENSE](LICENSE) file for more details.
