package solver

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"github.com/tylerdevx/geetest-solver/internal/crypto"
	"github.com/tylerdevx/geetest-solver/internal/utils"
	"github.com/tylerdevx/geetest-solver/solver/detection"
	"math/big"
)

var encryptionModulus, _ = new(big.Int).SetString("00C1E3934D1614465B33053E7F48EE4EC87B14B95EF88947713D25EECBFF7E74C7977D02DC1D9451F79DD5D1C10C29ACB6A9B4D6FB7D0A0279B6719E1772565F09AF627715919221AEF91899CAE08C0D686D748B20A3603BE2318CA6BC2B59706592A9219D0BF05C9F65023A21D2330807252AE0066D59CEEFA5F2748EA80BAB81", 16)
var encryptionExponent, _ = new(big.Int).SetString("10001", 16)

var encryptionPublicKey = &rsa.PublicKey{
	N: encryptionModulus,
	E: int(encryptionExponent.Int64()),
}

type V4PuzzlePayload struct {
	SetLeft      int     `json:"setLeft"`
	Passtime     int     `json:"passtime"`
	Userresponse float64 `json:"userresponse"`
	DeviceId     string  `json:"device_id"`
	LotNumber    string  `json:"lot_number"`
	PowMsg       string  `json:"pow_msg"`
	PowSign      string  `json:"pow_sign"`
	Geetest      string  `json:"geetest"`
	Lang         string  `json:"lang"`
	Ep           string  `json:"ep"`
	Biht         string  `json:"biht"`
	G9M2         string  `json:"G9M2"`
	F9293D       struct {
		F28A struct {
			A91474 string `json:"a91474"`
		} `json:"76f28a"`
	} `json:"f9293d"`
	Em struct {
		Ph int    `json:"ph"`
		Cp int    `json:"cp"`
		Ek string `json:"ek"`
		Wd int    `json:"wd"`
		Nt int    `json:"nt"`
		Si int    `json:"si"`
		Sc int    `json:"sc"`
	} `json:"em"`
}

func (c *GeetestSolverConfig) generateV4PuzzlePayload(captchaInfo *V4PuzzleCaptchaInfo) (*V4PuzzlePayload, error) {
	payload := &V4PuzzlePayload{}

	puzzleImage, err := c.FetchImage(geetestCaptchaImageUrl + captchaInfo.Data.Bg)

	if err != nil {
		return nil, err
	}

	pieceImage, err := c.FetchImage(geetestCaptchaImageUrl + captchaInfo.Data.Slice)

	if err != nil {
		return nil, err
	}

	puzzleImageBase64 := base64.StdEncoding.EncodeToString(puzzleImage)
	pieceImageBase64 := base64.StdEncoding.EncodeToString(pieceImage)
	puzzleSolver := detection.NewPuzzleSolver(puzzleImageBase64, pieceImageBase64)

	xPos, err := puzzleSolver.GetPosition(captchaInfo.Data.Ypos)

	if err != nil {
		return nil, err
	}

	powSolution, err := c.SolveV4PuzzlePow(captchaInfo)

	if err != nil {
		return nil, err
	}

	payload.SetLeft = xPos
	payload.Passtime = utils.RandomInt(500, 700)
	payload.Userresponse = (float64(xPos) / 1.0059466666666665) + 2
	payload.DeviceId = ""
	payload.LotNumber = captchaInfo.Data.LotNumber
	payload.PowMsg = powSolution.PowMessage
	payload.PowSign = powSolution.PowSign
	payload.Geetest = "captcha"
	payload.Lang = "zh"
	payload.Ep = "123"
	payload.Biht = "1426265548"
	payload.G9M2 = "1RZL"
	payload.F9293D.F28A.A91474 = "477ff9"
	payload.Em.Ph = 0
	payload.Em.Cp = 0
	payload.Em.Ek = "11"
	payload.Em.Wd = 1
	payload.Em.Nt = 0
	payload.Em.Si = 0
	payload.Em.Sc = 0

	return payload, nil
}

func encryptPayload(plainPayload []byte, publicKey *rsa.PublicKey) (string, error) {
	deviceGuid, err := crypto.RandomHex(16)

	if err != nil {
		return "", err
	}

	encryptedGuid, err := crypto.RsaEncrypt(deviceGuid, publicKey)

	if err != nil {
		return "", err
	}

	encryptedData, err := crypto.AesEncrypt(plainPayload, []byte(deviceGuid), []byte("0000000000000000"))

	if err != nil {
		return "", err
	}

	encryptedPayload := hex.EncodeToString(encryptedData) + encryptedGuid

	return encryptedPayload, nil
}
