-- =============================================================
-- SHEET-METAL-AUTOMAT-ON — Seed Data (Faz 12.2)
-- Production başlangıç verisi: Tenant, Makine, Malzeme, Admin
-- =============================================================
-- Kullanım: make seed
-- =============================================================

BEGIN;

-- 1. Demo Tenant (Fason Atölye)
INSERT INTO tenants (id, name, slug, plan, erp_type, is_active, created_at, updated_at)
VALUES 
    (1, 'Demo Atölye A.Ş.', 'demo-atolye', 'enterprise', 'logo', true, NOW(), NOW()),
    (2, 'Test Fabrikası Ltd.', 'test-fabrika', 'pro', 'sap', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. Admin Kullanıcı (her tenant için)
-- Şifre: admin123 (bcrypt hash - üretimde değiştirilmeli!)
INSERT INTO users (id, tenant_id, email, password_hash, role, full_name, is_active, created_at, updated_at)
VALUES 
    (1, 1, 'admin@demo-atolye.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/RK.s5uDSm', 'admin', 'Demo Admin', true, NOW(), NOW()),
    (2, 2, 'admin@test-fabrika.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/RK.s5uDSm', 'admin', 'Fabrika Admin', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. CNC Makineler
INSERT INTO machines (id, tenant_id, name, type, power_kw, max_thickness_mm, bed_size_x_mm, bed_size_y_mm, 
                      laser_speed_mm_per_min, pierce_time_seconds, setup_cost_usd, hourly_rate_usd, is_active, created_at, updated_at)
VALUES 
    -- Fiber Lazer
    (1, 1, 'Trumpf TruLaser 3030 (6kW)', 'fiber_laser', 6.0, 25.0, 3000.0, 1500.0, 
     15000.0, 1.2, 50.0, 85.0, true, NOW(), NOW()),
    -- CO2 Lazer
    (2, 1, 'Bystronic ByStar 4020 (4kW)', 'co2_laser', 4.0, 12.0, 4000.0, 2000.0, 
     8000.0, 2.0, 40.0, 65.0, true, NOW(), NOW()),
    -- Abkant Pres
    (3, 1, 'Amada HFE 100-3 (Press Brake)', 'press_brake', 18.5, 12.0, 3200.0, 0.0, 
     0.0, 0.0, 30.0, 45.0, true, NOW(), NOW()),
    -- Fiber Lazer (Tenant 2)
    (4, 2, 'Prima Power Platino Fiber 4kW', 'fiber_laser', 4.0, 20.0, 3000.0, 1500.0,
     12000.0, 1.5, 45.0, 75.0, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 4. Malzeme Listesi (ISO Standart)
INSERT INTO materials (id, tenant_id, name, code, density_kg_per_m3, 
                       price_per_kg_usd, min_thickness_mm, max_thickness_mm, is_active, created_at, updated_at)
VALUES 
    -- DKP (Soğuk Hadde Çelik)
    (1, 1, 'DKP Çelik (St37)', 'DKP-ST37', 7850.0, 1.15, 0.5, 12.0, true, NOW(), NOW()),
    -- Paslanmaz Çelik
    (2, 1, 'Paslanmaz 304 (2B)', 'SS-304-2B', 7930.0, 4.20, 0.5, 10.0, true, NOW(), NOW()),
    (3, 1, 'Paslanmaz 316L', 'SS-316L', 7980.0, 6.80, 0.5, 8.0, true, NOW(), NOW()),
    -- Alüminyum
    (4, 1, 'Alüminyum 5052 H32', 'AL-5052-H32', 2680.0, 3.45, 0.5, 8.0, true, NOW(), NOW()),
    (5, 1, 'Alüminyum 6061 T6', 'AL-6061-T6', 2700.0, 3.85, 0.8, 10.0, true, NOW(), NOW()),
    -- Galvanizli Çelik
    (6, 1, 'Galvanizli Çelik (Z275)', 'GAL-Z275', 7850.0, 1.35, 0.5, 6.0, true, NOW(), NOW()),
    -- Bakır
    (7, 1, 'Elektrolitik Bakır (Cu-ETP)', 'CU-ETP', 8960.0, 9.50, 0.3, 4.0, true, NOW(), NOW()),
    -- Tenant 2 malzemeleri
    (8, 2, 'DKP Çelik (St37)', 'DKP-ST37', 7850.0, 1.18, 0.5, 12.0, true, NOW(), NOW()),
    (9, 2, 'Paslanmaz 304', 'SS-304', 7930.0, 4.35, 0.5, 10.0, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 5. LME Başlangıç Fiyatları
INSERT INTO lme_price_logs (id, metal_code, price_usd_per_ton, source, recorded_at)
VALUES 
    (1, 'AL', 2485.50, 'seed_data', NOW()),
    (2, 'CU', 9250.00, 'seed_data', NOW()),
    (3, 'ZN', 2890.00, 'seed_data', NOW()),
    (4, 'NI', 16800.00, 'seed_data', NOW()),
    (5, 'SN', 26500.00, 'seed_data', NOW())
ON CONFLICT (id) DO NOTHING;

-- 6. API Anahtarları (White-Label entegrasyonu için)
INSERT INTO api_keys (id, tenant_id, name, key_hash, permissions, is_active, expires_at, created_at)
VALUES 
    (1, 1, 'Demo B2B Portal Key', 'demo_api_key_hash_v1', '["quote:read","quote:create","order:read"]', true, NOW() + INTERVAL '1 year', NOW()),
    (2, 2, 'Test Webhook Key', 'test_webhook_key_hash_v1', '["order:read","erp:sync"]', true, NOW() + INTERVAL '6 months', NOW())
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- Doğrulama
SELECT 'Tenants: ' || COUNT(*) FROM tenants;
SELECT 'Users: ' || COUNT(*) FROM users;
SELECT 'Machines: ' || COUNT(*) FROM machines;
SELECT 'Materials: ' || COUNT(*) FROM materials;
SELECT 'LME Prices: ' || COUNT(*) FROM lme_price_logs;
SELECT 'API Keys: ' || COUNT(*) FROM api_keys;
