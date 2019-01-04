package services

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestKeyValue(t *testing.T) {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"keyValueItem": KeyValueItemQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"setKeyValueItem": SetKeyValueItemMutation,
			},
		}),
	})

	addPermission(testUserUUID, testAPIKey, "KeyValue::*")

	// setKeyValueItem should make a new item and return it
	testKey := "testkey"
	testValue := "testValue1"
	reqStr := fmt.Sprintf(`
	mutation {
		setKeyValueItem(input: {
			key: "%s"
			value: "%s"
		}) {
			key
			value
		}
	}
	`, testKey, testValue)

	res := graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("setKeyValueItem returns error: %v\n", res.Errors)
		t.Fail()
		return
	}

	if data := res.Data.(map[string]interface{})["setKeyValueItem"].(map[string]interface{}); false ||
		data["key"].(string) != testKey ||
		data["value"].(string) != testValue {
		fmt.Printf("setKeyValueItem returns invalid data: %v\n", data)
	}

	// setKeyValueItem should update value for existing key
	testNewValue := "testValue2"
	reqStr = fmt.Sprintf(`
	mutation {
		setKeyValueItem(input: {
			key: "%s"
			value: "%s"
		}) {
			key
			value
		}
	}
	`, testKey, testNewValue)

	res = graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("setKeyValueItem returns error: %v\n", res.Errors)
		t.Fail()
		return
	}

	if data := res.Data.(map[string]interface{})["setKeyValueItem"].(map[string]interface{}); false ||
		data["key"].(string) != testKey ||
		data["value"].(string) != testNewValue {
		fmt.Printf("setKeyValueItem returns invalid data: %v\n", data)
	}

	// keyValueItem should return item
	reqStr = fmt.Sprintf(`
	query {
		keyValueItem(key: "%s") {
			key
			value
		}
	}
	`, testKey)

	res = graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("keyValueItem returns error: %v\n", res.Errors)
		t.Fail()
		return
	}

	if data := res.Data.(map[string]interface{})["keyValueItem"].(map[string]interface{}); false ||
		data["key"].(string) != testKey ||
		data["value"].(string) != testNewValue {
		fmt.Printf("keyValueItem returns invalid data: %v\n", data)
	}

	// keyValueItem should return nil for key which not exists
	reqStr = fmt.Sprintf(`
	query {
		keyValueItem(key: "%s") {
			key
			value
		}
	}
	`, testKey+"_invalid")

	res = graphql.Do(graphql.Params{Schema: schema, RequestString: reqStr, Context: testCtx})
	if res.HasErrors() {
		fmt.Printf("keyValueItem returns error: %v\n", res.Errors)
		t.Fail()
		return
	}

	if data := res.Data.(map[string]interface{})["keyValueItem"]; data != nil {
		fmt.Printf("keyValueItem should return nil: %v\n", data)
	}
}
