package cipher

//import "fmt"
//import "math"

type Cipher interface {
	Encode(string) string
	Decode(string) string
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
    }else if rune < 97 || rune > 122 {
			continue
		}else{
      res+=string(rune)
    }
	}
	return res
}

/* CAESAR CIPHER */

type CaesarCipher struct{
}

func encodeCaesar(str string) string{
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

func decodeCaesar(str string) string{
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
	encStr:=encodeCaesar(prepStr)
	return encStr 
}

func (ciph *CaesarCipher) Decode(str string) string {
	return decodeCaesar(str)
}

func NewCaesar() Cipher{
	var res Cipher = &CaesarCipher{} 
	return res
}

/* SHIFT CIPHER*/

type ShiftCipher struct{
	shift int
}

func encodeShift(str string, n int) string{
	res:=""
	for _, r := range str{
		r = r+rune(n)
		if(r >= 123){
		  r = r - 123 + 97
	  }else if(r <= 96){
       r = 122 - (96 - r)
     }

		res+=string(r)
	}
	return res
}

func decodeShift(str string, n int) string{
  res:=""
  for _, r := range str{
    r = r-rune(n)
    if(r <= 96){
      r = 122 - (96 - r)
    }else if(r >= 123){
      r = 97 + r - 123   
		}
    res+=string(r)
  }
  return res
 }

func (this *ShiftCipher) Encode(str string) string {
	prepStr:=prep(str)
	encStr:=encodeShift(prepStr, this.shift)
	return encStr 
}

func (this *ShiftCipher) Decode(str string) string {
	return decodeShift(str, this.shift)
}

func NewShift(n int) Cipher{
	if(!(n >= -25 && n <= 25) || n == 0){
		return nil
	}
	var res Cipher = &ShiftCipher{n}
	return res
}

/*VINGERIE CIPHER*/
type VingerieCipher struct{
	key string
}

func strIdx(i int, n int){

}

func extendKey(str string, n int) string{
	res := make([]rune, n)
	l := len(str)
	for i:= range res {
		res[i] = rune(str[i%l]);
	}
	return string(res)
}

func letterNum(c rune) int {
	return int(c - 97)
}
func encodeVingerie(str string, key string) string{
	res := make([]rune, len(str))
	extKey := extendKey(key, len(str))

	for i,r := range extKey{
		res[i]=[]rune(encodeShift(string(r),letterNum(rune(str[i]))))[0]
	}

	return string(res)
}

func decodeVingerie(str string, key string) string{
   res := make([]rune, len(str))
   extKey := extendKey(key, len(str))

   for i,r := range str{
		 sub := float64(r) - float64(extKey[i])
		  res[i] = rune('a' + sub)
			if(sub < 0){
				res[i] =rune (int('z') + int(sub) + 1)
			}

   }

   return string(res)
 }

func (this *VingerieCipher) Encode(str string) string {
  prepStr:=prep(str)
  encStr:=encodeVingerie(prepStr, this.key)
  return encStr
}

func (this *VingerieCipher) Decode(str string) string {
  return decodeVingerie(str, this.key)
}

func validate(str string) bool{
	if str == ""{
		return false
	}
	onlyA := true
  
	for _, rune := range str{
     switch rune{
       case ' ':
         return false
       case ',':
         return false
       case '.':
         return false
       case '!':
         return false
       case '-':
         return false
       case '@':
         return false
       case '#':
         return false
     }
     if 65 <= rune && rune <= 90{
			 return false
     }else if rune < 97 || rune > 122 {
       return false
     }
		 if(rune != 'a'){
		 	 onlyA = false
		 }
   }

	 if onlyA{
	 	 return false
	 }

	return true
}

func NewVigenere(str string) Cipher{
	if !validate(str){
		return nil
	}
	var res Cipher = &VingerieCipher{str}
	return res 
	
}
