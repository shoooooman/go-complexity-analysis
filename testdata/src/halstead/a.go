package halstead

func f1() { // want "Cyclomatic complexity: 1, Halstead difficulty: 2.500, volume: 18.095"
	print("Hello, World")
}

func f2() { // want "Cyclomatic complexity: 1, Halstead difficulty: 6.857, volume: 101.579"
	a := 2
	b := 1
	c := 3
	avg := (a + b + c) / 3
	println(avg)
}

func f3() { // want "Cyclomatic complexity: 2, Halstead difficulty: NaN, volume: 25.266"
	if false {

	} else {

	}
}

func f4() { // want "Cyclomatic complexity: 9, Halstead difficulty: 10.833, volume: 144.000"
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

type t1 struct {
}

func (t *t1) f5() { // want "Cyclomatic complexity: 1, Halstead difficulty: NaN, volume: 10.000"
}
