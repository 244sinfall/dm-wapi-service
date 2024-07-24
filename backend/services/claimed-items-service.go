package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	legendaryTableItem = "Легендарный"
	epicTableItem      = "Эпический"
	rareTableItem      = "Редкий"
	greenTableItem     = "Необычный"
	otherTableItem     = "Прочие"
)

type ClaimedItemsList struct {
	Legendary []ClaimedItem `json:"legendary"`
	Epic      []ClaimedItem `json:"epic"`
	Rare      []ClaimedItem `json:"rare"`
	Green     []ClaimedItem `json:"green"`
	Other     []ClaimedItem `json:"other"`
}

func (c *ClaimedItemsList) deleteFromNeededSlice(item ClaimedItem) {
	switch item.Quality {
	case "Легендарный":
		for i, v := range claimedItems.Legendary {
			if v.Id == item.Id {
				claimedItems.Legendary = append(claimedItems.Legendary[:i], claimedItems.Legendary[i+1:]...)
			}
		}
	case "Эпический":
		for i, v := range claimedItems.Epic {
			if v.Id == item.Id {
				claimedItems.Epic = append(claimedItems.Epic[:i], claimedItems.Epic[i+1:]...)
			}
		}
	case "Редкий":
		for i, v := range claimedItems.Rare {
			if v.Id == item.Id {
				claimedItems.Rare = append(claimedItems.Rare[:i], claimedItems.Rare[i+1:]...)
			}
		}
	case "Необычный":
		for i, v := range claimedItems.Green {
			if v.Id == item.Id {
				claimedItems.Green = append(claimedItems.Green[:i], claimedItems.Green[i+1:]...)
			}
		}
	default:
		for i, v := range claimedItems.Other {
			if v.Id == item.Id {
				claimedItems.Other = append(claimedItems.Other[:i], claimedItems.Other[i+1:]...)
			}
		}
	}
}

func (c *ClaimedItemsList) addToNeededSlice(item ClaimedItem) {
	switch item.Quality {
	case "Легендарный":
		claimedItems.Legendary = append(claimedItems.Legendary, item)
	case "Эпический":
		claimedItems.Epic = append(claimedItems.Epic, item)
	case "Редкий":
		claimedItems.Rare = append(claimedItems.Rare, item)
	case "Необычный":
		claimedItems.Green = append(claimedItems.Green, item)
	default:
		claimedItems.Other = append(claimedItems.Other, item)
	}
}

func (c *ClaimedItemsList) findNeededSlice(quality string) []ClaimedItem {
	switch quality {
	case "Легендарный":
		return claimedItems.Legendary
	case "Эпический":
		return claimedItems.Epic
	case "Редкий":
		return claimedItems.Rare
	case "Необычный":
		return claimedItems.Green
	default:
		return claimedItems.Other
	}
}

func (c *ClaimedItemsList) Add(i ClaimedItem, f *firestore.Client, ctx context.Context) error {
	if strings.TrimSpace(i.Owner) == "" || strings.TrimSpace(i.Link) == "" || strings.TrimSpace(i.Name) == "" || strings.TrimSpace(i.OwnerProfile) == "" {
		return errors.New("required fields are empty")
	}
	newDoc := f.Collection("claimedItems").NewDoc()
	i.Id = newDoc.ID
	i.AddedAt = time.Now()
	_, err := newDoc.Create(ctx, i)
	if err != nil {
		return err
	}
	claimedItems.addToNeededSlice(i)
	return nil
}

func (c *ClaimedItemsList) Delete(id string, f *firestore.Client, ctx context.Context) (*ClaimedItem, error) {
	docRef := f.Doc("claimedItems/" + id)
	data, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	itemToDelete := new(ClaimedItem)
	err = data.DataTo(itemToDelete)
	itemToDelete.Id = id
	if err != nil {
		return nil, err
	}
	claimedItems.deleteFromNeededSlice(*itemToDelete)
	_, err = docRef.Delete(ctx)
	if err != nil {
		return nil, err
	}
	return itemToDelete, nil
}

