package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"

	"github.com/jiazhoulvke/json_filter"
)

var (
	ErrSyntax = errors.New("syntax error")
)

var (
	sql          string
	sqlFile      string
	errorOutput  string
	resultOutput string
)

func init() {
	pflag.StringVarP(&sql, "sql", "q", "", "sql")
	pflag.StringVarP(&sqlFile, "sql_file", "f", "", "sql file")
	pflag.StringVarP(&errorOutput, "error_output", "", "", "error output")
	pflag.StringVarP(&resultOutput, "output", "o", "", "output")
}

func main() {
	pflag.Parse()

	// sql
	if sqlFile != "" {
		sqlFileData, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sql = string(sqlFileData)
	}
	if sql == "" {
		sql = "select * from t"
	}

	// input
	var r io.Reader = os.Stdin

	paths := pflag.Args()
	if len(paths) > 0 {
		files := make([]io.Reader, 0)
		for _, p := range paths {
			f, err := os.Open(p)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer f.Close()
			files = append(files, f)
		}
		r = io.MultiReader(files...)
	}

	// output
	var w io.Writer
	if resultOutput != "" {
		resultFile, err := os.Create(resultOutput)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resultFile.Close()
		w = resultFile
	} else {
		w = os.Stdout
	}

	// error output
	var errWriter io.Writer
	if errorOutput != "" {
		errFile, err := os.OpenFile(errorOutput, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer errFile.Close()
		errWriter = errFile
	} else {
		errWriter = os.Stderr
	}

	filter, err := json_filter.NewJSONFilterWithConfig(json_filter.FilterConfig{
		SQL:       sql,
		ErrWriter: errWriter,
		Reader:    r,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for filter.Next() {
		line, err := filter.GetData()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintln(w, string(line))
	}
}
