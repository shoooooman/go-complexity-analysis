package a

func f() { // want "branch cnt: 8"
	for true {
		if false {

		} else if false {

		} else if false {

		} else if false {

		} else if false {

		} else if false {

		} else {

		}
	}
}
