package claimeditems

import (
	"darkmoon-wapi-service/globals"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"google.golang.org/api/iterator"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func invalidateCache() {
	c.Delete("items")
}

func add(i claimedItem) error {
	if strings.TrimSpace(i.Owner) == "" || strings.TrimSpace(i.Link) == "" || strings.TrimSpace(i.Name) == "" || strings.TrimSpace(i.OwnerProfile) == "" {
		return errors.New("required fields are empty")
	}
	f := globals.GetFirestore()
	newDoc := f.Collection("claimedItems").NewDoc()
	i.Id = newDoc.ID
	i.AddedAt = time.Now()
	_, err := newDoc.Create(globals.GetGlobalContext(), i)
	if err != nil {
		return err
	}
	invalidateCache()
	return nil
}

func delete(id string) (*claimedItem, error) {
	f := globals.GetFirestore()
	docRef := f.Doc("claimedItems/" + id)
	data, err := docRef.Get(globals.GetGlobalContext())
	if err != nil {
		return nil, err
	}
	itemToDelete := new(claimedItem)
	err = data.DataTo(itemToDelete)
	if err != nil {
		return nil, err
	}
	_, err = docRef.Delete(globals.GetGlobalContext())
	if err != nil {
		return nil, err
	}
	invalidateCache()
	return itemToDelete, nil
}

func approve(id string, approveUser string) error {
	f := globals.GetFirestore()
	docRef := f.Doc("claimedItems/" + id)
	doc, err := docRef.Get(globals.GetGlobalContext())
	if err != nil {
		return err
	}
	var item claimedItem
	err = doc.DataTo(&item)
	if err != nil {
		return err
	}
	data := doc.Data()
	if data["Accepted"].(bool) {
		return errors.New("this item is already accepted")
	}
	item.Acceptor = approveUser
	item.Accepted = true
	item.AcceptedAt = time.Now()
	_, err = docRef.Set(globals.GetGlobalContext(), item)
	invalidateCache()
	return err
}

func update(id string, toUpdate claimedItem) (*claimedItem, error) {
	f := globals.GetFirestore()
	docRef := f.Doc("claimedItems/" + id)
	// doc, err := docRef.Get(globals.GetGlobalContext())
	// var item claimedItem
	// err = doc.DataTo(&item)
	// if err != nil {
	// 	return nil, err
	// }

	// var reviewerChange bool
	// if v.Name != toUpdate.Name || v.Link != toUpdate.Link || v.Reviewer != toUpdate.Reviewer {
	// 			if admin {
	// 				if v.Reviewer != toUpdate.Reviewer {
	// 					reviewerChange = true
	// 				}
	// 				v.Name = toUpdate.Name
	// 				v.Link = toUpdate.Link
	// 				v.Reviewer = toUpdate.Reviewer
	// 			}
	// 		}
	// 		v.Owner = toUpdate.Owner
	// 		v.OwnerProfile = toUpdate.OwnerProfile
	// 		v.OwnerProofLink = toUpdate.OwnerProofLink
	// 		v.AdditionalInfo = toUpdate.AdditionalInfo
	// if !reviewerChange {
	// 	if !strings.Contains(v.Reviewer, editorName) {
	// 		y, m, d := time.Now().Date()
	// 		v.Reviewer += fmt.Sprintf("\nИзменил: %v (%v.%v.%v)", editorName, d, int(m), y)
	// 	}
	// }
	rx, _ := regexp.Compile("[0-9]+")
	toUpdate.OwnerProofName = "№ " + rx.FindString(toUpdate.OwnerProofLink)
	_, err := docRef.Set(globals.GetGlobalContext(), toUpdate)
	return &toUpdate, err

}

func fetchClaimedItems() []claimedItem {
	claimedItems := make([]claimedItem, 0, 100)
	f := globals.GetFirestore()
	collection := f.Collection("claimedItems")
	iter := collection.Documents(globals.GetGlobalContext())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Data inconsistency for claimed item: %v\n", doc.Ref.ID)
			continue
		}
		item := new(claimedItem)
		err = doc.DataTo(item)
		if err != nil {
			fmt.Printf("Data inconsistency for claimed item: %v\n", doc.Ref.ID)
			continue
		}
		claimedItems = append(claimedItems, *item)
	}
	return claimedItems
}

func getClaimedItems() map[string][]claimedItem {
	items, found := c.Get("items")
	if !found {
		newItems := fetchClaimedItems()
		c.Set("items", newItems, cache.DefaultExpiration)
		newItemsMapped := make(map[string][]claimedItem, 100)
		for _, item := range newItems {
			if list, ok := newItemsMapped[item.GetKey()]; ok {
				list = append(list, item)
				newItemsMapped[item.GetKey()] = list
			} else {
				newArr := make([]claimedItem, 0, 100)
				newArr = append(newArr, item)
				newItemsMapped[item.GetKey()] = newArr
			}
		}
		return newItemsMapped
	}
	return items.(map[string][]claimedItem)
}
