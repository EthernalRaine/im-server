package encryption

import (
	"bytes"
	"chimera/utility/logging"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rc4"
	"encoding/base64"
	"io"
)

func SwapRC4State(pwd []byte, data []byte) []byte {
	c, err := rc4.NewCipher(pwd)
	if err != nil {
		logging.Error("RC4/SwapState", "Error creating RC4 Ciphertext")
		return nil
	}
	crypted := make([]byte, len(data))
	c.XORKeyStream(crypted, data)
	return crypted
}

func EncryptAES(key string, message string) string {
	plainText := []byte(message)
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		logging.Error("AES/Encryption", "Error Encryption AES #1: %s", err.Error())
		return ""
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		logging.Error("AES/Encryption", "Error Encryption AES #2: %s", err.Error())
		return ""
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText)
}

func DecryptAES(key string, secure string) string {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		logging.Error("AES/Decryption", "Error Encryption AES #1: %s", err.Error())
		return ""
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		logging.Error("AES/Decryption", "Error Encryption AES #2: %s", err.Error())
		return ""
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		logging.Error("AES/Decryption", "Error Encryption AES #3 (short)")
		return ""
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func EncryptTDES(data, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := key
	iv := ciphertext[:des.BlockSize]
	origData := pKCS5Padding(data, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	mode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func DecryptTDES(data, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := key
	iv := ciphertext[:des.BlockSize]

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	decrypter.CryptBlocks(decrypted, data)
	decrypted = pKCS5UnPadding(decrypted)
	return decrypted, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
