package config

import (
	"log"
	"os" 
	"strconv" 
	"fmt" 
)

type AppConfig struct {
	AppPort        string 
	DBDSN          string 
	JWTAccessKey   string 
	JWTRefreshKey  string
	AccessTTLMin   int 
	RefreshTTLDays int 
	SMTPHost       string 
	SMTPPort       int 
	SMTPUser       string 
	SMTPPass       string 
	FromEmail      string 
	Env            string 
}

var C AppConfig

func Init() {
	C = AppConfig{
		AppPort: getenv("APP_PORT", "8007"),
		DBDSN: getenv("DB_DSN", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)),
		JWTAccessKey:   must("JWT_ACCESS_SECRET"),
		JWTRefreshKey:  must("JWT_REFRESH_SECRET"),
		AccessTTLMin:   atoi(getenv("ACCESS_TTL_MIN", "15")),
		RefreshTTLDays: atoi(getenv("REFRESH_TTL_DAYS", "7")),
		SMTPHost:       getenv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       atoi(getenv("SMTP_PORT", "587")),
		SMTPUser:       must("SMTP_USER"),
		SMTPPass:       must("SMTP_PASS"),
		FromEmail:      getenv("FROM_EMAIL", "myappondev@gmail.com"),
		Env:            getenv("APP_ENV", "dev"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env: %s", k)
	}
	return v
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}