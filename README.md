# 🏭 Sheet Metal Automation Platform

> **CAD-to-MES Middleware** — Kurumsal B2B Sac Metal / CNC Üretim SaaS Platformu

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://golang.org)
[![Rust](https://img.shields.io/badge/Rust-2021-000000?logo=rust)](https://www.rust-lang.org)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-Proprietary-red)](./LICENSE)

Sac metal ve CNC imalat sektörü için geliştirilmiş, **gerçek zamanlı teklif motoru**, **AI destekli DFM analizi**, **ERP/MES entegrasyonu** ve **cloud-native** altyapıya sahip tam yığın B2B SaaS platformu.

---

## 🗺️ Mimari Genel Bakış

```
┌─────────────────────────────────────────────────────────┐
│  CAD Bridge (.NET / Python)                             │
│  SolidWorks ─ Inventor ─ SolidEdge ─ FreeCAD           │
└────────────────────┬────────────────────────────────────┘
                     │ WebSocket
┌────────────────────▼────────────────────────────────────┐
│  Frontend (Vue 3 + TresJS + Three.js)                   │
│  Dashboard ─ 3D Viewer ─ Teklif ─ Sipariş ─ Admin      │
└────────────────────┬────────────────────────────────────┘
                     │ HTTPS / Nginx
┌────────────────────▼────────────────────────────────────┐
│  Main API  (Go 1.26 + Fiber v2)                         │
│  22 Controller │ 32 Service │ 14 Middleware             │
│  JWT + API Key + Rate Limit + Circuit Breaker           │
└──────────┬──────────────┬──────────────┬────────────────┘
           │ AMQP         │ SQL          │ S3
┌──────────▼──────┐ ┌────▼────┐ ┌──────▼─────┐
│ Geometry Engine │ │ Postgres│ │   MinIO    │
│ (Rust)          │ │  15     │ │   S3       │
│ 36 CAM Modülü   │ └─────────┘ └────────────┘
│ 11 Math Modülü  │
│ K-Factor / DXF  │  ┌────────────┐  ┌──────────┐
│ Flat Pattern    │  │  RabbitMQ  │  │  Redis 7 │
│ Nesting / G-Code│  │  + DLQ     │  │  Cache   │
└─────────────────┘  └────────────┘  └──────────┘

ERP Adaptörleri: SAP ─ Logo Tiger ─ Mikro ─ Odoo ─ ERPNext
                 Netsis ─ Canias ─ ProManage ─ TrexDCAS
                 e-Fatura GİB ─ Dynamics ─ NetSuite
```

---

## ⚡ Hızlı Başlangıç

### Gereksinimler

| Araç | Versiyon | Zorunlu |
|------|----------|---------|
| Docker + Docker Compose | 24+ | ✅ |
| Go | 1.26+ | Geliştirme için |
| Rust (rustup) | 1.80+ | Geliştirme için |
| Node.js | 20 LTS | Frontend geliştirme |

### 1. Ortam Değişkenleri

```bash
cp .env.example .env
# .env dosyasını düzenle (DB parolası, JWT secret, MinIO credentialları)
```

### 2. Tüm Sistemi Başlat (Docker Compose)

```bash
docker compose up -d
```

**Servislerin hazır olması ~30 saniye sürer.**

### 3. Erişim Noktaları

| Servis | URL | Kimlik Bilgisi |
|--------|-----|----------------|
| **Frontend** | http://localhost:8080 | .env → USER credentials |
| **API** | http://localhost:8080/api/v1 | JWT Bearer Token |
| **Swagger** | http://localhost:8080/swagger | — |
| **Health Check** | http://localhost:8080/health | — |
| **RabbitMQ UI** | http://localhost:15672 | guest / guest |
| **MinIO Console** | http://localhost:9001 | .env → MINIO_ROOT_* |
| **Prometheus** | http://localhost:8080/metrics | — |

---

## 🔧 Servis Portları

```
8080  → Nginx (Frontend + API Gateway)
3000  → Main API (Go/Fiber) — internal
5432  → PostgreSQL
5672  → RabbitMQ AMQP
15672 → RabbitMQ Management UI
9000  → MinIO API
9001  → MinIO Console
6379  → Redis
```

---

## 🧪 Test Komutu

### Go Unit Testleri (Main API)

```bash
cd services/main-api
go test ./services/... -v -cover
```

### Rust Unit Testleri (Geometry Engine)

```bash
cd services/geometry-engine
cargo test -- --nocapture
```

### E2E Testleri (Playwright)

```bash
cd apps/frontend
npm run test:e2e
```

---

## 🔌 ERP Entegrasyon Rehberi

### SAP Bağlantısı

```env
# .env içinde:
SAP_HOST=sap.sirketiniz.com
SAP_CLIENT=100
SAP_USER=RFC_USER
SAP_PASSWORD=gizli_parola
SAP_SYSTEM_NUMBER=00
```

### Logo Tiger SQL Bağlantısı

```env
LOGO_SQL_SERVER=192.168.1.100
LOGO_DB=TIGERDB
LOGO_USER=sa
LOGO_PASSWORD=gizli_parola
LOGO_FIRM_NUMBER=1
```

### Odoo REST

```env
ODOO_URL=https://odoo.sirketiniz.com
ODOO_DB=production
ODOO_USERNAME=api_user
ODOO_API_KEY=xxxxx
```

---

## 🌐 API Kullanımı (B2B)

### Hızlı Teklif Oluşturma

```bash
# 1. Token al
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@demo.com","password":"demo123"}' | jq -r '.token')

# 2. Teklif iste
curl -X POST http://localhost:8080/api/quote \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "material": "steel",
    "thickness_mm": 2.0,
    "flat_area_mm2": 200000,
    "cutting_length_mm": 1500,
    "bend_count": 3,
    "quantity": 50,
    "bounding_box_x": 500,
    "bounding_box_y": 400,
    "white_label_margin": 0.20,
    "target_currency": "TRY"
  }'
```

### B2B API Key ile (Partner Entegrasyonu)

```bash
curl -X POST http://localhost:8080/api/external/b2b/quote \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{ ... }'
```

---

## 📁 Proje Yapısı

```
SHEET-METAL-AUTOMATION/
├── apps/
│   ├── frontend/          # Vue 3 + TresJS + Vite
│   └── cad-bridge/        # SolidWorks / Inventor / SolidEdge .NET Addinleri
├── services/
│   ├── main-api/          # Go 1.26 + Fiber v2 (Ana API Gateway)
│   │   ├── controllers/   # 22 HTTP Controller
│   │   ├── services/      # 32 Business Logic Servisi
│   │   │   └── erp/       # 14 ERP/MES Adaptörü
│   │   ├── middlewares/   # 14 Middleware (JWT, Rate Limit, Circuit Breaker...)
│   │   ├── models/        # GORM Veri Modelleri
│   │   └── routes/        # API Route Tanımları
│   └── geometry-engine/   # Rust CAM Motor (36 CAM + 11 Math modülü)
│       └── src/
│           ├── cam/       # K-Faktör, Nesting, DXF, G-Code, Bükme...
│           └── math/      # Flat Pattern, Kerf, Flange, Ekonomi Ölçeği...
├── k8s/                   # 37 Kubernetes Manifest (K3s production)
├── nginx/                 # Nginx Gateway Konfigürasyonu
└── docker-compose.yml     # Yerel geliştirme ortamı
```

---

## 🔒 Ortam Değişkenleri

| Değişken | Açıklama | Örnek |
|----------|----------|-------|
| `DATABASE_URL` | PostgreSQL bağlantı URI | `postgres://user:pass@host/db` |
| `JWT_SECRET` | JWT imzalama anahtarı | `gizli_256bit_anahtar` |
| `RABBITMQ_HOST` | RabbitMQ sunucu adresi | `localhost` |
| `MINIO_ROOT_USER` | MinIO admin kullanıcı | `minioadmin` |
| `MINIO_ROOT_PASSWORD` | MinIO admin parola | `gizliparola` |
| `REDIS_PASSWORD` | Redis auth parolası | `redisparola` |
| `STRIPE_SECRET_KEY` | Stripe API anahtarı | `sk_live_...` |
| `IYZICO_API_KEY` | İyzico API anahtarı | `sandbox-...` |
| `LME_API_KEY` | LME metal fiyat API | `...` |
| `OPENAI_API_KEY` | GPT-4 DFM Advisor | `sk-...` |

Tam liste için [`.env.example`](./.env.example) dosyasına bakın.

---

## 🤝 Katkı Rehberi

1. **Branch stratejisi:** `feature/faz-XX-açıklama` formatında dal oluştur
2. **Commit mesajları:** `Faz/Sprint Tamamlandi: Kısa açıklama` formatını kullan
3. **Test:** Yeni servis eklendiğinde ilgili `_test.go` dosyasını ekle
4. **IP Koruması:** Fiyatlandırma matematiği **sadece** `services/main-api/services/pricing.go` ve Rust engine'de kalacak — frontend'e taşıma

---

## 📄 Lisans

Bu yazılım özel mülkiyete aittir. Tüm hakları saklıdır.
© 2026 — Tüm haklar saklıdır.

