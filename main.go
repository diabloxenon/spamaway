package main

import (
	"./lib/bayesian"
	"./lib/utils"
	"fmt"
	"io/ioutil"
	// "os"
	"strings"
)

// BuildDictionary creates dictionary from all the emails in directory
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
		// This labels the data for training purposes for spam dataset
		if strings.Contains(email.Name(), "spms") {
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
					spamwordlist = append(spamwordlist, words...)
				}
			}
		} else {
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
					goodwordlist = append(goodwordlist, words...)
				}
			}
		}
	}

	// STATS: Wordcount -> 138777
	fmt.Println(len(goodwordlist))
	fmt.Println(len(spamwordlist))

	// We now have the dictionary of words, which may have duplicate entries
	goodworddict := utils.Set(goodwordlist) // Duplicates removed.
	spamworddict := utils.Set(spamwordlist) // Duplicates removed.

	// STATS: Wordcount -> 13397
	fmt.Println(len(goodworddict))
	fmt.Println(len(spamworddict))

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
	fmt.Println(len(goodworddict))
	fmt.Println(len(spamworddict))

	return []utils.WordDict{goodworddict, spamworddict}, nil
}

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

func BuildFmtDict(filename string) ([][]string, error) {
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

func ReadFmtDict(filename string) ([]string, error) {
	wordlist := []string{}
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
		wordlist = append(wordlist, words...)
	}
	return wordlist, nil
}

func main() {
	// trainDir := "dataset/train_data"

	fmt.Println("1. Building dictionary")
	// dict, err := BuildDictionary(trainDir)
	// utils.Check(err)
	dict, err := BuildFmtDict("dataset/spam_train.txt")
	utils.Check(err)

	// fmt.Println("2. Building training features and labels")
	// featTrain, err := BuildFeatures(trainDir, dict)
	// utils.Check(err)
	// labelTrain, err := BuildLabels(trainDir)
	// utils.Check(err)

	// fmt.Println("3. Training the Classifier")
	// featTrain, err := BuildFeatures(trainDir, dict)
	// utils.Check(err)

	var (
		Fam  bayesian.Class = "fam"  // The good ones
		Spam bayesian.Class = "spam" // The bad ones
	)
	// New Multinomial TF classifier
	classifier := bayesian.NewClassifier(bayesian.MultinomialTf)

	// famMails := utils.MapToArr(dict[0])
	// spamMails := utils.MapToArr(dict[1])

	famMails := &dict[0]
	spamMails := &dict[1]

	// classifier.Learn(famMails, Fam)
	// classifier.Learn(spamMails, Spam)

	// Do learning using two documents
	classifier.Learn(
		bayesian.NewDocument(Fam, *famMails),
		bayesian.NewDocument(Spam, *spamMails),
	)

	// classifier.ConvertTermsFreqToTfIdf()

	// dir := "dataset/test_data"
	// emailList, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	fmt.Printf("Directory not present %s", err)
	// }

	// // Collecting all words from those emails
	// for _, email := range emailList {
	// 	dat, _ := ReadEmail(fmt.Sprintf("%s/%s", dir, email.Name()))

	// // Classify tokens from a document
	// allScores, class, certain := classifier.Classify(dat)
	// fmt.Println(allScores, class, certain)
	// }

	data, err := ioutil.ReadFile("dataset/spam_test.txt")
	if err != nil {
		fmt.Printf("File not present %s", err)
	}

	// Breaks the email into lines.
	dat := strings.Split(string(data), "\n")

	success := 0.0
	failed := 0.0
	total := float64(len(dat) - 1)
	//Last line is blank
	for _, line := range dat[0 : len(dat)-1] {
		// Remove labels from start
		words := strings.Split(line[2:], " ")
		// Classify tokens from a document
		allScores, class, certain := classifier.Classify(words)
		// fmt.Println(allScores, class, certain)
		// fmt.Println(verifySpam(line[0], string(class)))
		// Rate of success
		isSpam := verifySpam(line[0], string(class))
		if (line[0] == 48 && !isSpam) || (line[0] == 49 && isSpam) {
			success++
		} else if !isSpam {
			fmt.Println(allScores, class, certain)
			fmt.Println(string(line[0]))
			failed++
		}
	}
	successPercent := (success / total) * 100
	failedPercent := (failed / total) * 100

	fmt.Printf("Success Rate is %0.3f\n", successPercent)
	fmt.Printf("Rate of false positive emails is %0.3f\n", failedPercent)

}

func verifySpam(label byte, class string) bool {
	if label == 49 && class == "spam" {
		return true
	} else if label == 48 && class == "fam" {
		return false
	}
	// false if no case matched
	return false
}
