package brackets

func Bracket(str string) (bool, error) {
	var s *Stack = New()

	for _, c := range str {

		if c == '{' {
			s.Push(1)
		} else if c == '[' {
			s.Push(2)
		} else if c == '(' {
			s.Push(3)
		} else if c == '}' {
			br, _ := s.Pop()
			if br != 1 {
				return false, nil
			}
		} else if c == ']' {
			br, _ := s.Pop()
			if br != 2 {
				return false, nil
			}
		} else if c == ')' {
			br, _ := s.Pop()
			if br != 3 {
				return false, nil
			}
		}
	}
	return len(*s) == 0, nil
}
