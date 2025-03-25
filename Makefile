migrate:
	goose -dir migrations postgres "postgres://postgres:password@localhost:5432/currency_db?sslmode=disable" up