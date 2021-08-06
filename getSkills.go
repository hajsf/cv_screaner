package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
func head(s string) bool {
	r, e := http.Head(s)
	return e == nil && r.StatusCode == 200
}
*/
func getSkills(w http.ResponseWriter, r *http.Request) {
	// define custom types
	type Input struct {
		Path   string `json:"path"`
		Skills string `json:"skills"`
	}
	type Output struct {
		MSG    string   `json:"msg"`
		File   string   `json:"file"`
		Skills []string `json:"skills"`
	}

	type Result struct {
		Rslt []Output `json:"result"`
	}

	// define vars
	var input Input
	var output Output
	var result Result

	// decode input or return error
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		//	fmt.Fprintf(w, "Decode error! please check your JSON formating.")
		return
	}
	/*
		routeRequired := head("http://localhost:3000/pdf/")
		if routeRequired {
			http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(input.Path))))
		}
	*/
	tmpDir := os.TempDir() + "\\cvsText"
	fmt.Println("temp is: ", tmpDir)
	/*	fmt.Println("cvTemp: ", exists)
		//	fmt.Println("tmpDir: ", tmpDir)
		if !exists {
			tmpDir := os.TempDir() + "\\cvscanner"
			_ = exec.Command(`SETX`, `cvTemp`, tmpDir).Run()
			/*	if err != nil {
						fmt.Printf("Error: %s\n", err)
				}
			/ //	fmt.Println("tmpDir: ", tmpDir)
		} */

	var files []string
	var paths []string

	//tmpDir := os.TempDir() + "\\cvscanner"
	root := input.Path
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.ToLower(filepath.Ext(path)) == ".pdf" {
			//copy(path, tmpDir)

			f := strings.TrimSuffix(info.Name(), filepath.Ext(path))
			files = append(files, f)
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for i, file := range files {
		// fmt.Println(file + ".pdf")
		fmt.Println("copy from: ", paths[i])
		source, err := os.Open(paths[i])
		if err != nil {
			panic(err)
		}
		defer source.Close()
		fmt.Println("copy to: ", tmpDir)
		destination, _ := os.Create(tmpDir + "/" + file + ".pdf")
		defer destination.Close()
		_, err = io.Copy(destination, source)
		//	_ = exec.Command(`COPY`, path, tmpDir).Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		if err != nil {
			panic(err)
		}
		str, err := ExtractTexts(file, paths[i])
		if err != nil {
			panic(err)
		}
		re := regexp.MustCompile(`(?i)` + input.Skills)
		matches := re.FindAllString(str, -1)
		uniqueSlice := unique(matches)
		if len(str) == 0 {
			uniqueSlice = []string{"Empty file or can not be read."}
		} else {
			if found := len(uniqueSlice); found == 0 {
				uniqueSlice = []string{"No single match had been found."}
			} else {
				uniqueSlice = append([]string{fmt.Sprintf("%02d skills found: ", len(uniqueSlice))}, uniqueSlice...)
			}
		}
		//	p := strings.Split(input.Skills, "|")
		//	fmt.Println("\nskills found:", len(uniqueSlice), "/", len(p))
		//	fmt.Println("'" + strings.Join(uniqueSlice, `', '`) + `'`)
		output = Output{MSG: "ok", File: file, Skills: uniqueSlice}
		result.Rslt = append(result.Rslt, output)
		// w.WriteHeader(http.StatusCreated)
		// json.NewEncoder(w).Encode(output)
		//	w.WriteHeader(http.StatusOK)
		//	w.Write(data)
		//json.NewEncoder(w).Encode(output)
	}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	// fmt.Println(output)
	// print user inputs
	//fmt.Fprintf(w, "Inputed name: %s", input.Path)
}
