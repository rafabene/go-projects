module github.com/rafabene/client-server-api/server

require github.com/rafabene/client-server-api/common v0.0.0

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	golang.org/x/text v0.23.0 // indirect
	gorm.io/driver/sqlite v1.5.7 // indirect
	gorm.io/gorm v1.25.12 // indirect
)

replace github.com/rafabene/client-server-api/common => ../common

go 1.23.6
