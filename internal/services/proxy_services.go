package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"math/rand"
	"regexp"
	"time"
)

// GetEnableProxyServerList 获取启用的代理服务器列表
func GetEnableProxyServerList() (res []models.ProxyServer, err error) {
	res, err = models.ProxyServer{}.GetProxyServerList()
	if err != nil {
		err = fmt.Errorf("failed to get proxy server list")
	}
	return
}

type AccountParams struct {
	Country         int    // 国家
	Province        int    // 省份
	City            int    // 城市
	ProxyProtocol   string // 代理协议
	ProxyURL        string // 代理服务地址
	SessionValidity int8   // Session有效性 节点掉线，Session时间内是否补充节点使用 1-是、0-否 int
	SessionTime     int    // Session时间 /分钟
	Account         string // 账号
	Password        string // 密码
	Format          string // 格式
}

// GenerateRandomSession 生成一个 16 位的随机 session 字符串
func GenerateRandomSession() (string, error) {
	// 创建一个长度为 8 的字节数组（16 位的字符串需要 8 个字节）
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// 将字节数组转换为 16 位的十六进制字符串
	return hex.EncodeToString(bytes), nil
}

// GenerateTimestampMD5 生成当前时间戳的 MD5 哈希值
func GenerateTimestampMD5() string {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 将时间戳转换为字符串
	timestampStr := fmt.Sprintf("%d", timestamp)

	// 生成 MD5 哈希值
	hash := md5.New()
	hash.Write([]byte(timestampStr))
	hashBytes := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashStr := hex.EncodeToString(hashBytes)

	return hashStr
}

func MakeAccountExtractionNodes(param AccountParams, num int) (cmd []string) {
	// curl -x http://account_country_province_city_session_SessionMinute_add:123-*-789@转发服务ip:12348 ipinfo.io
	// curl -x socks5://account_country_province_city_session2_SessionMinute_add:123-*-789@转发服务ip:12347 ipinfo.io
	// const CurlCmd = "curl -x %s://%s ipinfo.io"
	var session string
	var err error
	var port int
	port, err = models.ProxyServer{}.GetProxyServerPortByProtocol(param.ProxyProtocol)
	if err != nil {
		log.Error("failed to get proxy server port:%v", err.Error())
		return
	}

	param.ProxyURL = param.ProxyURL + fmt.Sprintf(":%d", port)
	country := param.Country
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		if param.Country == 0 {
			country = r.Intn(250) + 1
		}
		var nodeStr = param.Format
		if (param.SessionTime != 0 && param.SessionValidity == 1) || i == 0 {
			session, err = GenerateRandomSession()
			if err != nil {
				session = GenerateTimestampMD5()
			}
		}
		addr := regexp.MustCompile(`hostname:port`)
		nodeStr = addr.ReplaceAllString(nodeStr, param.ProxyURL)
		username := regexp.MustCompile(`username`)
		nodeStr = username.ReplaceAllString(nodeStr, fmt.Sprintf("%s_%d_%d_%d_%s_%d_%d", param.Account, country, param.Province, param.City, session, param.SessionTime, param.SessionValidity))
		password := regexp.MustCompile(`password`)
		nodeStr = password.ReplaceAllString(nodeStr, param.Password)
		cmd = append(cmd, nodeStr)
	}
	return
}

func GetAccountExtractionNodes(param AccountParams, num int, uid, subUserId int64) (cmd []string, err error) {
	subuser, err := models.SubUser{
		Id:           subUserId,
		ParentUserID: uid,
	}.GetSubUserById()
	if err != nil || subuser.Id <= 0 || subuser.Status != 1 {
		err = fmt.Errorf("failed to get sub user")
		return
	}
	param.Account = subuser.Username
	param.Password = subuser.Password
	return MakeAccountExtractionNodes(param, num), nil
}
