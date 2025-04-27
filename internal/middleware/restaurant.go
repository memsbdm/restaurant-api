package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/memsbdm/restaurant-api/internal/dto"
	"github.com/memsbdm/restaurant-api/internal/handler"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/memsbdm/restaurant-api/internal/service"
	"github.com/memsbdm/restaurant-api/pkg/keys"
)

func RestaurantMiddleware(appEnv string, restaurantSvc service.RestaurantService, restaurantUserSvc service.RestaurantUserService) Middle {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID, err := keys.GetUserIDFromContext(ctx)
			if err != nil {
				response.HandleError(w, response.ErrUnauthorized)
				return
			}

			// Check if the request has an active restaurant ID in the header or cookie
			restaurantID, err := extractActiveRestaurantIDFromRequest(r)
			if err == nil {
				// If the restaurant ID is present, check if the user belongs to that restaurant
				roleID, err := restaurantUserSvc.GetRestaurantUserRoleID(ctx, restaurantID, userID)
				if err == nil {
					// User belongs to the restaurant
					restaurant, err := restaurantSvc.GetByID(ctx, restaurantID)
					if err == nil {
						ctx = enrichContextWithRestaurantInfos(ctx, restaurant, roleID)
						sendActiveRestaurantToClient(w, r, appEnv, restaurant.ID)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			// If we reach here, it means we either didn't find the restaurant or the user doesn't belong to it
			// Check if the user has any restaurant linked
			restaurantUser, err := restaurantUserSvc.GetAnyRestaurantUserLinkByUserID(ctx, userID)
			if err != nil {
				if errors.Is(err, service.ErrRestaurantOrUserNotFound) {
					response.HandleError(w, service.ErrNoRestaurantFoundForUser)
					return
				}
				response.HandleError(w, err)
				return
			}

			// If the user has a restaurant linked, we proceed with that one
			restaurant, err := restaurantSvc.GetByID(ctx, restaurantUser.RestaurantID)
			if err != nil {
				response.HandleError(w, err)
				return
			}

			ctx = enrichContextWithRestaurantInfos(ctx, restaurant, restaurantUser.RoleID)
			sendActiveRestaurantToClient(w, r, appEnv, restaurant.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractActiveRestaurantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	var restaurantID uuid.UUID
	var err error

	if handler.IsMobileRequest(r) {
		restaurantID, err = getActiveRestaurantIDFromHeader(r)
	} else {
		restaurantID, err = getActiveRestaurantIDFromCookie(r)
	}

	return restaurantID, err
}

func getActiveRestaurantIDFromCookie(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie(keys.ActiveRestaurantCookieName)
	if err != nil {
		return uuid.Nil, errors.New("active restaurant cookie not found")
	}

	restaurantID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.Nil, errors.New("error during cookie restaurant id parsing to uuid")
	}

	return restaurantID, nil
}

func getActiveRestaurantIDFromHeader(r *http.Request) (uuid.UUID, error) {
	restaurantIDstr := r.Header.Get(keys.ActiveRestaurantHeaderName)
	if restaurantIDstr == "" {
		return uuid.Nil, errors.New("active restaurant header not found")
	}

	restaurantID, err := uuid.Parse(restaurantIDstr)
	if err != nil {
		return uuid.Nil, errors.New("error during header restaurant id parsing to uuid")
	}

	return restaurantID, nil
}

func enrichContextWithRestaurantInfos(ctx context.Context, restaurant *dto.Restaurant, userRoleID int) context.Context {
	ctx = context.WithValue(ctx, keys.RestaurantIDContextKey, restaurant.ID)
	ctx = context.WithValue(ctx, keys.RestaurantContextKey, restaurant)
	ctx = context.WithValue(ctx, keys.UserRoleIDContextKey, userRoleID)
	return ctx
}

func sendActiveRestaurantToClient(w http.ResponseWriter, r *http.Request, appEnv string, restaurantID uuid.UUID) {
	if handler.IsMobileRequest(r) {
		w.Header().Set(keys.ActiveRestaurantHeaderName, restaurantID.String())
		return
	}
	handler.SetActiveRestaurantCookie(w, restaurantID, appEnv)
}
