package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	"golang.org/x/crypto/sha3"
	"math/big"
	"math/rand"
	"time"
	"fmt"
)

// Terms represents an instance of the terms of service. There should only be
// one Current version of the Terms for a given language at a given time, but
// older versions may remain Valid. The difficulty indicates how many 1 bits
// a HashMask must have when signing these terms.
type Terms struct {
	gorm.Model
	Current    bool
	Lang       string
	Text       string   `sql:"type:text;"`
	Valid      bool
	Difficulty int
}

// TermsSig represent a user's signature of the terms of service. It tracks the
// address of the signer, when they signed it, the IP they submitted the
// signature from, the Nonce they used to make the HashMask match, the
// cryptographic signature, and a reference to the Terms object they signed.
type TermsSig struct {
	gorm.Model
	Signer    *types.Address `gorm:"index"`
	Timestamp string
	IP        string
	Nonce     []byte
	Banned    bool
	Signature *types.Signature
	TermsID   uint
}

// HashMask is a string of bytes. When a user signs the terms-of-service, they
// must provide a Nonce that makes the hash of the terms (along with the
// timestamp) match the hash mask. Matching the hash mask means that any 1 bits
// in the HashMask must also be 1 in the hash. 0s in the HashMask can be 1 or 0
// in the hash.
type HashMask struct {
	gorm.Model
	Mask       []byte
	Expiration time.Time
}

// TermsManager keeps track of a database instance and a cache and provides
// several interfaces for interacting with the Terms of Service in the
// database.
type TermsManager struct {
	db           *gorm.DB
	isTx          bool
	signers       map[string]struct{}
	hashMaskCache map[uint]*HashMask
}

// OnesCount returns how many 1s appear in the binary representation of a
// string of bytes. Similar to math/bits.OnesCount, but works for arbitrarily
// sized byte slices.
func OnesCount(data []byte) (int) {
	count := 0
	for i := 0; i < len(data); i++ {
		val := byte(1)
		for j := 0; j < 8; j++ {
			if data[i] & val != 0 {
				count += 1
			}
			val *= 2
		}
	}
	return count
}

// ClearExpiredHashMasks deletes any expired hash masks. These can be deleted
// even if a corresponding signature has been created; they're only needed
// during the sign up process.
func (tm *TermsManager) ClearExpiredHashMasks() {
	tm.db.Model(&HashMask{}).Where("expiration < NOW()").Delete(HashMask{})
}

// GetTerms returns the current Terms object for a given language.
func (tm *TermsManager) GetTerms(language string) (*Terms, error) {
	terms := &Terms{}
	err := tm.db.Model(&Terms{}).Where("lang = ? AND current = ?", language, true).First(terms).Error
	return terms, err
}

// GetNewHashMask generates a HashMask for a Terms object, considering the
// difficulty for those Terms. It saves the HashMask in the database, and
// returns the Bytes along with the byte string.
func (tm *TermsManager) GetNewHashMask(terms *Terms) ([]byte, uint, error) {
	hashMask, ok := tm.hashMaskCache[terms.ID]
	if ok && hashMask.Expiration.After(time.Now().Add(time.Hour / 2)) {
		return hashMask.Mask, hashMask.ID, nil
	}
	mask := new(big.Int)
	rand.Seed(time.Now().UTC().UnixNano())
	for OnesCount(mask.Bytes()) < terms.Difficulty {
		mask = mask.Or(mask, new(big.Int).Exp(big.NewInt(2), big.NewInt(rand.Int63n(255)), nil))
	}
	hashMask = &HashMask{
		Mask: mask.Bytes(),
		Expiration: time.Now().Add(time.Hour),
	}
	err := tm.db.Model(&HashMask{}).Create(hashMask).Error
	if err == nil {
		tm.hashMaskCache[terms.ID] = hashMask
	}
	return mask.Bytes(), hashMask.ID, err
}

// GetHashMaskById retrieves a HashMask from the database given its ID.
func (tm *TermsManager) GetHashMaskById(id uint) ([]byte, error) {
	hashMask := &HashMask{}
	err := tm.db.Model(&HashMask{}).First(hashMask, id).Error
	return hashMask.Mask, err
}

// CheckTerms verifies that a signature is valid for a given Terms of use,
// specified by ID.
func (tm *TermsManager) CheckTerms(id uint, sig *types.Signature, address *types.Address, timestamp string, nonce []byte, mask []byte) (bool, error) {
	terms := &Terms{}
	if err := tm.db.Model(&Terms{}).First(terms).Error; err != nil {
		return false, err
	}
	return terms.CheckSig(sig, address, timestamp, nonce, mask)
}

