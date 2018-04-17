package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"github.com/satori/go.uuid"
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

	var tempFileName string = getInputFileName(requestId)
	var outputFile string = getOutputFileName(requestId)

	//write file to tmp
	inputFile, createError := os.Create(tempFileName)
	bytes, pipeErr := io.Copy(inputFile, r.Body)
	if pipeErr != nil || createError != nil {
		fmt.Println("Error copying body to local file")
	}
	log.Println("Created file: ", tempFileName, "bytes:" , bytes)

	//convert
	toPdf(requestId)

	//Send pdf
	_, copyError2 := os.Stat(outputFile)
	if copyError2 != nil {
		fmt.Println("Output file does not exist!")
		http.Error(w, "Conversion Failed", 500)
		return;
	}

	log.Println("Serving file: ", outputFile, "bytes:" , bytes)

	http.ServeFile(w, r, outputFile);
}

func toPdf(requestId string) {
	//convert to pdf
	//var bash = "/bin/sh"
	var inputFilename = getInputFileName(requestId);
	var outputFilename = getOutputFileName(requestId);

	var chrome = os.Getenv("CHROME_LOCATION")
	if(chrome == "") {
		chrome = "/usr/local/Caskroom/google-chrome/latest/Google Chrome.app/Contents/MacOS/Google Chrome"; //location in homebrew (mac) for development/testing
	}

	var opts = []string{
		"-c",
		chrome +
		" --headless" +
		" --no-sandbox" +
		" --disable-gpu" +
		" --print-to-pdf=" + outputFilename +
		" file://" + inputFilename,
	}
	cmd := exec.Command("sh", opts...)

	out, cmdErr := cmd.CombinedOutput()

	if cmdErr != nil {
		log.Printf("Command finished with error: %v", cmdErr)
		log.Printf("Output: %s", out)

		foundPath, _ := exec.LookPath(chrome)
		log.Print("lookpath: " + foundPath);
		return;
	}

	log.Println("Waiting for command to finish...")
	cmd.Wait()
	log.Printf("Output: %s", out)
}

func getUuid() string {
	//out, err := exec.Command("uuidgen").Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return strings.Trim(string(out), "\n")

	return uuid.NewV4().String();

}

func getInputFileName(requestId string) string {
	return "/tmp/convert-" + requestId + "-input.html"
}

func getOutputFileName(requestId string) string {
	return "/tmp/convert-" + requestId + "-output.pdf"
}

func main() {
	http.HandleFunc("/", example)            // set router
	http.HandleFunc("/convert", handleConvertRequest)       // set router

	log.Print("Starting Server...")

	err := http.ListenAndServe(":9190", nil) // set listen port use listenandservetls when we have a cert.
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
