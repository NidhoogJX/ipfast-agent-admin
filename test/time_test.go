package test

import (
	"fmt"
	"testing"
)

// 测试时间格式转换
func TestTimeTransform(t *testing.T) {
	// var createdTime int64 = 1729069564
	// timeObj := time.Unix(createdTime, 0) //将时间戳转为时间格式
	// fmt.Println(timeObj)
	// year := timeObj.Year()     //年
	// month := timeObj.Month()   //月
	// day := timeObj.Day()       //日
	// hour := timeObj.Hour()     //小时
	// minute := timeObj.Minute() //分钟
	// second := timeObj.Second() //秒
	// time, err := fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	// if err != nil {
	// 	log.Printf("格式化时间失败: %v", err)
	// }
	// fmt.Println("输出格式化后的时间%w", time) // 输出格式化后的时间
	// log.Printf("时间格式:%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)

	// a := timeObj.Format("2006-01-02 15:04:05\n\n")         // 格式化当前时间
	// log.Printf("\x20\x20\x20\x20\x20\x20格式化后的时间: %s\n", a) // 输出格式化后的时间
	var a float64 = 1526.000000
	res := fmt.Sprintf("%.2f", a) // 输出格式化后的时间
	// res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", a), 64)
	fmt.Println("", res) // 输出格式化后的时间
}
