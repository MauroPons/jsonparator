package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CreateAndOpenDashboardInBrowser() {

	basePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	dashboardFilePath := basePath + "/static/index.html"
	dashboardFile, err := os.Open(dashboardFilePath)
	if err != nil {
		panic("Error getting index.html")
	}
	finalDashboardFilePath := options.BasePath +  "index.html"
	dashboardFileCopy, err := os.Create(finalDashboardFilePath)
	if err != nil {
		fmt.Println(err)
	}
	w := bufio.NewWriter(dashboardFileCopy)
	scanner := bufio.NewScanner(dashboardFile)
	count := 0
	for scanner.Scan() {
		count++
		line := scanner.Text()
		if strings.Contains(line, "let basePath = null;") {
			valueTemp := fmt.Sprintf( "			let basePath = '%s';", basePath)
			line = string(valueTemp)
		}
		if strings.Contains(line, "let filterParams = null;") {
			valueTemp := fmt.Sprintf(`			let filterParams = {"scope1": "%s", "scope2": "%s", "relativePath": "%s"}`, options.Hosts[0], options.Hosts[1], options.FilePathSource)
			line = valueTemp
		}
		fmt.Fprintln(w, line)
		_ = w.Flush()
	}

	url := "file://" + finalDashboardFilePath
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("x-www-browser", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
