package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store/teststore"
	"testing"
)

func TestServer_HandleCreateUser(t *testing.T) {
	s := newServer(teststore.New())

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"name": "Name",
				"email": "current@mail.com",
				"password": "123456",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "empty payload",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"name": "Name",
				"email": "nenormalmail",
				"password": "12345",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid password",
			payload: map[string]interface{}{
				"name": "Name",
				"email": "normal@mail.com",
				"password":"12345",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleAuthorizeUser(t *testing.T) {
	s := newServer(teststore.New())

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email": "test@example.com",
				"password": "password",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid payload",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid user",
			payload: map[string]interface{}{
				"email": "test@example.com",
				"password": "",
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := map[string]interface{}{
				"name": "User",
				"email": "test@example.com",
				"password": "password",
			}

			resp := struct {
				Token string `json:"token"`
			}{}

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(user)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)

			json.NewEncoder(b).Encode(tc.payload)
			rec = httptest.NewRecorder()
			req, _ = http.NewRequest(http.MethodPost, "/authorize", b)
			s.ServeHTTP(rec, req)

			json.NewDecoder(rec.Body).Decode(resp)
			recToken := resp.Token

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotNil(t, recToken)
		})
	}
}

func TestServer_HandleCreateArticle(t *testing.T) {
	s := newServer(teststore.New())

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"article_header": "Test Article",
				"article_text": "Article test text",
				"author_id": 1,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid article text",
			payload: map[string]interface{}{
				"article_header": "Test Article",
				"article_text": 7,
				"author_id": 1,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			us := map[string]interface{}{
				"name": "User",
				"email": "user@mail.com",
				"password": "123456",
			}

			resp := &struct {
				User *model.User `json:"user"`
				Token string `json:"token"`
			}{}

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(us)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)

			json.NewDecoder(rec.Body).Decode(resp)

			token := fmt.Sprintf("Baerer %s", resp.Token)

			json.NewEncoder(b).Encode(tc.payload)
			req, _ = http.NewRequest(http.MethodPost, "/private/create/article", b)
			rec = httptest.NewRecorder()
			req.Header.Add("Authorization", token)
			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleFindArticleByHeading(t *testing.T) {
	ts := teststore.New()
	ts.Article().CreateArticle(model.TestArticle(t, 5))

	s := newServer(ts)

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"article_heading": "TestArticle",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid article header",
			payload: map[string]interface{}{
				"article_heading": "Article",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid request",
			payload: map[string]interface{}{
				"article_heading": 1,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/find/article", b)
			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleChangeArticle(t *testing.T) {
	ts := teststore.New()
	article := model.TestArticle(t, 5)
	ts.Article().CreateArticle(article)
	s := newServer(ts)

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"id" : article.ID,
				"article_header": "Updated TestArticle",
				"article_text": "updated article text",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid id",
			payload: map[string]interface{}{
				"id" : 1,
				"article_header": "Updated TestArticle",
				"article_text": "updated article text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			us := map[string]interface{}{
				"name": "User",
				"email": "user@mail.com",
				"password": "123456",
			}

			resp := &struct {
				User *model.User `json:"user"`
				Token string `json:"token"`
			}{}

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(us)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)

			json.NewDecoder(rec.Body).Decode(resp)
			token := fmt.Sprintf("Baerer %s", resp.Token)

			json.NewEncoder(b).Encode(tc.payload)
			req, _ = http.NewRequest(http.MethodPut, "/private/change/article", b)
			rec = httptest.NewRecorder()
			req.Header.Add("Authorization", token)
			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleDeleteArticle(t *testing.T) {
	ts := teststore.New()
	article := model.TestArticle(t, 5)
	ts.Article().CreateArticle(article)
	s := newServer(ts)

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"id" : article.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid id",
			payload: map[string]interface{}{
				"id" : 1,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			us := map[string]interface{}{
				"name": "User",
				"email": "user@mail.com",
				"password": "123456",
			}

			resp := &struct {
				User *model.User `json:"user"`
				Token string `json:"token"`
			}{}

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(us)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)

			json.NewDecoder(rec.Body).Decode(resp)
			token := fmt.Sprintf("Baerer %s", resp.Token)

			json.NewEncoder(b).Encode(tc.payload)
			req, _ = http.NewRequest(http.MethodDelete, "/private/delete/article", b)
			rec = httptest.NewRecorder()
			req.Header.Add("Authorization", token)
			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}


