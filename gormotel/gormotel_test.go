package gormotel

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestTracePlugin(t *testing.T) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// register trace plugin
	err = db.Use(&TracePlugin{})
	if err != nil {
		t.Fatal(err)
	}
}
