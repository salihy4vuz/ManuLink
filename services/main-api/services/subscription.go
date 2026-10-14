package services

import (
	"errors"
	"sync"
)

type SubscriptionTier string

const (
	TierFree       SubscriptionTier = "Free"
	TierPro        SubscriptionTier = "Pro"
	TierEnterprise SubscriptionTier = "Enterprise"
)

type UserAccount struct {
	Username string
	Tier     SubscriptionTier
}

var (
	// In-memory mock DB for user accounts
	accounts = map[string]UserAccount{
		"admin": {Username: "admin", Tier: TierFree},
	}
	accountMutex sync.RWMutex
)

// GetUserTier returns the current subscription tier of the user.
func GetUserTier(username string) SubscriptionTier {
	accountMutex.RLock()
	defer accountMutex.RUnlock()

	if acc, exists := accounts[username]; exists {
		return acc.Tier
	}
	// Default to Free if unknown
	return TierFree
}

// UpgradeUserTier simulates a Stripe checkout success and upgrades the user.
func UpgradeUserTier(username string, newTier SubscriptionTier) error {
	if newTier != TierFree && newTier != TierPro && newTier != TierEnterprise {
		return errors.New("invalid subscription tier")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// If user exists, upgrade them. Otherwise create them.
	acc := accounts[username]
	acc.Username = username
	acc.Tier = newTier
	accounts[username] = acc

	return nil
}
