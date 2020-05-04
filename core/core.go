package core

import (
	"../utils"
	"fmt"
	"io/ioutil"
	"strings"
)

// BuildFeatures returns the feature matrix
func BuildFeatures(dir string, dictionary utils.WordDict) (utils.FeatMat, error) {
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Matrix to have features
	// Renders a matrix like this featMatrix => [len(emailList)][len(dictionary)]int
	featMatrix := make([][]int, len(emailList))
	for i := range featMatrix {
		featMatrix[i] = make([]int, len(dictionary))
	}

	// Collecting the number of occurences of each of the words in the emails.
	for emailI, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for lineI, line := range dat {
			// Body of email is only 3rd line of text file
			if lineI == 2 {
				words := strings.Split(line, " ")
				for word, wordI := range dictionary {
					featMatrix[emailI][wordI] = utils.Count(words, word)
				}
			}
		}
	}
	return featMatrix, nil
}

// BuildLabels for spam emails
func BuildLabels(dir string) (utils.LabelMat, error) {
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Label vector
	labelMat := make([]int, len(emailList))

	for i, email := range emailList {
		if strings.Contains(email.Name(), "spms") {
			labelMat[i] = 1
		} else {
			labelMat[i] = 0
		}
	}

	return labelMat, nil
}

// BuildDictionary creates dictionary from all the emails in directory.
// Configured for Enron Datasets only.
func BuildDictionary(dir string) ([]utils.WordDict, error) {
	var err error
	// Read the file names and sorts them.
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Slice to hold all the words in the emails
	goodwordlist := []string{}
	spamwordlist := []string{}

	// Collecting all words from those emails
	for _, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}

		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")

		for i, line := range dat {
			// Body of email is only 3rd line of text file
			if i == 2 {
				words := strings.Split(line, " ")
				if strings.Contains(email.Name(), "spms") {
					// This labels the data for training purposes for spam dataset
					spamwordlist = append(spamwordlist, words...)
				} else {
					goodwordlist = append(goodwordlist, words...)
				}
			}
		}
	}

	// STATS: Wordcount -> 138777
	// fmt.Println(len(goodwordlist))
	// fmt.Println(len(spamwordlist))

	// We now have the dictionary of words, which may have duplicate entries
	goodworddict := utils.Set(goodwordlist) // Duplicates removed.
	spamworddict := utils.Set(spamwordlist) // Duplicates removed.

	// STATS: Wordcount -> 13397
	// fmt.Println(len(goodworddict))
	// fmt.Println(len(spamworddict))

	// Removes punctuations and non-alphabets
	for word := range goodworddict {
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(goodworddict, word)
		}
	}
	for word := range spamworddict {
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(spamworddict, word)
		}
	}

	// STATS: Wordcount -> 11793
	// fmt.Println(len(goodworddict))
	// fmt.Println(len(spamworddict))

	return []utils.WordDict{goodworddict, spamworddict}, nil
}


// BuildFmtList preprocesses data to good and spam list for models to train.
func BuildFmtList(filename string) ([][]string, error) {
	goodwordlist := []string{}
	spamwordlist := []string{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("File opening failed %s", err)
	}
	// Breaks the email into lines.
	dat := strings.Split(string(data), "\n")
	//Last line is blank
	for _, line := range dat[0 : len(dat)-1] {
		// Remove labels from start
		words := strings.Split(line[2:len(line)], " ")
		switch line[0] {
		case 49:
			// Label of email is 1 i.e. (49 in ascii) then it is a spam mail
			spamwordlist = append(spamwordlist, words...)
		case 48:
			// Label of email is 0 i.e. (48 in ascii) then it is a fam mail
			goodwordlist = append(goodwordlist, words...)
		}
	}
	return [][]string{goodwordlist, spamwordlist}, nil
}

// ReadEmail reads emails for test data to perform.
func ReadEmail(emailName string) ([]string, error) {
	wordlist := []string{}
	data, err := ioutil.ReadFile(emailName)
	if err != nil {
		return nil, fmt.Errorf("File opening failed %s", err)
	}
	// Breaks the email into lines.
	dat := strings.Split(string(data), "\n")
	for i, line := range dat {
		// Body of email is only 3rd line of text file
		if i == 2 {
			words := strings.Split(line, " ")
			wordlist = append(wordlist, words...)
		}
	}
	// Remove Duplicates
	wordDict := utils.Set(wordlist)

	// Useless chars deleted
	for word := range wordDict {
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(wordDict, word)
		}
	}

	return utils.MapToArr(wordDict), nil
}