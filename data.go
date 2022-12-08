package main

import (
	"fmt"
	"sort"
)

type Rates map[string]float64

func (self Rates) getSortedCurrencies() []string {
	var keys []string
	for key := range self {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type DataError struct {
	Code int
	Type string
	Info string
}

func (self DataError) String() string {
	return fmt.Sprintf("[%d] %s: %s", self.Code, self.Type, self.Info)
}

type Data struct {
	Success    bool
	Historical bool
	Base       string
	Date       string
	Timestamp  int
	Rates      Rates
	Error      DataError
}

func (self Data) isValid() bool {
	return self.Success && self.Error.Code == 0 && len(self.Rates) > 0
}
