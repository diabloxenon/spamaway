package main

import (
	"./lib/utils"
	"fmt"
	"io/ioutil"
	// "os"
	// "reflect"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Create dictionary from all the emails in directory
func BuildDictionary(dir string) error {
	var (
		err error

		// Slice to hold all the words in the emails
		WordDict []string
	)
	// Read the file names and sorts them.
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Directory not present %s", err)
	}

	WordDict = []string{}

	// Collecting all words from those emails
	for _, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for i, line := range dat{
			// Body of email is only 3rd line of text file
			if i == 2{
				words := strings.Split(line, " ")
				WordDict = append(WordDict, words...)
			}
		}
	}
	fmt.Println(len(WordDict))
	// We now have the array of words, which may have duplicate entries
	WordDict = utils.Set(WordDict)
	fmt.Println(len(WordDict))
	return nil
}

func main() {
	err := BuildDictionary("dataset/train_data")
	if err != nil {
		fmt.Println(err)
	}

}
