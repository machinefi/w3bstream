package builder

import "sort"

type Addition interface {
	SqlExpr
	AdditionKind() AdditionType
}

type AdditionType int

const (
	AdditionJoin AdditionType = iota
	AdditionWhere
	AdditionGroupBy
	AdditionCombination
	AdditionOrderBy
	AdditionLimit
	AdditionOnConflict
	AdditionOther
	AdditionComment
)

func (t AdditionType) AdditionKind() AdditionType { return t }

var (
	_ Addition = (*join)(nil)
	_ Addition = (*where)(nil)
	_ Addition = (*groupby)(nil)
	_ Addition = (*orderby)(nil)
	_ Addition = (*addition)(nil)
	_ Addition = (*limit)(nil)
	_ Addition = (*comment)(nil)
)

type Additions []Addition

func (adds Additions) Len() int { return len(adds) }

func (adds Additions) Less(i, j int) bool {
	return adds[i].AdditionKind() < adds[j].AdditionKind()
}

func (adds Additions) Swap(i, j int) { adds[i], adds[j] = adds[j], adds[i] }

func WriteAdditions(e *Ex, adds ...Addition) {
	final := make(Additions, 0, len(adds))
	for i := range adds {
		if IsNilExpr(adds[i]) {
			continue
		}
		final = append(final, adds[i])
	}
	if len(final) == 0 {
		return
	}
	sort.Sort(final)
	for i := range final {
		e.WriteQueryByte('\n')
		e.WriteExpr(final[i])
	}
}
