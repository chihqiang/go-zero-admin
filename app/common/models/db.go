package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&Account{},
		&Role{},
		&Menu{},
	)
}

func MigrateData(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	if err := db.Model(&Role{}).Count(&count).Error; err != nil {
		return fmt.Errorf("检查数据失败: %w", err)
	}
	if count > 0 {
		return nil
	}

	// 1. 创建超级管理员角色
	adminRole := Role{
		Name:   "超级管理员",
		Sort:   1,
		Status: true,
		Remark: "超级管理员角色，拥有所有权限",
	}
	if err := db.Create(&adminRole).Error; err != nil {
		return fmt.Errorf("创建超级管理员角色失败: %w", err)
	}

	// 2. 创建菜单结构
	menus := []Menu{
		// 仪表盘目录
		{Name: "仪表盘", MenuType: 1, Path: "/admin/dashboard", Component: "admin/dashboard/page", Icon: "Dashboard", Sort: 1, Pid: 0, Remark: "仪表盘目录"},
		// 仪表盘菜单
		{Name: "数据概览", MenuType: 2, Path: "/admin/dashboard", Component: "admin/dashboard/page", Icon: "Dashboard", Sort: 1, Pid: 0, Remark: "仪表盘页面"},
		// 系统管理目录
		{Name: "系统管理", MenuType: 1, Path: "/admin/sys", Component: "admin/sys/page", Icon: "Setting", Sort: 2, Pid: 0, Remark: "系统管理目录"},
	}

	// 保存一级菜单获取ID
	if err := db.Create(&menus).Error; err != nil {
		return fmt.Errorf("创建一级菜单失败: %w", err)
	}

	// 获取系统管理目录ID
	var sysMenu Menu
	if err := db.Where("name = ?", "系统管理").First(&sysMenu).Error; err != nil {
		return fmt.Errorf("获取系统管理菜单失败: %w", err)
	}

	// 添加二级菜单
	level2Menus := []Menu{
		{Name: "账号管理", MenuType: 2, Path: "/admin/sys/account", Component: "admin/sys/account/page", Icon: "Users", Sort: 1, Pid: sysMenu.ID, Remark: "账号管理菜单"},
		{Name: "角色管理", MenuType: 2, Path: "/admin/sys/roles", Component: "admin/sys/roles/page", Icon: "UserRole", Sort: 2, Pid: sysMenu.ID, Remark: "角色管理菜单"},
		{Name: "菜单管理", MenuType: 2, Path: "/admin/sys/menu", Component: "admin/sys/menu/page", Icon: "Menu", Sort: 3, Pid: sysMenu.ID, Remark: "菜单管理菜单"},
	}
	if err := db.Create(&level2Menus).Error; err != nil {
		return fmt.Errorf("创建二级菜单失败: %w", err)
	}

	// 获取账号管理和角色管理菜单ID
	var accountMenu, roleMenu, menuMenu Menu
	if err := db.Where("name = ?", "账号管理").First(&accountMenu).Error; err != nil {
		return fmt.Errorf("获取账号管理菜单失败: %w", err)
	}
	if err := db.Where("name = ?", "角色管理").First(&roleMenu).Error; err != nil {
		return fmt.Errorf("获取角色管理菜单失败: %w", err)
	}
	if err := db.Where("name = ?", "菜单管理").First(&menuMenu).Error; err != nil {
		return fmt.Errorf("获取菜单管理菜单失败: %w", err)
	}

	// 添加按钮
	buttons := []Menu{
		// 账号管理按钮
		{Name: "账号列表", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 1, ApiUrl: "/api/v1/account/list", ApiMethod: "GET", Pid: accountMenu.ID, Remark: "获取账号列表"},
		{Name: "账号详情", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 2, ApiUrl: "/api/v1/account/detail", ApiMethod: "GET", Pid: accountMenu.ID, Remark: "获取账号详情"},
		{Name: "创建账号", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 3, ApiUrl: "/api/v1/account/create", ApiMethod: "POST", Pid: accountMenu.ID, Remark: "创建账号"},
		{Name: "更新账号", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 4, ApiUrl: "/api/v1/account/update", ApiMethod: "PUT", Pid: accountMenu.ID, Remark: "更新账号"},
		{Name: "删除账号", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 5, ApiUrl: "/api/v1/account/delete", ApiMethod: "DELETE", Pid: accountMenu.ID, Remark: "删除账号"},
		// 角色管理按钮
		{Name: "角色列表", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 1, ApiUrl: "/api/v1/role/list", ApiMethod: "GET", Pid: roleMenu.ID, Remark: "获取角色列表"},
		{Name: "角色详情", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 2, ApiUrl: "/api/v1/role/detail", ApiMethod: "GET", Pid: roleMenu.ID, Remark: "获取角色详情"},
		{Name: "创建角色", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 3, ApiUrl: "/api/v1/role/create", ApiMethod: "POST", Pid: roleMenu.ID, Remark: "创建角色"},
		{Name: "更新角色", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 4, ApiUrl: "/api/v1/role/update", ApiMethod: "PUT", Pid: roleMenu.ID, Remark: "更新角色"},
		{Name: "删除角色", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 5, ApiUrl: "/api/v1/role/delete", ApiMethod: "DELETE", Pid: roleMenu.ID, Remark: "删除角色"},
		{Name: "所有角色", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 6, ApiUrl: "/api/v1/role/all", ApiMethod: "GET", Pid: roleMenu.ID, Remark: "获取所有角色列表"},
		{Name: "关联菜单", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 7, ApiUrl: "/api/v1/role/associate-menus", ApiMethod: "POST", Pid: roleMenu.ID, Remark: "关联角色和菜单"},
		// 菜单管理按钮
		{Name: "菜单列表", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 1, ApiUrl: "/api/v1/menu/list", ApiMethod: "GET", Pid: menuMenu.ID, Remark: "获取菜单列表"},
		{Name: "所有菜单", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 2, ApiUrl: "/api/v1/menu/all", ApiMethod: "GET", Pid: menuMenu.ID, Remark: "获取所有菜单列表"},
		{Name: "菜单详情", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 3, ApiUrl: "/api/v1/menu/detail", ApiMethod: "GET", Pid: menuMenu.ID, Remark: "获取菜单详情"},
		{Name: "创建菜单", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 4, ApiUrl: "/api/v1/menu/create", ApiMethod: "POST", Pid: menuMenu.ID, Remark: "创建菜单"},
		{Name: "更新菜单", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 5, ApiUrl: "/api/v1/menu/update", ApiMethod: "PUT", Pid: menuMenu.ID, Remark: "更新菜单"},
		{Name: "删除菜单", MenuType: 3, Path: "", Component: "", Icon: "", Sort: 6, ApiUrl: "/api/v1/menu/delete", ApiMethod: "DELETE", Pid: menuMenu.ID, Remark: "删除菜单"},
	}
	if err := db.Create(&buttons).Error; err != nil {
		return fmt.Errorf("创建按钮菜单失败: %w", err)
	}

	// 3. 将所有菜单分配给超级管理员角色
	var allMenus []Menu
	if err := db.Find(&allMenus).Error; err != nil {
		return fmt.Errorf("获取所有菜单失败: %w", err)
	}
	for _, m := range allMenus {
		db.Model(&adminRole).Association("Menus").Append(&m)
	}

	// 4. 创建超级管理员账户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	adminAccount := Account{
		Name:     "超级管理员",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Status:   true,
	}
	if err := db.Create(&adminAccount).Error; err != nil {
		return fmt.Errorf("创建超级管理员账户失败: %w", err)
	}

	// 5. 将超级管理员角色分配给超级管理员账户
	db.Model(&adminAccount).Association("Roles").Append(&adminRole)

	return nil
}
