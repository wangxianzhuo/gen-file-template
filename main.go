// Copyright 2021 XianZhuo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var templateFile = flag.String("temp-file", "", "template file path")
var storePath = flag.String("store-path", "./", "the path to store generated file")
var useCurrentDateAsFileName = flag.Bool("file-name-use-date", true, "use current date like yyyyMMdd.md as the generated file name")
var fileExtension = flag.String("file-extension", "md", "the generated file extension")

const fileNamePatternAsDate = "20060102"
const defaultFileName = "default"

func main() {
	paramsHandle()

	err := createPathIfNotExist(*storePath)
	if err != nil {
		log.Printf("Create store path error: %v.", err)
		os.Exit(-1)
	}

	f, err := createFileIfNotExist(*storePath)
	if err != nil {
		log.Printf("Create file error: %v.", err)
		os.Exit(-1)
	}

	templateFile, err := os.Open(*templateFile)
	if err != nil {
		log.Printf("Open template file error: %v.", err)
		os.Exit(-1)
	}

	_, err = io.Copy(f, templateFile)
	if err != nil {
		log.Printf("Generate file %q error: %v.", f.Name(), err)
		os.Exit(-1)
	}

	log.Printf("generate file %q.\nPress 'Ctrl + c' to exit.", f.Name())

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT)

	for q := range quit {
		log.Println("exit", q)
		os.Exit(1)
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}

func createPathIfNotExist(filePath string) error {
	if !isExist(filePath) {
		return os.MkdirAll(filePath, os.ModePerm)
	}
	return nil
}

func createFileIfNotExist(filePath string) (*os.File, error) {
	var fileName string

	if *useCurrentDateAsFileName {
		fileName = time.Now().Format(fileNamePatternAsDate)
	} else {
		fileName = defaultFileName
	}

	fileFullName := path.Join(filePath, fileName) + "." + *fileExtension

	if isExist(fileFullName) {
		return nil, fmt.Errorf("file %q exist", fileFullName)
	}

	return os.Create(fileFullName)
}

func paramsHandle() {
	flag.Parse()

	if !isExist(*templateFile) {
		fmt.Printf("template file %q not exist\n", *templateFile)
		flag.Usage()

		os.Exit(-1)
	}
}
