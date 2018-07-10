package downcase

func Downcase(str string) (string, error) {
	res:= ""
	for _, rune := range str{
		if 65 <= rune && rune <= 90{
			res+=string(rune + 32) 
		}else{
			res+=string(rune)
		}
  }

	return res, nil
}
