package esutil

import (
	"encoding/json"
	"log"
	"os"
)

// OutDSL OutDSL(query.Source())
func OutDSL(data interface{}, err error) {
	if err != nil {
		log.Printf("query.Source() has err, err is %v \n", err)
		return
	}

	m := map[string]interface{}{"query": data}

	b, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		log.Printf("OutElasticDSL MarshalIndent fail, err is %v \n", err)
		return
	}
	_, _ = os.Stdout.Write(b)
}
