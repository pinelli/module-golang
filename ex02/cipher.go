package cipher

//import "fmt"

type Cipher interface {
	Encode(string) string
	Decode(string) string
}

type CaesarCipher struct{
}

func prep(str string) string{
	res:= ""
  for _, rune := range str{
		switch rune{
			case ' ':
				continue
			case ',':
				continue
			case '.':
				continue
			case '!':
				continue
			case '-':
				continue
  		case '@':
				continue
			case '#':
				continue
		}
    if 65 <= rune && rune <= 90{
      res+=string(rune + 32)
    }else{
      res+=string(rune)
    }
	}
	return res
}

func encode(str string) string{
	res:=""
	for _, r := range str{
		r = r+3
		if(r >= 123){
		  r = r - 123 + 97
	  }
		res+=string(r)
	}
	return res
}

func decode(str string) string{
  res:=""
  for _, r := range str{
    r = r-3
    if(r <= 96){
      r = 122 - (96 - r) 
    }
    res+=string(r)
  }
  return res
}


func (ciph *CaesarCipher) Encode(str string) string {
	prepStr:=prep(str)
	encStr:=encode(prepStr)
	return encStr 
}

func (ciph *CaesarCipher) Decode(str string) string {
	return decode(str)
}

func NewCaesar() Cipher{
	var res Cipher = &CaesarCipher{} 
	return res
}


/*
func NewShift(n int) Cipher{
	return nil
}

func NewVigenere(str string) Cipher{
	return nil
}
*/
