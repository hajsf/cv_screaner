package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

func ExtractTexts(file, path string) (string, error) {
	doc, _ := fitz.New(path)
	/*	if err != nil {
			 fmt.Printf("err 1: %v", err)
		}
	*/
	defer doc.Close()
	tmpDir := os.TempDir() + "\\cvsText"
	fmt.Println("temp is: ", tmpDir)
	//_ = exec.Command(`SETX`, `cvTemp`, tmpDir).Run()
	//	if err != nil {
	//		fmt.Printf("Error: %s\n", err)
	//	}
	//	tmpDir, err := ioutil.TempDir(os.TempDir(), "fitz")
	//	tmpDir := os.TempDir() + "\\cvsText"
	/*	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err = os.Mkdir(tmpDir, 0755)
		if err != nil {
			fmt.Printf("err 2: %v", err)
		} /*else {
				fmt.Printf("tmpDir: %s created", tmpDir)
		} */
	//	} /* else {
	/*		 fmt.Printf("\ntmpDir: %v already exixted", tmpDir)
	} */

	var str string
	// Extract pages as text
	for n := 0; n < doc.NumPage(); n++ {
		text, err := doc.Text(n)
		if err != nil {
			// fmt.Printf("err 3: %v", err)
			return "", err
		}
		str += text
	}
	str = trimEmptyLines(str)

	f, err := os.Create(filepath.Join(tmpDir, file+".txt"))
	if err != nil {
		fmt.Printf("err 4: %v", err)
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		fmt.Printf("err 5: %v", err)
		return "", err
	}
	return str, nil
}
