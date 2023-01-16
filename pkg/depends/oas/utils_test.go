package oas_test

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func NewCaseGroup(name string) *group {
	return &group{
		name: name,
	}
}

type group struct {
	name string
	list []*caseItem
}

func (g *group) It(desc string, result string, v interface{}) {
	g.list = append(g.list, &caseItem{
		desc:   desc,
		result: result,
		value:  v,
	})
}

func (g *group) ItCase(item *caseItem) {
	g.list = append(g.list, item)
}

func (g *group) Run(t *testing.T) {
	for i := range g.list {
		item := g.list[i]
		data, errForMarshal := json.Marshal(item.value)
		if errForMarshal != nil {
			t.Logf("[%s] %s, marshal failed, results: %s", g.name, item.desc, string(data))
		}
		NewWithT(t).Expect(errForMarshal).To(BeNil())
		NewWithT(t).Expect(item.result).To(Equal(string(data)))
		// assert.Nil(t, errForMarshal)
		// assert.Equal(t,
		// 	item.result,
		// 	string(data),
		// 	fmt.Sprintf("[%s] %s, marshal failed, results: %s", g.name, item.desc, string(data)),
		// )

		expectRv := reflect.Indirect(reflect.ValueOf(item.value))
		value := reflect.New(expectRv.Type()).Interface()
		errForUnmarshal := json.Unmarshal([]byte(item.result), value)
		// assert.Nil(t, errForUnmarshal)
		if errForUnmarshal != nil {
			t.Logf("%s", errForUnmarshal.Error())
		}
		NewWithT(t).Expect(errForMarshal).To(BeNil())
		value = reflect.Indirect(reflect.ValueOf(value)).Interface()
		NewWithT(t).Expect(expectRv.Interface()).To(Equal(value))
		// assert.Equal(t,
		// 	expectRv.Interface(),
		// 	value,
		// 	fmt.Sprintf("[%s] %s, unmarshal failed)", g.name, item.desc),
		// )
	}
}

type caseItem struct {
	value  interface{}
	result string
	desc   string
}
