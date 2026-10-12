// Package config resolves process-local settings. It never persists credentials.
package config

import "os"

const (
	BaseURL    = "http://pan.dlut.edu.cn/v1"
	AuthURL    = "http://202.118.66.144:8088/v1/auth/sign_in"
	NAOSIRoot  = "ak7an4chdvtq"
	NAOSIGroup = "ej7f8s176d8f"
)

func Value(flag, name, fallback string) string {
	if flag != "" {
		return flag
	}
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
