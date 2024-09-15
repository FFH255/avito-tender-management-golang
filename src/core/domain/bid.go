package domain

import (
	"github.com/pkg/errors"
	"slices"
	"time"
)

// BidName Полное название предложения
type BidName string

func NewBidName(str string) (BidName, error) {
	if len(str) == 0 {
		return "", errors.Wrap(ErrValidation, "bid name must not be empty")
	}
	if len(str) > 100 {
		return "", errors.Wrap(ErrValidation, "bid name must not exceed 100 characters")
	}
	return BidName(str), nil
}

// BidDescription Описание предложения
type BidDescription string

func NewBidDescription(str string) (BidDescription, error) {
	if len(str) == 0 {
		return "", errors.Wrap(ErrValidation, "bid description must not be empty")
	}
	if len(str) > 500 {
		return "", errors.Wrap(ErrValidation, "bid description must not exceed 500 characters")
	}
	return BidDescription(str), nil
}

// BidStatus Статус предложения
type BidStatus string

const (
	BidCreatedStatus   BidStatus = "Created"
	BidPublishedStatus BidStatus = "Published"
	BidCanceledStatus  BidStatus = "Canceled"
)

func NewBidStatus(str string) (BidStatus, error) {
	switch str {
	case string(BidCreatedStatus):
		return BidCreatedStatus, nil

	case string(BidPublishedStatus):
		return BidPublishedStatus, nil

	case string(BidCanceledStatus):
		return BidCanceledStatus, nil

	default:
		return "", errors.Wrapf(ErrValidation, "invalid BidStatus: %s", str)
	}
}

// BidAuthorType Тип автора предложения
type BidAuthorType string

const (
	BidAuthorOrganizationType BidAuthorType = "Organization"
	BidAuthorUserType         BidAuthorType = "User"
)

func NewBidAuthorType(str string) (BidAuthorType, error) {
	switch str {
	case string(BidAuthorOrganizationType):
		return BidAuthorOrganizationType, nil

	case string(BidAuthorUserType):
		return BidAuthorUserType, nil

	default:
		return "", errors.Wrapf(ErrValidation, "invalid BidAuthorType: %s", str)
	}
}

// BidVersion Номер версии после правок
type BidVersion int

func NewBidVersion(n int) BidVersion {
	if n <= 0 {
		return 1
	}
	return BidVersion(n)
}

// BidSnapshot Снимок состояния Bid
type BidSnapshot struct {
	ID          ID
	Name        BidName
	Description BidDescription
	Version     BidVersion
}

// Bid Предложение
type Bid struct {
	ID          ID             `json:"id"`
	Name        BidName        `json:"name"`
	Description BidDescription `json:"description"`
	Status      BidStatus      `json:"status"`
	TenderID    ID             `json:"-"`
	AuthorType  BidAuthorType  `json:"authorType"`
	AuthorID    ID             `json:"authorId"`
	Version     BidVersion     `json:"version"`
	CreatedAt   time.Time      `json:"createdAt"`
	Snapshots   []BidSnapshot  `json:"-"`
}

func (b *Bid) ChangeStatus(status string) error {
	s, err := NewBidStatus(status)
	if err != nil {
		return err
	}
	b.Status = s
	return nil
}

func (b *Bid) takeSnapshot() {
	snapshot := BidSnapshot{
		ID:          NewID(),
		Name:        b.Name,
		Description: b.Description,
		Version:     b.Version,
	}
	b.Snapshots = append(b.Snapshots, snapshot)
	b.Version++
}

func (b *Bid) Edit(name, description *string) error {
	b.takeSnapshot()

	if name != nil {
		bidName, err := NewBidName(*name)
		if err != nil {
			return err
		}
		b.Name = bidName
	}

	if description != nil {
		bidDescription, err := NewBidDescription(*description)
		if err != nil {
			return err
		}
		b.Description = bidDescription
	}

	return nil
}

func (b *Bid) Rollback(version int) error {
	v := NewBidVersion(version)

	i := slices.IndexFunc(b.Snapshots, func(s BidSnapshot) bool {
		return s.Version == v
	})

	if i == -1 {
		return errors.Wrapf(ErrNotFound, "Bid snapshot %d not found", v)
	}

	snapshot := b.Snapshots[i]

	b.Snapshots = append(b.Snapshots, BidSnapshot{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		Version:     b.Version,
	})
	b.Name = snapshot.Name
	b.Description = snapshot.Description
	b.Version++

	return nil
}

func NewBid(name, description, authorType string, tenderID, authorID ID) (*Bid, error) {
	id := NewID()

	bidName, err := NewBidName(name)
	if err != nil {
		return nil, err
	}

	bidDescription, err := NewBidDescription(description)
	if err != nil {
		return nil, err
	}

	status := BidCreatedStatus

	bidAuthorType, err := NewBidAuthorType(authorType)
	if err != nil {
		return nil, err
	}

	version := NewBidVersion(1)
	createdAt := time.Now()

	return &Bid{
		ID:          id,
		Name:        bidName,
		Description: bidDescription,
		Status:      status,
		TenderID:    tenderID,
		AuthorType:  bidAuthorType,
		AuthorID:    authorID,
		Version:     version,
		CreatedAt:   createdAt,
		Snapshots:   make([]BidSnapshot, 0),
	}, nil
}
