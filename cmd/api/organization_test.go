package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/LincolnG4/Haku/internal/auth"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func TestOrganizationAuthorization(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	// Setup users and organization in mock store
	org := &store.Organization{ID: 1, Name: "TestOrg"}
	adminUser := &store.User{ID: 10, Username: "admin", Email: "admin@example.com"}
	memberUser := &store.User{ID: 20, Username: "member", Email: "member@example.com"}
	nonMemberUser := &store.User{ID: 30, Username: "nonmember", Email: "nonmember@example.com"}

	app.store.Organizations.Create(nil, org)
	app.store.Users.Create(nil, adminUser)
	app.store.Users.Create(nil, memberUser)
	app.store.Users.Create(nil, nonMemberUser)

	// Add admin and member to org
	app.store.Organizations.AddMember(nil, &store.OrganizationMember{UserID: adminUser.ID, OrganizationID: org.ID, RoleID: store.AdminRole})
	app.store.Organizations.AddMember(nil, &store.OrganizationMember{UserID: memberUser.ID, OrganizationID: org.ID, RoleID: 2}) // 2 = not admin

	// Helper to get token for a user
	getToken := func(user *store.User) string {
		claims := &auth.MyClaims{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:  fmt.Sprintf("%d", user.ID),
				Issuer:   "test-aud",
				Audience: []string{"test-aud"},
			},
		}
		token, _ := app.authenticator.GenerateToken(claims)
		return token
	}

	// Test: Only members can get members
	t.Run("forbid non-member from getting members", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/organizations/1/members", nil)
		req.Header.Set("Authorization", "Bearer "+getToken(nonMemberUser))
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("allow member to get members", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/organizations/1/members", nil)
		req.Header.Set("Authorization", "Bearer "+getToken(memberUser))
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusOK, rr.Code)
	})

	// Test: Only admins can add members
	addPayload := map[string]interface{}{"user_id": nonMemberUser.ID, "role_id": 2}
	addBody, _ := json.Marshal(addPayload)

	t.Run("forbid non-admin from adding member", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/v1/organizations/1/members", bytes.NewReader(addBody))
		req.Header.Set("Authorization", "Bearer "+getToken(memberUser))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusForbidden, rr.Code)
	})

	t.Run("allow admin to add member", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/v1/organizations/1/members", bytes.NewReader(addBody))
		req.Header.Set("Authorization", "Bearer "+getToken(adminUser))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusCreated, rr.Code)
	})
}
