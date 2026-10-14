<template>
  <div class="landing">
    <!-- Promo Bar (Xometry'de üst bildirim çubuğu) -->
    <div class="promo-bar">
      <span>🔬 Türkiye'nin Yapay Zeka Destekli Anında Teklif Motorunu Keşfedin — </span>
      <a href="/quote">Hemen Deneyin &rarr;</a>
    </div>

    <!-- Navigation (Xometry tarzı beyaz navbar) -->
    <nav class="site-nav">
      <div class="nav-inner">
        <router-link to="/" class="brand-logo">
          <Hexagon size="28" class="brand-icon" />
          <span>Apex<strong>FactoryOS</strong></span>
        </router-link>

        <div class="nav-menu">
          <div class="nav-dropdown">
            <button class="nav-btn">Yetenekler <ChevronDown size="16" /></button>
            <div class="dropdown-panel">
              <a href="#services">CNC Freze & Torna</a>
              <a href="#services">Lazer & Sac Metal</a>
              <a href="#services">Abkant Büküm</a>
              <a href="#services">Yüzey İşlemleri</a>
            </div>
          </div>
          <div class="nav-dropdown">
            <button class="nav-btn">Malzemeler <ChevronDown size="16" /></button>
            <div class="dropdown-panel">
              <a href="#materials">Alüminyum</a>
              <a href="#materials">Paslanmaz Çelik</a>
              <a href="#materials">Karbon Çelik</a>
              <a href="#materials">Titanyum</a>
            </div>
          </div>
          <a href="#how-it-works" class="nav-link">Nasıl Çalışır?</a>
          <a href="#faq" class="nav-link">SSS</a>
          <div class="nav-dropdown">
            <button class="nav-btn">Şirket <ChevronDown size="16" /></button>
            <div class="dropdown-panel">
              <a href="#">Hakkımızda</a>
              <a href="#">Medya</a>
              <a href="#">Kariyer</a>
              <a href="#">İletişim</a>
            </div>
          </div>
        </div>

        <div class="nav-actions">
          <div class="lang-toggle">
            <span class="active">TR</span><span class="sep">/</span><span>EN</span>
          </div>
          <router-link to="/login" class="btn-ghost">Giriş Yap</router-link>
          <router-link to="/quote" class="btn-primary-blue">Anında Fiyat Al</router-link>
        </div>
      </div>
    </nav>

    <main>
      <!-- HERO (Xometry'nin ikiye bölünmüş hero alanı) -->
      <section class="hero">
        <div class="hero-left">
          <div class="hero-tag">AI-Powered Manufacturing Platform</div>
          <h1>Büyük Fikirlerin<br /><span class="text-blue">Üretildiği Yer.</span></h1>
          <p>
            CNC İşleme, Sac Metal, Lazer Kesim ve Abkant Büküm hizmetlerini tek bir platformdan alın. STEP veya DXF dosyanızı yükleyin, saniyeler içinde anında teklif alın, üretime gönderin.
          </p>
          <div class="hero-cta">
            <router-link to="/quote" class="btn-cta-large">
              <UploadCloud size="22" />  3D Dosya Yükle & Fiyat Al
            </router-link>
            <a href="#how-it-works" class="btn-link-arrow">Nasıl çalışır? <ArrowRight size="18" /></a>
          </div>
          <div class="hero-stats">
            <div class="stat"><strong>5.000+</strong><span>Üretilen Parça</span></div>
            <div class="stat-divider"></div>
            <div class="stat"><strong>±0.01mm</strong><span>Tolerans</span></div>
            <div class="stat-divider"></div>
            <div class="stat"><strong>24 Saat</strong><span>İlk Sevkiyat</span></div>
          </div>
        </div>
        <div class="hero-right">
          <div class="hero-visual">
            <div class="visual-card card-1">
              <div class="card-icon"><Cpu size="28" /></div>
              <div class="card-body">
                <div class="card-label">AI Analizi</div>
                <div class="card-value text-green">Tamamlandı ✓</div>
              </div>
            </div>
            <div class="visual-card card-2">
              <div class="card-icon"><CircleDollarSign size="28" /></div>
              <div class="card-body">
                <div class="card-label">Birim Fiyat</div>
                <div class="card-value">₺487.50</div>
              </div>
            </div>
            <div class="visual-card card-3">
              <div class="card-icon"><Clock size="28" /></div>
              <div class="card-body">
                <div class="card-label">Üretim Süresi</div>
                <div class="card-value">3-5 İş Günü</div>
              </div>
            </div>
            <div class="hero-3d-bg">
              <div class="rotating-hex">
                <Hexagon size="300" class="hex-outline" />
                <Hexagon size="220" class="hex-outline-2" />
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- TRUSTED BY (Güven bandı) -->
      <section class="trusted-section">
        <div class="container">
          <p class="trusted-label">500+ Türk ve Uluslararası Şirket Tarafından Tercih Ediliyor</p>
          <div class="logo-row">
            <span class="trust-logo">LOGO Yazılım</span>
            <span class="trust-logo">Netsis ERP</span>
            <span class="trust-logo">SAP S/4HANA</span>
            <span class="trust-logo">SolidWorks</span>
            <span class="trust-logo">Canias ERP</span>
            <span class="trust-logo">Mikro ERP</span>
          </div>
        </div>
      </section>

      <!-- CAPABILITIES (Xometry'nin hizmet kartları) -->
      <section id="services" class="capabilities-section">
        <div class="container">
          <div class="section-eyebrow">ÜRETİM YETENEKLERİMİZ</div>
          <h2>Tek Platform. Sonsuz Kapasite.</h2>
          <p class="section-sub">Prototipten seri üretime, ihtiyacınız olan her türlü hassas işleme hizmetini sunuyoruz.</p>

          <div class="caps-grid">
            <div class="cap-card">
              <div class="cap-image cnc-bg">
                <div class="cap-overlay">
                  <span class="tolerance-badge">±0.01mm</span>
                </div>
              </div>
              <div class="cap-body">
                <h3>CNC Talaşlı İmalat</h3>
                <p>3, 4 ve 5 Eksen CNC Freze & Torna. Alüminyum'dan Titanyum'a geniş malzeme seçeneği.</p>
                <a href="#" class="cap-link">Detaylı İncele <ArrowRight size="14" /></a>
              </div>
            </div>
            <div class="cap-card">
              <div class="cap-image sheet-bg">
                <div class="cap-overlay">
                  <span class="tolerance-badge">0.5-20mm</span>
                </div>
              </div>
              <div class="cap-body">
                <h3>Sac Metal & Lazer Kesim</h3>
                <p>Fiber Lazer Kesim, CNC Abkant Büküm ve Punch işlemleri. Tam otomatik büküm analizi.</p>
                <a href="#" class="cap-link">Detaylı İncele <ArrowRight size="14" /></a>
              </div>
            </div>
            <div class="cap-card">
              <div class="cap-image weld-bg">
                <div class="cap-overlay">
                  <span class="tolerance-badge">ISO 9001</span>
                </div>
              </div>
              <div class="cap-body">
                <h3>Kaynak & Birleştirme</h3>
                <p>TIG, MIG ve Robotik Kaynak. Kompleks montaj ve yapısal parçalar için kapsamlı hizmet.</p>
                <a href="#" class="cap-link">Detaylı İncele <ArrowRight size="14" /></a>
              </div>
            </div>
            <div class="cap-card">
              <div class="cap-image quality-bg">
                <div class="cap-overlay">
                  <span class="tolerance-badge">CMM Raporu</span>
                </div>
              </div>
              <div class="cap-body">
                <h3>Kalite & Sertifikasyon</h3>
                <p>CMM Boyutsal Denetim, EN 10204 Malzeme Sertifikası ve tam kalite traceability.</p>
                <a href="#" class="cap-link">Detaylı İncele <ArrowRight size="14" /></a>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- HOW IT WORKS (Xometry'nin süreç adımları) -->
      <section id="how-it-works" class="how-section">
        <div class="container">
          <div class="section-eyebrow">SÜREÇ</div>
          <h2>Üretim Saniyeler İçinde Başlar</h2>
          <div class="steps">
            <div class="step">
              <div class="step-num">1</div>
              <div class="step-content">
                <h3><UploadCloud size="20" class="inline-icon" /> Yükle</h3>
                <p>STEP, STP veya DXF dosyanızı platforma yükleyin. Uçtan uca şifreleme ile güvendedir.</p>
              </div>
            </div>
            <div class="step-arrow"><ArrowRight size="28" /></div>
            <div class="step">
              <div class="step-num">2</div>
              <div class="step-content">
                <h3><Cpu size="20" class="inline-icon" /> Fiyat Al</h3>
                <p>Rust tabanlı AI motorumuz parçayı analiz eder, lazer gaz maliyetinden nesting veriminize kadar hesaplar.</p>
              </div>
            </div>
            <div class="step-arrow"><ArrowRight size="28" /></div>
            <div class="step">
              <div class="step-num">3</div>
              <div class="step-content">
                <h3><Zap size="20" class="inline-icon" /> Üret</h3>
                <p>Siparişi onaylayın. MES sistemimiz parçayı anında üretim parkurundaki en uygun makineye atar.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- WHY US (Xometry'nin "Why Xometry?" bölümü) -->
      <section class="why-section">
        <div class="container">
          <div class="why-grid">
            <div class="why-content">
              <div class="section-eyebrow">NEDEN BİZ?</div>
              <h2>Geleneksel Tedarik Zincirinin Sonu.</h2>
              <p>Haftalar süren e-posta trafiği, belirsiz fiyatlar ve geç teslimatlarla dolu tedarikçi yönetimi devri kapandı. Apex FactoryOS, üretim sürecinizi tamamen dijitalleştiriyor.</p>
              <router-link to="/quote" class="btn-primary-blue mt-6">Platformu Keşfedin</router-link>
            </div>
            <div class="why-pillars">
              <div class="pillar">
                <div class="pillar-icon"><Zap size="28" /></div>
                <h4>Anında Fiyat, Sıfır Bekleme</h4>
                <p>Geleneksel tedarik süreçlerinde 3-4 gün olan teklif alma süresi, AI motorumuzla saniyeye iniyor.</p>
              </div>
              <div class="pillar">
                <div class="pillar-icon"><ShieldCheck size="28" /></div>
                <h4>ISO 9001 Kalite Güvencesi</h4>
                <p>Her parça CMM ölçüm cihazından geçer. Malzeme sertifikası (EN 10204 3.1) standart olarak sunulur.</p>
              </div>
              <div class="pillar">
                <div class="pillar-icon"><Globe size="28" /></div>
                <h4>ERP Entegrasyonu</h4>
                <p>Logo, Netsis, SAP ve SolidWorks ile doğrudan bağlantı. Mevcut iş akışlarınızı bozmadan entegre olun.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- PRICING (Fiyatlandırma şeffaflığı) -->
      <section id="pricing" class="pricing-section">
        <div class="container">
          <div class="section-eyebrow">FİYATLANDIRMA</div>
          <h2>Şeffaf, Sürprizsiz Maliyet.</h2>
          <p class="section-sub">Sabit abonelik ücreti yok. Gizli komisyon yok. Sadece ürettiğiniz kadar ödeyin.</p>
          <div class="pricing-breakdown">
            <div class="breakdown-left">
              <h3>AI Maliyet Algoritması Nasıl Çalışır?</h3>
              <div class="cost-line">
                <div class="cost-icon"><Layers size="18" /></div>
                <div>
                  <strong>Malzeme Maliyeti</strong>
                  <p>Piyasa fiyatlarına anlık bağlı, fire oranı optimize edilmiş hammadde bedeli.</p>
                </div>
              </div>
              <div class="cost-line">
                <div class="cost-icon"><Settings size="18" /></div>
                <div>
                  <strong>İşleme Süresi (Makine Dakikası)</strong>
                  <p>Takım yolu simülasyonu ile CNC veya lazer kesim süresinin milimetrik hesabı.</p>
                </div>
              </div>
              <div class="cost-line">
                <div class="cost-icon"><Activity size="18" /></div>
                <div>
                  <strong>Lazer Gaz & Enerji Maliyeti</strong>
                  <p>O2, N2 ve Basınçlı Hava gazının kesim hızı ve kalınlığa göre hesabı.</p>
                </div>
              </div>
            </div>
            <div class="breakdown-right">
              <div class="quote-preview">
                <div class="qp-header">
                  <span>Örnek Teklif Analizi</span>
                  <span class="qp-badge">Canlı AI</span>
                </div>
                <div class="qp-part">Parça: Flanş (AlSi 6061-T6)</div>
                <div class="qp-lines">
                  <div class="qp-line"><span>Malzeme (0.8 kg)</span><span>₺224.00</span></div>
                  <div class="qp-line"><span>CNC İşleme (38 dk)</span><span>₺190.00</span></div>
                  <div class="qp-line"><span>Kurulum & Fikstür</span><span>₺35.00</span></div>
                  <div class="qp-line"><span>Kalite Kontrolü</span><span>₺38.50</span></div>
                </div>
                <div class="qp-total"><span>Birim Fiyat</span><span>₺487.50</span></div>
                <router-link to="/quote" class="qp-btn">Parçanızı Fiyatlandırın</router-link>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- MATERIALS (Malzeme kütüphanesi) -->
      <section id="materials" class="materials-section">
        <div class="container">
          <div class="section-eyebrow">MALZEME KÜTÜPHANESİ</div>
          <h2>50+ Metal ve Alaşım</h2>
          <div class="mat-grid">
            <div class="mat-col">
              <h4><span class="mat-dot al"></span>Alüminyum</h4>
              <ul><li>Al 6061-T6</li><li>Al 7075-T651</li><li>Al 5083</li><li>Al 2024</li></ul>
            </div>
            <div class="mat-col">
              <h4><span class="mat-dot ss"></span>Paslanmaz Çelik</h4>
              <ul><li>AISI 304 / 304L</li><li>AISI 316 / 316L</li><li>AISI 430</li><li>Duplex 2205</li></ul>
            </div>
            <div class="mat-col">
              <h4><span class="mat-dot cs"></span>Karbon & Yapı Çeliği</h4>
              <ul><li>ST37 / ST52</li><li>S235 / S355</li><li>Hardox 400/500</li><li>Ck45 Sertleştirilebilir</li></ul>
            </div>
            <div class="mat-col">
              <h4><span class="mat-dot sp"></span>Özel Alaşımlar</h4>
              <ul><li>Titanyum Grade 2/5</li><li>Pirinç CuZn37</li><li>Bakır C11000</li><li>Inconel 625</li></ul>
            </div>
          </div>
        </div>
      </section>

      <!-- FAQ -->
      <section id="faq" class="faq-section">
        <div class="container">
          <div class="section-eyebrow">SIKÇA SORULAN SORULAR</div>
          <h2>Aklınızdaki Sorular</h2>
          <div class="faq-cols">
            <div class="faq-col">
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>Hangi CAD formatlarını destekliyorsunuz?</h4>
                  <p>Yapay zeka motorumuz .STEP (.STP) ve sac açınımları için .DXF formatlarını tam destekler. Dönüştürme aracımızla diğer formatları da işleyebilirsiniz.</p>
                </div>
              </div>
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>Tolerans standartlarınız nedir?</h4>
                  <p>CNC işlemelerde standart olarak ISO 2768-m (±0.1mm) uygulanır. Kritik yüzeyler için ±0.01mm hassas tolerans belirtilebilir.</p>
                </div>
              </div>
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>Minimum sipariş miktarı (MOQ) var mı?</h4>
                  <p>Hayır. Tek adet prototipten 50.000 adetlik seri üretime kadar her miktarda sipariş kabul ediyoruz.</p>
                </div>
              </div>
            </div>
            <div class="faq-col">
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>Fikri mülkiyetim (IP) güvende mi?</h4>
                  <p>Tüm CAD ve teknik dosyalarınız MinIO Secure Vault altyapısında uçtan uca şifrelenerek depolanır. 6698 sayılı KVKK ve ISO 27001 kapsamında 3. taraflarla kesinlikle paylaşılmaz.</p>
                </div>
              </div>
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>Kalite belgesi (sertifika) alabilir miyim?</h4>
                  <p>Evet. Her sipariş için isteğe bağlı olarak EN 10204 3.1 Malzeme Sertifikası ve CMM Boyutsal Ölçüm Raporu sunulmaktadır.</p>
                </div>
              </div>
              <div class="faq-q">
                <ChevronRight size="18" class="faq-icon" />
                <div>
                  <h4>ERP sistemimle entegre olabilir miyim?</h4>
                  <p>Evet. Logo, Netsis, SAP ve diğer popüler ERP sistemleriyle API üzerinden doğrudan entegrasyon sağlıyoruz. Webhook desteği mevcuttur.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- CTA BAND -->
      <section class="cta-band">
        <div class="container">
          <h2>Üretim sürecinizi bugün dönüştürün.</h2>
          <p>CAD dosyanızı yükleyin, yapay zeka motorumuz saniyeler içinde eksiksiz bir teklif hazırlasın.</p>
          <router-link to="/quote" class="btn-cta-large white">
            <UploadCloud size="22" /> Ücretsiz Teklif Al
          </router-link>
        </div>
      </section>

      <!-- FOOTER -->
      <footer class="site-footer">
        <div class="container">
          <div class="footer-top">
            <div class="footer-brand">
              <div class="footer-logo">
                <Hexagon size="28" class="brand-icon" />
                <strong>Apex FactoryOS</strong>
              </div>
              <p>Türkiye'nin yapay zeka destekli B2B bulut üretim ağı. Prototipten seri üretime, tek platform.</p>
              <div class="footer-certs">
                <span>ISO 9001:2015</span>
                <span>ISO 27001</span>
                <span>KVKK Uyumlu</span>
              </div>
            </div>
            <div class="footer-col">
              <h5>Üretim Hizmetleri</h5>
              <a href="#">CNC Freze</a>
              <a href="#">CNC Torna</a>
              <a href="#">Fiber Lazer Kesim</a>
              <a href="#">Abkant Büküm</a>
              <a href="#">Robotik Kaynak</a>
              <a href="#">Yüzey Kaplama</a>
            </div>
            <div class="footer-col">
              <h5>Malzemeler</h5>
              <a href="#">Alüminyum Alaşımları</a>
              <a href="#">Paslanmaz Çelik</a>
              <a href="#">Karbon Çelik</a>
              <a href="#">Titanyum</a>
              <a href="#">Pirinç & Bakır</a>
              <a href="#">Tam Malzeme Listesi</a>
            </div>
            <div class="footer-col">
              <h5>Platform</h5>
              <a href="#">Nasıl Çalışır?</a>
              <a href="#">API Dökümantasyonu</a>
              <a href="#">ERP Entegrasyonları</a>
              <a href="#">Tasarım Rehberi (DFM)</a>
              <a href="#">Müşteri Portalı</a>
            </div>
            <div class="footer-col">
              <h5>Şirket</h5>
              <a href="#">Hakkımızda</a>
              <a href="#">Kariyer</a>
              <a href="#">Basın</a>
              <a href="#">İletişim</a>
              <a href="#">Bayilik</a>
            </div>
          </div>
          <div class="footer-bottom">
            <p>&copy; 2026 Apex Factory OS A.Ş. Tüm hakları saklıdır.</p>
            <div class="legal-links">
              <a href="#"><strong>KVKK Aydınlatma Metni</strong></a>
              <a href="#">Gizlilik Politikası</a>
              <a href="#">Çerez Politikası</a>
              <a href="#">Kullanım Koşulları</a>
              <a href="#">Mesafeli Satış Sözleşmesi</a>
              <a href="#">İptal & İade</a>
            </div>
          </div>
        </div>
      </footer>
    </main>
  </div>
</template>

<script setup>
import {
  Hexagon, ChevronDown, ChevronRight, ArrowRight, UploadCloud,
  Cpu, CircleDollarSign, Clock, Zap, ShieldCheck, Globe,
  Layers, Settings, Activity
} from 'lucide-vue-next'
</script>

<style scoped>
/* ============================================================
   GLOBAL RESET & VARIABLES
   ============================================================ */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.landing {
  font-family: -apple-system, BlinkMacSystemFont, 'Inter', 'Segoe UI', sans-serif;
  color: #0f172a;
  background: #fff;
  font-size: 16px;
  line-height: 1.6;
}

.container { max-width: 1280px; margin: 0 auto; padding: 0 2rem; }
.section-eyebrow {
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 2px;
  color: #1d6bf3;
  text-transform: uppercase;
  margin-bottom: 0.75rem;
}
.section-sub { font-size: 1.15rem; color: #475569; max-width: 600px; margin: 0.75rem auto 0; }
.text-blue { color: #1d6bf3; }
.text-green { color: #10b981; }
.mt-6 { margin-top: 1.5rem; }
.inline-icon { display: inline-block; vertical-align: middle; margin-right: 0.4rem; }

/* ============================================================
   PROMO BAR
   ============================================================ */
.promo-bar {
  background: #1d6bf3;
  color: #fff;
  text-align: center;
  font-size: 0.875rem;
  padding: 0.6rem 1rem;
}
.promo-bar a { color: #fff; font-weight: 700; margin-left: 0.25rem; }

/* ============================================================
   NAVIGATION
   ============================================================ */
.site-nav {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
}

.nav-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 2rem;
  height: 68px;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
  color: #0f172a;
  font-size: 1.2rem;
  white-space: nowrap;
}
.brand-logo strong { color: #1d6bf3; font-weight: 900; }
.brand-icon { color: #1d6bf3; }

.nav-menu {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  flex: 1;
}

.nav-dropdown { position: relative; }
.nav-btn {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.925rem;
  font-weight: 500;
  color: #334155;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  transition: all 0.15s;
}
.nav-btn:hover { background: #f1f5f9; color: #1d6bf3; }

.nav-link {
  font-size: 0.925rem;
  font-weight: 500;
  color: #334155;
  text-decoration: none;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  transition: all 0.15s;
}
.nav-link:hover { background: #f1f5f9; color: #1d6bf3; }

.dropdown-panel {
  display: none;
  position: absolute;
  top: calc(100% + 0.5rem);
  left: 0;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
  min-width: 200px;
  padding: 0.5rem;
  z-index: 200;
}
.nav-dropdown:hover .dropdown-panel { display: block; }
.dropdown-panel a {
  display: block;
  padding: 0.6rem 1rem;
  color: #334155;
  text-decoration: none;
  font-size: 0.9rem;
  border-radius: 6px;
  transition: background 0.15s;
}
.dropdown-panel a:hover { background: #f1f5f9; color: #1d6bf3; }

.nav-actions { display: flex; align-items: center; gap: 1rem; margin-left: auto; }
.lang-toggle { display: flex; gap: 0.3rem; font-size: 0.85rem; font-weight: 600; color: #94a3b8; cursor: pointer; }
.lang-toggle .active { color: #0f172a; }
.lang-toggle .sep { color: #cbd5e1; }

.btn-ghost {
  text-decoration: none;
  font-size: 0.925rem;
  font-weight: 600;
  color: #334155;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  transition: all 0.15s;
}
.btn-ghost:hover { background: #f1f5f9; }

.btn-primary-blue {
  text-decoration: none;
  font-size: 0.925rem;
  font-weight: 700;
  color: #fff;
  background: #1d6bf3;
  padding: 0.6rem 1.5rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
}
.btn-primary-blue:hover { background: #1557d0; box-shadow: 0 4px 12px rgba(29,107,243,0.3); }

/* ============================================================
   HERO
   ============================================================ */
.hero {
  display: grid;
  grid-template-columns: 1fr 1fr;
  min-height: 80vh;
  align-items: center;
  gap: 4rem;
  padding: 5rem 2rem 5rem;
  max-width: 1280px;
  margin: 0 auto;
}

.hero-tag {
  display: inline-block;
  background: #eff6ff;
  color: #1d6bf3;
  border: 1px solid #bfdbfe;
  padding: 0.4rem 1rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  margin-bottom: 1.5rem;
}

.hero-left h1 {
  font-size: 3.8rem;
  font-weight: 900;
  line-height: 1.05;
  letter-spacing: -2px;
  color: #0f172a;
  margin-bottom: 1.5rem;
}

.hero-left p {
  font-size: 1.15rem;
  color: #475569;
  line-height: 1.7;
  margin-bottom: 2.5rem;
  max-width: 520px;
}

.hero-cta { display: flex; align-items: center; gap: 1.5rem; margin-bottom: 3rem; }

.btn-cta-large {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  background: #1d6bf3;
  color: #fff;
  font-size: 1rem;
  font-weight: 700;
  padding: 0.9rem 2rem;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.2s;
  box-shadow: 0 4px 15px rgba(29,107,243,0.35);
}
.btn-cta-large:hover { background: #1557d0; transform: translateY(-2px); box-shadow: 0 6px 20px rgba(29,107,243,0.4); }
.btn-cta-large.white { background: #fff; color: #1d6bf3; box-shadow: none; border: 2px solid rgba(255,255,255,0.5); }
.btn-cta-large.white:hover { background: #f0f7ff; }

.btn-link-arrow {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  color: #1d6bf3;
  font-weight: 600;
  font-size: 0.95rem;
  text-decoration: none;
  transition: gap 0.2s;
}
.btn-link-arrow:hover { gap: 0.7rem; }

.hero-stats {
  display: flex;
  align-items: center;
  gap: 2rem;
}
.stat strong { display: block; font-size: 1.6rem; font-weight: 900; color: #0f172a; }
.stat span { font-size: 0.85rem; color: #64748b; font-weight: 500; }
.stat-divider { width: 1px; height: 40px; background: #e2e8f0; }

/* Hero visual */
.hero-right { position: relative; height: 500px; display: flex; align-items: center; justify-content: center; }
.hero-visual { position: relative; width: 100%; height: 100%; }

.visual-card {
  position: absolute;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 200px;
  z-index: 3;
  animation: float 3s ease-in-out infinite;
}
.card-1 { top: 15%; left: 5%; animation-delay: 0s; }
.card-2 { top: 45%; right: 5%; animation-delay: 1s; }
.card-3 { bottom: 15%; left: 20%; animation-delay: 0.5s; }
.card-icon { color: #1d6bf3; }
.card-label { font-size: 0.8rem; color: #64748b; font-weight: 600; }
.card-value { font-size: 1.1rem; font-weight: 800; color: #0f172a; }

.hero-3d-bg {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1;
  opacity: 0.08;
}
.rotating-hex { position: relative; animation: rotate 20s linear infinite; }
.hex-outline { color: #1d6bf3; }
.hex-outline-2 { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: #1d6bf3; }

@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-8px); } }
@keyframes rotate { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

/* ============================================================
   TRUSTED BY
   ============================================================ */
.trusted-section {
  padding: 2.5rem 0;
  border-top: 1px solid #f1f5f9;
  border-bottom: 1px solid #f1f5f9;
  background: #f8fafc;
}
.trusted-label {
  text-align: center;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #94a3b8;
  margin-bottom: 1.5rem;
}
.logo-row { display: flex; justify-content: center; align-items: center; gap: 3rem; flex-wrap: wrap; }
.trust-logo {
  font-size: 1.1rem;
  font-weight: 800;
  color: #cbd5e1;
  letter-spacing: -0.5px;
  transition: color 0.2s;
}
.trust-logo:hover { color: #94a3b8; }

/* ============================================================
   CAPABILITIES
   ============================================================ */
.capabilities-section { padding: 6rem 0; }
.capabilities-section h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; }
.capabilities-section > .container > .section-sub { text-align: center; }
.capabilities-section > .container > h2 { text-align: center; }

.caps-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1.5rem; margin-top: 3rem; }
.cap-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.25s, box-shadow 0.25s;
}
.cap-card:hover { transform: translateY(-6px); box-shadow: 0 12px 30px rgba(0,0,0,0.1); border-color: #bfdbfe; }

.cap-image {
  height: 160px;
  position: relative;
  background-size: cover;
  background-position: center;
}
.cnc-bg { background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%); }
.sheet-bg { background: linear-gradient(135deg, #1e293b 0%, #0c4a6e 100%); }
.weld-bg { background: linear-gradient(135deg, #1c1917 0%, #451a03 100%); }
.quality-bg { background: linear-gradient(135deg, #0f172a 0%, #042f2e 100%); }

.cap-overlay {
  position: absolute;
  bottom: 1rem;
  left: 1rem;
}
.tolerance-badge {
  background: rgba(29,107,243,0.9);
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  letter-spacing: 0.5px;
}

.cap-body { padding: 1.5rem; }
.cap-body h3 { font-size: 1.1rem; font-weight: 800; margin-bottom: 0.5rem; color: #0f172a; }
.cap-body p { font-size: 0.9rem; color: #64748b; line-height: 1.5; margin-bottom: 1rem; }
.cap-link {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.9rem;
  font-weight: 700;
  color: #1d6bf3;
  text-decoration: none;
  transition: gap 0.2s;
}
.cap-link:hover { gap: 0.7rem; }

/* ============================================================
   HOW IT WORKS
   ============================================================ */
.how-section { padding: 6rem 0; background: #f8fafc; }
.how-section h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; text-align: center; margin-bottom: 4rem; }
.steps { display: flex; align-items: center; gap: 1.5rem; }
.step { flex: 1; background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 2.5rem; }
.step-num { font-size: 3rem; font-weight: 900; color: #eff6ff; line-height: 1; margin-bottom: 1rem; font-variant-numeric: tabular-nums; }
.step-content h3 { font-size: 1.2rem; font-weight: 800; margin-bottom: 0.75rem; color: #0f172a; display: flex; align-items: center; gap: 0.5rem; }
.step-content p { font-size: 0.95rem; color: #64748b; line-height: 1.6; }
.step-arrow { color: #cbd5e1; flex-shrink: 0; }

/* ============================================================
   WHY SECTION
   ============================================================ */
.why-section { padding: 6rem 0; }
.why-grid { display: grid; grid-template-columns: 1fr 1.2fr; gap: 6rem; align-items: center; }
.why-content h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; margin-bottom: 1rem; }
.why-content p { font-size: 1.05rem; color: #475569; line-height: 1.7; }
.why-pillars { display: flex; flex-direction: column; gap: 2rem; }
.pillar { display: flex; gap: 1.25rem; }
.pillar-icon { flex-shrink: 0; width: 52px; height: 52px; background: #eff6ff; border-radius: 10px; display: flex; align-items: center; justify-content: center; color: #1d6bf3; }
.pillar h4 { font-size: 1rem; font-weight: 800; margin-bottom: 0.4rem; color: #0f172a; }
.pillar p { font-size: 0.9rem; color: #64748b; line-height: 1.6; }

/* ============================================================
   PRICING
   ============================================================ */
.pricing-section { padding: 6rem 0; background: #f8fafc; }
.pricing-section h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; text-align: center; }
.pricing-section .section-sub { text-align: center; }

.pricing-breakdown { display: grid; grid-template-columns: 1fr 1fr; gap: 4rem; margin-top: 3rem; align-items: start; }
.breakdown-left h3 { font-size: 1.3rem; font-weight: 800; margin-bottom: 2rem; color: #0f172a; }
.cost-line { display: flex; gap: 1.2rem; margin-bottom: 2rem; align-items: flex-start; }
.cost-icon { flex-shrink: 0; width: 40px; height: 40px; background: #eff6ff; border-radius: 8px; display: flex; align-items: center; justify-content: center; color: #1d6bf3; margin-top: 2px; }
.cost-line strong { font-size: 1rem; color: #0f172a; display: block; margin-bottom: 0.25rem; }
.cost-line p { font-size: 0.9rem; color: #64748b; }

.quote-preview {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 20px rgba(0,0,0,0.06);
}
.qp-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; font-weight: 700; color: #0f172a; }
.qp-badge { background: #d1fae5; color: #059669; font-size: 0.75rem; font-weight: 700; padding: 0.2rem 0.75rem; border-radius: 20px; }
.qp-part { font-size: 0.9rem; color: #64748b; margin-bottom: 1.5rem; padding-bottom: 1rem; border-bottom: 1px solid #f1f5f9; }
.qp-lines { margin-bottom: 1rem; }
.qp-line { display: flex; justify-content: space-between; font-size: 0.95rem; color: #475569; padding: 0.5rem 0; border-bottom: 1px solid #f8fafc; }
.qp-line span:last-child { font-family: 'JetBrains Mono', monospace; font-weight: 600; }
.qp-total { display: flex; justify-content: space-between; font-weight: 900; font-size: 1.2rem; color: #0f172a; padding: 1rem 0; border-top: 2px solid #e2e8f0; margin-bottom: 1.5rem; }
.qp-total span:last-child { color: #1d6bf3; font-family: 'JetBrains Mono', monospace; }
.qp-btn {
  display: block;
  text-align: center;
  background: #1d6bf3;
  color: #fff;
  padding: 0.9rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 700;
  transition: all 0.2s;
}
.qp-btn:hover { background: #1557d0; }

/* ============================================================
   MATERIALS
   ============================================================ */
.materials-section { padding: 6rem 0; border-top: 1px solid #f1f5f9; }
.materials-section h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; text-align: center; margin-bottom: 3rem; }

.mat-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 3rem; }
.mat-col h4 { display: flex; align-items: center; gap: 0.5rem; font-size: 1rem; font-weight: 800; color: #0f172a; margin-bottom: 1rem; }
.mat-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.mat-dot.al { background: #3b82f6; }
.mat-dot.ss { background: #94a3b8; }
.mat-dot.cs { background: #64748b; }
.mat-dot.sp { background: #f59e0b; }
.mat-col ul { list-style: none; }
.mat-col li { font-size: 0.9rem; color: #64748b; padding: 0.4rem 0; border-bottom: 1px solid #f1f5f9; }
.mat-col li:last-child { border-bottom: none; }

/* ============================================================
   FAQ
   ============================================================ */
.faq-section { padding: 6rem 0; background: #f8fafc; }
.faq-section h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; text-align: center; margin-bottom: 3rem; }
.faq-cols { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; max-width: 1100px; margin: 0 auto; }
.faq-q { display: flex; gap: 1rem; background: #fff; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1.5rem; margin-bottom: 1rem; }
.faq-q:hover { border-color: #bfdbfe; }
.faq-icon { color: #1d6bf3; flex-shrink: 0; margin-top: 2px; }
.faq-q h4 { font-size: 0.95rem; font-weight: 700; color: #0f172a; margin-bottom: 0.5rem; }
.faq-q p { font-size: 0.875rem; color: #64748b; line-height: 1.6; }

/* ============================================================
   CTA BAND
   ============================================================ */
.cta-band {
  padding: 6rem 2rem;
  background: linear-gradient(135deg, #1d6bf3 0%, #0b47c0 100%);
  text-align: center;
  color: white;
}
.cta-band h2 { font-size: 2.5rem; font-weight: 900; letter-spacing: -1px; margin-bottom: 1rem; }
.cta-band p { font-size: 1.15rem; opacity: 0.85; margin-bottom: 2.5rem; }

/* ============================================================
   FOOTER
   ============================================================ */
.site-footer {
  background: #0f172a;
  color: #94a3b8;
  padding: 4rem 0 0;
}

.footer-top { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr 1fr; gap: 3rem; padding-bottom: 4rem; border-bottom: 1px solid rgba(255,255,255,0.08); }

.footer-logo { display: flex; align-items: center; gap: 0.5rem; color: #fff; margin-bottom: 1rem; }
.footer-logo .brand-icon { color: #60a5fa; }
.footer-logo strong { font-size: 1.1rem; }
.footer-brand p { font-size: 0.875rem; line-height: 1.6; max-width: 280px; margin-bottom: 1.5rem; }

.footer-certs { display: flex; gap: 0.5rem; flex-wrap: wrap; }
.footer-certs span { border: 1px solid rgba(255,255,255,0.15); padding: 0.25rem 0.6rem; border-radius: 4px; font-size: 0.75rem; color: #cbd5e1; font-weight: 600; }

.footer-col h5 { color: #fff; font-size: 0.875rem; font-weight: 700; margin-bottom: 1.5rem; letter-spacing: 0.5px; }
.footer-col a { display: block; font-size: 0.875rem; color: #94a3b8; text-decoration: none; margin-bottom: 0.75rem; transition: color 0.15s; }
.footer-col a:hover { color: #60a5fa; }

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  padding: 1.5rem 0;
  font-size: 0.8rem;
}

.legal-links { display: flex; gap: 1.5rem; flex-wrap: wrap; }
.legal-links a { color: #64748b; text-decoration: none; font-size: 0.8rem; transition: color 0.15s; }
.legal-links a:hover { color: #94a3b8; }
.legal-links a strong { color: #94a3b8; }

/* ============================================================
   RESPONSIVE
   ============================================================ */
@media (max-width: 1024px) {
  .hero { grid-template-columns: 1fr; min-height: auto; padding: 4rem 2rem; }
  .hero-right { display: none; }
  .caps-grid { grid-template-columns: 1fr 1fr; }
  .mat-grid { grid-template-columns: 1fr 1fr; }
  .why-grid { grid-template-columns: 1fr; gap: 3rem; }
  .pricing-breakdown { grid-template-columns: 1fr; }
  .footer-top { grid-template-columns: 1fr 1fr; }
}

@media (max-width: 768px) {
  .hero-left h1 { font-size: 2.8rem; }
  .caps-grid { grid-template-columns: 1fr; }
  .mat-grid { grid-template-columns: 1fr; }
  .faq-cols { grid-template-columns: 1fr; }
  .steps { flex-direction: column; }
  .step-arrow { transform: rotate(90deg); }
  .footer-bottom { flex-direction: column; text-align: center; }
  .footer-top { grid-template-columns: 1fr; }
}
</style>
