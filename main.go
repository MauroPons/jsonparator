package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli/v2"
)

var options *InputData
var fileRelativePathSource *os.File
var fileRelativePathError *os.File
var progressBar *ProgressBar
var arrayAllParamsSorted []string
var mapFileParams = make(map[string]*os.File)
var fieldErrorCounter = FieldErrorCounter{mapField: make(map[string]int)}

type InputData struct {
	Hosts                          []string
	Headers                        []string
	Exclude                        []string
	SeparateInFilesNotContainParam []string
	FilePathSource                 string
	FilePathError                  string
	BasePath                       string
	FileName                       string
	CallerScope                    string
	Velocity                       int
	FilePathTotalLines             int
	FilePathTotalLinesError        int
	Currency                       int
}

type FieldErrorCounter struct {
	mapField map[string]int
	mux      sync.Mutex
}

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "JsonParator"
	app.Usage = "Compares Api's Json responses"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"P"},
			Usage:   "Specifies the file from which to read targets. It should contain one column only with a rel RelativePath. eg: /v1/cards?query=123",
		},
		&cli.StringSliceFlag{
			Name:  "host",
			Usage: "Targeted hosts. Exactly 2 must be specified. eg: --host 'http://host1.com --host 'http://host2.com'",
		},
		&cli.StringSliceFlag{
			Name:    "header",
			Aliases: []string{"H"},
			Usage:   "Headers to be used in the http call",
		},
		&cli.IntFlag{
			Name:    "velocity",
			Aliases: []string{"V"},
			Value:   4,
			Usage:   "Set comparators velocity in K RPM",
		},
		&cli.StringSliceFlag{
			Name:    "exclude",
			Aliases: []string{"E"},
			Usage:   "Excludes a value from both json for the specified RelativePath. A RelativePath is a series of keys separated by a dot or #",
		},
		&cli.StringSliceFlag{
			Name:    "requestNotContainParam",
			Aliases: []string{"M"},
			Usage:   "Separate into files the requests that not contain parameter. At set this value the separateFileByParameters is turned on automatically#",
		},
	}
	app.Action = action
	return app
}

func action(contextClient *cli.Context) error {
	options = parseFlags(contextClient)

	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileRelativePathSource = openFile()
	defer fileRelativePathSource.Close()

	options.FilePathTotalLines = GetTotalLines(fileRelativePathSource)

	fileRelativePathError = createFileError()
	defer fileRelativePathError.Close()

	fmt.Println("Separating files by parameter's cuts")
	arrayAllParamsSorted = getAllParamsSorted(fileRelativePathSource)
	createFilesByParams(arrayAllParamsSorted, "source")
	createFilesByParams(arrayAllParamsSorted, "error")
	loadFileByParams(fileRelativePathSource, "source")

	progressBar = NewProgressBar()
	progressBar.Start()

	pipeline := NewPipeline()
	pipeline.Run(context)

	progressBar.Stop()

	summaryFieldError()
	separateFilesByParams()

	return nil
}

func separateFilesByParams() {
	fmt.Println("Summary param's cuts: ")
	dir := options.BasePath + "/table-summary.csv"
	fileSummary, _ := os.Create(dir)
	w := bufio.NewWriter(fileSummary)
	fmt.Fprintln(w, fmt.Sprintln("PARAMS,TOTAL,CORRECT,INCORRECT,%CORRECT,%INCORRECT"))
	matchTotal := int(math.Abs(float64(options.FilePathTotalLines - options.FilePathTotalLinesError)))
	percentMatch, percentError := getPercentValues(options.FilePathTotalLines, matchTotal, options.FilePathTotalLinesError)
	fmt.Fprintln(w, fmt.Sprintln("GENERAL,", options.FilePathTotalLines, ",", matchTotal, ",", options.FilePathTotalLinesError, ",", percentMatch, ",", percentError))
	w.Flush()

	for _, param := range arrayAllParamsSorted {
		total, match, error := getCountRowsByParams(param)
		percentMatch, percentError := getPercentValues(total, match, error)
		fmt.Fprintln(w, fmt.Sprintln(strings.ToUpper(param), ",", total, ",", match, ",", error, ",", percentMatch, ",", percentError))
		w.Flush()
	}

	table, _ := tablewriter.NewCSV(os.Stdout, dir, true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})
	table.Render()

}

func summaryFieldError() {
	var arrayFieldError []string
	for key := range fieldErrorCounter.mapField {
		arrayFieldError = append(arrayFieldError, key)
	}
	sort.Strings(arrayFieldError)
	dir := options.BasePath + "/common-error-summary.csv"
	fileSummary, err := os.Create(dir)
	fi, err := fileSummary.Stat()
	if err == nil && fi.Size() > 0 {
		fmt.Println("Summary field's error: ")
		w := bufio.NewWriter(fileSummary)
		fmt.Fprintln(w, fmt.Sprintln("ATTRIBUTES,INCORRECT"))
		for index := range arrayFieldError {
			total := fieldErrorCounter.mapField[arrayFieldError[index]]
			fmt.Fprintln(w, fmt.Sprintln(arrayFieldError[index], ",", total))
			w.Flush()
		}
		table, _ := tablewriter.NewCSV(os.Stdout, dir, true)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
		table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor})
		table.Render()
	}
}

