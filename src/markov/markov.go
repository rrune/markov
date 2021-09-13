package markov

import (
	"encoding/gob"
	//. "fmt"
	"io/ioutil"
	"math/rand"
	"os"
	. "rune/markov/util"
	"strings"
)

type chain struct {
	Links map[string]link
}

type link struct {
	CanFollow     []string
	EndOfSentence bool
}

func New() *chain {
	return &chain{make(map[string]link)}
}

func (c chain) Read(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	Check(err)
	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(&c)
	Check(err)
}

func (c chain) Write(filePath string) {
	err := os.Remove(filePath)
	Check(err)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0)
	Check(err)
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(c)
	Check(err)
}

func (c chain) Train(filePath string) {
	file, err := ioutil.ReadFile(filePath)
	Check(err)
	words := strings.Fields(string(file))

	for i, word := range words {
		var next string
		if len(words) > i+1 {
			next = words[i+1]
			next = strings.TrimSuffix(next, ".")
		}

		endOfSentence := false
		if word[len(word)-1] == '.' {
			endOfSentence = true
			word = strings.TrimSuffix(word, ".")
			next = ""
		}

		if _, ok := c.Links[word]; !ok {
			arr := []string{next}
			c.Links[word] = link{arr, endOfSentence}
		} else {
			tempLink := link{c.Links[word].CanFollow, c.Links[word].EndOfSentence}
			tempLink.CanFollow = append(tempLink.CanFollow, next)
			c.Links[word] = tempLink
		}
	}
}

func (c chain) getWord() string {
	keys := make([]string, 0, len(c.Links))
	for key := range c.Links {
		keys = append(keys, key)
	}
	return keys[rand.Intn(len(keys))]
}

func (c chain) GetChain(limit int) string {
	return c.GetChainStartWord(c.getWord(), limit)
}

func (c chain) GetChainStartWord(startWord string, limit int) string {
	result := "" + startWord
	for i := 0; i < limit; i++ {
		if (len(c.Links[startWord].CanFollow)) != 0 {
			nextWord := c.Links[startWord].CanFollow[rand.Intn(len(c.Links[startWord].CanFollow))]
			if nextWord != "" {
				result += " " + nextWord
			} else {
				result += nextWord
			}

			if c.Links[startWord].EndOfSentence && rand.Intn(2) == 1 {
				result += "."
				startWord = c.getWord()
			} else {
				startWord = nextWord
			}
		} else {
			//force new sentence
			if result[len(result)-1] != '.' {
				result += "."
			}
			startWord = c.getWord()
		}
	}
	return result
}
