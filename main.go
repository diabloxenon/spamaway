package main

import (
	"./core"
	"./lib/bayesian"
	"./utils"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	trainDir := "dataset/spam_train.txt"
	fmt.Println("1. Building dictionary")
	dict, err := core.BuildFmtList(trainDir)
	utils.Check(err)

	fmt.Println("2. Building training features and labels")
	famMails := &dict[0]
	spamMails := &dict[1]

	fmt.Println("3. Training the Classifier")
	var (
		Fam  bayesian.Class = "fam"  // The good ones
		Spam bayesian.Class = "spam" // The bad ones
	)
	// New Multinomial TF classifier
	classifier := bayesian.NewClassifier(bayesian.MultinomialTf)
	// Train the model
	classifier.Learn(
		bayesian.NewDocument(Fam, *famMails),
		bayesian.NewDocument(Spam, *spamMails),
	)

	fmt.Println("4. Testing the Model")
	testDir := "dataset/spam_test.txt"
	data, err := ioutil.ReadFile(testDir)
	if err != nil {
		fmt.Printf("File not present %s", err)
	}
	// Breaks the email into lines.
	dat := strings.Split(string(data), "\n")

	// Performance variables
	success := 0.0
	failed := 0.0
	total := float64(len(dat) - 1)
	conmat := [][]int{[]int{0, 0}, []int{0, 0}}

	//Last line is blank
	for _, line := range dat[0 : len(dat)-1] {
		// Remove labels from start
		words := strings.Split(line[2:], " ")
		// Classify tokens from a document
		_, class, _ := classifier.Classify(words)
		// fmt.Println(allScores, class, certain)
		utils.SpamFamConfusionMatrix(conmat, line[0], string(class))
	}

	fmt.Println("5. Benchmarking...")
	fmt.Println("Confusion matrix generated\n")
	fmt.Println(" Fam Fam | Fam Spam")
	fmt.Println("         |         ")
	fmt.Printf("   %d   |   %d    \n", conmat[0][0], conmat[0][1])
	fmt.Println("---------|---------")
	fmt.Printf("   %d    |   %d    \n", conmat[1][0], conmat[1][1])
	fmt.Println("         |         ")
	fmt.Println("Spam Fam | Spam Spam\n")

	success = float64(conmat[0][0] + conmat[1][1])
	failed = total - success

	successPercent := (success / total) * 100
	failedPercent := (failed / total) * 100

	fmt.Printf("\nSuccess Rate of Multinomial NB Classifier (in percent): %0.3f\n", successPercent)
	fmt.Printf("False Positive/Negative Rate of MultiNBC (in percent): %0.3f\n", failedPercent)
}
