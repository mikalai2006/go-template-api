package service

import (
	"errors"
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/auths"
	"github.com/mikalai2006/go-template-api/pkg/hasher"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	hasher       hasher.PasswordHasher
	tokenManager auths.TokenManager

	repository   repository.Authorization
	otpGenerator utils.Generator

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	verificationCodeLength int
}

func NewAuthService(
	repo repository.Authorization,
	hasherPkg hasher.PasswordHasher,
	tokenManager auths.TokenManager,
	refreshTokenTTL time.Duration,
	accessTokenTTL time.Duration,
	otp utils.Generator,
	verificationCodeLength int,
) *AuthService {
	return &AuthService{
		hasher:       hasherPkg,
		tokenManager: tokenManager,

		repository:   repo,
		otpGenerator: otp,

		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,

		verificationCodeLength: verificationCodeLength,
	}
}

func (s *AuthService) CreateAuth(auth *domain.SignInInput) (string, error) {
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return "", err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	authData := &domain.Auth{
		VkID:      auth.VkID,
		GoogleID:  auth.GoogleID,
		GithubID:  auth.GithubID,
		AppleID:   auth.AppleID,
		Roles:     []string{"user"},
		Login:     auth.Login,
		Password:  passwordHash,
		Email:     auth.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Strategy:  "local",
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	id, err := s.repository.CreateAuth(authData)
	if err != nil {
		return "", err
	}

	// if created auth, send email with verification code

	return id, nil
}

func (s *AuthService) ExistAuth(auth *domain.SignInInput) (domain.Auth, error) {
	return s.repository.CheckExistAuth(auth)
}

func (s *AuthService) SignIn(auth *domain.SignInInput) (domain.ResponseTokens, error) {
	var result domain.ResponseTokens
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return result, err
	}
	auth.Password = passwordHash

	// fmt.Println("sign in ", auth)
	user, err := s.repository.GetByCredentials(auth)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, err
		}
		return result, err
	}

	return s.CreateSession(&user)
}

func (s *AuthService) CreateSession(auth *domain.Auth) (domain.ResponseTokens, error) {
	var (
		res domain.ResponseTokens
		err error
	)

	claims := domain.DataForClaims{
		Roles:  auth.Roles,
		UserID: auth.ID.Hex(),
		UID:    auth.UserData.ID.Hex(),
	}

	res.AccessToken, err = s.tokenManager.NewJWT(claims, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repository.SetSession(auth.ID, session)

	return res, err
}

func (s *AuthService) VerificationCode(userID, hash string) error {
	err := s.repository.VerificationCode(userID, hash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (domain.ResponseTokens, error) {
	var result domain.ResponseTokens

	user, err := s.repository.RefreshToken(refreshToken)
	if err != nil {
		return result, err
	}

	return s.CreateSession(&user)
}
