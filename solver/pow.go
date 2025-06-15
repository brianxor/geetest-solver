package solver

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/tylerdevx/geetest-solver/internal/crypto"
)

var hashFunctions = map[string]func(string) string{
	"md5":    md5Hash,
	"sha1":   sha1Hash,
	"sha256": sha256Hash,
}

type PowSolution struct {
	PowMessage string
	PowSign    string
}

func (c *GeetestSolverConfig) SolveV4PuzzlePow(captchaInfo *V4PuzzleCaptchaInfo) (*PowSolution, error) {
	deviceGuid, err := crypto.RandomHex(16)

	if err != nil {
		return nil, err
	}

	hashFunc, exists := hashFunctions[captchaInfo.Data.PowDetail.Hashfunc]

	if !exists {
		return nil, fmt.Errorf("pow hash function not found")
	}

	powMessage := fmt.Sprintf("%s|%d|%s|%s|%s|%s||%s",
		captchaInfo.Data.PowDetail.Version,
		captchaInfo.Data.PowDetail.Bits,
		captchaInfo.Data.PowDetail.Hashfunc,
		captchaInfo.Data.PowDetail.Datetime,
		c.CaptchaId,
		captchaInfo.Data.LotNumber,
		deviceGuid,
	)

	powSign := hashFunc(powMessage)

	return &PowSolution{
		PowMessage: powMessage,
		PowSign:    powSign,
	}, nil
}

func md5Hash(text string) string {
	sum := md5.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

func sha1Hash(text string) string {
	sum := sha1.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

func sha256Hash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}
