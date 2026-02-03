module github.com/rootix/portfolio

go 1.22

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/gorilla/mux v1.8.1
	github.com/uptrace/bun v1.2.5
	github.com/uptrace/bun/dialect/pgdialect v1.2.5
	github.com/uptrace/bun/driver/pgdriver v1.2.5
	github.com/yuin/goldmark v1.5.6
	golang.org/x/crypto v0.26.0
)

replace github.com/uptrace/bun => github.com/uptrace/bun v1.2.5
replace github.com/uptrace/bun/dialect/pgdialect => github.com/uptrace/bun/dialect/pgdialect v1.2.5
replace github.com/uptrace/bun/driver/pgdriver => github.com/uptrace/bun/driver/pgdriver v1.2.5
