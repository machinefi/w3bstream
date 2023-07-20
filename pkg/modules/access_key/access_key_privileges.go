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
	OperatorID  string
	Summary     string
	Method      string
	Attr        enums.ApiOperatorAttr
	MinimalPerm enums.AccessPermission
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
	Operators map[string]*OperatorMeta
}

type GroupAccessPrivilege struct {
	Name string                 `json:"name"`
	Perm enums.AccessPermission `json:"perm"`
}

type GroupAccessPrivileges []GroupAccessPrivilege

func (gaps GroupAccessPrivileges) ConvToPrivilegeModel() models.GroupAccessPrivileges {
	ret := make(models.GroupAccessPrivileges)

	for k := range gOperatorGroups {
		ret[k] = enums.ACCESS_PERMISSION__NO_ACCESS
	}

	for i := range gaps {
		p := &gaps[i]
		if _, ok := ret[p.Name]; !ok {
			continue
		}

		switch p.Perm {
		case enums.ACCESS_PERMISSION__READ_WRITE, enums.ACCESS_PERMISSION__READONLY, enums.ACCESS_PERMISSION__NO_ACCESS:
			ret[p.Name] = p.Perm
		default:
			ret[p.Name] = enums.ACCESS_PERMISSION__NO_ACCESS
		}
	}
	return ret
}

func ConvToGroupMetaWithPrivileges(privileges models.GroupAccessPrivileges) []*GroupMetaWithPrivilege {
	ret := make([]*GroupMetaWithPrivilege, 0, len(privileges))
	for name, perm := range privileges {
		ret = append(ret, &GroupMetaWithPrivilege{
			GroupMetaBase: GroupMetaBase{
				Name: name,
				Desc: gOperatorGroups[name].Desc,
			},
			Perm: perm,
		})
	}
	return ret
}

type GroupMetaWithPrivilege struct {
	GroupMetaBase
	Perm enums.AccessPermission `json:"perm"`
}

var (
	// gOperatorGroups mapping group name and  group meta
	gOperatorGroups = map[string]*GroupMeta{}
	// gOperators mapping operator id and group name
	gOperators = map[string]string{}

	// gOperatorGroupMetas group meta list
	gOperatorGroupMetas         []*GroupMetaBase
	gOperatorGroupMetasInitOnce = &sync.Once{}
)

func RouterRegister(r *kit.Router, name, desc string) {
	if _, ok := gOperatorGroups[name]; ok {
		panic(errors.Errorf("group already registered: group[%s]", name))
	}

	routes := r.Routes()
	group := &GroupMeta{
		GroupMetaBase: GroupMetaBase{
			Name: name,
			Desc: desc,
		},
		Operators: map[string]*OperatorMeta{},
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

		switch op.Method {
		case http.MethodGet:
			op.MinimalPerm = enums.ACCESS_PERMISSION__READONLY
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			op.MinimalPerm = enums.ACCESS_PERMISSION__READ_WRITE
		default:
			continue
		}

		if with, ok := fact.Operator.(WithOperatorAttr); ok {
			op.Attr = with.OperatorAttr()
		}

		if groupName, ok := gOperators[op.OperatorID]; ok {
			panic(errors.Errorf("operator id already registered in group: operator[%s] group[%s]", op.OperatorID, groupName))
		}

		gOperators[op.OperatorID] = name
		group.Operators[op.OperatorID] = op
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
