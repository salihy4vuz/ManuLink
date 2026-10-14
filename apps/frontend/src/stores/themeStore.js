import { defineStore } from 'pinia';
import { fetchTenantConfig } from '../api/tenant';

export const useThemeStore = defineStore('theme', {
  state: () => ({
    tenantId: null,
    name: 'Sheet Metal Automation',
    primaryColor: '#3498db',
    secondaryColor: '#2c3e50',
    logoUrl: '/default-logo.svg',
    faviconUrl: '/favicon.ico',
    isLoading: false,
    error: null
  }),
  actions: {
    async loadTenantTheme(hostname) {
      this.isLoading = true;
      try {
        const config = await fetchTenantConfig(hostname);
        if (config) {
          this.tenantId = config.id;
          this.name = config.name || this.name;
          this.primaryColor = config.primaryColor || this.primaryColor;
          this.secondaryColor = config.secondaryColor || this.secondaryColor;
          this.logoUrl = config.logoUrl || this.logoUrl;
          this.faviconUrl = config.faviconUrl || this.faviconUrl;
          this.applyTheme();
        }
      } catch (err) {
        this.error = err.message;
        console.error('Failed to load tenant theme:', err);
      } finally {
        this.isLoading = false;
      }
    },
    applyTheme() {
      const root = document.documentElement;
      root.style.setProperty('--color-primary', this.primaryColor);
      root.style.setProperty('--color-secondary', this.secondaryColor);
      
      // Update Title
      document.title = `${this.name} - B2B Portal`;
      
      // Update Favicon
      let link = document.querySelector("link[rel~='icon']");
      if (!link) {
        link = document.createElement('link');
        link.rel = 'icon';
        document.getElementsByTagName('head')[0].appendChild(link);
      }
      link.href = this.faviconUrl;
    },
    updateTheme(newConfig) {
      this.primaryColor = newConfig.primaryColor;
      this.secondaryColor = newConfig.secondaryColor;
      this.logoUrl = newConfig.logoUrl;
      this.applyTheme();
    }
  }
});
