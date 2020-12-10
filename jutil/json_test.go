package util

func ExampleOutJSON() {
	testMap := map[string]string{"hello": "world"}
	OutJSON(testMap)
	// Output: {
	//	"hello": "world"
	// }
}


