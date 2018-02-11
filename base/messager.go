package base

import (
	"bytes"
	"container/list"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"strconv"
	"time"

	"gopkg.in/configo.v2"
)

type AesEncrypt struct {
	Key string
	Iv  string
}

func NewEnc() *AesEncrypt {
	return &AesEncrypt{}
}

func init() {
	InitConfig()
}

func (this *AesEncrypt) getKey() []byte {

	keyLen := len(this.Key)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(this.Key)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//加密字符串
func (this *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	plantText := []byte(strMesg)

	key := this.getKey()
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = this.PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, []byte(this.Iv)[:aes.BlockSize])

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (strDesc string, err error) {

	defer func() {
		//错误处理
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()

	key := this.getKey()
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return "", err
	}
	blockModel := cipher.NewCBCDecrypter(block, []byte(this.Iv)[:aes.BlockSize])
	plantText := make([]byte, len(src))
	blockModel.CryptBlocks(plantText, src)
	plantText = this.PKCS7UnPadding(plantText, block.BlockSize())
	return string(plantText), nil
}

//补位
func (this *AesEncrypt) PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//补位
func (this *AesEncrypt) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

type MessageType int

const (
	MT_SEND MessageType = iota
	MT_CHECK
)

var SEND_URL string = "/sms/send"
var CHECK_URL string = "/sms/check"

func (m *MessageType) StringUri() string {
	if *m == MT_CHECK {
		return CHECK_URL
	}
	return SEND_URL
}

func SendMessage(m MessageType, data map[string]string) *map[string]string {

	aesEnc := new(AesEncrypt)
	aesEnc.Iv = configsms.Iv
	aesEnc.Key = configsms.Key

	rlt := ListToString(TransToListByDesc(data), "&")

	b, _ := aesEnc.Encrypt(rlt)
	str := base64.StdEncoding.EncodeToString(b)

	senddata := strings.Join([]string{
		"sign=" + url.QueryEscape(str),
		rlt,
	}, "&")

	c := &http.Client{}

	postUrl := LoadUrl(m)
	log.Println("msg post: " + postUrl)
	req, err := http.NewRequest("POST", postUrl, strings.NewReader(senddata))
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	r := make(map[string]string)
	log.Println("msg body: " + string(body))
	_ = json.Unmarshal(body, &r)
	return &r
}

func TransToListByDesc(data map[string]string) *list.List {
	l := list.New()
	for k, v := range data {
		inval := fmt.Sprintf("%s=%s", k, v)

		ele := l.Front()
		var last *list.Element = nil
		for ; ele != nil; ele = ele.Next() {
			tmp := ele.Value.(string)

			if strings.Compare(tmp, inval) > 0 {

				if last != nil {
					l.InsertBefore(inval, ele)
					break
				} else {
					l.PushFront(inval)
					break
				}
			}
			last = ele
		}

		if ele == nil {
			l.PushBack(inval)
		}

	}
	return l
}

func ListToString(list *list.List, con string) string {
	ele := list.Front()
	rlt := ""
	for ; ele != nil; ele = ele.Next() {
		if rlt == "" {
			rlt = ele.Value.(string)
		} else {
			rlt = strings.Join([]string{rlt, ele.Value.(string)}, con)
		}

	}

	return rlt
}

func LoadUrl(messageType MessageType) string {
	sms, e := configo.Get("sms")
	if e != nil {
		return "http://localhost:8006" + messageType.StringUri()
	}
	addr := sms.MustGet("addr", "localhost")
	port := sms.MustGet("port", "8006")

	return "http://" + strings.Join([]string{addr, port}, ":") + messageType.StringUri()
}

type ConfigSMS struct {
	AppId string
	Key   string
	Iv    string
}

var configsms = ConfigSMS{}

func InitConfig() error {
	conf, e := configo.Get("configsms")
	if e != nil {
		log.Println(e.Error())
		return e
	}
	configsms.AppId = conf.MustGet("app_id", "5f9c2d45f20b789e")
	configsms.Iv = conf.MustGet("iv", "aaC5p6c5L2g6KeJ5")
	configsms.Key = conf.MustGet("key", "sdf234wef34efrfT")
	return nil
}

func RegisterSend(mobile string) *map[string]string {
	m := make(map[string]string)
	m["mobile"] = mobile
	m["servicename"] = "mana_register"
	m["clientid"] = configsms.AppId
	m["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["content"] = ""
	m["type"] = "1"
	return SendMessage(MT_SEND, m)
}

func MessageCheck(mobile, code string) *map[string]string {
	m := make(map[string]string)
	m["mobile"] = mobile
	//m["servicename"] = "bingo_register"
	m["clientid"] = configsms.AppId
	m["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	//m["content"] = "{\"code\":\"1234\",\"product\":\"玛娜爱情\"}"
	m["code"] = code
	//CheckMessage(m)
	return SendMessage(MT_CHECK, m)

}
