package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func fetchRawForDay(day string, accessKey string) (body []byte, err error) {
	body, err = util.HttpGetJson(
		// https is not supported on the free plan.
		fmt.Sprintf("http://api.exchangeratesapi.io/%s", day),
		map[string]string{"access_key": accessKey},
	)

	return
}

func getAccessKey() string {
	accessKey, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".config", ProgramName, "api_key"))
	check(err)
	return strings.TrimSpace(string(accessKey))
}

func main() {
	log.SetFlags(0)
	var opts Opts
	check(duckopt.Parse(getUsage(), "$0").Bind(&opts))
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

	raw, err := ratesCache.ReadBytes(func() []byte {
		body, err := fetchRawForDay(day, getAccessKey())
		check(err)
		return body
	})

	check(err)

	var data Data
	check(json.Unmarshal(raw, &data))

	if !data.isValid() {
		fmt.Println(col.Red("Unable to fetch rates"))
		log.Fatal(data.Error)
	}

	check(ratesCache.WriteBytes(raw))

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
