package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/mailer"
	"github.com/memsbdm/restaurant-api/internal/repository"
	"github.com/memsbdm/restaurant-api/pkg/keys"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

var (
	ErrEmailConflict        = errors.New("email already taken")
	ErrEmailAlreadyVerified = errors.New("email already verified")
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (dto.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (dto.UserDTO, error)
	Create(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, error)
	Update(ctx context.Context, user *dto.UserDTO) (dto.UserDTO, error)
	SendVerificationEmail(ctx context.Context, user dto.UserDTO) error
	VerifyEmail(ctx context.Context, token string) (dto.UserDTO, error)
	ResendVerificationEmail(ctx context.Context, userID uuid.UUID) error
}

type userService struct {
	cfg       *config.App
	userRepo  repository.UserRepository
	mailerSvc MailerService
	tokenSvc  TokenService
}

func NewUserService(cfg *config.App, userRepo repository.UserRepository, tokenSvc TokenService, mailerSvc MailerService) *userService {
	return &userService{
		cfg:       cfg,
		userRepo:  userRepo,
		mailerSvc: mailerSvc,
		tokenSvc:  tokenSvc,
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (dto.UserDTO, error) {
	dbUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (dto.UserDTO, error) {
	dbUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}

func (s *userService) Create(ctx context.Context, user *dto.CreateUserDto) (dto.UserDTO, error) {
	emailTaken, err := s.userRepo.EmailTaken(ctx, user.Email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	if emailTaken {
		return dto.UserDTO{}, ErrEmailConflict
	}

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return dto.UserDTO{}, err
	}

	user.Password = hashedPassword

	dbUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}

func (s *userService) Update(ctx context.Context, user *dto.UserDTO) (dto.UserDTO, error) {
	dbUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(dbUser), nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, user dto.UserDTO) error {
	spt, err := s.tokenSvc.GenerateSPT(ctx, keys.EmailVerification, user.ID.String(), keys.EmailVerificationTokenDuration)
	if err != nil {
		return err
	}

	emailtmpl, err := s.mailerSvc.RenderTemplate("verify_email.tmpl", map[string]any{
		"Host":  s.cfg.Host,
		"User":  user,
		"Token": spt,
	})
	if err != nil {
		return err
	}

	return s.mailerSvc.Send(&mailer.Mail{
		To:      []string{user.Email},
		Subject: "Verify your email",
		Body:    emailtmpl,
	})
}

func (s *userService) VerifyEmail(ctx context.Context, token string) (dto.UserDTO, error) {
	decodedToken, err := s.tokenSvc.VerifySPT(ctx, keys.EmailVerification, token)
	if err != nil {
		return dto.UserDTO{}, err
	}

	userID := decodedToken
	dbUser, err := s.userRepo.GetByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return dto.UserDTO{}, err
	}

	dbUser.IsEmailVerified = true
	dbUserDto := dto.NewUserDTO(dbUser)
	updatedUser, err := s.userRepo.Update(ctx, &dbUserDto)
	if err != nil {
		return dto.UserDTO{}, err
	}

	err = s.tokenSvc.RevokeSPT(ctx, keys.EmailVerification, userID)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.NewUserDTO(updatedUser), nil
}

func (s *userService) ResendVerificationEmail(ctx context.Context, userID uuid.UUID) error {
	dbUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if dbUser.IsEmailVerified {
		return ErrEmailAlreadyVerified
	}

	return s.SendVerificationEmail(ctx, dto.NewUserDTO(dbUser))
}
