package vivo

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

func int64toString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func hMACSHA256HEX(data string, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	hash := h.Sum(nil)
	return hash
}

func base64encode(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

func GenerateRequestID() string {
	s := uuid.New().String()
	return s
}

func GenerateSessionID() string {
	s := uuid.New().String()
	return s
}

func PcmToWav(dst []byte) []byte {
	numchannel := 2
	saplerate := 16000
	byteDst := dst
	longSampleRate := saplerate
	byteRate := 16 * saplerate * numchannel / 8
	totalAudioLen := len(byteDst)
	totalDataLen := totalAudioLen + 36
	var header = make([]byte, 44)
	// RIFF/WAVE header
	header[0] = 'R'
	header[1] = 'I'
	header[2] = 'F'
	header[3] = 'F'
	header[4] = byte(totalDataLen & 0xff)
	header[5] = byte((totalDataLen >> 8) & 0xff)
	header[6] = byte((totalDataLen >> 16) & 0xff)
	header[7] = byte((totalDataLen >> 24) & 0xff)
	//WAVE
	header[8] = 'W'
	header[9] = 'A'
	header[10] = 'V'
	header[11] = 'E'
	// 'fmt ' chunk
	header[12] = 'f'
	header[13] = 'm'
	header[14] = 't'
	header[15] = ' '
	// 4 bytes: size of 'fmt ' chunk
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	// format = 1
	header[20] = 1
	header[21] = 0
	header[22] = byte(numchannel)
	header[23] = 0
	header[24] = byte(longSampleRate & 0xff)
	header[25] = byte((longSampleRate >> 8) & 0xff)
	header[26] = byte((longSampleRate >> 16) & 0xff)
	header[27] = byte((longSampleRate >> 24) & 0xff)
	header[28] = byte(byteRate & 0xff)
	header[29] = byte((byteRate >> 8) & 0xff)
	header[30] = byte((byteRate >> 16) & 0xff)
	header[31] = byte((byteRate >> 24) & 0xff)
	// block align
	header[32] = byte(2 * 16 / 8)
	header[33] = 0
	// bits per sample
	header[34] = 16
	header[35] = 0
	//data
	header[36] = 'd'
	header[37] = 'a'
	header[38] = 't'
	header[39] = 'a'
	header[40] = byte(totalAudioLen & 0xff)
	header[41] = byte((totalAudioLen >> 8) & 0xff)
	header[42] = byte((totalAudioLen >> 16) & 0xff)
	header[43] = byte((totalAudioLen >> 24) & 0xff)

	headerDst := string(header)
	resDst := headerDst + string(dst)
	return []byte(resDst)
}

func GenerateVisionChatImage(file []byte) string {
	return "data:image/JPEG;base64," + base64encode(file)
}

func pcmIntToBytes(data []int, bitDepth int) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	bytesPerSample := bitDepth / 8
	if bitDepth%8 != 0 {
		// 对于非标准的、非字节对齐的位深度
		bytesPerSample = (bitDepth + 7) / 8
	}
	if bytesPerSample <= 0 {
		return nil, fmt.Errorf("无效的位深度 %d", bitDepth)
	}

	buf := new(bytes.Buffer)
	buf.Grow(len(data) * bytesPerSample) // 预分配容量

	for _, sample := range data {
		switch bitDepth {
		case 8:
			// 假设 8bit WAV 是 unsigned (0-255)
			if sample < 0 || sample > 255 {
				// Clamp or return error? Clamping might hide issues.
				// return nil, fmt.Errorf("8-bit sample %d out of range [0, 255]", sample)
				if sample < 0 {
					sample = 0
				}
				if sample > 255 {
					sample = 255
				}
			}
			if err := buf.WriteByte(byte(sample)); err != nil {
				return nil, err
			}
		case 16:
			// int -> int16 -> uint16 (for bit pattern) -> LittleEndian bytes
			if err := binary.Write(buf, binary.LittleEndian, int16(sample)); err != nil {
				return nil, err
			}
		case 24:
			// 手动写入 3 字节小端序
			b := []byte{
				byte(sample & 0xFF),
				byte((sample >> 8) & 0xFF),
				byte((sample >> 16) & 0xFF),
			}
			if _, err := buf.Write(b); err != nil {
				return nil, err
			}
		case 32:
			// int -> int32 -> uint32 (for bit pattern) -> LittleEndian bytes
			// 假设是整数 PCM，非浮点
			if err := binary.Write(buf, binary.LittleEndian, int32(sample)); err != nil {
				return nil, err
			}
		default:
			// 尝试通用小端写入，可能不适用于所有非标准格式
			temp := make([]byte, bytesPerSample)
			for j := 0; j < bytesPerSample; j++ {
				temp[j] = byte((sample >> (j * 8)) & 0xFF)
			}
			if _, err := buf.Write(temp); err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}
