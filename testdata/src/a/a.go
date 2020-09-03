package a

func f0() { // want "Cyclomatic complexity: 1"
}

func f1() { // want "Cyclomatic complexity: 2"
	if false {

	} else {

	}
}

func f2() { // want "Cyclomatic complexity: 11"
	for true {
		if false {

		} else if false {

		} else if false {

		} else if false {

		} else if false {
			n := 0
			switch n {
			case 0:
			case 1:
			case 2:
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
