package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type node struct {
	val  string
	next *node
}

type linkedList struct {
	head   *node
	length int
}

func (l *linkedList) prepend(n *node) {
	second := l.head
	l.head = n
	l.head.next = second
	l.length++
}

func main() {
	// open the dictionary file ==========================
	filetoCheck := os.Args[1]
	readDictionary, err := os.Open("./dictionaries/large")
	if err != nil {
		log.Fatal("something went wrong reading the dictionary\n", err)
	}
	defer readDictionary.Close()

	// create a filescanner for the dictionary to scan each line ==========================
	dicScanner := bufio.NewScanner(readDictionary)
	dicScanner.Split(bufio.ScanLines)

	// scan each line into the correct spot in the hash table ==========================
	dicHash := make(map[string]linkedList)
	for dicScanner.Scan() {
		// look at the first letter of the thing
		word := dicScanner.Text()
		firstLetter := word[0:1]
		// if theres a key for that letter, then prepend it to the linked list there
		if existingList, ok := dicHash[firstLetter]; ok {
			newNode := &node{val: word}
			existingList.prepend(newNode)
			dicHash[firstLetter] = existingList
		} else {
			// if not then create the key and start the linked list
			newList := linkedList{}
			newNode := &node{val: word}
			newList.prepend(newNode)
			dicHash[firstLetter] = newList
		}
	}

	// open the input file - at the moment just a hard coded test file ==========================
	readInput, err := os.Open("./texts/" + filetoCheck)
	if err != nil {
		log.Fatal("error opening the input file", err)
	}
	defer readInput.Close()

	// create a filescanner for the input to scan each word ==========================
	inputScanner := bufio.NewScanner(readInput)
	inputScanner.Split(bufio.ScanWords)

	mispelled := []string{}
	_ = mispelled

	//regex thing to remove unwatned chars
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}

	// scan each word and check the hash table ==========================
	for inputScanner.Scan() {
		// firstLetter := rune(inputScanner.Text()[0])
		word := strings.ToLower(reg.ReplaceAllString(inputScanner.Text(), ""))
		if len(word) > 0 {
			key := string(word[0])
			doesExist := searchLinkedList(key, word, dicHash)
			// if the word is there then continue.
			if doesExist {
				continue
			} else {
				// if not then it must be mispelt and therefore gets appended to the mispell arr
				mispelled = append(mispelled, word)
			}
		}
	}

	fmt.Println("MISSPELLED WORDS:")
	for _, word := range mispelled {
		if len(word) > 1 {
			fmt.Println(word)
		}
	}
	fmt.Println(strconv.Itoa(len(mispelled)) + " words mispelled")
}

// test function to iterate over dictionary and print words ==========================
func printLinkedList(firstLetter string, dic map[string]linkedList) {
	list := dic[firstLetter]
	for i := 0; i < list.length; i++ {
		fmt.Println(list.head.val)
		list.head = list.head.next
	}
	return
}

// function that takes searches the linked list corresponding to the input key ==========================
func searchLinkedList(firstLetter string, word string, dictionary map[string]linkedList) bool {
	list := dictionary[firstLetter]
	for i := 0; i < list.length; i++ {
		if list.head.val == word {
			return true
		} else {
			list.head = list.head.next
		}
	}
	return false
}
