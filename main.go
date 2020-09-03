package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var options 			*InputData
var fileSource 			*os.File
var fileError 			*os.File
var progressBar			*ProgressBar

type InputData struct {
	FilePathSource     string
	FilePathTotalLines int
	FilePathError      string
	BasePath           string
	FileName           string
	Velocity           int
	Hosts              []string
	Headers            []string
	ParametersCutting  []string
	Exclude            string
	CallerScope        string
	Currency 		   int
}

func main() {

	port := os.Getenv("PORT") // To build use WJC_PORT instead of PORT
	if port == "" {
		port = "8080" // Dejar en 8080 sino no sube en el scope
	}
	//router := mux.NewRouter()
	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	//router.HandleFunc("/", IndexHandler)
	//fmt.Println("Starting up on " + port)
	//http.ListenAndServe(":" + port, router)

	app := newApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

//func IndexHandler(w http.ResponseWriter, r *http.Request){
//	tmpl := template.Must(template.ParseFiles("static/index.html"))
//	tmpl.Execute(w, nil)
//}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "JsonParator"
	app.Usage = "Compares Api's Json responses"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "path",
			Usage: "specifies the file from which to read targets. It should contain one column only with a rel RelativePath. eg: /v1/cards?query=123",
		},
		&cli.StringSliceFlag{
			Name:  "host",
			Usage: "targeted hosts. Exactly 2 must be specified. eg: --host 'http://host1.com --host 'http://host2.com'",
		},
		&cli.StringSliceFlag{
			Name:    "header",
			Aliases: []string{"H"},
			Usage:   "headers to be used in the http call",
		},
		&cli.IntFlag{
			Name:    "velocity",
			Aliases: []string{"k"},
			Value:   4,
			Usage:   "Set comparators velocity in K RPM",
		},
		&cli.StringFlag{
			Name:  "exclude",
			Usage: "excludes a value from both json for the specified RelativePath. A RelativePath is a series of keys separated by a dot or #",
			Value: "results.#.payer_costs.#.payment_method_option_id",
		},
		&cli.StringFlag{
			Name:  "parametersCutting",
			Usage: "check is request contains the params parametrized #",
			Value: "caller_id,display_filtered,differential_pricing_id,bins,public_key",
		},
	}

	app.Action = action
	return app
}

func action(contextClient *cli.Context) error {
	options = parseFlags(contextClient)

	CreateAndOpenDashboardInBrowser()

	context, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	fileSource = openFile()
	defer fileSource.Close()

	fileError = createFileError()
	defer fileError.Close()

	calculateTotalLinesOfSourceFile()

	progressBar = NewProgressBar()
	progressBar.Start()

	pipeline := NewPipeline()
	pipeline.Run(context)

	progressBar.Stop()

	return nil
}

func calculateTotalLinesOfSourceFile() {
	out, err0 := exec.Command("wc", "-l", options.FilePathSource).Output()
	if err0 != nil {
		panic("Error getting total line of file: " + options.FilePathSource + "\n")
	}
	outString := strings.Trim(string(out), " ")
	valueString := strings.Split(outString, " ")[0]
	var total, err1 = strconv.ParseUint(valueString, 10, 32)
	if err1 != nil {
		panic("Error getting total line of file: " + options.FilePathSource + "\n")
	}
	options.FilePathTotalLines = int(total)
}

func createFileError() *os.File {
	fmt.Println("Creating error file in %s", options.FilePathError)
	file, err := os.Create(options.FilePathError)
	if err != nil {
		panic("Can not read file: " + options.FilePathError)
	}
	return file
}

func openFile() *os.File {
	fmt.Println("Reading file from %s", options.FilePathSource)
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
	opts.Exclude = context.String("exclude")
	opts.ParametersCutting = strings.Split(context.String("parametersCutting"), ",")
	dir, fileName := filepath.Split(opts.FilePathSource)
	path := dir + time.Now().Format("20060102150405")
	_ = os.Mkdir(path , 0777)
	opts.CallerScope = getCallerScope(opts)
	opts.BasePath = path + "/" + opts.CallerScope + "/"
	opts.FileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	_ = os.Mkdir(opts.BasePath, 0777)
	opts.FilePathError = opts.BasePath + opts.FileName + ".error"
	return &opts
}

func formatPath(relativePath string) string {
	var finalPath = strings.Trim(strings.ReplaceAll(relativePath, "'", ""), "")
	firstCharacter :=  finalPath[0:1]
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
			}else if h[0] == "X-Caller-Scopes"{
				callerScopeTemp = h[1]
			}
		}
	}
	return callerScopeTemp
}
