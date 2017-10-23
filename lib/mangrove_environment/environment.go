package mangrove_environment

import (
	"os"
)

const defaultEnvironment = "development"
const environmentVariableKey = "ENV"

// Get 実行環境の取得
func Get() string {
	env := os.Getenv(environmentVariableKey)
	if env == "" {
		return defaultEnvironment
	}

	return env
}

// Set 実行環境を設定
func Set(env string) {
	os.Setenv(environmentVariableKey, env)
}