// CheckSig verifies that a signature is valid for a given Terms of use object
func (terms *Terms) CheckSig(sig *types.Signature, address *types.Address, timestamp string, nonce []byte, mask []byte) (bool, error) {
	signedMessage := []byte(fmt.Sprintf("%v\n%v\n%#x", terms.Text, timestamp, nonce))
	termsSha := sha3.NewLegacyKeccak256()
	termsSha.Write(signedMessage)
	hash := termsSha.Sum(nil)
	if !CheckMask(mask, hash) {
		return false, fmt.Errorf("Hash must match mask: %#x, got %#x", mask, hash)
	}
	return sig.Verify(address, signedMessage), nil
}

// SaveSig verifies that a signature is valid for a given Terms, then saves it
// to the database
func (tm *TermsManager) SaveSig(id uint, sig *types.Signature, address *types.Address, timestamp, host_ip string, nonce []byte, mask []byte) (error) {
	terms := &Terms{}
	if err := tm.db.Model(&Terms{}).First(terms).Error; err != nil {
		return err
	}
	if ok, err := terms.CheckSig(sig, address, timestamp, nonce, mask); !ok || err != nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("Signature invalid")
	}
	terms_sig := &TermsSig{
		Signer: address,
		Timestamp: timestamp,
		IP: host_ip,
		Nonce: nonce,
		Signature: sig,
		TermsID: terms.ID,
		Banned: false,
	}
	return tm.db.Model(&TermsSig{}).Create(terms_sig).Error
}

// FindValidNonce finds a valid hash for a given hash mask, Terms, and
// Timestamp, and returns the Nonce used to generate thas hashmask.
func FindValidNonce(text, timestamp string, mask []byte) (<-chan []byte) {
	ch := make(chan []byte)
	go func(text, timestamp string, mask []byte, ch chan []byte) {
		rand.Seed(time.Now().UTC().UnixNano())
		hash := []byte{}
		nonce := make([]byte, 32)
		for !CheckMask(mask, hash) {
			nonce = make([]byte, 32)
			rand.Read(nonce[:])
			termsSha := sha3.NewLegacyKeccak256()
			termsSha.Write([]byte(fmt.Sprintf("%v\n%v\n%#x", text, timestamp, nonce)))
			hash = termsSha.Sum(nil)
		}
		ch <- nonce
	}(text, timestamp, mask, ch)
	return ch
}

// CheckMask verifies that a given hash matches the provided hashmask
func CheckMask(mask, hash []byte) bool {
	maskInt := new(big.Int).SetBytes(mask)
	hashInt := new(big.Int).SetBytes(hash)
	return new(big.Int).And(hashInt, maskInt).Cmp(maskInt) == 0
}

// Check ensures that a given address has signed the terms
func (tm *TermsManager) CheckAddress(address *types.Address) (<-chan bool) {
	result := make(chan bool)
	go func(address *types.Address, result chan bool) {
		if _, ok := tm.signers[address.String()]; ok {
			result <- ok
			return
		}
		var count int
		err := tm.db.Table("terms_sigs").Joins("LEFT JOIN terms ON terms.id = terms_sigs.terms_id").Where("terms_sigs.signer = ? AND terms.valid = ? AND terms_sigs.banned = ?", address, true, false).Count(&count).Error
		if err != nil && (count >= 1) {
			tm.signers[address.String()] = struct{}{}
		}
		result <- (count >= 1)
	}(address, result)
	return result
}

// UpdateTerms creates a new Terms object in the database, deprecating the old
// terms
func (tm *TermsManager) UpdateTerms(lang, text string) error {
	oldTerms, oldTermsErr := tm.GetTerms(lang)
	tx := tm.db
	if !tm.isTx {
		tx = tm.db.Begin()
	}
	if tx.Error != nil {
		return tx.Error
	}
	difficulty := 3
	if oldTermsErr == nil {
		oldTerms.Current = false
		if err := tx.Save(oldTerms).Error; err != nil {
			tx.Rollback()
			return err
		}
		difficulty = oldTerms.Difficulty
	}
	if err := tx.Create(&Terms{
		Current: true,
		Lang: lang,
		Text: text,
		Valid: true,
		Difficulty: difficulty,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if !tm.isTx {
		return tx.Commit().Error
	}
	return nil
}

func NewTermsManager(db *gorm.DB) (*TermsManager) {
	return &TermsManager{db, false, make(map[string]struct{}), make(map[uint]*HashMask)}
}

func NewTxTermsManager(db *gorm.DB) (*TermsManager) {
	return &TermsManager{db, true, make(map[string]struct{}), make(map[uint]*HashMask)}
}
