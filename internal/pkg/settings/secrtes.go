package settings

import "os"

func CheckSecretes(necessary []string) bool {
	for _, key := range necessary {
		_, ok := os.LookupEnv(key)
		if !ok {
			return false
		}
	}
	return true
}
