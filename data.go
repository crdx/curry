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
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func (self DataError) String() string {
	return fmt.Sprintf("[%d] %s: %s", self.Code, self.Type, self.Info)
}

type Data struct {
	Success    bool      `json:"success"`
	Historical bool      `json:"historical"`
	Base       string    `json:"base"`
	Date       string    `json:"date"`
	Timestamp  int       `json:"timestamp"`
	Rates      Rates     `json:"rates"`
	Error      DataError `json:"error"`
}

func (self Data) isValid() bool {
	return self.Success && self.Error.Code == 0 && len(self.Rates) > 0
}
