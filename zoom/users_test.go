package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersService_List(t *testing.T) {
	assert := assert.New(t)

	roleID := "foo"
	userType := 2
	usersRes := &UsersListResponse{
		Users: []*UsersListItem{
			{
				ID:     "id1",
				RoleID: roleID,
				Email:  "user1@example.com",
				Type:   userType,
			},
			{
				ID:     "id2",
				RoleID: roleID,
				Email:  "user2@example.com",
				Type:   userType,
			},
		},
	}

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodGet)
		assert.Equal(req.URL.String(), "/users?role_id=foo")

		b, _ := json.Marshal(usersRes)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	u := &UsersService{
		client: zoomClient,
	}

	listOpts := &UsersListOptions{
		RoleID: Ptr(roleID),
	}
	actualUsersRes, res, err := u.List(context.Background(), listOpts)

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Len(actualUsersRes.Users, 2)
	assert.Equal(actualUsersRes.Users[0].ID, usersRes.Users[0].ID)
	assert.Equal(actualUsersRes.Users[1].ID, usersRes.Users[1].ID)
}

func TestUsersService_Create(t *testing.T) {
	assert := assert.New(t)

	userType := 2
	email := "eleanor@example.com"
	firstName := "Eleanor"
	lastName := "Roosevelt"

	usersRes := &UsersCreateResponse{
		ID:        "id1",
		Email:     email,
		Type:      userType,
		FirstName: firstName,
		LastName:  lastName,
	}

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodPost)
		assert.Equal(req.URL.String(), "/users")

		b, _ := json.Marshal(usersRes)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	u := &UsersService{
		client: zoomClient,
	}

	createOpts := &UsersCreateOptions{
		UserInfo: &UsersCreateOptionsUserInfo{
			Email:     email,
			Type:      userType,
			FirstName: Ptr(firstName),
			LastName:  Ptr(lastName),
		},
	}
	actualUsersRes, res, err := u.Create(context.Background(), createOpts)

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal(actualUsersRes.ID, usersRes.ID)
	assert.Equal(actualUsersRes.Email, usersRes.Email)
	assert.Equal(actualUsersRes.FirstName, usersRes.FirstName)
	assert.Equal(actualUsersRes.LastName, usersRes.LastName)
}

func TestUsersService_Delete(t *testing.T) {
	assert := assert.New(t)

	userID := "id1"

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodDelete)
		assert.Equal(req.URL.String(), fmt.Sprintf("/users/%s", userID))

		w.Header().Add("Content-Type", "application/json")
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	u := &UsersService{
		client: zoomClient,
	}

	res, err := u.Delete(context.Background(), userID, &UsersDeleteOptions{})

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}
