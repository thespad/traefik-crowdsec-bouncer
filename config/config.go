package config

import (
  "log"
  "os"
)

// OptionalEnv returns the value of an environment variable, or the provided default value if the variable is not set.
func OptionalEnv(varName string, optional string) string {
  envVar := os.Getenv(varName)
  if envVar == "" {
    return optional
  }
  return envVar
}

// NullableEnv returns the value of an environment variable. If the variable is not set, an empty string is returned.
func NullableEnv(varName string) string {
  envVar := os.Getenv(varName)
  return envVar
}

// RequiredEnv returns the value of an environment variable. If the variable is not set, an error is returned.
func RequiredEnv(varName string) string {
  envVar := os.Getenv(varName)
  if envVar == "" {
    log.Fatalf("The required env var %s is not provided. Exiting", varName)
  }
  return envVar
}

// ExpectedEnv returns the value of an environment variable if it is one of the expected values. If the variable is not set, an error is returned.
func ExpectedEnv(varName string, expected []string) string {
  envVar := RequiredEnv(varName)
  if !contains(expected, envVar) {
    log.Fatalf("The value for env var %s is not expected. Expected values are %v", varName, expected)
  }
  return envVar
}

func contains(source []string, target string) bool {
  for _, a := range source {
    if a == target {
      return true
    }
  }
  return false
}
