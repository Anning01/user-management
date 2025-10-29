module github.com/Anning01/user-management

go 1.24.4

// Web框架
//go get -u github.com/gin-gonic/gin

// ORM
//go get -u gorm.io/gorm
//go get -u gorm.io/driver/mysql

// 环境变量处理
//go get -u github.com/spf13/viper

// JWT认证
//go get -u github.com/golang-jwt/jwt/v5

// 密码加密
//go get -u golang.org/x/crypto

// 日志
//go get -u github.com/sirupsen/logrus

// 数据库迁移
//go get -u github.com/go-gormigrate/gormigrate/v2

// 数据验证
//go get -u github.com/go-playground/validator/v10

require (
	github.com/spf13/viper v1.21.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
