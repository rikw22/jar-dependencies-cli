package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Dependency struct {
	Name     string `json:"Name"`
	Version  string `json:"Version"`
	FullName string `json:"FullName"`
}

// Example: dom4j-2.1.3
var regex1 = regexp.MustCompile(`(.*)(?:-)([0-9]+\.[0-9]+(\.[0-9]+)?(\.[0-9]+)?)$`)

// Examples:
// 	- hibernate-core-5.4.32.Final
// 	- javassist-3.27.0-GA
var regex2 = regexp.MustCompile(`(.*)(?:-)([0-9]+\.[0-9]+(\.[0-9]+)?)(\.|-)+[a-zA-Z]+$`)

func processDependencyFilename(filename string) Dependency {
	packageName := strings.Replace(filename, "BOOT-INF/lib/", "", -1)
	packageName = strings.Replace(packageName, "WEB-INF/lib/", "", -1)
	packageName = strings.Replace(packageName, "WEB-INF/lib-provided/", "", -1)

	FullName := packageName

	packageName = strings.TrimSuffix(packageName, ".jar")

	Version := ""

	matched1 := regex1.FindStringSubmatch(packageName)
	if matched1 != nil {
		packageName = matched1[1]
		Version = matched1[2]
	}

	matched2 := regex2.FindStringSubmatch(packageName)
	if matched2 != nil {
		packageName = matched2[1]
		Version = matched2[2]
	}

	return Dependency{Name: packageName, Version: Version, FullName: FullName}
}

func processJarFile(reader *zip.ReadCloser) []Dependency {
	var dependencies []Dependency

	// 3. Iterate over jar/war files inside the archive and unzip each of them
	for _, f := range reader.File {
		if strings.HasPrefix(f.Name, "BOOT-INF/lib/") && strings.HasSuffix(f.Name, ".jar") {
			dependency := processDependencyFilename(f.Name)
			dependencies = append(dependencies, dependency)
		}
	}
	return dependencies
}

func processWarFile(reader *zip.ReadCloser) []Dependency {
	var dependencies []Dependency

	// 3. Iterate over jar/war files inside the archive and unzip each of them
	for _, f := range reader.File {
		if strings.HasPrefix(f.Name, "WEB-INF/lib/") && strings.HasSuffix(f.Name, ".jar") {
			dependency := processDependencyFilename(f.Name)
			dependencies = append(dependencies, dependency)
		}

		if strings.HasPrefix(f.Name, "WEB-INF/lib-provided/") && strings.HasSuffix(f.Name, ".jar") {
			dependency := processDependencyFilename(f.Name)
			dependencies = append(dependencies, dependency)
		}
	}
	return dependencies
}

func processFile(source string) error {
	// 1. Open the jar/war file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var dependencies []Dependency
	if strings.HasSuffix(source, ".jar") {
		dependencies = processJarFile(reader)
	} else if strings.HasSuffix(source, ".war") {
		dependencies = processWarFile(reader)
	} else {
		return fmt.Errorf("File %s is not a jar or war file", source)
	}

	result, _ := json.Marshal(dependencies)
	fmt.Println(string(result))

	return nil
}

/**
Inspiration:
	- https://gosamples.dev/unzip-file/#:~:text=To%20unzip%20a%20compressed%20archive,through%20all%20the%20archive%20files.
	- https://www.sohamkamani.com/golang/json/
*/
func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "Jar/war file")
	flag.Parse()

	if len(filename) == 0 {
		fmt.Println("Usage: main -f file.jar")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := processFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}
