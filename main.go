package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

// Вызов go run . -decode="[5 #][5 -_]-[5 #]"
// Вызов go run . -encode="###---121212"
func main() {
	//fmt.Println(encodev3("-#-#-#=="))
	// Определяем флаг `task` с типом string и значением по умолчанию "World"
	decodeStr := flag.String("decode", "", "set")
	encodeStr := flag.String("encode", "", "set")

	// Обязательно вызываем функцию `flag.Parse()` для обработки переданных аргументов
	flag.Parse()

	//log.Printf("decodeStr=%#v", *decodeStr)
	//log.Printf("encodeStr=%#v", *encodeStr)

	if decodeStr != nil && *decodeStr != "" {
		fmt.Println(decode(*decodeStr))
	}

	if encodeStr != nil && *encodeStr != "" {
		fmt.Println(encodev3(*encodeStr))
	}
}

// 1. "[5 #][5 -_]-[5 #] => "#####-----_-#####"
func decode(input string) (output string) {
	inputRune := []rune(input)
	var inBacket bool //Внутри квадратных скобок
	for i := 0; i < len(inputRune); i++ {
		switch {
		case inputRune[i] == '[':
			inBacket = true
			backet, _, good := strings.Cut(string(inputRune[i+1:]), "]") //12 _#
			if !good {
				log.Fatalf("found [ but not found ]")
			}
			numStr, subStr, good := strings.Cut(backet, " ")
			if !good {
				log.Fatalf("\" \" not found in \"%v\"", backet)
			}
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Fatalf("%v is not a number", subStr)
			}
			for j := 0; j < num; j++ {
				output += subStr
			}
		case inputRune[i] == ']':
			inBacket = false
		case !inBacket:
			output += string(inputRune[i])
		}
	}
	return
}

// 2. "#####--_-_-_-_-_-#####" => "[5 #][5 -_]-[5 #]
// ######=######=  => "[2 ######=] [6 #]=[6 #]=
// #_$#_$  => "[2 #_$]
//func encode(input string) (output string) {
//	if len(input) == 0 {
//		return
//	}
//
//	inputRune := []rune(input)
//	count := 1
//	simbol := inputRune[0]
//	for i := 0; i < len(inputRune); i++ {
//		if i >= 1 {
//			if inputRune[i] == inputRune[i-1] {
//				count++
//			} else {
//				output += getOutputIter(count, simbol)
//				count = 1
//				simbol = inputRune[i]
//			}
//			if i == len(inputRune)-1 {
//				output += getOutputIter(count, simbol)
//			}
//		}
//	}
//	return
//}

//func getOutputIter(count int, simbol rune) string {
//	if count > 1 {
//		return "[" + strconv.Itoa(count) + " " + string(simbol) + "]"
//	} else {
//		return string(simbol)
//	}
//}

// И так, наш энкодер будет работать по следующему принципу:
// В первом цикле перебираем начальный символ
// Во втором цикле увеличиваем количество символов, детектируя дубли
// В третьем цикле детектируем наличие дублей в последующих символах
// При успехе во втором и третьем цикле - смещаем индекс первого цикла количество элементов в которых детектированы дубли, а в результат пишем количество и значение
// Если дошли до конца цикла - сохраняем символ в выход и переходим к следующему
func encodev3(input string) (output string) {
	inputRune := []rune(input)
	if len(inputRune) <= 1 { //Тут уверяемся что у нас не менее 2 символов
		return input
	}
	count := 1
	lMax := 10
	//log.Printf("len(inputRune)=%v", len(inputRune))
	for i := 0; i < len(inputRune); {
		//log.Printf("i=%v", i)
		var foundedPairs bool
		for j := i + 1; j+(j-i) <= len(inputRune) && j-i <= lMax; { //Тут делим пополам, потому что нам нужно искать пары: этот срез и следующий
			a1, b1 := i, j
			a2, b2 := j, j+(j-i)
			var iNext int
			//log.Printf("i=%v, j=%v, inputRune[%v:%v]=%#v, inputRune[%v:%v]=%#v", i, j, a1, b1, string(inputRune[a1:b1]), a2, b2, string(inputRune[a2:b2]))
			//Если нашли кобинацию при которой слайсы совпадают
			if slices.Equal(inputRune[a1:b1], inputRune[a2:b2]) {
				//log.Printf("equal_l1")
				count = 2 //Сразу записываем 2
				iNext = b2
				//И идем циклом по следующим значениям в поисках совпадений
				for k := b2; k <= len(inputRune); k += j - i { //Смещаемся на размер фрагмента, повторы которого ищем
					a3, b3 := k, k+(j-i)
					//log.Printf("i=%v, j=%v, k=%v, inputRune[%v:%v]=%#v, inputRune[%v:%v]=%#v", i, j, k, a1, b1, string(inputRune[a1:b1]), a3, b3, string(inputRune[a3:b3]))
					if slices.Equal(inputRune[a1:b1], inputRune[a3:b3]) {
						//log.Printf("equal_l2")
						count++
						continue
					} else {
						iNext = a3
						break
					}
				}
				output += fmt.Sprintf("[%v %v]", count, string(inputRune[i:j]))
				foundedPairs = true
				i = iNext
				break
			} else {
				j++
			}

		}
		if !foundedPairs {
			//Вот это пишем если ничего не нашли
			output += string(inputRune[i])
			i++
		}
	}
	return
}
