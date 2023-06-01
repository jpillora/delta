package delta_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/jpillora/delta"
)

type Foo struct {
	ID     string `json:"id"`
	Wizzle bool   `json:"wizzle"`
}
type Order struct {
	ID    string `json:"id"`
	Items []int  `json:"items"`
	Foos  []Foo  `json:"foos"`
}

func simpleOrders() (*Order, *Order) {
	a := &Order{
		ID:    "1234",
		Items: []int{1, 2, 3, 4, 5},
	}
	b := &Order{
		ID:    "1234",
		Items: []int{1, 2, 4, 5},
	}
	return a, b
}

func TestSimplePatch(t *testing.T) {
	a, b := simpleOrders()
	_, err := delta.CopyPatch(b, a)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("a != b,\n a=%+v,\n b=%+v", a, b)
	}
}

func TestSimpleJSONPatch(t *testing.T) {
	a, b := simpleOrders()
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	ops, err := delta.CopyPatch(b, a, delta.OptionApply(false))
	if err != nil {
		t.Fatal(err)
	}
	opsJSON, _ := json.Marshal(ops)
	patch, err := jsonpatch.DecodePatch(opsJSON)
	if err != nil {
		panic(err)
	}
	jp, err := patch.Apply(ja)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(jb, jp) {
		t.Fatalf("jb != jp,\n jb=%+v,\n jp=%+v", jb, jp)
	}
}

func TestTable(t *testing.T) {

	type tc struct {
		src, dst *Order
	}

	a, b := simpleOrders()

	for i, tc := range []tc{
		// test case
		{
			src: a,
			dst: b,
		},
		// test case
		{
			src: &Order{
				ID:    "1234",
				Items: []int{1, 2, 3, 4},
				Foos: []Foo{
					{ID: "a", Wizzle: true},
					{ID: "b", Wizzle: false},
				},
			},
			dst: &Order{
				ID:    "1234",
				Items: []int{1, 2, 4},
				Foos: []Foo{
					{ID: "c", Wizzle: true},
					{ID: "d", Wizzle: false},
				},
			},
		},
		// test case
		{
			src: &Order{
				ID:    "1",
				Items: []int{1, 2, 3, 4, 6, 7},
				Foos: []Foo{
					{ID: "a", Wizzle: true},
				},
			},
			dst: &Order{
				ID:    "55",
				Items: []int{1, 2, 4, 5, 7},
				Foos: []Foo{
					{ID: "aa", Wizzle: true},
				},
			},
		},
	} {
		// marshal json before we mutate
		js, _ := json.Marshal(tc.src)
		jd, _ := json.Marshal(tc.dst)
		// copy src into dst, and compute the patch
		ops, err := delta.CopyPatch(tc.dst, tc.src)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tc.src, tc.dst) {
			t.Fatalf("[#%d] struct-not-equal\n src=%+v,\n dst=%+v", i, tc.src, tc.dst)
		}
		// we will marshal ops to send to the client
		jops, err := json.Marshal(ops)
		if err != nil {
			t.Fatal(err)
		}
		// client will perform the json patch
		patch, err := jsonpatch.DecodePatch(jops)
		if err != nil {
			t.Fatal(err)
		}
		jp, err := patch.Apply(js)
		if err != nil {
			panic(err)
		}
		// patched json should match original dst json
		if !bytes.Equal(jd, jp) {
			t.Fatalf("[#%d] json-not-equal\n dst=%s,\n jp=%s", i, string(jd), string(jp))
		}
	}
}
