package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMeetingsService_List(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	userID := "userID123"
	var meetingID1 int64 = 1
	var meetingID2 int64 = 2
	meetingType := 2
	meetingsRes := &MeetingsListResponse{
		Meetings: []*MeetingsListItem{
			{
				ID:        meetingID1,
				JoinURL:   "https://zoom.us/j/123?pwd=xyz",
				StartTime: now,
				Type:      meetingType,
			},
			{
				ID:        meetingID2,
				JoinURL:   "https://zoom.us/j/456?pwd=abc",
				StartTime: now,
				Type:      meetingType,
			},
		},
	}

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodGet)
		assert.Equal(req.URL.String(), fmt.Sprintf("/users/%s/meetings?type=%v", userID, meetingType))

		b, _ := json.Marshal(meetingsRes)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	m := &MeetingsService{
		client: zoomClient,
	}

	listOpts := &MeetingsListOptions{
		Type: Ptr(strconv.Itoa(meetingType)),
	}
	actualMeetingsRes, res, err := m.List(context.Background(), userID, listOpts)

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Len(actualMeetingsRes.Meetings, 2)
	assert.Equal(actualMeetingsRes.Meetings[0].ID, meetingsRes.Meetings[0].ID)
	assert.Equal(actualMeetingsRes.Meetings[1].ID, meetingsRes.Meetings[1].ID)
}

func TestMeetingsService_Create(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	userID := "userID123"
	var meetingID int64 = 1
	meetingsRes := &MeetingsCreateResponse{
		ID:        meetingID,
		JoinURL:   "https://zoom.us/j/123?pwd=xyz",
		Password:  "xyz",
		StartTime: now,
		Type:      2,
	}

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodPost)
		assert.Equal(req.URL.String(), fmt.Sprintf("/users/%s/meetings", userID))

		b, _ := json.Marshal(meetingsRes)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	m := &MeetingsService{
		client: zoomClient,
	}

	createOpts := &MeetingsCreateOptions{
		DefaultPassword: Ptr(true),
		Duration:        Ptr(int(time.Duration.Minutes(30))),
		Settings: &MeetingsCreateOptionsSettings{
			JBHTime:        Ptr(0),
			JoinBeforeHost: Ptr(true),
		},
		StartTime: Ptr(MeetingsCreateOptionsStartTime(now)),
		Type:      Ptr(2),
	}
	actualMeetingsRes, res, err := m.Create(context.Background(), userID, createOpts)

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal(meetingID, actualMeetingsRes.ID)
	assert.Equal(meetingsRes.JoinURL, actualMeetingsRes.JoinURL)
	assert.Equal(meetingsRes.Password, actualMeetingsRes.Password)
	assert.Equal(createOpts.Type, Ptr(actualMeetingsRes.Type))
	assert.WithinDuration(time.Time(*createOpts.StartTime), actualMeetingsRes.StartTime, time.Duration(1*time.Second))
}

func TestMeetingsService_Delete(t *testing.T) {
	assert := assert.New(t)

	var meetingID int64 = 1

	h := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(req.Method, http.MethodDelete)
		assert.Equal(req.URL.String(), fmt.Sprintf("/meetings/%v", meetingID))

		w.Header().Add("Content-Type", "application/json")
	}

	zoomClient, testServer := testClient(h)
	defer testServer.Close()

	m := &MeetingsService{
		client: zoomClient,
	}

	res, err := m.Delete(context.Background(), meetingID, &MeetingsDeleteOptions{})

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}
