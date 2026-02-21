create-migration:
	migrate create -ext sql -dir internal/infra/persistence/database/migrations -seq $(name)