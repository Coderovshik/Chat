package user

import (
	"context"
	"strconv"
	"time"

	"github.com/Coderovshik/chat_server/internal/config"
	"github.com/Coderovshik/chat_server/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       UserRepository
	timeout    time.Duration
	tokenTTL   time.Duration
	signingKey string
}

func NewService(repo UserRepository, cfg *config.Config) *Service {
	return &Service{
		repo:       repo,
		timeout:    cfg.Timeout,
		tokenTTL:   cfg.TokenTTL,
		signingKey: cfg.SigningKey,
	}
}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	passhash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Passhash: string(passhash),
	}

	u, err = s.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:       strconv.Itoa(int(u.ID)),
		Email:    u.Email,
		Username: u.Username,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Passhash), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	claims := &util.UserClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "chat_app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
		},
	}
	ss, err := util.NewJWTSignedString(claims, s.signingKey)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		accessToken: ss,
	}, nil
}
