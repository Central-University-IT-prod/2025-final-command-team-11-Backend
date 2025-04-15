package guest

import (
	"time"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

type Guest struct {
	guest      GuestStorage
	bookEntity EntityStorage
	book       BookingStorage
	id         IdUseCase
}

func New(guest GuestStorage, bookEntity EntityStorage, book BookingStorage, id IdUseCase) *Guest {
	return &Guest{
		guest:      guest,
		book:       book,
		bookEntity: bookEntity,
		id:         id,
	}
}

func (g *Guest) Create(c ctx.Context, bookId, email, owner string) e.Error {
	user, err := g.id.GetUser(c, email)
	if err != nil {
		return err
	}

	booking, err := g.book.GetById(c, bookId)
	if err != nil {
		return err
	}

	if booking.UserId != owner {
		return e.New("Forbidden.", e.Forbidden)
	}

	ent, err := g.bookEntity.GetEntity(c, booking.EntityId)
	if err != nil {
		return err
	}

	if ent.Type == types.OPENSPACE {
		return e.New("You can`t invite peoples to open space", e.BadInput)
	}

	_, err = g.guest.GetById(c, bookId, user.Id)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err == nil {
		return nil
	}

	guests, err := g.guest.Get(c, booking.Id)
	if err != nil {
		return err
	}

	if len(guests)+1 >= ent.Capacity {
		return e.New("Room is full.", e.BadInput)
	}

	guest := &entity.Guest{
		UserId:    user.Id,
		BookingId: bookId,
		CreatedAt: time.Now().UTC(),
	}

	return g.guest.Create(c, guest)
}

func (g *Guest) Get(c ctx.Context, bookingId string, userId string) ([]*entity.Guest, e.Error) {
	booking, err := g.book.GetById(c, bookingId)
	if err != nil {
		return nil, err
	}

	if booking.UserId != userId {
		return nil, e.New("Forbidden.", e.Forbidden)
	}

	guests, err := g.guest.Get(c, bookingId)
	if err != nil {
		return nil, err
	}

	for _, guest := range guests {
		user, err := g.id.GetUserById(c, guest.UserId)
		if err != nil {
			return nil, err
		}

		guest.UserId = user.Email
	}

	return guests, nil
}

func (g *Guest) Delete(c ctx.Context, bookId, email, owner string) e.Error {
	user, err := g.id.GetUser(c, email)
	if err != nil {
		return err
	}

	booking, err := g.book.GetById(c, bookId)
	if err != nil {
		return err
	}

	if booking.UserId != owner {
		return e.New("Forbidden.", e.Forbidden)
	}

	return g.guest.Delete(c, bookId, user.Id)
}
