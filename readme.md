## Template golang siap pakai

### Cara pakai

1. Clone repo ini
2. Ubah nama folder `golang-template` menjadi nama project
3. Ubah nama modulde  `golang-template` menjadi nama project
4. rename `example.local.env` atau apapun itu menjadi `*.env`
5. seluruh configurasi ada di file `*.env`
6. jalankan `go mod download` untuk mendownload dependency

### Cara menjalankan

1. Jalankan `go run cmd/main/main.go`

### Cara menjalankan test

1. Jalankan `go test ./...`

### Cara menjalankan test coverage

1. Jalankan `go test ./... -coverprofile=coverage.out`

### Package yang digunakan

1. [x] [gin-gonic](https://gin-gonic.com) sebagai http server
2. [x] [gorm](https://gorm.io) sebagai orm dan migration
3. [x] [go-redis](https://redis.uptrace.dev/) sebagai cache
4. [x] [golang-jwt](https://github.com/golang-jwt/jwt) sebagai jwt
5. [x] [godotenv](https://github.com/joho/godotenv) sebagai env loader
6. [x] [gin-swagger](https://github.com/swaggo/gin-swagger) sebagai swagger generator
7. [ ] [casbin](https://casbin.org/) sebagai acl

### Database yang bisa digunakan

1. [x] [mysql](https://www.mysql.com/)
2. [x] [postgres](https://www.postgresql.org/)
3. [x] [sqlite](https://www.sqlite.org/index.html)
4. [x] [mssql](https://www.microsoft.com/en-us/sql-server/sql-server-2019)

### Struktur folder

```
├── cmd
│   └── main // entry point
│       └── main.go
├── etc
│   ├── config // config file
│   │   ├── *.env
|
├── pkg
│   ├── utils // custom package
│   │
│   ├── auth // auth package
│   │

```

### note: template dan dokumentasi masi tahap pengembangan