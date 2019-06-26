package dbgu_test

import (
	"os"
	"time"

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
	// C.E: (*time.Time)(nil)
	// D["a"].A: 2
	// D["a"].B: ""
	// D["a"].C: (*dbgu_test.Thing)(nil)
	// D["a"].E: (*time.Time)(nil)
	// E: (*time.Time)(nil)

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
	//       "D": null,
	//       "E": null
	//     },
	//     "D": {
	//       "a": {
	//         "A": 2,
	//         "B": "",
	//         "C": null,
	//         "D": null,
	//         "E": null
	//       }
	//     },
	//     "E": null
	//   }

}

func ExampleDumper_Diff() {
	now := time.Date(2019, 01, 01, 10, 10, 10, 10, &time.Location{})
	a := Thing{}
	b := Thing{
		A: 1,
		B: "2",
		C: &Thing{},
		D: map[string]Thing{
			"a": Thing{A: 2},
		},
		E: &now,
	}
	dbgu.Dump(dbgu.OptNoColor).Diff(os.Stdout, a, b)

	// Output:
	// A        From 0 <- 1
	// B        From "" <- "2"
	// C.A      From nil <- 0
	// C.B      From nil <- ""
	// D["a"].A From nil <- 2
	// D["a"].B From nil <- ""
	// E        From nil <- 2019-01-01 10:10:10.00000001 +0000 UTC
}

type Thing struct {
	A int
	B string
	C *Thing
	D map[string]Thing
	E *time.Time
}