func getPercentValues(total int, match int, error int) (string, string) {
	percentMatch := float64(0)
	percentError := float64(0)
	if total != 0 {
		percentMatch = float64(match) * 100 / float64(total)
		percentError = float64(error) * 100 / float64(total)
	}
	return fmt.Sprintf("%.2f", percentMatch), fmt.Sprintf("%.2f", percentError)
}

func getCountRowsByParams(param string) (int, int, int) {
	total := GetTotalLines(mapFileParams["source-"+param])
	error := GetTotalLines(mapFileParams["error-"+param])
	match := int(math.Abs(float64(total - error)))
	return total, match, error
}

func loadFileByParams(file *os.File, s string) {
	_, _ = file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		relativePath := scanner.Text()
		addRelativePathToFileParam(strings.Trim(relativePath, "\""), "source")
	}
}

func createFilesByParams(array []string, prefix string) {
	for _, param := range array {
		nameFile := prefix + "-" + param
		relativePaths := options.BasePath + "relative-paths/"
		_ = os.Mkdir(relativePaths, 0777)
		pathFileParam := relativePaths + nameFile + ".txt"
		file, _ := os.Create(pathFileParam)
		mapFileParams[nameFile] = file
	}
}

func getAllParamsSorted(file *os.File) []string {
	var arrayParamsSorted []string
	scanner := bufio.NewScanner(file)
	_, _ = file.Seek(0, 0)
	mapParams := make(map[string]bool)
	for scanner.Scan() {
		relativePath := scanner.Text()
		urlParse, err := url.Parse(relativePath)
		if err != nil {
			panic(err)
		}
		mapValues, _ := url.ParseQuery(urlParse.RawQuery)
		for key := range mapValues {
			mapParams[key] = true
		}
	}
	for key := range mapParams {
		arrayParamsSorted = append(arrayParamsSorted, key)
	}

	for _, param := range options.SeparateInFilesNotContainParam {
		arrayParamsSorted = append(arrayParamsSorted, param+"-no")
	}
	sort.Strings(arrayParamsSorted)
	return arrayParamsSorted
}

func GetTotalLines(file *os.File) int {
	count := 0
	_, _ = file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count
}

func createFileError() *os.File {
	fmt.Println("Creating error file in", options.FilePathError)

	file, err := os.Create(options.FilePathError)
	if err != nil {
		panic("Can not read file: " + options.FilePathError)
	}
	return file
}

func openFile() *os.File {
	fmt.Println("Reading file", options.FilePathSource)
	file, err := os.Open(options.FilePathSource)
	if err != nil {
		panic("Can not read file: " + options.FilePathSource)
	}
	return file
}

func parseFlags(context *cli.Context) *InputData {
	fmt.Println("\nGetting data...")
	opts := InputData{}
	opts.Currency = 700
	opts.Hosts = context.StringSlice("host")
	if len(opts.Hosts) != 2 {
		panic("Invalid numbers the hosts provided")
	}
	opts.FilePathSource = formatPath(context.String("path"))
	opts.Headers = context.StringSlice("header")
	opts.Velocity = context.Int("velocity") * 1000
	opts.Exclude = context.StringSlice("exclude")
	opts.SeparateInFilesNotContainParam = context.StringSlice("requestNotContainParam")

	dir, fileName := filepath.Split(opts.FilePathSource)
	path := dir + time.Now().Format("20060102150405")
	_ = os.Mkdir(path, 0777)
	opts.CallerScope = getCallerScope(opts)
	opts.BasePath = path + "/" + opts.CallerScope + "/"
	opts.FileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	_ = os.Mkdir(opts.BasePath, 0777)
	opts.FilePathError = opts.BasePath + opts.FileName + ".error"
	return &opts
}

func formatPath(relativePath string) string {
	var finalPath = strings.Trim(strings.ReplaceAll(relativePath, "'", ""), "")
	firstCharacter := finalPath[0:1]
	if firstCharacter != "/" {
		finalPath = "/" + finalPath
	}
	return finalPath
}

func getCallerScope(options InputData) string {
	callerScopeTemp := "no-caller-scope"
	if options.Headers != nil {
		for _, header := range options.Headers {
			if header == "" {
				continue
			}
			h := strings.Split(header, ":")
			if len(h) != 2 {
				panic("Invalid header")
			} else if h[0] == "X-Caller-Scopes" && h[1] != "" {
				callerScopeTemp = h[1]
			}
		}
	}
	return callerScopeTemp
}

func addRelativePathToFileParam(relativePath string, prefix string) {
	urlParse, err := url.Parse(relativePath)
	if err != nil {
		panic(err)
	}
	mapValues, _ := url.ParseQuery(urlParse.RawQuery)
	for key := range mapValues {
		file := mapFileParams[prefix+"-"+key]
		w := bufio.NewWriter(file)
		fmt.Fprintln(w, relativePath)
		_ = w.Flush()
	}

	for _, param := range options.SeparateInFilesNotContainParam {
		value := mapValues[param]
		if value == nil {
			key := prefix + "-" + param + "-no"
			file := mapFileParams[key]
			w := bufio.NewWriter(file)
			fmt.Fprintln(w, relativePath)
			_ = w.Flush()
		}
	}
}

func (c *FieldErrorCounter) Add(key string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.mapField[key]++
	c.mux.Unlock()
}
