package modelgen

import (
	"go/types"
	"reflect"
	"regexp"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

func forEachStructField(st *types.Struct, each func(v *types.Var, name, tag string)) {
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Exported() {
			continue
		}
		tag, exists := reflect.StructTag(st.Tag(i)).Lookup("db")
		if exists && tag != "-" {
			each(f, builder.GetColumnName(f.Name(), tag), tag)
			continue
		}
		if f.Anonymous() {
			if embed, ok := f.Type().Underlying().(*types.Struct); ok {
				forEachStructField(embed, each)
			}
		}
	}
}

var regexpRelAnnotate = regexp.MustCompile(`@rel ([^\n]+)`)

func parseColRelFromDoc(doc string) (string, []string) {
	others := make([]string, 0)
	rel := ""
	for _, line := range strings.Split(doc, "\n") {
		if len(line) == 0 {
			continue
		}
		matches := regexpRelAnnotate.FindAllStringSubmatch(line, 1)
		if matches == nil {
			others = append(others, line)
			continue
		}
		if len(matches) == 1 {
			rel = matches[0][1]
		}
	}
	return rel, others
}

var regexpDefAnnotate = regexp.MustCompile(`@def ([^\n]+)`)

func parseKeysFromDoc(doc string) (*Keys, []string) {
	keys := &Keys{}
	others := make([]string, 0)

	for _, line := range strings.Split(doc, "\n") {
		if len(line) == 0 {
			continue
		}

		matches := regexpDefAnnotate.FindAllStringSubmatch(line, -1)

		if matches == nil {
			others = append(others, line)
			continue
		}

		for _, sub := range matches {
			if len(sub) == 2 {
				def := builder.ParseIndexDefine(sub[1])
				switch def.Kind {
				case "primary":
					keys.Primary = def.ToDefs()
				case "unique_index":
					if keys.UniqueIndexes == nil {
						keys.UniqueIndexes = builder.Indexes{}
					}
					keys.UniqueIndexes[def.ID()] = def.ToDefs()
				case "index":
					if keys.Indexes == nil {
						keys.Indexes = builder.Indexes{}
					}
					keys.Indexes[def.ID()] = def.ToDefs()
				}
			}
		}
	}

	return keys, others
}

func uniqueStrings(lst []string) (res []string) {
	m := make(map[string]bool)
	for _, s := range lst {
		m[s] = true
	}
	for _, s := range lst {
		if _, ok := m[s]; ok {
			delete(m, s)
			res = append(res, s)
		}
	}
	return
}

func filterStrings(lst []string, chk func(s string, i int) bool) []string {
	left, _ := partedStrings(lst, chk)
	return left
}

func partedStrings(lst []string, chk func(s string, i int) bool) ([]string, []string) {
	newLs, newRs := make([]string, 0), make([]string, 0)
	for i, s := range lst {
		if chk(s, i) {
			newLs = append(newLs, s)
		} else {
			newRs = append(newRs, s)
		}
	}
	return newLs, newRs
}
