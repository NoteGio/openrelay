package db_test

import (
	"crypto/rand"
	"crypto/ecdsa"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/notegio/openrelay/types"
	"testing"
	"fmt"
	"bytes"
	// "log"
)

func getTermsManager() (*dbModule.TermsManager, func(), error) {
	db, err := getDb()
	if err != nil {
		return nil, func(){}, fmt.Errorf("Could not get db: '%v'", err.Error())
	}
	tx := db.Begin()
	if err := tx.AutoMigrate(&dbModule.Terms{}).Error; err != nil {
		tx.Rollback()
		return nil, func(){}, fmt.Errorf("Could not migrate Terms: '%v'", err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.TermsSig{}).Error; err != nil {
		tx.Rollback()
		return nil, func(){}, fmt.Errorf("Could not migrate TermsSig: '%v'", err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.HashMask{}).Error; err != nil {
		tx.Rollback()
		return nil, func(){}, fmt.Errorf("Could not migrate HashMask: '%v'", err.Error())
	}
	return dbModule.NewTxTermsManager(tx), func() { tx.Rollback() }, nil
}

func TestUpdateTerms(t *testing.T) {
	tm, revert, err := getTermsManager()
	if err != nil {
		t.Fatalf("Could not get db: %v", err.Error())
	}
	defer revert()
	if err := tm.UpdateTerms("en", "Don't break the law"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	terms, err := tm.GetTerms("en")
	if err != nil {
		t.Fatalf("Error getting terms: '%v'", err.Error())
	}
	if terms.Text != "Don't break the law" {
		t.Errorf("Unexpected terms text: '%v'", terms.Text)
	}
}

func TestUpdateTermsRepeat(t *testing.T) {
	tm, revert, err := getTermsManager()
	if err != nil {
		t.Fatalf("Could not get db: %v", err.Error())
	}
	defer revert()
	if err := tm.UpdateTerms("en", "Don't break the law"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	if err := tm.UpdateTerms("en", "Don't break the law please"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	terms, err := tm.GetTerms("en")
	if err != nil {
		t.Fatalf("Error getting terms: '%v'", err.Error())
	}
	if terms.Text != "Don't break the law please" {
		t.Errorf("Unexpected terms text: '%v'", terms.Text)
	}
}

func TestHashMask(t *testing.T) {
	tm, revert, err := getTermsManager()
	if err != nil {
		t.Fatalf("Could not get db: %v", err.Error())
	}
	defer revert()
	if err := tm.UpdateTerms("en", "Don't break the law"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	terms, err := tm.GetTerms("en")
	if err != nil {
		t.Fatalf("Error getting terms: '%v'", err.Error())
	}
	mask, id, err := tm.GetNewHashMask(terms)
	if err != nil {
		t.Fatalf("Error creating mask: '%v'", err.Error())
	}
	loadedMask, err := tm.GetHashMaskById(id)
	if err != nil {
		t.Fatalf("Error creating mask: '%v'", err.Error())
	}
	if count := dbModule.OnesCount(loadedMask); count != 3 {
		t.Errorf("Expected 3 bits in mask, got %v", count)
	}
	if !bytes.Equal(loadedMask, mask) {
		t.Errorf("Expected '%#x' = '%#x'", loadedMask, mask)
	}
}

func TestFindValidNonce(t *testing.T) {
	tm, revert, err := getTermsManager()
	if err != nil {
		t.Fatalf("Could not get db: %v", err.Error())
	}
	defer revert()
	if err := tm.UpdateTerms("en", "Don't break the law"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	terms, err := tm.GetTerms("en")
	if err != nil {
		t.Fatalf("Error getting terms: '%v'", err.Error())
	}
	mask, _, err := tm.GetNewHashMask(terms)
	if err != nil {
		t.Fatalf("Error creating mask: '%v'", err.Error())
	}
	timestamp := "1543351413"
	nonce := <-dbModule.FindValidNonce(terms, timestamp, mask)
	termsSha := sha3.NewKeccak256()
	termsSha.Write([]byte(fmt.Sprintf("%v\n%v\n%#x", terms.Text, timestamp, nonce)))
	hash := termsSha.Sum(nil)
	if !dbModule.CheckMask(mask, hash) {
		t.Errorf("Mask does not match: %#x <> %#x", mask, hash)
	}

}


func TestCheckTerms(t *testing.T) {
	tm, revert, err := getTermsManager()
	if err != nil {
		t.Fatalf("Could not get db: %v", err.Error())
	}
	defer revert()
	if err := tm.UpdateTerms("en", "Don't break the law"); err != nil {
		t.Fatalf("Error setting terms: '%v'", err.Error())
	}
	terms, err := tm.GetTerms("en")
	if err != nil {
		t.Fatalf("Error getting terms: '%v'", err.Error())
	}
	mask, _, err := tm.GetNewHashMask(terms)
	if err != nil {
		t.Fatalf("Error creating mask: '%v'", err.Error())
	}
	timestamp := "1543351413"
	nonce := <-dbModule.FindValidNonce(terms, timestamp, mask)
	termsSha := sha3.NewKeccak256()
	termsSha.Write([]byte(fmt.Sprintf("%v\n%v\n%#x", terms.Text, timestamp, nonce)))
	hash := termsSha.Sum(nil)

	key, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	address := crypto.PubkeyToAddress(key.PublicKey)
	signer := &types.Address{}
	copy(signer[:], address[:])

	hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), hash[:]...)
	signedBytes := crypto.Keccak256(hashedBytes)

	sig, _ := crypto.Sign(signedBytes, key)
	signature := make(types.Signature, 66)
	signature[0] = sig[64] + 27
	copy(signature[1:33], sig[0:32])
	copy(signature[33:65], sig[32:64])
	signature[65] = types.SigTypeEthSign

	if ok, err := tm.CheckTerms(terms.ID, &signature, signer, timestamp, nonce, mask); !ok || err != nil {
		if err != nil {
			t.Errorf("Error checking signature: '%v'", err.Error())
		} else {
			t.Errorf("Signature did not validate")
		}
	}
}
