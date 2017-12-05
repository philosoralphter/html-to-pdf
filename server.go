package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func example(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprint(w, "Did you mean to /convert a file?") // send data to client side
}

func handleConvertRequest(w http.ResponseWriter, r *http.Request) {
	var requestId string = getUuid()

	var tempFileName string = "/tmp/convert-" + requestId + "-input.html"
	//write file to tmp
	inputFile, createError := os.Create(tempFileName)
	bytes, pipeErr := io.Copy(inputFile, r.Body)
	if pipeErr != nil || createError != nil {
		fmt.Println("Error copying body to local file")
	}
	log.Println("Created file: ", tempFileName, "bytes:" , bytes)

	//convert
	toPdf(tempFileName, requestId)

	//Sned pdf
	_, copyError2 := os.Stat(tempFileName)
	if copyError2 != nil {
		fmt.Println("Output file does not exist!")
		http.Error(w, "Conversion Failed", 500)
		return;
	}

	http.ServeFile(w, r, tempFileName);
}

func toPdf(inputFilename string,requestId string) {
	//convert to pdf
	var chrome = "/bin/bash"
	opts := []string{
		"-c",
		//"ls",
		" /usr/local/Caskroom/google-chrome/latest/Google\\ Chrome.app/Contents/MacOS/Google\\ Chrome"+
		" --headless"+
		" --disable-gpu"+
		" --print-to-pdf=/tmp/convert-" + requestId + "-output.pdf" +
		" file://" + inputFilename,
	}
	cmd := exec.Command(chrome, opts...)

	out, cmdErr := cmd.CombinedOutput()

	if cmdErr != nil {
		log.Printf("Command finished with error: %v", cmdErr)
		log.Printf("Output: %s", out)
		log.Print(cmdErr)
		return;
	}

	log.Println("Waiting for command to finish...")
	cmd.Wait()
	log.Printf("Output: %s", out)
}

func getUuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(out), "\n")
}

func main() {
	http.HandleFunc("/", example)            // set router
	http.HandleFunc("/convert", handleConvertRequest)       // set router
	err := http.ListenAndServe(":9190", nil) // set listen port use listenandservetls when we haver a cert.
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
