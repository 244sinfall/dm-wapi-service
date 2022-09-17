package economics

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const cacheFrequency = 30 * time.Minute
const fieldsInOneDocumentDb = 500

type cachedChecks struct {
	checks    []APIResponseCheck
	updatedAt time.Time
	updating  bool
}

var CachedChecks cachedChecks

func UploadChecksToDatabase(f *firestore.Client, ctx context.Context) error {
	checkMap := make(map[string]APIResponseCheck, fieldsInOneDocumentDb)
	counter := 1
	higherBound := CachedChecks.checks[0].Id //12620
	documents, err := f.Collection("checks").Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, doc := range documents {
		if !strings.Contains(doc.Ref.Path, "info") {
			_, err := f.Collection("checks").Doc(doc.Ref.ID).Delete(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("Old cache deleted.")
	batch := f.Batch()
	for idx, check := range CachedChecks.checks {
		checkMap[strconv.Itoa(check.Id)] = check
		counter++
		if counter == fieldsInOneDocumentDb || idx == len(CachedChecks.checks)-1 { //12120
			lowerBound := CachedChecks.checks[idx].Id
			docName := fmt.Sprintf("%d-%d", lowerBound, higherBound)
			batch.Set(f.Collection("checks").Doc(docName), checkMap)
			higherBound = lowerBound - 1
			checkMap = make(map[string]APIResponseCheck)
			counter = 1
		}
	}
	batch.Set(f.Doc("checks/info"), map[string]interface{}{
		"updatedAt": time.Now(),
		"count":     CachedChecks.checks[0].Id,
	})
	_, err = batch.Commit(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Pushing cache completed.")
	return nil
}

func getChecksFromDatabase(f *firestore.Client, ctx context.Context) error {
	fmt.Println("Receiving checks from database...")
	checkInfoSnapshot, err := f.Doc("checks/info").Get(ctx)
	if err != nil {
		return err
	}
	checkInfo := checkInfoSnapshot.Data()
	updatedAt := checkInfo["updatedAt"].(time.Time)
	checkCount := checkInfo["count"].(int64)
	CachedChecks = cachedChecks{
		checks:    make([]APIResponseCheck, checkCount),
		updatedAt: updatedAt,
		updating:  true,
	}
	documents, err := f.Collection("checks").Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, doc := range documents {
		if !strings.Contains(doc.Ref.Path, "info") {
			rawChecks := make(map[string]APIResponseCheck, fieldsInOneDocumentDb)
			err := doc.DataTo(&rawChecks)
			if err != nil {
				return err
			}

			for k, v := range rawChecks {
				checkId, err := strconv.Atoi(k)
				if err != nil {
					fmt.Println(err)
					continue
				}
				CachedChecks.checks[checkId-1] = v

			}
		}
	}
	CachedChecks.updating = false
	fmt.Println("Received checks from database successfully")
	return nil
}

func ParseAndDeployNewChecks(f *firestore.Client, ctx context.Context) error {
	CachedChecks = cachedChecks{
		checks:    nil,
		updatedAt: time.Now(),
		updating:  true,
	}
	parsedChecks, err := ParseChecksFromDarkmoon()
	if err != nil {
		fmt.Println("Unable to parse checks from Darkmoon. Error: " + err.Error())
		return err
	} else {
		CachedChecks = cachedChecks{
			checks:    parsedChecks,
			updatedAt: time.Now(),
			updating:  false,
		}
		err := UploadChecksToDatabase(f, ctx)
		if err != nil {
			fmt.Println("Unable to upload checks to database with error: ", err.Error())
			return err
		}
	}
	return nil
}

func ChecksScheduler(f *firestore.Client, ctx context.Context, ping bool) {
	fmt.Println("Scheduler just started. Looking for retrieve cached checks from database")
	err := getChecksFromDatabase(f, ctx)
	if err != nil {
		fmt.Println("Unable to fetch checks from DB with error: " + err.Error())
	}
	for {
		if time.Now().Sub(CachedChecks.updatedAt) > cacheFrequency || err != nil {
			err := ParseAndDeployNewChecks(f, ctx)
			if err != nil {
				fmt.Println("Unable to parse new checks! " + err.Error())
			}
			err = nil
		} else {
			if ping {
				return
			}
			schedule := cacheFrequency - time.Now().Sub(CachedChecks.updatedAt) + (2 * time.Minute)
			time.Sleep(schedule)
		}
	}
}