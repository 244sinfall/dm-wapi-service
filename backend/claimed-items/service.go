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

func add(i claimedItem) error {
	if strings.TrimSpace(i.Owner) == "" || strings.TrimSpace(i.Link) == "" || strings.TrimSpace(i.Name) == "" || strings.TrimSpace(i.OwnerProfile) == "" {
		return errors.New("required fields are empty")
	}
	f := globals.GetFirestore()
	newDoc := f.Collection("claimedItems").NewDoc()
	i.Id = newDoc.ID
	i.AddedAt = time.Now()
	_, err := newDoc.Create(globals.GetGlobalContext(), i)
	itemsI, _ := c.Get("items")
	items := itemsI.(map[string][]claimedItem)
	items[i.GetKey()] = append(items[i.GetKey()], i)
	c.Set("items", items, cache.NoExpiration)
	if err != nil {
		return err
	}
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
	itemsI, _ := c.Get("items")
	items := itemsI.(map[string][]claimedItem)
	var indexToDelete int
	for index, item := range items[itemToDelete.GetKey()] {
		if item.Id == itemToDelete.Id {
			indexToDelete = index
		}
	}
	items[itemToDelete.GetKey()] = append(items[itemToDelete.GetKey()][:indexToDelete], items[itemToDelete.GetKey()][indexToDelete+1:]...)
	c.Set("items", items, cache.NoExpiration)
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
	itemsI, _ := c.Get("items")
	items := itemsI.(map[string][]claimedItem)
	for _, itemIter := range items[item.GetKey()] {
		if itemIter.Id == item.Id {
			itemIter = item
		}
	}
	c.Set("items", items, cache.NoExpiration)
	return err
}

func update(id string, toUpdate claimedItem) (*claimedItem, error) {
	f := globals.GetFirestore()
	docRef := f.Doc("claimedItems/" + id)
	rx, _ := regexp.Compile("[0-9]+")
	toUpdate.OwnerProofName = "№ " + rx.FindString(toUpdate.OwnerProofLink)
	_, err := docRef.Set(globals.GetGlobalContext(), toUpdate)
	itemsI, _ := c.Get("items")
	items := itemsI.(map[string][]claimedItem)
	for _, itemIter := range items[toUpdate.GetKey()] {
		if itemIter.Id == toUpdate.Id {
			itemIter = toUpdate
		}
	}
	c.Set("items", items, cache.NoExpiration)
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
		c.Set("items", newItemsMapped, cache.NoExpiration)
		return newItemsMapped
	}
	return items.(map[string][]claimedItem)
}
