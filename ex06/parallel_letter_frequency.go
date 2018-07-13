package letter

import "sync"

func add(c rune, letters *map[rune]int, mutex *sync.Mutex) {
	(*mutex).Lock()

	_, ok := (*letters)[c]
	if !ok {
		(*letters)[c] = 1
	} else {
		(*letters)[c] += 1
	}

	(*mutex).Unlock()
}

func Count(text string, letters *map[rune]int, mutex *sync.Mutex, chanel chan bool) {
	for _, c := range []rune(text) {
		add(rune(c), letters, mutex)
	}
	chanel <- true
}

func Frequency(text string) map[rune]int {
	var letters map[rune]int = make(map[rune]int)
	var mutex = &sync.Mutex{}
	chanel := make(chan bool, 3)

	Count(text, &letters, mutex, chanel)
	return letters
}

func ConcurrentFrequency(texts []string) map[rune]int {
	var letters map[rune]int = make(map[rune]int)
	var mutex = sync.Mutex{}

	l := len(texts)
	chanel := make(chan bool, 3)

	for _, str := range texts {
		go Count(str, &letters, &mutex, chanel)
	}

	for i := 0; i < l; i++ {
		<-chanel
	}
	return letters
}
