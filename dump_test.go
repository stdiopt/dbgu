package dbgu_test

import (
	"os"

	"github.com/stdiopt/dbgu"
)

func ExampleDumper_Var() {
	a := Thing{
		A: 1,
		B: "2",
		C: &Thing{},
		D: map[string]Thing{
			"a": Thing{A: 2},
		},
	}
	dbgu.Dump(dbgu.OptNoColor).Var(os.Stdout, a)
	// Output:
	// A: 1
	// B: "2"
	// C.A: 0
	// C.B: ""
	// C.C: (*dbgu_test.Thing)(nil)
	// C.D: map[string]dbgu_test.Thing(nil)
	// D: map[string]dbgu_test.Thing{"a":dbgu_test.Thing{A:2, B:"", C:(*dbgu_test.Thing)(nil), D:map[string]dbgu_test.Thing(nil)}}

}

func ExampleDumper_Var_json() {
	a := Thing{
		A: 1,
		B: "2",
		C: &Thing{},
		D: map[string]Thing{
			"a": Thing{A: 2},
		},
	}
	dbgu.Dump(dbgu.OptJSON, dbgu.OptNoColor).Var(os.Stdout, a)
	// Output:
	// {
	//     "A": 1,
	//     "B": "2",
	//     "C": {
	//       "A": 0,
	//       "B": "",
	//       "C": null,
	//       "D": null
	//     },
	//     "D": {
	//       "a": {
	//         "A": 2,
	//         "B": "",
	//         "C": null,
	//         "D": null
	//       }
	//     }
	//   }

}

func ExampleDumper_Diff() {
	a := Thing{}
	b := Thing{
		A: 1,
		B: "2",
		C: &Thing{},
		D: map[string]Thing{
			"a": Thing{A: 2},
		},
	}
	dbgu.Dump(dbgu.OptNoColor).Diff(os.Stdout, a, b)
	// Output:
	// A        From 0 <- 1
	// B        From "" <- "2"
	// C.A      From nil <- 0
	// C.B      From nil <- ""
	// D["a"].A From nil <- 2
	// D["a"].B From nil <- ""
}

type Thing struct {
	A int
	B string
	C *Thing
	D map[string]Thing
}
