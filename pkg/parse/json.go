package parse

import (
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/gunparse/pkg"
	"github.com/mattfenwick/gunparse/pkg/example"
)

func JsonObject(input string) pkg.Result[example.ParseError, *pkg.Pair[int, int], rune, *example.Object] {
	return example.ObjectParser.Parse(example.StringToRunes(input), pkg.NewPair[int, int](1, 1))
}

func SerializeObject(obj *example.Object) string {
	//simpleTree := SerializeObjectHelper(obj)
	return json.MustMarshalToString(obj)
}

//func SerializeObjectHelper(obj *example.Object) map[string]interface{} {
//	out := map[string]interface{}{}
//	for _, kv := range obj.Body {
//		out[*kv.Key] = SerializeValue(kv.Value)
//	}
//	return out
//}
