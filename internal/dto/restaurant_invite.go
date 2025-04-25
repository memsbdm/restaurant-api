package dto

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantInvite struct {
	ID               int        `json:"id"`
	InvitedByUserID  uuid.UUID  `json:"invited_by_user_id"`
	CanceledByUserID *uuid.UUID `json:"canceled_by_user_id"`
	RoleID           int        `json:"role_id"`
	Email            string     `json:"email"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	AcceptedAt       *time.Time `json:"accepted_at"`
	CanceledAt       *time.Time `json:"canceled_at"`
	InvitedByUser    User       `json:"invited_by_user"`
	CanceledByUser   *User      `json:"canceled_by_user"`
	Role             Role       `json:"role"`
	RestaurantID     uuid.UUID  `json:"restaurant_id"`
	Restaurant       Restaurant `json:"restaurant"`
}
