package middlewares

import (
	"crm-service-go/app/clients"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(resource clients.PolicyResource, act clients.AuthorizationAction) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		user := LoggedUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User hasn't logged in yet"})
			return
		}

		entitlementClient := clients.NewEntitlementClient()
		entitlementClient.Client.SetToken(c.GetHeader("Authorization"))
		entitlements, err := entitlementClient.FindPolicy(user.Sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
			return
		}

		policies := entitlements.Policies
		policy := findPolicyByResource(policies, resource)
		subjects := findSubjectByAction(policy, act)
		if len(subjects) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
			return
		}

		storeCodes := utils.Map[*clients.SubjectSchema, string](subjects, func(item *clients.SubjectSchema) string {
			return item.Name
		})

		c.Set(pkg.ProjectKeyUserStoreCodes, storeCodes)
		c.Set(pkg.ProjectKeyUserRole, entitlements.Roles)
		c.Set(pkg.ProjectKeyUserPolicy, subjects)
		c.Next()
	}
}

func findPolicyByResource(policies []*clients.PolicySchema, resource clients.PolicyResource) *clients.PolicySchema {
	for _, policy := range policies {
		if policy.Resource == resource {
			return policy
		}
	}
	return nil
}

func findSubjectByAction(policy *clients.PolicySchema, action clients.AuthorizationAction) []*clients.SubjectSchema {
	var subjects []*clients.SubjectSchema
	for _, subject := range policy.Subject {
		for _, item := range subject.Actions {
			if item == action {
				if subject.Name == "*" {
					return []*clients.SubjectSchema{&subject}
				}
				subjects = append(subjects, &subject)
			}
		}
	}
	return subjects
}
