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

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)
	fmt.Println("URL", r.URL)

	for k, v := range r.Form {
		fmt.Println("Query Params:")
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprint(w, "Did you mean to /convert a file?")
}

func handleConvertRequest(w http.ResponseWriter, r *http.Request) {
	var requestId string = getUuid()

	var tempFileName string = getInputFileName(requestId)
	var outputFile string = getOutputFileName(requestId)

	r.ParseForm();

	//write file to tmp
	var bytes int64;
	inputFile, createError := os.Create(tempFileName)
	defer inputFile.Close();

	//Determine if using body or external URL
	if targetURL, urlIsSet := r.Form["from-url"]; urlIsSet {

		log.Println("Using URL to fetch Original HTML");
		// Get the page
		resp, err := http.Get(strings.Join(targetURL, ""))
		if err != nil {
			w.WriteHeader(404)
			http.Se
			//return err
		}
		defer resp.Body.Close()

		// Write the body to file
		bytes1, err := io.Copy(inputFile, resp.Body)
		bytes = bytes1;
		if err != nil {
			//return err
		}

	} else {
		//Use posted body
		bytes2, pipeErr := io.Copy(inputFile, r.Body)
		bytes = bytes2
		if pipeErr != nil || createError != nil {
			fmt.Println("Error copying body to local file")
		}
	}

	log.Println("Created file: ", tempFileName)

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
		chrome = "/Applications/Google\\ Chrome.app/Contents/MacOS/Google\\ Chrome"; //location in homebrew (mac) for development/testing
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
	http.HandleFunc("/", defaultRoute)                // set router
	http.HandleFunc("/convert", handleConvertRequest) // set router

	log.Print("Starting Server...")

	err := http.ListenAndServe(":9190", nil) // set listen port use listenandservetls when we have a cert.
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
