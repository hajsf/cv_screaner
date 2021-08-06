package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/zserge/lorca"
)

func init() {
	//
}

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

//go:embed static
var staticFiles embed.FS

func main() {
	tmpDir := os.TempDir() + "\\cvsText"
	if p, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err = os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("err 2: %v", err)
		} else {
			fmt.Println("te,p created at:", p)
			_, exists := os.LookupEnv("cv")
			if !exists {
				//
				//err = os.Setenv(`cv`, tmpDir)
				_ = exec.Command(`SETX`, `cv`, tmpDir).Run()
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				}
				//	fmt.Println("tmpDir: ", tmpDir) */
			} else {
				fmt.Println("Env exisit")
			}
		}
	} else {
		fmt.Println("checking Env ")
		_, exists := os.LookupEnv("cv")
		if !exists {
			//
			//err = os.Setenv(`cv`, tmpDir)
			_ = exec.Command(`SETX`, `cv`, tmpDir).Run()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			} else {
				fmt.Println("Env created")
			}
			//	fmt.Println("tmpDir: ", tmpDir) */
		} else {
			fmt.Println("Env exisit")
		}
	}
	go func() {
		// http.FS can be used to create a http Filesystem
		var staticFS = http.FS(staticFiles)
		fs := http.FileServer(staticFS) // embeded static files
		// Serve static files, to be embedded in the binary
		http.Handle("/static/", fs)

		http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
			http.ServeFile(rw, r, "http://localhost:3000/static/favicon.ico")
		})

		//	www := http.FileServer(http.Dir("/files/")) // side static files
		// Serve public files, to be beside binary
		http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./files"))))

		http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(tmpDir))))

		//	defer os.RemoveAll(tempDir)
		/*	tmpDir := os.TempDir() + "\\scanner"
			if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
				err = os.Mkdir(tmpDir, 0755)
				if err != nil {
					fmt.Printf("err 2: %v", err)
				} else {
					fmt.Println(tmpDir)
					http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(tmpDir))))
				}
			}  else {
				 fmt.Printf("\ntmpDir: %v already exixted", tmpDir)
			} */

		http.HandleFunc("/getSkills", getSkills)

		http.ListenAndServe(":3000", nil)

	}()
	// Start UI
	date := "2021-10-30"
	dateStamp, _ := time.Parse(layoutISO, date)
	today := time.Now()
	var url string
	if today.After(dateStamp) {
		url = "http://localhost:3000/static/expired.html"
	} else {
		url = "http://localhost:3000/static/"
	}
	ui, err := lorca.New(url, "", 1200, 800)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer ui.Close()

	// Bind Go function to be available in JS. Go function may be long-running and
	// blocking - in JS it's represented with a Promise.
	ui.Bind("add", func(a, b int) int { return a + b })

	// Call JS function from Go. Functions may be asynchronous, i.e. return promises
	n := ui.Eval(`Math.random()`).Float()
	fmt.Println(n)

	// Call JS that calls Go and so on and so on...
	m := ui.Eval(`add(2, 3)`).Int()
	fmt.Println(m)

	// Wait for the browser window to be closed
	<-ui.Done()
}

/*
To return JSON to the client instead of text
	data := SomeStruct{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
*/
