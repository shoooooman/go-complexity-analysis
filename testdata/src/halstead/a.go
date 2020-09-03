package a

func f1() { // Want "Halstead complexity: 4.754888, 1.000000"
	print("Hello, World")
}

func f2() {
	a := 2
	b := 1
	c := 3
	avg := (a + b + c) / 3
	println(avg)
}

func f3() {
	if false {

	} else {

	}
}

func f4() {
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
