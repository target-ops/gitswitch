package cmd

import "os"

// osEnv is a thin indirection so tests can stub the environment without
// pulling `os` into every file under cmd/.
func osEnv(key string) string { return os.Getenv(key) }
