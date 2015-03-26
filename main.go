package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	twilio "github.com/carlosdp/twiliogo"
)

var from = flag.String("from", "", "the sender phone number")
var to = flag.String("to", "", "the recipient phone number")
var msg = flag.String("msg", "", "the contents of the msg")
var sms = flag.Bool("sms", true, "send message as an SMS")

func main() {
	flag.Parse()
	pw, err := speakeasy.Ask("Password for Twilio credentials: ")
	if err != nil {
		log.Fatal("Invalid input for password.")
	}

	client := twilio.NewClient(decrypt(pw, os.Getenv("TWILIO_SID")), decrypt(pw, os.Getenv("TWILIO_SECRET")))
	message, err := twilio.NewMessage(client, *from, *to, twilio.Body(*msg))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(message.Status)
	}

}

func decrypt(key string, ciphertext string) string {
	access_token, _ := hex.DecodeString(ciphertext)
	key += strings.Repeat(" ", 16-(len(key)%16))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	if len(access_token) < aes.BlockSize {
		panic("access_token too short")
	}
	iv := access_token[:aes.BlockSize]
	access_token = access_token[aes.BlockSize:]
	if len(access_token)%aes.BlockSize != 0 {
		panic("access_token is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(access_token, access_token)
	return strings.TrimSpace(fmt.Sprintf("%s", access_token))
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say string `xml:",omitempty"`
}
