package esutil

func ExampleOutDSL() {
	testMap := map[string]string{"hello": "world"}
	OutDSL(testMap, nil)
	// Output: {
	//    "query": {
	//       "hello": "world"
	//    }
	// }
}
