package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"math/big"
	"math/rand"
	"time"
	"fmt"
)

type Terms struct {
	gorm.Model
	Current    bool
	Lang       string
	Text       string
	Valid      bool
	Difficulty int
}

type TermsSig struct {
	gorm.Model
	Signer    *types.Address `gorm:"index"`
	Timestamp string
	IP        string
	Nonce     *types.Uint256
	Signature *types.Signature
	Terms     Terms
}

type HashMask struct {
	gorm.Model
	Mask       []byte
	Expiration time.Time
}

type TermsManager struct {
	db      *gorm.DB
	isTx    bool
	signers map[string]struct{}
}

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

func (tm *TermsManager) ClearExpiredHashMasks() {
	tm.db.Model(&HashMask{}).Where("expiration < NOW()").Delete(HashMask{})
}

func (tm *TermsManager) GetTerms(language string) (*Terms, error) {
	terms := &Terms{}
	err := tm.db.Model(&Terms{}).Where("lang = ? AND current = ?", language, true).First(terms).Error
	return terms, err
}

func (tm *TermsManager) GetNewHashMask(terms *Terms) ([]byte, uint, error) {
	mask := new(big.Int)
	rand.Seed(time.Now().UTC().UnixNano())
	for OnesCount(mask.Bytes()) < terms.Difficulty {
		mask = mask.Or(mask, new(big.Int).Exp(big.NewInt(2), big.NewInt(rand.Int63n(255)), nil))
	}
	hashMask := &HashMask{
		Mask: mask.Bytes(),
		Expiration: time.Now().Add(time.Hour),
	}
	err := tm.db.Model(&HashMask{}).Create(hashMask).Error
	return mask.Bytes(), hashMask.ID, err
}

func (tm *TermsManager) GetHashMaskById(id uint) ([]byte, error) {
	hashMask := &HashMask{}
	err := tm.db.Model(&HashMask{}).First(hashMask, id).Error
	return hashMask.Mask, err
}

func (tm *TermsManager) CheckTerms(id uint, sig *types.Signature, address *types.Address, timestamp string, nonce []byte, mask []byte) (bool, error) {
	terms := &Terms{}
	if err := tm.db.Model(&Terms{}).First(terms).Error; err != nil {
		return false, err
	}
	return terms.CheckSig(sig, address, timestamp, nonce, mask)
}

func (terms *Terms) CheckSig(sig *types.Signature, address *types.Address, timestamp string, nonce []byte, mask []byte) (bool, error) {
	termsSha := sha3.NewKeccak256()
	termsSha.Write([]byte(fmt.Sprintf("%v\n%v\n%#x", terms.Text, timestamp, nonce)))
	hash := termsSha.Sum(nil)
	if !CheckMask(mask, hash) {
		return false, fmt.Errorf("Hash must match mask: %v", mask)
	}
	return sig.Verify(address, hash), nil
}

func FindValidNonce(terms *Terms, timestamp string, mask []byte) (<-chan []byte) {
	ch := make(chan []byte)
	go func(terms *Terms, timestamp string, mask []byte, ch chan []byte) {
		rand.Seed(time.Now().UTC().UnixNano())
		hash := []byte{}
		nonce := make([]byte, 32)
		for !CheckMask(mask, hash) {
			nonce = make([]byte, 32)
			rand.Read(nonce[:])
			termsSha := sha3.NewKeccak256()
			termsSha.Write([]byte(fmt.Sprintf("%v\n%v\n%#x", terms.Text, timestamp, nonce)))
			hash = termsSha.Sum(nil)
		}
		ch <- nonce
	}(terms, timestamp, mask, ch)
	return ch
}

func CheckMask(mask, hash []byte) bool {
	maskInt := new(big.Int).SetBytes(mask)
	hashInt := new(big.Int).SetBytes(hash)
	return new(big.Int).And(hashInt, maskInt).Cmp(maskInt) == 0
}

func (tm *TermsManager) Check(address *types.Address) (bool, error) {
	if _, ok := tm.signers[address.String()]; ok {
		return true, nil
	}
	var count int
	err := tm.db.Joins("LEFT JOIN terms ON terms.id = terms_sig.terms").Where("terms_sig.signer = ? AND terms.valid = 1", address).Count(&count).Error
	if err != nil {
		tm.signers[address.String()] = struct{}{}
	}
	return (count >= 1), err
}

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
	return &TermsManager{db, false, make(map[string]struct{})}
}

func NewTxTermsManager(db *gorm.DB) (*TermsManager) {
	return &TermsManager{db, true, make(map[string]struct{})}
}
