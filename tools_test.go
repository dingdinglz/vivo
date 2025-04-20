package vivo

import (
	"fmt"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	fmt.Println(generateRandomString(8))
}

func TestGenerateRequestID(t *testing.T) {
	fmt.Println(GenerateRequestID())
}

func TestSign(t *testing.T) {
	fmt.Println(hMACSHA256HEX("GET\n/search/geo\ncity=%E6%B7%B1%E5%9C%B3&keywords=%E4%B8%8A%E6%A2%85%E6%9E%97&page_num=1&page_size=3\n1080389454\n1629255133\nx-ai-gateway-app-id:1080389454\nx-ai-gateway-timestamp:1629255133\nx-ai-gateway-nonce:le1qqjex", "XpurLJTrKSuAGoIq"))
	fmt.Println(base64encode(hMACSHA256HEX("GET\n/search/geo\ncity=%E6%B7%B1%E5%9C%B3&keywords=%E4%B8%8A%E6%A2%85%E6%9E%97&page_num=1&page_size=3\n1080389454\n1629255133\nx-ai-gateway-app-id:1080389454\nx-ai-gateway-timestamp:1629255133\nx-ai-gateway-nonce:le1qqjex", "XpurLJTrKSuAGoIq")))
}
