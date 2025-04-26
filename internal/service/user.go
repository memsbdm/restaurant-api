package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/database"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/mailer"
	"github.com/memsbdm/restaurant-api/pkg/keys"
	"github.com/memsbdm/restaurant-api/pkg/security"
)

var (
	ErrEmailConflict        = errors.New("email already taken")
	ErrEmailAlreadyVerified = errors.New("email already verified")
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (dto.User, error)
	GetByEmail(ctx context.Context, email string) (dto.User, error)
	Create(ctx context.Context, user *dto.CreateUser) (dto.User, error)
	Update(ctx context.Context, user *dto.User) (dto.User, error)
	SendVerificationEmail(ctx context.Context, user dto.User) error
	VerifyEmail(ctx context.Context, token string) (dto.User, error)
	ResendVerificationEmail(ctx context.Context, userID uuid.UUID) error
	GetRestaurantsByUserID(ctx context.Context, userID uuid.UUID) ([]dto.Restaurant, error)
}

type userService struct {
	cfg       *config.App
	db        *database.DB
	mailerSvc MailerService
	tokenSvc  TokenService
}

func NewUserService(cfg *config.App, db *database.DB, tokenSvc TokenService, mailerSvc MailerService) *userService {
	return &userService{
		cfg:       cfg,
		db:        db,
		mailerSvc: mailerSvc,
		tokenSvc:  tokenSvc,
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (dto.User, error) {
	dbUser, err := s.db.Queries.GetUserByID(ctx, id)
	if err != nil {
		return dto.User{}, err
	}

	return dto.NewUser(&dbUser), nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (dto.User, error) {
	dbUser, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.User{}, err
	}

	return dto.NewUser(&dbUser), nil
}

func (s *userService) Create(ctx context.Context, user *dto.CreateUser) (dto.User, error) {
	emailTaken, err := s.db.Queries.UserEmailTaken(ctx, user.Email)
	if err != nil {
		return dto.User{}, err
	}

	if emailTaken {
		return dto.User{}, ErrEmailConflict
	}

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return dto.User{}, err
	}

	user.Password = hashedPassword

	dbUser, err := s.db.Queries.CreateUser(ctx, user.ToParams())
	if err != nil {
		return dto.User{}, err
	}

	return dto.NewUser(&dbUser), nil
}

func (s *userService) Update(ctx context.Context, user *dto.User) (dto.User, error) {
	dbUser, err := s.db.Queries.UpdateUser(ctx, user.ToUpdateParams())
	if err != nil {
		return dto.User{}, err
	}

	return dto.NewUser(&dbUser), nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, user dto.User) error {
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

func (s *userService) VerifyEmail(ctx context.Context, token string) (dto.User, error) {
	decodedToken, err := s.tokenSvc.VerifySPT(ctx, keys.EmailVerification, token)
	if err != nil {
		return dto.User{}, err
	}

	userID := decodedToken
	dbUser, err := s.db.Queries.GetUserByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return dto.User{}, err
	}

	dbUser.IsEmailVerified = true
	dbUserDTO := dto.NewUser(&dbUser)
	updatedUser, err := s.db.Queries.UpdateUser(ctx, dbUserDTO.ToUpdateParams())
	if err != nil {
		return dto.User{}, err
	}

	err = s.tokenSvc.RevokeSPT(ctx, keys.EmailVerification, userID)
	if err != nil {
		return dto.User{}, err
	}

	return dto.NewUser(&updatedUser), nil
}

func (s *userService) ResendVerificationEmail(ctx context.Context, userID uuid.UUID) error {
	dbUser, err := s.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if dbUser.IsEmailVerified {
		return ErrEmailAlreadyVerified
	}

	return s.SendVerificationEmail(ctx, dto.NewUser(&dbUser))
}

func (s *userService) GetRestaurantsByUserID(ctx context.Context, userID uuid.UUID) ([]dto.Restaurant, error) {
	dbRestaurants, err := s.db.Queries.GetRestaurantsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	restaurants := make([]dto.Restaurant, len(dbRestaurants))
	for i, dbRestaurant := range dbRestaurants {
		restaurants[i] = dto.NewRestaurant(&dbRestaurant)
	}

	return restaurants, nil
}
