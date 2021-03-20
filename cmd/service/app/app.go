package app

import (
	"context"
	"encoding/json"
	"github.com/artrey/payments/cmd/service/app/dto"
	"github.com/artrey/payments/cmd/service/app/middleware/authenticator"
	"github.com/artrey/payments/cmd/service/app/middleware/authorizator"
	"github.com/artrey/payments/cmd/service/app/middleware/identificator"
	"github.com/artrey/payments/pkg/business"
	"github.com/artrey/payments/pkg/security"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	securitySvc *security.Service
	businessSvc *business.Service
	router      chi.Router
}

func NewServer(securitySvc *security.Service, businessSvc *business.Service, router chi.Router) *Server {
	return &Server{securitySvc: securitySvc, businessSvc: businessSvc, router: router}
}

func (s *Server) Init() error {
	s.router.Post("/users", s.handleRegister)
	s.router.Put("/users", s.handleLogin)

	identificatorMd := identificator.Identificator
	authenticatorMd := authenticator.Authenticator(
		identificator.Identifier, s.securitySvc.UserDetails,
	)

	// функция-связка между middleware и security service
	// (для чистоты: security service ничего не знает об http)
	roleChecker := func(ctx context.Context, roles ...string) bool {
		userDetails, err := authenticator.Authentication(ctx)
		if err != nil {
			return false
		}
		return s.securitySvc.HasAnyRole(ctx, userDetails, roles...)
	}
	adminRoleMd := authorizator.Authorizator(roleChecker, security.RoleAdmin)
	userRoleMd := authorizator.Authorizator(roleChecker, security.RoleUser)

	s.router.Get("/public", s.handlePublic)
	s.router.With(identificatorMd, authenticatorMd, adminRoleMd).Get("/admin", s.handleAdmin)
	s.router.With(identificatorMd, authenticatorMd, userRoleMd).Get("/user", s.handleUser)

	s.router.With(identificatorMd, authenticatorMd, userRoleMd).Post("/user/payments", s.handleCreatePayment)
	s.router.With(identificatorMd, authenticatorMd, userRoleMd).Get("/user/payments", s.handleUserPayments)
	s.router.With(identificatorMd, authenticatorMd, adminRoleMd).Get("/admin/payments", s.handleAllPayments)

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleRegister(writer http.ResponseWriter, request *http.Request) {
	login := request.PostFormValue("login")
	if login == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	password := request.PostFormValue("password")
	if password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.securitySvc.Register(request.Context(), login, password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &dto.TokenDTO{Token: token}

	writeJson(writer, data, http.StatusCreated)
}

func (s *Server) handleLogin(writer http.ResponseWriter, request *http.Request) {
	login := request.PostFormValue("login")
	if login == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	password := request.PostFormValue("password")
	if password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.securitySvc.Login(request.Context(), login, password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &dto.TokenDTO{Token: token}

	writeJson(writer, data, http.StatusOK)
}

// Доступно всем
func (s *Server) handlePublic(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("public"))
	if err != nil {
		log.Print(err)
	}
}

// Только пользователям с ролью ADMIN
func (s *Server) handleAdmin(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("admin"))
	if err != nil {
		log.Print(err)
	}
}

// Только пользователям с ролью USER
func (s *Server) handleUser(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("user"))
	if err != nil {
		log.Print(err)
	}
}

// Только пользователям с ролью USER
func (s *Server) handleCreatePayment(writer http.ResponseWriter, request *http.Request) {
	amountString := request.PostFormValue("amount")
	if amountString == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseInt(amountString, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	auth, err := authenticator.Authentication(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	userDetail, ok := auth.(*security.UserDetails)
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	id, err := s.businessSvc.CreatePayment(request.Context(), userDetail.ID, amount)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := &dto.PaymentDTO{
		Id:       id,
		SenderId: userDetail.ID,
		Amount:   amount,
	}

	writeJson(writer, data, http.StatusCreated)
}

// Только пользователям с ролью USER
func (s *Server) handleUserPayments(writer http.ResponseWriter, request *http.Request) {
	auth, err := authenticator.Authentication(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	userDetail, ok := auth.(*security.UserDetails)
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	payments, err := s.businessSvc.GetUserPayments(request.Context(), userDetail.ID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make([]*dto.PaymentDTO, len(payments))
	for i, p := range payments {
		data[i] = dto.FromPaymentModel(p)
	}

	writeJson(writer, data, http.StatusOK)
}

// Только пользователям с ролью ADMIN
func (s *Server) handleAllPayments(writer http.ResponseWriter, request *http.Request) {
	payments, err := s.businessSvc.GetAllPayments(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make([]*dto.PaymentDTO, len(payments))
	for i, p := range payments {
		data[i] = dto.FromPaymentModel(p)
	}

	writeJson(writer, data, http.StatusOK)
}

func writeJson(writer http.ResponseWriter, data interface{}, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_, err = writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}
