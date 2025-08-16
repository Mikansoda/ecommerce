package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"
	"fmt"

	"marketplace/entity"
	"marketplace/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password, role string) error
	VerifyOTP(ctx context.Context, email, otp string) error
	Login(ctx context.Context, email, password string) (access string, refresh string, err error)
	Refresh(ctx context.Context, refreshToken string) (access string, newRefresh string, err error)
	Logout(ctx context.Context, accessToken string) error
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, username, email, password, role string) error {
	// cek sudah ada belum
	_, err := s.repo.FindByEmail(ctx, email)
	switch {
	case err == nil:
		return errors.New("email already registered")
	case !errors.Is(err, gorm.ErrRecordNotFound):
    return err
}

// cek username udah kepake belum
	_, err = s.repo.FindByUsername(ctx, username)
	if err == nil {
		return errors.New("username already taken")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// hash password
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// generate OTP plaintext & hash
	otpPlain := generateNumericOTP(6)
	otpHash, _ := bcrypt.GenerateFromPassword([]byte(otpPlain), bcrypt.DefaultCost)
	exp := time.Now().Add(10 * time.Minute)

	u := &entity.UsersAndAdmins{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(pwHash),
		Role:         entity.RoleUser,
		IsActive:     false,
		OTPHash:      string(otpHash),
		OTPExpiresAt: &exp,
	}

	fmt.Printf("DEBUG: Saving user -> Email: %s, Username: %s, Role: %s\n", u.Email, u.Username, u.Role)
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}

	// kirim email OTP
	return SendEmailOTP(email, otpPlain)
}

func (s *authService) VerifyOTP(ctx context.Context, email, otp string) error {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}
	if u.OTPExpiresAt == nil || time.Now().After(*u.OTPExpiresAt) {
		return errors.New("otp expired")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.OTPHash), []byte(otp)); err != nil {
		return errors.New("invalid otp")
	}
	u.IsActive = true
	u.OTPHash = ""
	u.OTPExpiresAt = nil
	return s.repo.Update(ctx, u)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if !u.IsActive {
		return "", "", errors.New("account not active, verify otp")
	}
	access, _, err := GenerateAccessToken(u.ID.String(), u.Email, string(u.Role))
	if err != nil {
		return "", "", err
	}
	refresh, rexp, err := GenerateRefreshToken(u.ID.String(), u.Email, string(u.Role))
	if err != nil {
		return "", "", err
	}
	// simpan refresh hash
	rh, _ := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	u.RefreshTokenHash = string(rh)
	u.RefreshExpiresAt = &rexp
	if err := s.repo.Update(ctx, u); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, email, _, err := parseRefresh(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}
	if u.RefreshExpiresAt == nil || time.Now().After(*u.RefreshExpiresAt) {
		return "", "", errors.New("refresh expired")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.RefreshTokenHash), []byte(refreshToken)); err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// generate token baru
	access, _, err := GenerateAccessToken(claims.UserID, u.Email, string(u.Role))
	if err != nil {
		return "", "", err
	}
	newRefresh, rexp, err := GenerateRefreshToken(claims.UserID, u.Email, string(u.Role))
	if err != nil {
		return "", "", err
	}
	rh, _ := bcrypt.GenerateFromPassword([]byte(newRefresh), bcrypt.DefaultCost)
	u.RefreshTokenHash = string(rh)
	u.RefreshExpiresAt = &rexp
	if err := s.repo.Update(ctx, u); err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

// In-memory blacklist
var accessBlacklist = make(map[string]time.Time)

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	// blacklist sampai expiry claimnya
	claims, err := parseAccess(accessToken)
	if err != nil {
		return errors.New("invalid token")
	}
	exp := claims.ExpiresAt.Time
	accessBlacklist[accessToken] = exp
	return nil
}

func generateNumericOTP(n int) string {
	// pakai base32 untuk randomness, terus ambil digit
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "000000"
	}
	code := strings.ToUpper(base32.StdEncoding.EncodeToString(b))
	// ambil digit aja
	digs := make([]rune, 0, n)
	for _, c := range code {
		if c >= '0' && c <= '9' {
			digs = append(digs, c)
			if len(digs) == n {
				break
			}
		}
	}
	// fallback kalau kurang digit
	for len(digs) < n {
		digs = append(digs, '0')
	}
	return string(digs)
}
