.PHONY: up down logs build rebuild db-shell clean smoke migrate seed lint fmt test help ps redis-cli logs-api logs-rust

# =============================================================
# SHEET-METAL-AUTOMAT-ON -- Production Makefile
# =============================================================

## help: Bu yardim mesajini goster
help:
	@echo ================================================================
	@echo  SHEET-METAL-AUTOMAT-ON -- Makefile Komutlari
	@echo ================================================================
	@echo  make up Tum servisleri baslat (detached)
	@echo  make down Tum servisleri durdur
	@echo  make rebuild Yeniden build alip baslat
	@echo  make logs Canli log takibi
	@echo  make logs-api Sadece Go API loglarini takip et
	@echo  make logs-rust Sadece Rust Engine loglarini takip et
	@echo  make smoke Smoke testleri calistir
	@echo  make migrate Veritabani migrationlarini calistir
	@echo  make seed Veritabanina baslangic verisi yukle
	@echo  make test Tum testleri calistir (Go + Rust)
	@echo  make lint Go ve Rust linter calistir
	@echo  make fmt Go ve Rust formatter calistir
	@echo  make ps Calisan servislerin durumunu goster
	@echo  make db-shell PostgreSQL terminaline gir
	@echo  make redis-cli Redis CLIye baglan
	@echo  make clean Durdurulmus container ve cache temizle
	@echo  make clean-all Volumeler dahil her seyi sil (DIKKAT!)
	@echo ================================================================

## up: Tum servisleri arka planda baslat
up:
	@echo Starting all services...
	docker-compose up -d
	@echo System is up. Run 'make smoke' to verify.

## down: Tum servisleri durdur (volumeler korunur)
down:
	@echo Stopping all services...
	docker-compose down

## rebuild: Yeniden build alarak baslat
rebuild:
	@echo Rebuilding all services...
	docker-compose up --build -d
	@echo Rebuild complete.

## logs: Canli log takibi
logs:
	docker-compose logs -f

## logs-api: Sadece Go API loglarini takip et
logs-api:
	docker-compose logs -f main-api

## logs-rust: Sadece Rust Engine loglarini takip et
logs-rust:
	docker-compose logs -f geometry-engine

## smoke: Temel HTTP smoke testlerini calistir
smoke:
	@echo Running smoke tests...
	@bash scripts/smoke_test.sh

## migrate: Veritabani migrationlarini calistir
migrate:
	@echo Running database migrations...
	docker-compose exec main-api ./main-api migrate
	@echo Migration complete.

## seed: Veritabanina baslangic verisi yukle
seed:
	@echo Seeding database...
	docker-compose exec -T postgres psql -U sheetmetal_admin -d sheetmetal_prod < scripts/seed_data.sql
	@echo Seed complete.

## test: Tum testleri calistir
test:
	@echo Running Go tests...
	cd services/main-api && go test -v -race -coverprofile=coverage.out ./...
	@echo Running Rust tests...
	cd services/geometry-engine && cargo test --verbose
	@echo All tests complete.

## lint: Go ve Rust linter calistir
lint:
	@echo Running Go linter...
	cd services/main-api && go vet ./...
	@echo Running Rust Clippy...
	cd services/geometry-engine && cargo clippy -- -D warnings

## fmt: Go ve Rust formatter calistir
fmt:
	@echo Formatting Go code...
	cd services/main-api && go fmt ./...
	@echo Formatting Rust code...
	cd services/geometry-engine && cargo fmt

## ps: Calisan servislerin durumunu goster
ps:
	docker-compose ps

## db-shell: PostgreSQL interaktif terminale gir
db-shell:
	docker-compose exec postgres psql -U sheetmetal_admin -d sheetmetal_prod

## redis-cli: Redis CLI baglantisi
redis-cli:
	docker-compose exec redis redis-cli

## clean: Durdurulmus containerlar ve cache temizle
clean:
	@echo Cleaning up stopped containers...
	docker system prune -f
	@echo Cleanup complete.

## clean-all: Volumeler dahil tum Docker kaynaklarini sil (DIKKAT!)
clean-all:
	@echo WARNING: This will delete ALL data including volumes!
	docker-compose down -v
	docker system prune -af
