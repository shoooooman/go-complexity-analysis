package a

func f1() { // want "Cyclomatic complexity: 1"
	print("Hello, World")
}

func f2() { // want "Cyclomatic complexity: 1"
	a := 2
	b := 1
	c := 3
	avg := (a + b + c) / 3
	println(avg)
}

func f3() { // want "Cyclomatic complexity: 2"
	if false {

	} else {

	}
}

func f4() { // want "Cyclomatic complexity: 9"
	for true {
		if false {

		} else if false {

		} else if false {

		} else if false {
			n := 0
			switch n {
			case 0:
			case 1:
			default:
			}
		} else {

		}
	}
}
