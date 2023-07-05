package client

import g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"

func MethodReceiver(typ string, rcv string) *g.SnippetField {
	if rcv == "" {
		return g.Var(g.Type(typ))
	}
	return g.Var(g.Type(typ), rcv)
}

func MethodStarReceiver(typ string, rcv string) *g.SnippetField {
	if rcv == "" {
		return g.Var(g.Star(g.Type(typ)))
	}
	return g.Var(g.Star(g.Type(typ)), rcv)
}

func SnippetReturnListOfInvokeMethod(f *g.File, rt g.SnippetType) []*g.SnippetField {
	lst := make([]*g.SnippetField, 0, 3)
	if rt != nil {
		lst = append(lst, g.Var(g.Star(rt)))
	}
	lst = append(lst,
		g.Var(g.Type(f.Use(PkgKit, "Metadata"))),
		g.Var(g.Error),
	)
	return lst
}
