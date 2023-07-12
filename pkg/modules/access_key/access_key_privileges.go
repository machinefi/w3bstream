package access_key

import (
	"net/http"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
)

type OperatorMeta struct {
	OperatorID string
	Summary    string
	Method     string
	Attr       enums.ApiOperatorAttr
}

type WithOperatorAttr interface {
	OperatorAttr() enums.ApiOperatorAttr
}

type GroupMetaBase struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type GroupMeta struct {
	GroupMetaBase
	Operators []*OperatorMeta
}

type GroupAccessPrivilege struct {
	Name string                 `json:"name"`
	Perm enums.AccessPermission `json:"perm"`
}

type GroupAccessPrivileges []GroupAccessPrivilege

func (gaps GroupAccessPrivileges) AccessPrivileges() models.AccessPrivileges {
	ret := make(models.AccessPrivileges)

	for i := range gaps {
		p := &gaps[i]
		meta, ok := gOperatorGroups[p.Name]
		if !ok {
			continue
		}

		switch p.Perm {
		case enums.ACCESS_PERMISSION__READ_WRITE:
			for _, op := range meta.Operators {
				ret[op.OperatorID] = struct{}{}
			}
		case enums.ACCESS_PERMISSION__READONLY:
			for _, op := range meta.Operators {
				if op.Method == http.MethodGet {
					ret[op.OperatorID] = struct{}{}
				}
			}
		default:
			continue
		}
	}
	return ret
}

var (
	gOperatorGroups = map[string]*GroupMeta{}

	gOperatorGroupMetas         []*GroupMetaBase
	gOperatorGroupMetasInitOnce = &sync.Once{}
)

func RouterRegister(r *kit.Router, name, desc string) {
	if _, ok := gOperatorGroups[name]; ok {
		panic(errors.Errorf("operator group: %s already registered", name))
	}

	routes := r.Routes()
	group := &GroupMeta{
		GroupMetaBase: GroupMetaBase{
			Name: name,
			Desc: desc,
		},
	}

	for _, route := range routes {
		factories := httptransport.NewHttpRouteMeta(route).Metas

		fact := factories[len(factories)-1]
		op := &OperatorMeta{
			OperatorID: fact.Type.Name(),
			Summary:    fact.Summary,
			Method:     fact.Method,
			Attr:       enums.API_OPERATOR_ATTR__COMMON,
		}

		if with, ok := fact.Operator.(WithOperatorAttr); ok {
			op.Attr = with.OperatorAttr()
		}
		group.Operators = append(group.Operators, op)
	}
	gOperatorGroups[name] = group
}

func OperatorGroupMetaList() []*GroupMetaBase {
	gOperatorGroupMetasInitOnce.Do(func() {
		for _, meta := range gOperatorGroups {
			v := meta.GroupMetaBase
			gOperatorGroupMetas = append(gOperatorGroupMetas, &v)
		}
		sort.Slice(gOperatorGroupMetas, func(i, j int) bool {
			return gOperatorGroupMetas[i].Name < gOperatorGroupMetas[j].Name
		})
	})
	return gOperatorGroupMetas
}
