package services

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestStorage(t *testing.T) {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"storageBucketList": StorageBucketListQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"createStorageBucket": CreateStorageBucketMutation,
			},
		}),
	})

	addPermission(testUserUUID, testAPIKey, "Storage::*")

	// createStorageBucket should make a new bucket and return it
	testBucketName := "my_bucket"
	reqStr := fmt.Sprintf(`
	mutation {
		createStorageBucket(input: {
			name: "%s",
			isPublic: false
		}) {
			name
			isPublic
		}
	}
	`, testBucketName)

	res := graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("createStorageBucket returns error: %v\n", res.Errors)
		t.Fail()
		return
	}
	if data := res.Data.(map[string]interface{})["createStorageBucket"].(map[string]interface{}); false ||
		data["name"].(string) != testBucketName ||
		data["isPublic"].(bool) {
		fmt.Printf("createStorageBucket returns invalid data: %v\n", data)
		t.Fail()
		return
	}

	// storageBucketList should return list of all buckets
	reqStr = `
	query {
		storageBucketList {
			name
			isPublic
		}
	}
	`

	res = graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("storageBucketList returns error: %v\n", res.Errors)
		t.Fail()
		return
	}
	if data := res.Data.(map[string]interface{})["storageBucketList"].([]interface{}); false ||
		len(data) != 1 ||
		data[0].(map[string]interface{})["name"].(string) != testBucketName {
		fmt.Printf("storageBucketList returns invalid data: %v\n", data)
		t.Fail()
		return
	}

	// checkStorageBucketPermission should return true if the api key is valid
	if !checkStorageBucketPermission(testBucketName, testAPIKey) ||
		checkStorageBucketPermission(testBucketName+"_invalid", testAPIKey) ||
		checkStorageBucketPermission(testBucketName, testAPIKey+"_invalid") {
		fmt.Printf("checkStorageBucketPermission returns invalid result")
		t.Fail()
		return
	}
}
