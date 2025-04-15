package models

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidRole  = errors.New("invalid role")

	ErrUserNotFound = errors.New("user not found")

	ErrFloorNotFound = errors.New("floor not found")

	ErrBookingEntityNotFound = errors.New("booking entity not found")

	ErrBookingNotFound    = errors.New("booking not found")
	ErrAlreadyHaveBooking = errors.New("already have booking")
	ErrNoFreePlaces       = errors.New("no free places")
	ErrNoAccessToBooking  = errors.New("no access to booking")
	ErrInvalidBookingTime = errors.New("invalid booking time")

	ErrNoRights = errors.New("no rights")

	ErrOrderNotFound = errors.New("order not found")

	ErrGuestNotFounc       = errors.New("guest not found")
	ErrGuestsLimitAchieved = errors.New("guests limit achieved")
)
