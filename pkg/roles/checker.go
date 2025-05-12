package role

import (
    "context"
    "google.golang.org/grpc/metadata"
    "log"
)

type Checker interface {
    HasRole(ctx context.Context, allowedRoles ...string) bool
}

type checkerImpl struct{}

func NewChecker() Checker {
    return &checkerImpl{}
}

func (c *checkerImpl) HasRole(ctx context.Context, allowedRoles ...string) bool {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        log.Println("No metadata found in context")
        return false
    }
    
    roles := md.Get("role")
    if len(roles) == 0 {
        log.Println("No role found in metadata")
        return false
    }
    
    userRole := roles[0]
    log.Printf("User role from metadata: %s", userRole)
    
    for _, allowed := range allowedRoles {
        if userRole == allowed { 
            log.Printf("Role match found: %s", allowed)
            return true
        }
    }
    
    log.Printf("No matching role found. User role: %s, Allowed roles: %v", userRole, allowedRoles)
    return false
}