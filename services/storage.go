package services

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/allegro/bigcache"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"

	"luppiter/components/database"
)

var gcpBucket *storage.BucketHandle
var cache *bigcache.BigCache

type StorageBucket struct {
	Name     string `gorm:"type:varchar(255);primary_key"`
	UserUUID string `gorm:"type:varchar(40);index"`
	IsPublic bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type StorageItem struct {
	Name       string `gorm:"varchar(255);primary_key"`
	BucketName string `gorm:"varchar(255);primary_key"`
	UUID       string `gorm:"varchar(40);unique_index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var storageBucketType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StorageBucket",
	Fields: graphql.Fields{
		"name":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"isPublic":  &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var storageItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StorageItem",
	Fields: graphql.Fields{
		"name":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"bucketName": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt":  &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var (
	ErrStorageNameDuplicated = errors.New("requested name already exists")
	ErrStorageBucketNotFound = errors.New("bucket not found")
	ErrStorageItemNotFound   = errors.New("file not found")
)

// StorageBucketListQuery returns [StorageBucket!]!
var StorageBucketListQuery = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(storageBucketType))),
	Description: "[Storage] Get the list of buckets.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		key, err := checkPermission(p.Context.Value("APIKey").(string), "Storage")
		if err != nil {
			return nil, err
		}

		var buckets []StorageBucket
		database.DB.Where(&StorageBucket{UserUUID: key.UserUUID}).Find(&buckets)
		return buckets, nil
	},
}

// CreateStorageBucketMutation returns StorageBucket!
var CreateStorageBucketMutation = &graphql.Field{
	Type:        graphql.NewNonNull(storageBucketType),
	Description: "[Storage] Create a new bucket.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "CreateStorageBucketInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"name":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"isPublic": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Boolean)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		key, err := checkPermission(p.Context.Value("APIKey").(string), "Storage")
		if err != nil {
			return nil, err
		}

		input := p.Args["input"].(map[string]interface{})

		var bucket StorageBucket
		database.DB.Where(&StorageBucket{Name: input["name"].(string)}).First(&bucket)
		if bucket.UserUUID != "" {
			return nil, ErrStorageNameDuplicated
		}

		bucket = StorageBucket{
			Name:     input["name"].(string),
			UserUUID: key.UserUUID,
			IsPublic: input["isPublic"].(bool),
		}
		errs := database.DB.Create(&bucket).GetErrors()
		if len(errs) > 0 {
			return nil, ErrInternal
		}

		return bucket, nil
	},
}

// UploadStorageItem create new file in the bucket.
func UploadStorageItem(ctx context.Context, bucketName, itemName string, reader io.Reader) error {
	if !checkStorageBucketPermission(bucketName, ctx.Value("APIKey").(string)) {
		return ErrStorageBucketNotFound
	}

	u := uuid.NewV4().String()
	item := StorageItem{
		BucketName: bucketName,
		Name:       itemName,
		UUID:       u,
	}
	errs := database.DB.Create(&item).GetErrors()
	if len(errs) > 0 {
		return ErrInternal
	}

	// Upload to google cloud storage.
	w := gcpBucket.Object(u).NewWriter(context.Background())
	_, err := io.Copy(w, reader)
	if err != nil {
		return err
	}

	return w.Close()
}

// DownloadStorageItem returns a file.
func DownloadStorageItem(ctx context.Context, bucketName, itemName string) (io.Reader, string, error) {
	if !checkStorageBucketPermission(bucketName, ctx.Value("APIKey").(string)) {
		return nil, "", ErrStorageBucketNotFound
	}

	var item StorageItem
	database.DB.Where(&StorageItem{BucketName: bucketName, Name: itemName}).First(&item)
	if item.UUID == "" {
		return nil, "", ErrStorageItemNotFound
	}

	obj := gcpBucket.Object(item.UUID)
	c := context.Background()
	file, err := obj.NewReader(c)
	if err != nil {
		return nil, "", ErrInternal
	}
	attrs, err := obj.Attrs(c)
	if err != nil {
		return nil, "", ErrInternal
	}

	return file, attrs.ContentType, nil
}

func checkStorageBucketPermission(bucketName string, key string) bool {
	apiKey, ok := getAPIKey(key)

	var bucket StorageBucket
	database.DB.Where(&StorageBucket{Name: bucketName}).First(&bucket)
	return bucket.UserUUID != "" && (bucket.IsPublic || (ok && apiKey.UserUUID == bucket.UserUUID))
}

func init() {
	// Initialize Google cloud stroage
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "secrets/service-account.json")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic("failed to create gcp client")
	}

	gcpBucket = client.Bucket("luppiter.lynlab.co.kr")
}
