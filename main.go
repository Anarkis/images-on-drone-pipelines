package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gopkg.in/yaml.v3"
)

func readCredentials(file string) map[interface{}]interface{} {
	yfile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Reading credentials from %s\n", file)

	data := make(map[interface{}]interface{})

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}

	return data
}

func findImages(dir string, repoImages map[string]void) {
	file, err := os.Open(dir + "/.drone.yml")
	var stepS390xFlag bool = false

	reImage := regexp.MustCompile(`image:[a-zA-Z \/:0-9.-]*`)
	reStep := regexp.MustCompile(`^name: `)
	reStepS390x := regexp.MustCompile(`^name: [a-zA-Z \/:0-9.-]*s390x`)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		step := reStep.FindAll([]byte(text), -1)

		if len(step) > 0 {
			stepS390x := reStepS390x.FindAll([]byte(text), -1)
			if len(stepS390x) > 0 {
				stepS390xFlag = true
			} else {
				stepS390xFlag = false
			}
		}
		if stepS390xFlag {
			found := reImage.FindAll([]byte(text), -1)
			if len(found) > 0 {
				fmt.Printf("\t%q\n", found[0])
				repoImages[string(found[0])] = member
			}
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func clone(repo string, credentials map[interface{}]interface{}) string {

	dirname := strings.Split(repo, "/")

	dir, err := ioutil.TempDir("./temp/", dirname[len(dirname)-1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Repo %s\n", repo)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: repo,
		Auth: &http.BasicAuth{
			Username: credentials["user"].(string),
			Password: credentials["pwd"].(string),
		},
		Depth: 1,
	})

	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getRepos(reposFile string) []string {
	file, err := os.Open(reposFile)
	repoList := []string{}

	fmt.Println("Fetching the repos to analyze")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repoList = append(repoList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return repoList
}

type void struct{}

var member void

func main() {

	repoImages := make(map[string]void)

	repoList := getRepos("repositories.txt")
	credentials := readCredentials("credentials.yaml")

	for _, repo := range repoList {
		dir := clone(repo, credentials)
		findImages(dir, repoImages)
		defer os.RemoveAll(dir)
	}

	fmt.Println("\nUnique images used in all the repos:")
	for k := range repoImages {
		fmt.Println(k)
	}
}
