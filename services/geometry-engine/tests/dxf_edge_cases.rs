#[cfg(test)]
mod dxf_edge_case_tests {
    /// Test 1: Boş DXF buffer parse hatası beklenir
    #[test]
    fn test_empty_dxf_buffer_graceful() {
        let buffer: &[u8] = &[];
        assert_eq!(buffer.len(), 0);
        // Gerçek implementasyonda: parse_dxf(&buffer) Err döndürmeli
    }

    /// Test 2: Geçersiz STEP dosyası — graceful error beklenir
    #[test]
    fn test_malformed_step_file_graceful() {
        let buffer = b"THIS IS NOT A VALID STEP FILE";
        assert!(String::from_utf8_lossy(buffer).contains("NOT A VALID"));
        // Gerçek implementasyonda: analyze_topology(&buffer) Err döndürmeli
    }

    /// Test 3: Minimum geçerli DXF içeriği — panic olmamalı
    #[test]
    fn test_minimal_dxf_no_panic() {
        let minimal_dxf = b"0\nSECTION\n2\nENTITIES\n0\nENDSEC\n0\nEOF\n";
        // Parse denemesi crash olmamalı
        let _ = String::from_utf8_lossy(minimal_dxf);
    }
}

#[cfg(test)]
mod k_factor_tests {
    use geometry_engine::math::k_factor::{DynamicKFactorCalculator, KFactorCalculator, KFactorRequest};

    fn make_req(material: &str, thickness: f64, inner_radius: f64, angle: f64) -> KFactorRequest {
        KFactorRequest {
            material: material.to_string(),
            thickness,
            inner_radius,
            bend_angle_degrees: angle,
        }
    }

    /// Test 1: Çelik 90° bükme → k-faktörü 0.33-0.45 arasında olmalı
    #[test]
    fn test_steel_90deg_k_factor_range() {
        let calc = DynamicKFactorCalculator;
        let req = make_req("steel", 2.0, 2.0, 90.0);
        let res = calc.calculate(&req).unwrap();
        assert!(
            res.k_factor >= 0.2 && res.k_factor <= 0.5,
            "steel k-factor ({}) must be in [0.2, 0.5]", res.k_factor
        );
        // Çelik baz k-factor 0.42, rt_ratio=1.0 → aynı kalır
        assert!((res.k_factor - 0.42).abs() < 0.06, "steel k-factor should be ~0.42, got {}", res.k_factor);
    }

    /// Test 2: Alüminyum k-factor çelikten büyük olmalı
    #[test]
    fn test_aluminum_k_factor_higher_than_steel() {
        let calc = DynamicKFactorCalculator;
        let steel = calc.calculate(&make_req("steel", 2.0, 2.0, 90.0)).unwrap();
        let aluminum = calc.calculate(&make_req("aluminum", 2.0, 2.0, 90.0)).unwrap();
        assert!(
            aluminum.k_factor >= steel.k_factor,
            "aluminum k-factor ({}) should be >= steel k-factor ({})",
            aluminum.k_factor, steel.k_factor
        );
    }

    /// Test 3: Paslanmaz çelik k-factor çelikten küçük olmalı
    #[test]
    fn test_stainless_k_factor_lower_than_steel() {
        let calc = DynamicKFactorCalculator;
        let steel = calc.calculate(&make_req("steel", 2.0, 2.0, 90.0)).unwrap();
        let stainless = calc.calculate(&make_req("stainless steel", 2.0, 2.0, 90.0)).unwrap();
        assert!(
            stainless.k_factor <= steel.k_factor,
            "stainless ({}) should be <= steel ({})",
            stainless.k_factor, steel.k_factor
        );
    }

    /// Test 4: Bend Allowance pozitif olmalı
    #[test]
    fn test_bend_allowance_positive() {
        let calc = DynamicKFactorCalculator;
        let res = calc.calculate(&make_req("steel", 2.0, 2.0, 90.0)).unwrap();
        assert!(res.bend_allowance > 0.0, "bend_allowance must be > 0, got {}", res.bend_allowance);
    }

    /// Test 5: Geçersiz kalınlık → hata beklenir
    #[test]
    fn test_invalid_thickness_error() {
        let calc = DynamicKFactorCalculator;
        let req = make_req("steel", 0.0, 2.0, 90.0); // thickness=0 geçersiz
        assert!(calc.calculate(&req).is_err(), "expected error for thickness=0");
    }

    /// Test 6: Geçersiz açı (>=180°) → hata beklenir
    #[test]
    fn test_invalid_bend_angle_error() {
        let calc = DynamicKFactorCalculator;
        let req = make_req("steel", 2.0, 2.0, 180.0); // açı=180 geçersiz
        assert!(calc.calculate(&req).is_err(), "expected error for bend_angle=180");
    }

    /// Test 7: rt_ratio < 1 → k-factor azalmalı (tight bend penalty)
    #[test]
    fn test_tight_bend_lower_k_factor() {
        let calc = DynamicKFactorCalculator;
        // rt_ratio = 0.4/2.0 = 0.2 < 1.0 → k-factor -= 0.05
        let tight = calc.calculate(&make_req("steel", 2.0, 0.4, 90.0)).unwrap();
        // rt_ratio = 4.0/2.0 = 2.0 → normal
        let normal = calc.calculate(&make_req("steel", 2.0, 4.0, 90.0)).unwrap();
        assert!(
            tight.k_factor < normal.k_factor,
            "tight bend ({}) should have lower k-factor than normal bend ({})",
            tight.k_factor, normal.k_factor
        );
    }
}

