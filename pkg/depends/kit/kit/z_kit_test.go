package kit

import "fmt"

func ExampleMetadata() {
	ma := Metadata{}
	ma.Add("a", "Av1")
	ma.Add("a", "Av2")

	mb := Metadata{}
	mb.Set("b", "Bv1")

	all := FromMetas(ma, mb)

	results := []interface{}{
		all.String(),
		all.Has("a"),
		all.Get("a"),
	}

	all.Del("b")

	results = append(results,
		all.Get("b"),
		all.String(),
	)

	for _, r := range results {
		fmt.Printf("%v\n", r)
	}
	// Output:
	// a=Av1&a=Av2&b=Bv1
	// true
	// Av1
	//
	// a=Av1&a=Av2
}
