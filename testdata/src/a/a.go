package a

func f0() { // want "Cyclomatic complexity: 1"
}

func f1() { // want "Cyclomatic complexity: 2"
	if false {

	} else {

	}
}

func f2() { // want "Cyclomatic complexity: 9"
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

func f3() { // want "Cyclomatic complexity: 4"
	if false || true {
		if false {

		}
	}
}

func f4() { // want "Cyclomatic complexity: 11"
	n := 0
	switch n {
	case 0:
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	}
}
