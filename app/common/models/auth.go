package models

// Account 账号表
type Account struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Status   bool   `gorm:"default:true" json:"status"`
	Roles    []Role `gorm:"many2many:sys_account_roles;" json:"roles,omitempty"`
}

func (Account) TableName() string {
	return "sys_accounts"
}

// Role 角色表
type Role struct {
	ID       int64     `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	Sort     int       `gorm:"default:0" json:"sort"`
	Status   bool      `gorm:"default:true" json:"status"`
	Remark   string    `gorm:"type:varchar(500)" json:"remark"`
	Menus    []Menu    `gorm:"many2many:sys_role_menus;" json:"menus,omitempty"`
	Accounts []Account `gorm:"many2many:sys_account_roles;" json:"-"`
}

func (Role) TableName() string {
	return "sys_roles"
}

// Menu 菜单表
type Menu struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	Pid       int64  `gorm:"default:0" json:"pid"`
	MenuType  int    `gorm:"not null" json:"menu_type"` // 1-目录 2-菜单 3-按钮
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Path      string `gorm:"type:varchar(255);default:''" json:"path"`
	Component string `gorm:"type:varchar(255);default:''" json:"component"`
	Icon      string `gorm:"type:varchar(100);default:''" json:"icon"`
	Sort      int    `gorm:"default:0" json:"sort"`
	ApiUrl    string `gorm:"type:varchar(255);default:''" json:"api_url"`
	ApiMethod string `gorm:"type:varchar(10);default:'*'" json:"api_method"`
	Visible   bool   `gorm:"default:true" json:"visible"`
	Status    bool   `gorm:"default:true" json:"status"`
	Remark    string `gorm:"type:varchar(500);default:''" json:"remark"`
	Roles     []Role `gorm:"many2many:sys_role_menus;" json:"-"`
}
