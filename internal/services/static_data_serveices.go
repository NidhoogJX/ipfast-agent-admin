package services

import "ipfast_server/internal/db/models"

// 分页查询 用户的静态IP
func StaticIPPaginate(page, pageSize int, ip, addressStr string, status int8, country_id, uid int64) (results []models.IpRecordByCountry, total int64, err error) {
	return models.IpRecord{}.Paginate(page, pageSize, ip, addressStr, status, country_id, uid, int8(2))
}

// 分页查询 用户的静态IP
func DataIPPaginate(page, pageSize int, ip, port string, status int8, country_id, uid int64) (results []models.IpRecordByCountry, total int64, err error) {
	return models.IpRecord{}.Paginate(page, pageSize, ip, port, status, country_id, uid, int8(3))
}
