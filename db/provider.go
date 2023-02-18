package db

import (
	"sync"
)

type Provider func() Database

func GetProvider() Provider {
	var v Database
	var once sync.Once

	return func() Database {
		once.Do(func() {
			v = createDatabase()
		})
		return v
	}
}
