package middleware

import (
	"context"
)

func ProjectProviderFromContext(ctx context.Context) string {
	return ctx.Value("ProjectProvider").(*ProjectProvider).ProjectName
}

// ProjectProvider with account id prefix
type ProjectProvider struct {
	ProjectName string `name:"projectName" in:"path" validate:"@projectName"`
}

func (ProjectProvider) ContextKey() string { return "ProjectProvider" }

func (ProjectProvider) Path() string { return "/:id" }

func (r *ProjectProvider) Output(ctx context.Context) (interface{}, error) {
	a := CurrentAccountFromContext(ctx)
	return &ProjectProvider{
		ProjectName: a.AccountID.String() + "_" + r.ProjectName,
	}, nil
}
