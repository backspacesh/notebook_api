package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store"
	"strings"
)

type server struct {
	router *mux.Router
	store store.Store
}

func newServer(store store.Store) *server {
	srv := &server {
		router: mux.NewRouter(),
		store: store,
	}

	srv.configureRouter()

	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/hello", s.hello()).Methods("GET")
	s.router.HandleFunc("/create", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/find/article", s.handleFindArticleByHeading()).Methods("GET")
	s.router.HandleFunc("/show_all_articles", s.handleShowAllArticles()).Methods("GET")
	s.router.HandleFunc("/authorize", s.handleAuthorizeUser()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.JwtAuthentication)
	private.HandleFunc("/create/article", s.handleCreateArticle()).Methods("POST")
	private.HandleFunc("/delete/article", s.handleDeleteArticle()).Methods("DELETE")
	private.HandleFunc("/change/article", s.handleChangeArticle()).Methods("PUT")
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User *model.User `json:"user"`
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name: req.Name,
			Email: req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		tk := &model.Token{ID: u.ID}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

		resp := response{
			User: u,
			Token: tokenString,
		}

		s.respond(w, r, http.StatusCreated, resp)
	}
}

func (s *server) handleAuthorizeUser() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := u.ComparePassword(req.Password); err != nil {
			s.error(w, r, http.StatusForbidden, errors.New("invalid email or password"))
			return
		}

		tk := &model.Token{ID: u.ID}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_string")))

		s.respond(w, r, http.StatusOK, tokenString)
	}
}

func (s *server) handleCreateArticle() http.HandlerFunc {
	type request struct {
		ArticleHeader string `json:"article_header"`
		ArticleText string `json:"article_text"`
		AuthorID int `json:"author_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		a := &model.Article{
			Heading: req.ArticleHeader,
			Text: req.ArticleText,
			AuthorID: req.AuthorID,
			Date: "",
		}

		if err := s.store.Article().CreateArticle(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, a)
	}
}

func (s *server) handleFindArticleByHeading() http.HandlerFunc {
	type request struct {
		Header string `json:"article_heading"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		ar, err := s.store.Article().FindByHeading(req.Header)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		s.respond(w, r, http.StatusOK, ar)
	}
}

func (s *server) handleShowAllArticles() http.HandlerFunc {
	type articles struct {
		AricleList []*model.Article `json:"articles"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ars, err := s.store.Article().ShowAllArticles()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		arts := &articles{
			AricleList: ars,
		}

		s.respond(w, r, http.StatusOK, arts)
	}
}

func (s *server) handleChangeArticle() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
		ArticleHeader string `json:"article_header"`
		ArticleText string `json:"article_text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		ar := &model.Article{
			ID: req.ID,
			Heading: req.ArticleHeader,
			Text: req.ArticleText,
		}

		if err := s.store.Article().ChangeArticleById(ar); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		ar, err := s.store.Article().FindByHeading(ar.Heading)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, ar)
	}
}

func (s *server) handleDeleteArticle() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
	}

	type response struct {
		Message string `json:"message"`
	}

	return func (w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		h, err := s.store.Article().DeleteArticle(req.ID)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		resp := &response{
			Message: fmt.Sprintf("Deleted article: %s", h),
		}

		s.respond(w, r, http.StatusOK, resp)
	}
}

func (s *server) hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	}
}

func (s *server) JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			s.respond(w, r, http.StatusUnauthorized, map[string]interface{}{"error": "Missing auth token!"})
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			s.respond(w, r, http.StatusBadRequest, map[string]interface{}{"error": "Invalid/Malformed auth token!"})
			return
		}

		tokenPart := splitted[1]

		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}

		if !token.Valid {
			s.respond(w, r, http.StatusUnauthorized, map[string]interface{}{"error": "Token is not valid!"})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error":err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
