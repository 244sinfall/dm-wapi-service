package claimed_items

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"time"
)

const (
	LegendaryTableItem = "Легендарный"
	EpicTableItem      = "Эпический"
	RareTableItem      = "Редкий"
	GreenTableItem     = "Необычный"
	OtherTableItem     = "Прочие"
)

const adminPermission = 4
const reviewerPermission = 3

type ClaimedItemsResponse struct {
	legendary []ClaimedItem
	epic      []ClaimedItem
	rare      []ClaimedItem
	green     []ClaimedItem
	other     []ClaimedItem
}

func (c *ClaimedItemsResponse) deleteFromNeededSlice(item ClaimedItem) {
	switch item.Quality {
	case "Легендарный":
		for i, v := range ClaimedItems.legendary {
			if v.Id == item.Id {
				ClaimedItems.legendary = append(ClaimedItems.legendary[:i], ClaimedItems.legendary[i+1:]...)
			}
		}
	case "Эпический":
		for i, v := range ClaimedItems.epic {
			if v.Id == item.Id {
				ClaimedItems.epic = append(ClaimedItems.epic[:i], ClaimedItems.epic[i+1:]...)
			}
		}
	case "Редкий":
		for i, v := range ClaimedItems.rare {
			if v.Id == item.Id {
				ClaimedItems.rare = append(ClaimedItems.rare[:i], ClaimedItems.rare[i+1:]...)
			}
		}
	case "Необычный":
		for i, v := range ClaimedItems.green {
			if v.Id == item.Id {
				ClaimedItems.green = append(ClaimedItems.green[:i], ClaimedItems.green[i+1:]...)
			}
		}
	case "Прочие":
		for i, v := range ClaimedItems.other {
			if v.Id == item.Id {
				ClaimedItems.other = append(ClaimedItems.other[:i], ClaimedItems.other[i+1:]...)
			}
		}
	default:
		for i, v := range ClaimedItems.other {
			if v.Id == item.Id {
				ClaimedItems.other = append(ClaimedItems.other[:i], ClaimedItems.other[i+1:]...)
			}
		}
	}
}

func (c *ClaimedItemsResponse) addToNeededSlice(item ClaimedItem) {
	switch item.Quality {
	case "Легендарный":
		ClaimedItems.legendary = append(ClaimedItems.legendary, item)
	case "Эпический":
		ClaimedItems.epic = append(ClaimedItems.epic, item)
	case "Редкий":
		ClaimedItems.rare = append(ClaimedItems.rare, item)
	case "Необычный":
		ClaimedItems.green = append(ClaimedItems.green, item)
	case "Прочие":
		ClaimedItems.other = append(ClaimedItems.other, item)
	default:
		ClaimedItems.legendary = append(ClaimedItems.other, item)
	}
}

func (c *ClaimedItemsResponse) findNeededSlice(quality string) []ClaimedItem {
	switch quality {
	case "Легендарный":
		return ClaimedItems.legendary
	case "Эпический":
		return ClaimedItems.epic
	case "Редкий":
		return ClaimedItems.rare
	case "Необычный":
		return ClaimedItems.green
	case "Прочие":
		return ClaimedItems.other
	default:
		return ClaimedItems.other
	}
}

func (c *ClaimedItemsResponse) add(i ClaimedItem, f *firestore.Client, ctx context.Context) error {
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
	ClaimedItems.addToNeededSlice(i)
	return nil
}

func (c *ClaimedItemsResponse) delete(id string, f *firestore.Client, ctx context.Context) error {
	docRef := f.Doc("claimedItems/" + id)
	data, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	itemToDelete := new(ClaimedItem)
	err = data.DataTo(itemToDelete)
	itemToDelete.Id = id
	if err != nil {
		return err
	}
	ClaimedItems.deleteFromNeededSlice(*itemToDelete)
	_, err = docRef.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClaimedItemsResponse) approve(id string, approveUser string, f *firestore.Client, ctx context.Context) error {
	docRef := f.Doc("claimedItems/" + id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	data := doc.Data()
	if data["Accepted"].(bool) == true {
		return errors.New("this item is already accepted")
	}
	quality := doc.Data()["Quality"].(string)
	selectedSlice := ClaimedItems.findNeededSlice(quality)
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

func (c *ClaimedItemsResponse) update(id string, permission int64, toUpdate ClaimedItem, editorName string, f *firestore.Client, ctx context.Context) error {
	docRef := f.Doc("claimedItems/" + id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	quality := doc.Data()["Quality"].(string)
	selectedSlice := ClaimedItems.findNeededSlice(quality)
	for i, v := range selectedSlice {
		if v.Id == id {
			var reviewerChange bool
			if v.Name != toUpdate.Name || v.Link != toUpdate.Link || v.Reviewer != toUpdate.Reviewer {
				if permission >= adminPermission {
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
			return err
		}
	}
	return errors.New("id not found in the database")
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

var ClaimedItems ClaimedItemsResponse

func GetClaimedItemsFromDatabase(f *firestore.Client, ctx context.Context) {
	fmt.Println("Receiving claimed items from database...")
	ClaimedItems = ClaimedItemsResponse{
		legendary: make([]ClaimedItem, 0, 100),
		epic:      make([]ClaimedItem, 0, 100),
		rare:      make([]ClaimedItem, 0, 100),
		green:     make([]ClaimedItem, 0, 100),
		other:     make([]ClaimedItem, 0, 100),
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
			case LegendaryTableItem:
				ClaimedItems.legendary = append(ClaimedItems.legendary, *item)
			case EpicTableItem:
				ClaimedItems.epic = append(ClaimedItems.epic, *item)
			case RareTableItem:
				ClaimedItems.rare = append(ClaimedItems.rare, *item)
			case GreenTableItem:
				ClaimedItems.green = append(ClaimedItems.green, *item)
			case OtherTableItem:
				ClaimedItems.other = append(ClaimedItems.other, *item)
			default:
				ClaimedItems.other = append(ClaimedItems.other, *item)
			}
		}
	}
	fmt.Println("Received claimed items from database.")
}

func AddClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < reviewerPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}
	claimedItem := new(ClaimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItem)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = ClaimedItems.add(*claimedItem, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func ApproveClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < adminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	user, err := a.GetUser(ctx, token.UID)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}
	err = ClaimedItems.approve(id, user.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func UpdateClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < reviewerPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	user, _ := a.GetUser(ctx, token.UID)
	claimedItemMock := new(ClaimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItemMock)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	err = ClaimedItems.update(id, permission, *claimedItemMock, user.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func DeleteClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < adminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	err = ClaimedItems.delete(id, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
	return
}

func ReceiveClaimedItems(c *gin.Context) {
	c.JSON(200, gin.H{"result": gin.H{"legendary": ClaimedItems.legendary, "epic": ClaimedItems.epic, "rare": ClaimedItems.rare, "green": ClaimedItems.green, "other": ClaimedItems.other}})
}
