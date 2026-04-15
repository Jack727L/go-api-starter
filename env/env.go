package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppEnvironment string

const (
	EnvironmentLocalDev AppEnvironment = "local-dev"
	EnvironmentDev      AppEnvironment = "dev"
	EnvironmentStaging  AppEnvironment = "staging"
	EnvironmentProd     AppEnvironment = "prod"
	EnvironmentTest     AppEnvironment = "test"
)

func IsLocalDevMode() bool { return GetApplicationEnv() == EnvironmentLocalDev }
func IsDevMode() bool      { return GetApplicationEnv() == EnvironmentDev }
func IsStagingMode() bool  { return GetApplicationEnv() == EnvironmentStaging }
func IsProdMode() bool     { return GetApplicationEnv() == EnvironmentProd }
func IsTestMode() bool     { return GetApplicationEnv() == EnvironmentTest }

// IsAsyncMode reports whether async Redis job processing is enabled.
// Only relevant when IsTestMode() is true (tests run sync by default).
func IsAsyncMode() bool {
	return os.Getenv("ASYNC_MODE") == "true"
}

func GetApplicationEnv() AppEnvironment {
	_ = godotenv.Load()
	envValue := os.Getenv("BACKEND_ENV")
	switch AppEnvironment(envValue) {
	case EnvironmentLocalDev, EnvironmentDev, EnvironmentStaging, EnvironmentProd, EnvironmentTest:
		return AppEnvironment(envValue)
	default:
		log.Fatalf("Invalid BACKEND_ENV value: %q (must be one of: local-dev, dev, staging, prod, test)", envValue)
	}
	return EnvironmentLocalDev // unreachable
}
