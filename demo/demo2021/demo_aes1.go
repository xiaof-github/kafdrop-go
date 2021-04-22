package agent

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {

	var key = []byte("WEtB1q8qOP3fYOM330BnZw==")
	fmt.Println("keySize: ", len(key))

	//每次得到的结果都不同，但是都可以解密
	fmt.Println("oracle=", []byte("oracle"))
	msg, ok := AesCBCEncrypt([]byte("oracle"), key)
	if ok != nil {
		fmt.Println("encrypt is failed.")
	}
	fmt.Println("msg=", msg)

	text, ok := AesCBCDncrypt(msg, key)
	if ok != nil {
		fmt.Println("decrypt is failed.")
	}

	fmt.Println("text=", text)
	fmt.Println("text=", string(text))
}

/*CBC加密 按照golang标准库的例子代码
不过里面没有填充的部分,所以补上
*/
//使用PKCS5进行填充

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//使用PKCS7进行填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//填充原文
	blockSize := block.BlockSize()
	fmt.Println("blockSize: ", blockSize)

	rawData = PKCS5Padding(rawData, blockSize)

	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))

	//block大小 16
	iv := cipherText[:blockSize]

	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	panic(err)
	// }
	copy(iv, key)
	fmt.Println("IV: ", iv)

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func AesCBCDncrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]

	encryptData = encryptData[blockSize:]
	fmt.Println("encryptData=", encryptData)
	fmt.Println(len(encryptData) % blockSize)

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS5UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData, key []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Dncrypt(rawData string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := AesCBCDncrypt(data, key)
	if err != nil {
		return "", err
	}
	return string(dnData), nil //
}