func (c *ClaimedItemsList) Approve(id string, approveUser string, f *firestore.Client, ctx context.Context) error {
	docRef := f.Doc("claimedItems/" + id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	data := doc.Data()
	if data["Accepted"].(bool) {
		return errors.New("this item is already accepted")
	}
	quality := doc.Data()["Quality"].(string)
	selectedSlice := claimedItems.findNeededSlice(quality)
	for i, v := range selectedSlice {
		if v.Id == id {
			v.Acceptor = approveUser
			v.Accepted = true
			v.AcceptedAt = time.Now()
			selectedSlice[i] = v
			_, err := docRef.Set(ctx, v)
			return err
		}
	}
	return errors.New("id not found in the database")
}

func (c *ClaimedItemsList) Update(id string, admin bool, toUpdate ClaimedItem, editorName string, f *firestore.Client, ctx context.Context) (newItem *ClaimedItem, oldItem *ClaimedItem, err error) {
	docRef := f.Doc("claimedItems/" + id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	quality := doc.Data()["Quality"].(string)
	selectedSlice := claimedItems.findNeededSlice(quality)
	for i, v := range selectedSlice {
		if v.Id == id {
			oldItem := v
			var reviewerChange bool
			if v.Name != toUpdate.Name || v.Link != toUpdate.Link || v.Reviewer != toUpdate.Reviewer {
				if admin {
					if v.Reviewer != toUpdate.Reviewer {
						reviewerChange = true
					}
					v.Name = toUpdate.Name
					v.Link = toUpdate.Link
					v.Reviewer = toUpdate.Reviewer
				}
			}
			v.Owner = toUpdate.Owner
			v.OwnerProfile = toUpdate.OwnerProfile
			v.OwnerProofLink = toUpdate.OwnerProofLink
			v.AdditionalInfo = toUpdate.AdditionalInfo
			if !reviewerChange {
				if !strings.Contains(v.Reviewer, editorName) {
					y, m, d := time.Now().Date()
					v.Reviewer += fmt.Sprintf("\nИзменил: %v (%v.%v.%v)", editorName, d, int(m), y)
				}
			}
			rx, _ := regexp.Compile("[0-9]+")
			v.OwnerProofName = "№ " + rx.FindString(toUpdate.OwnerProofLink)
			selectedSlice[i] = v
			_, err := docRef.Set(ctx, v)
			return &v, &oldItem, err
		}
	}
	return nil, nil, errors.New("id not found in the database")
}

type ClaimedItem struct {
	Id             string    `json:"id"`
	Quality        string    `json:"quality"`
	Name           string    `json:"name"`
	Link           string    `json:"link"`
	Owner          string    `json:"owner"`
	OwnerProfile   string    `json:"ownerProfile"`
	OwnerProofName string    `json:"ownerProof"`
	OwnerProofLink string    `json:"ownerProofLink"`
	Reviewer       string    `json:"reviewer"`
	Accepted       bool      `json:"accepted"`
	Acceptor       string    `json:"acceptor"`
	AddedAt        time.Time `json:"addedAt"`
	AcceptedAt     time.Time `json:"acceptedAt"`
	AdditionalInfo string    `json:"additionalInfo"`
}

var claimedItems ClaimedItemsList

func GetClaimedItems() ClaimedItemsList {
	return claimedItems
}

func GetClaimedItemsFromDatabase(f *firestore.Client, ctx context.Context) {
	fmt.Println("Receiving claimed items from database...")
	claimedItems = ClaimedItemsList{
		Legendary: make([]ClaimedItem, 0, 100),
		Epic:      make([]ClaimedItem, 0, 100),
		Rare:      make([]ClaimedItem, 0, 100),
		Green:     make([]ClaimedItem, 0, 100),
		Other:     make([]ClaimedItem, 0, 100),
	}
	documents, err := f.Collection("claimedItems").Documents(ctx).GetAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, doc := range documents {
		item := new(ClaimedItem)
		err := doc.DataTo(item)
		item.Id = doc.Ref.ID
		if err != nil {
			fmt.Println(err)
		} else {
			switch item.Quality {
			case legendaryTableItem:
				claimedItems.Legendary = append(claimedItems.Legendary, *item)
			case epicTableItem:
				claimedItems.Epic = append(claimedItems.Epic, *item)
			case rareTableItem:
				claimedItems.Rare = append(claimedItems.Rare, *item)
			case greenTableItem:
				claimedItems.Green = append(claimedItems.Green, *item)
			default:
				claimedItems.Other = append(claimedItems.Other, *item)
			}
		}
	}
	fmt.Println("Received claimed items from database.")
}
