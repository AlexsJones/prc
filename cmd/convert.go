/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	template "github.com/AlexsJones/prc/templates"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	inputFiles []string
	outputPath string
)

func generateRule(file string, b []byte) {
	//Filename
	extension := filepath.Ext(file)
	file = file[0 : len(file)-len(extension)]
	file = filepath.Base(file)
	file = strings.Replace(file, "-", " ", -1)
	// e.g foo bar
	file = strings.Replace(file, "_", " ", -1)
	// Remove any - or _ from the filename
	file = strings.Title(strings.ToLower(file))
	// Set to titlecase e.g "Foo Bar"
	file = strings.Replace(file, " ", "", -1)
	// Remove any spare white space e.g FooBar
	var rule template.RecordingRule
	err := yaml.Unmarshal(b, &rule)
	if err != nil {
		log.Fatal(err)
	}

	// Create the prometheus rule
	var prometheusRule template.PrometheusRule
	prometheusRule.Spec.Groups = rule.Groups
	prometheusRule.Metadata.Name = file
	prometheusRule.APIVersion = "monitoring.coreos.com/v1"
	prometheusRule.Kind = "PrometheusRule"

	bytes, err := yaml.Marshal(prometheusRule)

	outputP := "."
	if outputPath != "" {
		outputP = outputPath
	}
	if err := ioutil.WriteFile(path.Join(outputP, fmt.Sprintf("prometheusrule-%s.yaml", file)), bytes, 0644); err != nil {
		log.Fatal(err)
	}
	color.Green("Created new PrometheusRule %s", path.Join(outputP, fmt.Sprintf("prometheusrule-%s.yaml", file)))
}

func loadFiles(inFiles []string) {
	for _, file := range inFiles {
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			log.Fatalf("File does not exist:%s", file)
		}
		if info.IsDir() {
			log.Debugf("Searching %s", file)
			files, err := ioutil.ReadDir(file)
			if err != nil {
				log.Fatal(err)
			}
			flist := []string{}
			for _, f := range files {
				flist = append(flist, fmt.Sprintf("%s/%s", file, f.Name()))
			}
			loadFiles(flist)
		} else {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			generateRule(file, b)
		}
	}

}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert from prometheus recording rules into PrometheusRules",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		loadFiles(inputFiles)
	},
}

func init() {

	convertCmd.Flags().StringArrayVarP(&inputFiles, "from-files", "f", []string{}, `Input prometheus recording rules file(s) e.g. --from-files=rule.yaml.
Also supports directory paths e.g. --from-files=./localrules/`)
	convertCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "Output path for input files when converted into PrometheusRules e.g. --output-path=../")

	err := convertCmd.MarkFlagRequired("from-files")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
