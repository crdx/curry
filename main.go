package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/crdx/curry/api"
	"github.com/crdx/curry/cache"
	"github.com/crdx/curry/util"

	"github.com/crdx/col"
	"github.com/crdx/duckopt"
)

const ProgramName = "curry"

func getUsage() string {
	return `
		Usage:
		    $0 raw
		    $0 list
		    $0 clear
		    $0 [options] <value>

		Convert currencies.

		Commands:
		    raw      Print raw json
		    clear    Clear the cache
		    list     List available currencies

		Options:
		    -f, --from TYPE    From this currency [default: GBP]
		    -t, --to TYPE      To this currency [default: GBP]
		    -q, --quiet        Show only the value
		    -C, --no-color     Disable colours
		    -h, --help         Show help
	`
}

type Opts struct {
	Clear        bool    `docopt:"clear"`
	Raw          bool    `docopt:"raw"`
	List         bool    `docopt:"list"`
	ValueFrom    float32 `docopt:"<value>"`
	CurrencyFrom string  `docopt:"--from"`
	CurrencyTo   string  `docopt:"--to"`
	Quiet        bool    `docopt:"--quiet"`
	NoColor      bool    `docopt:"--no-color"`
}

type Rates map[string]float32

func (self Rates) getSortedCurrencies() []string {
	var keys []string
	for key := range self {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type HistoricDataError struct {
	Code int
	Type string
	Info string
}

func (self HistoricDataError) String() string {
	return fmt.Sprintf("[%d] %s: %s", self.Code, self.Type, self.Info)
}

type HistoricData struct {
	Success    bool
	Historical bool
	Base       string
	Date       string
	Timestamp  int
	Rates      Rates
	Error      HistoricDataError
}

func (self HistoricData) isValid() bool {
	return self.Success && self.Error.Code == 0 && len(self.Rates) > 0
}

func getAccessKey() string {
	accessKey, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".config", ProgramName, "api_key"))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(accessKey))
}

func main() {
	log.SetFlags(0)
	var opts Opts
	if err := duckopt.Parse(getUsage(), "$0").Bind(&opts); err != nil {
		panic(err)
	}
	col.InitUnless(opts.NoColor)

	day := util.GetYesterday()
	ratesCache := cache.New(path.Join(os.Getenv("HOME"), ".cache", ProgramName, day+".json"))

	if opts.Clear {
		if ratesCache.Delete() {
			fmt.Println(col.Green("Cache file cleared"))
		} else {
			fmt.Println(col.Yellow("No cache file found"))
		}
		os.Exit(0)
	}

	raw := ratesCache.ReadBytes(func() []byte {
		return api.FetchRawForDay(day, getAccessKey())
	})

	var data HistoricData
	if err := json.Unmarshal(raw, &data); err != nil {
		log.Fatal(err)
	}

	if !data.isValid() {
		fmt.Println(col.Red("Unable to fetch rates"))
		log.Fatal(data.Error)
	}

	ratesCache.WriteBytes(raw)

	if opts.Raw {
		fmt.Printf("%s\n", raw)
		os.Exit(0)
	}

	if opts.List {
		for _, rate := range data.Rates.getSortedCurrencies() {
			fmt.Println(rate)
		}
		os.Exit(0)
	}

	if opts.ValueFrom == 0 {
		log.Fatal(col.Red("Supply a valid non-zero float"))
	}

	rateFrom := data.Rates[opts.CurrencyFrom]
	if rateFrom == 0 {
		log.Fatalf(
			col.Red("Source %s is not supported. Run \"%s list\" to see available currencies."),
			opts.CurrencyFrom,
			ProgramName,
		)
	}

	rateTo := data.Rates[opts.CurrencyTo]
	if rateTo == 0 {
		log.Fatalf(
			col.Red("Target %s is not supported. Run \"%s list\" to see available currencies."),
			opts.CurrencyTo,
			ProgramName,
		)
	}

	valueTo := opts.ValueFrom / (rateFrom / rateTo)

	if opts.Quiet {
		fmt.Printf("%.2f\n", valueTo)
	} else {
		fmt.Printf(
			"%.2f %s is %.2f %s (as of %s)\n",
			opts.ValueFrom,
			opts.CurrencyFrom,
			valueTo,
			opts.CurrencyTo,
			day,
		)
	}
}
