package util

func ExampleOutJSON() {
	testMap := map[string]string{"hello": "world"}
	OutJSON(testMap)
	// Output: {
	//	"hello": "world"
	// }
}

func ExampleOutElasticDSL() {
	testMap := map[string]string{"hello": "world"}
	OutElasticDSL(testMap, nil)
	// Output: {
	//	"query": {
	//		"hello": "world"
	//	}
	// }
}