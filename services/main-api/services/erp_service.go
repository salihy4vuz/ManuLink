package services

import (
	"log"
	"strings"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services/erp"
)

// GetERPAdapter acts as a Factory returning the correct strategy based on Tenant settings
func GetERPAdapter(erpSystem string) erp.ERPAdapter {
	switch strings.ToLower(erpSystem) {
	case "logo":
		return &erp.LogoAdapter{}
	case "mikro":
		return &erp.MikroAdapter{}
	case "netsis":
		return &erp.NetsisAdapter{}
	case "canias":
		return &erp.CaniasAdapter{}
	// --- KOBI Optimizasyonu: Kurumsal / Global ERP'ler devredisi ---
	// case "sap":
	// 	return &erp.SAPAdapter{}
	// case "odoo":
	// 	return &erp.OdooAdapter{}
	// case "dynamics":
	// 	return &erp.DynamicsAdapter{}
	// case "netsuite":
	// 	return &erp.NetSuiteAdapter{}
	// case "epicor":
	// 	return &erp.EpicorAdapter{}
	// -------------------------------------------------------------
	case "promanage":
		return &erp.ProManageAdapter{}
	case "trexdcas":
		return &erp.TrexAdapter{}
	default:
		return &erp.DefaultAdapter{}
	}
}

// SyncOrderToERP acts as the main entry point for the controllers
func SyncOrderToERP(tenant models.Tenant, order models.Order, items []models.OrderItem) error {
	adapter := GetERPAdapter(tenant.ERPSystem)
	
	log.Printf("[ERP Router] Routing order %d for Tenant '%s' to %s adapter...\n", order.ID, tenant.Name, tenant.ERPSystem)
	
	return adapter.SyncOrder(tenant, order, items)
}
