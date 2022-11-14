package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// ---------------------------------------------------FUNCTIONS FOR HANGMAN GAME----------------------------------------------------------------
func splitHangman(file string, n int) []string { // Split each character every 7 ligns in the hangman file
	listlign := []string{}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("File ", file, "not fond")
	}
	content2 := string(content)
	comp := 0
	inf := 0
	for i, j := range content2 {
		if j == '\n' { // compt carriages return to reach the end of a hangman postion
			comp++
		}
		if comp == n+1 {
			listlign = append(listlign, content2[inf:i])
			inf = i
			comp = 0
		}
	}
	return listlign
}

type HangManData struct {
	Word             string   // Word composed of '_', ex: H_ll_
	ToFind           string   // Final word chosen by the program at the beginning. It is the word to find
	Attempts         int      // Number of attempts left
	HangmanPositions []string // It can be the array where the positions parsed in "hangman.txt" are stored
	BasicLetter      []string //letter given at the beggining
	TriedLetter      []string // Letter which were already tried (success or not)
}

func SplitWords(s string) []string { //Split words in the txt file
	list := []int{0}
	list2 := []string{}
	listfinal := []string{}
	for i := range s {
		if string(s[i]) == "\n" {
			list = append(list, i)
		}
	}
	list = append(list, len(s))
	list2 = append(list2, s[list[0]:list[1]])
	for i := 1; i <= len(list)-2; i++ {
		list2 = append(list2, s[list[i]+1:list[i+1]])
	}
	for i := range list2 {
		if list2[i] != "" {
			listfinal = append(listfinal, list2[i])
		}
	}
	return listfinal
}
func end(etat string) { //print a message at the end of the game
	switch etat {
	case "win":
		fmt.Println("CONGRATS !")
	case "loose":
		fmt.Println("GAME OVER !")
	}
}
func Read(file string) []string { // read the file, stock his content and split each words in a list
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("The file ", file, " wasn't found in the current directory")
		return []string{}
	}
	words := SplitWords(string(content))
	return words
}
func ReadJson(file string) []byte { //read an encoded json file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return content
}
func ChoseWord(words []string) string { // choose a word in a list
	rand.Seed(time.Now().UnixNano())
	InDWord := rand.Intn(len(words))
	return string(words[InDWord])
}
func Hangman(game HangManData) { //main function of the game, which deals with user input and its consequences
	if game.Attempts == 0 {
		end("loose")
		return
	}
	if game.Word == game.ToFind {
		end("win")
		return
	}
	var attempt string
	fmt.Printf("choose: ")
	fmt.Scanln(&attempt)
	Research(game, attempt)
}
func Research(game HangManData, letter string) { // research if the given letter is present in the word, or if the word given is the same, or end the game while saving its state if the key-word STOP is entered
	if letter == "STOP" {
		b, _ := json.Marshal(game)
		message := []byte(b)
		ioutil.WriteFile("save", message, 0644)
		return
	}
	for _, i := range game.TriedLetter {
		if i == letter {
			fmt.Println("the letter ", letter, " was already tried")
			Hangman(game)
			return
		}
	}
	game.TriedLetter = append(game.TriedLetter, letter)
	listtmp := []rune(game.Word)
	flag := false
	if len(letter) > 1 && letter == game.ToFind {
		end("win")
		return
	} else if len(letter) > 1 {
		game.Attempts -= 2
		if game.Attempts > 0 {
			fmt.Println("The word is wrong,", game.Attempts, "attempts remaining")
			fmt.Println(game.HangmanPositions[len(game.HangmanPositions)-game.Attempts-1])
			ToPrint := ""
			for i := range game.Word {
				ToPrint += string(game.Word[i])
				ToPrint += " "
			}
			fmt.Print(ToPrint)
		} else {
			fmt.Println("The word is wrong,", "0", "attempts remaining")
			fmt.Println(game.HangmanPositions[9])
			fmt.Println()
		}
		if game.Attempts >= 1 {
			Hangman(game)
			return
		} else {
			end("loose")
			return
		}
	}
	for i := range game.ToFind {
		if string(game.ToFind[i]) == string(letter) {
			listtmp[i] = rune(letter[0])
			game.TriedLetter = append(game.TriedLetter, string(letter[0]))
			flag = true
		}
	}
	game.Word = ""
	for i := range listtmp {
		game.Word += string(listtmp[i])
	}
	wordtemp := ""
	if flag {
		for i := range game.Word {
			wordtemp += string(game.Word[i])
			wordtemp += " "
		}
		fmt.Print(wordtemp)
	} else {
		game.Attempts--
		fmt.Println("Not present in the word,", game.Attempts, "attempts remaining")
		fmt.Println(game.HangmanPositions[len(game.HangmanPositions)-game.Attempts-1])
		fmt.Println()
		ToPrint := ""
		for i := range game.Word {
			ToPrint += string(game.Word[i])
			ToPrint += " "
		}
		fmt.Print(ToPrint)
	}
	Hangman(game)
}
func duplicateInArray(arr []int) int { //find if there is a duplicate in an array and return the index if so, -1 if not
	visited := make(map[int]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
			return arr[i]
		} else {
			visited[arr[i]] = true
		}
	}
	return -1
}
func main() { // initialise the hangman, with the first user's input and continue the game with the hangman function
	var game HangManData
	args := os.Args[1:]
	file := args[0]
	printtmp := ""
	for i := range args {
		if args[i] == "--startWith" {
			formate := ReadJson(string(args[i+1]))
			json.Unmarshal(formate, &game)
			for i := range game.Word {
				printtmp += string(game.Word[i])
				printtmp += " "
			}
			fmt.Print(printtmp)
			fmt.Print("\n")
			Hangman(game)
			return
		}
	}
	words := Read(file)
	if len(words) == 0 {
		return
	}
	FinalWord := ChoseWord(words)
	tmp := FinalWord
	Basic := []string{}
	ListIndex := []int{}
	ceuil := len(FinalWord)/2 - 1
	FlagTmp := false
	tmp2 := 0
	for !FlagTmp { // choose random letter to show at the beggining, while preventing duplicated letters
		Index := rand.Intn(len(tmp))
		if duplicateInArray(ListIndex) != -1 {
			ListIndex = ListIndex[:len(ListIndex)-1]
		}
		ListIndex = append(ListIndex, Index)
		tmp2++
		if tmp2 >= ceuil && duplicateInArray(ListIndex) == -1 {
			FlagTmp = true
		}
	}
	for i := range ListIndex {
		FindDuplicate := FinalWord[ListIndex[i]]
		for j, h := range FinalWord {
			if h == rune(FindDuplicate) {
				ListIndex = append(ListIndex, j)
			}
		}
	}
	WordStep1 := "" //will be the word with found letters and "_" for unknown letters
	flag := false
	for i := range FinalWord {
		for j := range ListIndex {
			if ListIndex[j] == i {
				WordStep1 += string(FinalWord[i])
				Basic = append(Basic, string(FinalWord[i]))
				flag = true
				break
			}
		}
		if !flag {
			WordStep1 += "_"
		}
		flag = false
	}
	wordtemp := ""
	for i := range WordStep1 { //crate a more readable word in ASCII because of the spaces
		wordtemp += string(WordStep1[i])
		wordtemp += " "
	}
	fmt.Print(wordtemp)
	fmt.Println("")
	if game.Word == "" {
		game = HangManData{Word: WordStep1, ToFind: FinalWord, Attempts: 10, HangmanPositions: splitHangman("hangman.txt", 7), BasicLetter: Basic, TriedLetter: []string{}}
		Hangman(game)
	}
}
