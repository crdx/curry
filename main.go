package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"crdx.org/curry/internal/cache"
	"crdx.org/curry/internal/util"

	"crdx.org/col"
	"crdx.org/duckopt"
)

const ProgramName = "curry"

func getUsage() string {
	return `
		Usage:
		    $0 raw
		    $0 ls
		    $0 clean
		    $0 [options] <value>

		Convert currencies.

		Commands:
		    raw      Print raw json
		    ls       List available currencies
		    clean    Clean cache

		Options:
		    -f, --from TYPE    From this currency [default: GBP]
		    -t, --to TYPE      To this currency [default: GBP]
		    -q, --quiet        Show only the value
		    -C, --no-color     Disable colours
		    -h, --help         Show help
	`
}

type Opts struct {
	Clean        bool   `docopt:"clean"`
	Raw          bool   `docopt:"raw"`
	List         bool   `docopt:"ls"`
	ValueFrom    string `docopt:"<value>"`
	CurrencyFrom string `docopt:"--from"`
	CurrencyTo   string `docopt:"--to"`
	Quiet        bool   `docopt:"--quiet"`
	NoColor      bool   `docopt:"--no-color"`
}

func check(err error) {
	if err != nil {
		log.Fatal(col.Red("Error: " + err.Error()))
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

	if opts.Clean {
		if ratesCache.Delete() {
			fmt.Println(col.Green("Cache file removed"))
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
		log.Fatal(col.Red("Error: unable to fetch rates: " + data.Error.String()))
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

	valueFrom, err := strconv.ParseFloat(opts.ValueFrom, 64)
	if err != nil {
		log.Fatalf(col.Red("Error: %s is not a number"), opts.ValueFrom)
	}

	if valueFrom == 0 {
		log.Fatal(col.Red("Error: cannot convert from 0"))
	}

	rateFrom := data.Rates[opts.CurrencyFrom]
	if rateFrom == 0 {
		log.Fatalf(
			col.Red("Error: source %s is not supported — run \"%s ls\" to see available currencies"),
			opts.CurrencyFrom,
			ProgramName,
		)
	}

	rateTo := data.Rates[opts.CurrencyTo]
	if rateTo == 0 {
		log.Fatalf(
			col.Red("Error: target %s is not supported — run \"%s ls\" to see available currencies"),
			opts.CurrencyTo,
			ProgramName,
		)
	}

	valueTo := valueFrom / (rateFrom / rateTo)

	if opts.Quiet {
		fmt.Printf("%.2f\n", valueTo)
	} else {
		fmt.Printf(
			"%.2f %s is %.2f %s (as of %s)\n",
			valueFrom,
			opts.CurrencyFrom,
			valueTo,
			opts.CurrencyTo,
			day,
		)
	}
}
