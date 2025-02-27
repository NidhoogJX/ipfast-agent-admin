package iplocation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"ipfast_server/pkg/util/log"
	"math/bits"
	"strconv"

	"os"
	//protoc --go_out=. cidr.proto

	"github.com/kentik/patricia"
	"github.com/kentik/patricia/string_tree"
	"github.com/spf13/viper"
)

/*
Geo地理位置信息

	CIDR: IP地址段
	CountryISO2: 国家二字母代码 例如 CN 中国
	CountryISO3: 国家三字母代码 例如 CHN 中国
	ProvinceCode: 第一级行政区划代码 例如 BJ 北京
	City: 城市名称 例如 Beijing 北京
*/
type GeoInfo struct {
	CountryISO3 string
	CountryName string
}

/*
IP库配置

	Language: 语言 zh-cn,de,en,es,fr,ja,pt-br,ru
*/
// var language = "en"

/*
全局IP前缀树
*/
var ipTree string_tree.TreeV4

/*
全局地理位置信息
*/
var locationMap map[int64]GeoInfo

/*
初始化自动加载IP库到内存
*/
func Setup() error {
	log.Info("初始化IP库中...")
	ipTree = *string_tree.NewTreeV4()
	locationMap = make(map[int64]GeoInfo)
	err := loadIpTree(viper.GetString("ipdata.path"))
	if err != nil {
		log.Error("IP库初始化失败:%v", err)
		return fmt.Errorf("IP库初始化失败:%v", err)
	}
	log.Info("初始化IP库成功")
	return nil
}

func GetGeoInfo(cityId int64) (geoInfo GeoInfo) {
	return locationMap[cityId]
}

/*
FindIpLoction 根据IP查找IP归属地

	param:
		ip: IP地址
	return:
		ipinfo: IP归属地信息
		err: 错误信息
*/
func FindIpLoction(ip string) (ipinfo GeoInfo, err error) {
	a, _, err := patricia.ParseIPFromString(ip)
	if err != nil {
		return
	}
	if a == nil {
		err = fmt.Errorf("invalid IP address: %s", ip)
		return
	}
	ok, v := ipTree.FindDeepestTag(*a)
	if ok {
		err = json.Unmarshal([]byte(v), &ipinfo)
		if err != nil {
			err = fmt.Errorf("failed to unmarshal json data %s", v)
			return
		}
	} else {
		err = fmt.Errorf("failed to find IP %s in the tree", ip)
	}
	return
}

/*
读取CSV文件

	param:
		filepath: 文件路径
	return:
		reader: 读取器
		err: 错误信息
*/
func loadCsvFile(filepath string) (reader *csv.Reader, file *os.File, err error) {
	file, err = os.Open(filepath)
	if err != nil {
		return
	}
	reader = csv.NewReader(file) // 读取CSV文件
	reader.Comma = ','           // 设置分隔符为制表符
	_, _ = reader.Read()         // 跳过第一行
	return
}

/*
加载ipv4.csv文件到内存
*/
func loadIpTree(filePath string) (err error) {
	reader, file, err := loadCsvFile(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	defer file.Close()
	if err != nil {
		return
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if record[1] == "" {
			continue
		}
		ipstart, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return fmt.Errorf("ip库文件格式错误 ipstart:%v,err :%v", record[1], err)
		}
		ipend, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return fmt.Errorf("ip库文件格式错误 ipend:%v", record[2])
		}
		jsonData, err := json.Marshal(GeoInfo{
			CountryISO3: record[2],
			CountryName: record[3],
		})
		if err != nil {
			return err
		}
		cidr, err := ipRangeToCIDR(ipstart, ipend)
		if err != nil {
			return err
		}
		for _, cidrStr := range cidr {
			ip, _, err := patricia.ParseIPFromString(cidrStr)
			if err != nil {
				return fmt.Errorf("无法解析 CIDR: %v", cidrStr)
			}
			ok, _ := ipTree.Add(*ip, string(jsonData), nil)
			if !ok {
				return fmt.Errorf("failed to add jsonData %s to the tree", record[0])
			}
		}
	}
	return
}

func ipRangeToCIDR(ipStart, ipEnd int64) ([]string, error) {

	var cidrs []string
	for ipStart <= ipEnd {
		maxSize := 32 - bits.TrailingZeros32(uint32(ipStart))
		maxDiff := 32 - bits.Len32(uint32(ipEnd-ipStart+1))
		size := maxSize
		if maxDiff > maxSize {
			size = maxDiff
		}
		cidr := fmt.Sprintf("%s/%d", long2ip(uint64(ipStart)), size)
		cidrs = append(cidrs, cidr)
		ipStart += 1 << (32 - size)
	}

	return cidrs, nil
}

func long2ip(ip uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d", (ip>>24)&0xFF, (ip>>16)&0xFF, (ip>>8)&0xFF, ip&0xFF)
}
