package models

// Announcement 系统公告表
type Announcement struct {
	ID          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:公告ID"`
	Title       string `json:"title" gorm:"type:varchar(255);comment:标题;not null"`
	Content     string `json:"content" gorm:"type:text;comment:内容;not null"`
	Enable      int8   `json:"enable" gorm:"type:tinyint;comment:是否启用(1是;0否);default:0;not null"`
	TopUp       int8   `json:"top_up" gorm:"type:tinyint;comment:是否置顶(1是;0否);default:0;not null"`
	Sort        int64  `json:"sort" gorm:"type:int;comment:排序(数字越小越靠前);default:1;not null"`
	CreatedTime int64  `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdatedTime int64  `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
}

var UserFiled = []string{
	"title",
	"content",
	"enable",
	"top_up",
	"sort",
	"updated_time",
}

type AnnouncementRes struct {
	ID          int64  `json:"id" `
	Title       string `json:"title" `
	Content     string `json:"content"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}

// TableName 设置表名
func (Announcement) TableName() string {
	return "ip_announcements"
}

// TableName 设置表名
func (AnnouncementRes) TableName() string {
	return "ip_announcements"
}

// SelectByPage 分页查询公告
func (model Announcement) SelectByPage(page int, pageSize int) (announcements []Announcement, total int64, err error) {
	err = DB.Model(&Announcement{}).
		Select("id, enable, top_up, sort, title, content, created_time, updated_time").
		Order("top_up DESC, sort ASC").
		Count(&total).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&announcements).Error
	return
}

// 添加公告信息
func (model Announcement) Create() (err error) {
	err = DB.Model(&Announcement{}).Create(&model).Error
	return
}

// 修改公告信息
func (model Announcement) Update() (err error) {
	err = DB.Model(&Announcement{}).Select(UserFiled).Where("id = ?", model.ID).Updates(model).Error
	return
}

// 删除公告信息
func (model Announcement) Delete() (err error) {
	err = DB.Model(&Announcement{}).Where("id = ?", model.ID).Delete(model).Error
	return
}
