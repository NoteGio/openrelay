package accounts_test

import (
	"github.com/notegio/0xrelay/accounts"
	"math/big"
	"testing"
	"time"
)

func TestBlacklisted(t *testing.T) {
	account := accounts.NewAccount(true, new(big.Int), 0, 0)
	if !account.Blacklisted() {
		t.Errorf("Expected account blacklisted")
	}
}
func TestExpiredDiscount(t *testing.T) {
	account := accounts.NewAccount(false, new(big.Int), 0, 0)
	if discount := account.Discount(); discount.Cmp(new(big.Int)) != 0 {
		t.Errorf("Expected 0 discount, got '%v'", discount)
	}
}

func TestFiftyPercentDiscount(t *testing.T) {
	baseFee := big.NewInt(10000)
	account := accounts.NewAccount(false, baseFee, 50, time.Now().Unix()+5)
	if discount := account.Discount(); discount.Cmp(big.NewInt(5000)) != 0 {
		t.Errorf("Expected 5000 discount, got '%v'", discount)
	}
}
