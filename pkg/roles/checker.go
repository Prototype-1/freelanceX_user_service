package role

import (
	"context"
)
type Checker interface {
	HasRole(ctx context.Context, allowedRoles ...string) bool
}

type checkerImpl struct{}

func NewChecker() Checker {
	return &checkerImpl{}
}

func (c *checkerImpl) HasRole(ctx context.Context, allowedRoles ...string) bool {
	roleVal := ctx.Value("role")
	if roleVal == nil {
		return false
	}
	userRole, ok := roleVal.(string)
	if !ok {
		return false
	}

	for _, allowed := range allowedRoles {
		if userRole == allowed {
			return true
		}
	}
	return false
}