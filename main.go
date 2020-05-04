package main

import (
	"./lib/utils"
	"fmt"
	"io/ioutil"
	// "os"
	// "reflect"
	"strings"
)

type WordDict map[string]bool
type FeatMat [][]int
type LabelMat []int

// BuildDictionary creates dictionary from all the emails in directory
func BuildDictionary(dir string) (WordDict, error) {
	var err error
	// Read the file names and sorts them.
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}
	
	// Slice to hold all the words in the emails
	wordlist := []string{}

	// Collecting all words from those emails
	for _, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for i, line := range dat{
			// Body of email is only 3rd line of text file
			if i == 2{
				words := strings.Split(line, " ")
				wordlist = append(wordlist, words...)
			}
		}
	}

	// STATS: Wordcount -> 138777
	fmt.Println(len(wordlist))
	
	// We now have the dictionary of words, which may have duplicate entries
	worddict := utils.Set(wordlist) // Duplicates removed.
	
	// STATS: Wordcount -> 13397
	fmt.Println(len(worddict))

	// Removes punctuations and non-alphabets
	for word := range worddict{
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(worddict, word)
		}
	}

	// STATS: Wordcount -> 11793
	fmt.Println(len(worddict))

	return worddict, nil
}

// BuildFeatures returns the feature matrix
func BuildFeatures(dir string, dictionary WordDict) (FeatMat, error) {
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Matrix to have features
	featMatrix := [len(emailList)][len(dictionary)]int

	// Collecting the number of occurences of each of the words in the emails.
	for email_i, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for line_i, line := range dat{
			// Body of email is only 3rd line of text file
			if line_i == 2{
				words := strings.Split(line, " ")
				for word, word_i := range dictionary{
					featMatrix[email_i][word_i] = utils.Count(words, word)
				}
			}
		}
	}
	return featMatrix, nil
}

func BuildLabels(dir string) (LabelMat, error){
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Label vector
	labelMat = [len(emails)]int

	for i, email := range emailList{
		labelMat[i] = 1 ? strings.Contains(email.Name(), "spms") : 0
	}

	return labelMat, nil
}

func main() {
	trainDir = "dataset/train_data"
	
	fmt.Println("1. Building dictionary")
	dict, err := BuildDictionary(trainDir)
	utils.Check(err)

	fmt.Println("2. Building training features and labels")
	featTrain, err := BuildFeatures(trainDir, dict)
	utils.Check(err)
	labelTrain, err := BuildLabels(trainDir)
	utils.Check(err)

	fmt.Println("3. Training the Classifier")
	featTrain, err := BuildFeatures(trainDir, dict)
	utils.Check(err)

}