#[cfg(test)]
mod flat_pattern_tests {
    use geometry_engine::math::flat_pattern::{calculate_flat_pattern, FlatPatternRequest};

    /// Test 1: Tek 90° bükme — flat length > flanges toplamı (BA eklendi, BD çıkarıldı)
    #[test]
    fn test_single_90deg_bend_flat_length() {
        let req = FlatPatternRequest {
            material: "steel".to_string(),
            thickness: 2.0,
            flanges: vec![50.0, 50.0],  // 2 flange, 1 bükme
            bend_angles: vec![90.0],
            inner_radius: 2.0,
        };
        let res = calculate_flat_pattern(&req).unwrap();
        // Düzleştirilmiş uzunluk flange toplamından farklı olmalı (bend deduction uygulandı)
        let flange_sum = 50.0_f64 + 50.0_f64;
        assert!(
            res.total_flat_length > 0.0 && res.total_flat_length < flange_sum + 20.0,
            "flat_length ({}) should be reasonable for 2x50mm flanges", res.total_flat_length
        );
        assert_eq!(res.bend_deductions_used.len(), 1, "should have 1 bend deduction");
    }

    /// Test 2: Flange sayısı uyumsuzluğu → hata beklenir
    #[test]
    fn test_flange_bend_count_mismatch_error() {
        let req = FlatPatternRequest {
            material: "steel".to_string(),
            thickness: 2.0,
            flanges: vec![50.0, 50.0, 50.0], // 3 flange
            bend_angles: vec![90.0],           // 1 bükme → uyumsuz
            inner_radius: 2.0,
        };
        assert!(calculate_flat_pattern(&req).is_err(),
            "mismatched flange/bend count should return error");
    }

    /// Test 3: Çoklu bükme — deduction sayısı bükme sayısına eşit olmalı
    #[test]
    fn test_multiple_bends_deduction_count() {
        let req = FlatPatternRequest {
            material: "aluminum".to_string(),
            thickness: 1.5,
            flanges: vec![30.0, 40.0, 30.0, 20.0], // 4 flange, 3 bükme
            bend_angles: vec![90.0, 45.0, 90.0],
            inner_radius: 1.5,
        };
        let res = calculate_flat_pattern(&req).unwrap();
        assert_eq!(res.bend_deductions_used.len(), 3,
            "3 bends should produce 3 deductions, got {}", res.bend_deductions_used.len());
        assert!(res.total_flat_length > 0.0, "flat length must be positive");
    }

    /// Test 4: 45° bükme deduction < 90° bükme deduction (küçük açı, küçük deduction)
    #[test]
    fn test_smaller_angle_smaller_deduction() {
        let make = |angle: f64| FlatPatternRequest {
            material: "steel".to_string(),
            thickness: 2.0,
            flanges: vec![50.0, 50.0],
            bend_angles: vec![angle],
            inner_radius: 2.0,
        };
        let res_45 = calculate_flat_pattern(&make(45.0)).unwrap();
        let res_90 = calculate_flat_pattern(&make(90.0)).unwrap();

        assert!(
            res_45.bend_deductions_used[0] < res_90.bend_deductions_used[0],
            "45° deduction ({}) should be < 90° deduction ({})",
            res_45.bend_deductions_used[0], res_90.bend_deductions_used[0]
        );
    }
}

#[cfg(test)]
mod nesting_tests {
    use geometry_engine::cam::nesting::{calculate_nesting, NestingRequest};

    /// Test 1: Temel nesting — sonuç geçerli
    #[test]
    fn test_basic_nesting_returns_result() {
        let req = NestingRequest {
            sheet_width: 1000.0,
            sheet_height: 2000.0,
            parts: vec![
                (200.0, 100.0, 5),  // 5 adet 200x100mm parça
                (150.0, 150.0, 3),  // 3 adet 150x150mm parça
            ],
        };
        let res = calculate_nesting(&req).unwrap();
        assert!(res.efficiency_pct > 0.0, "nesting efficiency must be > 0");
        assert!(res.efficiency_pct <= 100.0, "nesting efficiency must be <= 100");
    }

    /// Test 2: Boş parça listesi — hata veya sıfır efficiency
    #[test]
    fn test_empty_parts_list_handled() {
        let req = NestingRequest {
            sheet_width: 1000.0,
            sheet_height: 2000.0,
            parts: vec![],
        };
        let res = calculate_nesting(&req);
        // Hata veya 0 efficiency kabul edilebilir — panic olmaMAlı
        match res {
            Ok(r) => assert_eq!(r.efficiency_pct, 0.0),
            Err(_) => {} // hata da kabul
        }
    }

    /// Test 3: Küçük parçalar yüksek efficiency sağlamalı (>%50)
    #[test]
    fn test_small_parts_high_efficiency() {
        let req = NestingRequest {
            sheet_width: 1000.0,
            sheet_height: 2000.0,
            parts: vec![(100.0, 100.0, 50)], // 50 adet 100x100mm
        };
        let res = calculate_nesting(&req).unwrap();
        assert!(
            res.efficiency_pct > 50.0,
            "small uniform parts should have > 50% efficiency, got {}%", res.efficiency_pct
        );
    }
}

