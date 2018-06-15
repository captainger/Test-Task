package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func fileHanlder(wg *sync.WaitGroup, fileName string, ch chan int, out io.Writer) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		ch <- 0
		fmt.Fprintln(out, "File ", fileName, " can not be open")
		wg.Done()
		return
	}
	number := strings.Count(string(fileContent), "Go")
	ch <- number
	fmt.Fprintf(out, "Count for %s: %d\n", fileName, number)
	wg.Done()
}

func urlHandler(wg *sync.WaitGroup, url string, ch chan int, out io.Writer) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- 0
		fmt.Fprintln(out, "URL ", url, " can not be open")
		wg.Done()
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	number := strings.Count(string(respBody), "Go")
	ch <- number
	fmt.Fprintf(out, "Count for %s: %d\n", url, number)
	wg.Done()
}

func GoSearcher(input io.Reader, out io.Writer, origin string) {
	in := bufio.NewScanner(input)
	var wg sync.WaitGroup
	numberOfGoroutines := 0
	total := 0
	ch1 := make(chan int)

	for in.Scan() {
		name := in.Text()
		if numberOfGoroutines >= 5 { //Если горутин 5, то ждем пока одна не закончит работу и уменьшаем счетчик горутин
			total += <-ch1
			numberOfGoroutines--
		}
		wg.Add(1) // Если горутин меньше пяти, то добавляем новую в вейт группу, увеличиваем счетчик и запускаем обработчик
		numberOfGoroutines++
		if origin == "url" {
			go urlHandler(&wg, name, ch1, out)
		} else if origin == "file" {
			go fileHanlder(&wg, name, ch1, out)
		}
	}
	for numberOfGoroutines > 0 { //обрабатываем ответы от горутин
		total += <-ch1
		numberOfGoroutines--
	}
	close(ch1)
	wg.Wait()
	fmt.Fprint(out, "Total: ", total)
}

func main() {
	out := os.Stdout
	in := os.Stdin
	if len(os.Args) != 3 || os.Args[1] != "-type" {
		log.Fatal("Usage: go run main.go -type [url/file]")
	}
	if !(os.Args[2] == "url" || os.Args[2] == "file") {
		log.Fatal("Invalid type. Type can be ony url or file")
	}
	GoSearcher(in, out, os.Args[2])
}
