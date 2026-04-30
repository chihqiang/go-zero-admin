

generate-api:
	@echo "格式化 API 文件"
	goctl api format --dir .
	@echo "根据 admin.api 生成 Go 代码"
	goctl api go -api app/admin/admin.api -dir app/admin -style gozero