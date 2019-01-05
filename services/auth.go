package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/lib/pq"

	"luppiter/components/database"
	"luppiter/components/random"
)

type APIKey struct {
	UserUUID    string         `gorm:"varchar(40);index"`
	Key         string         `gorm:"varchar(255);primary_key"`
	Comment     string         `gorm:"varchar(255)"`
	Permissions pq.StringArray `gorm:"type:varchar(255)[]"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var apiKeyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "APIKey",
	Fields: graphql.Fields{
		"key":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"comment":     &graphql.Field{Type: graphql.String},
		"permissions": &graphql.Field{Type: graphql.NewNonNull(graphql.NewList(graphql.String))},
		"createdAt":   &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var (
	ErrUnauthorized     = errors.New("authorization failed")
	ErrInsuffPermission = errors.New("insufficient permission")
)

func getAPIKey(key string) (*APIKey, bool) {
	if key == "" {
		return nil, false
	}

	var apiKey APIKey
	database.DB.Where(&APIKey{Key: key}).First(&apiKey)
	if apiKey.UserUUID == "" {
		return nil, false
	}
	return &apiKey, true
}

func checkPermission(key string, serviceName string) (*APIKey, error) {
	apiKey, ok := getAPIKey(key)
	if !ok {
		return nil, ErrUnauthorized
	}

	for _, p := range apiKey.Permissions {
		if p == (serviceName + "::*") {
			return apiKey, nil
		}
	}
	return apiKey, ErrInsuffPermission
}

// APIKeysQuery returns [APIKey!]!
var APIKeysQuery = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(apiKeyType))),
	Description: "[Auth] Get api keys",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := getUserInfo(p.Context.Value("Authorization").(string))
		if err != nil {
			return nil, err
		}

		var keys []APIKey
		database.DB.Where(&APIKey{UserUUID: user.UUID}).Find(&keys)
		return keys, nil
	},
}

// CreateAPIKeyMutation returns APIKey!
var CreateAPIKeyMutation = &graphql.Field{
	Type:        graphql.NewNonNull(apiKeyType),
	Description: "[Auth] Create a new api key",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "CreateAPIKeyInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"comment": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := getUserInfo(p.Context.Value("Authorization").(string))
		if err != nil {
			return nil, err
		}

		input := p.Args["input"].(map[string]interface{})
		key := APIKey{
			UserUUID:    user.UUID,
			Key:         random.GenerateRandomString(40),
			Comment:     input["comment"].(string),
			Permissions: []string{},
		}

		errs := database.DB.Save(&key).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return key, nil
	},
}

// AddPermissionToAPIKeyMutation returns APIKey!
var AddPermissionToAPIKeyMutation = &graphql.Field{
	Type:        graphql.NewNonNull(apiKeyType),
	Description: "[Auth] Add a permission to the api key",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "AddPermissionToAPIKeyInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"key":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"permission": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := getUserInfo(p.Context.Value("Authorization").(string))
		if err != nil {
			return nil, err
		}

		input := p.Args["input"].(map[string]interface{})
		return addPermission(user.UUID, input["key"].(string), input["permission"].(string))
	},
}

func addPermission(userUUID string, key string, permission string) (*APIKey, error) {
	var k APIKey
	database.DB.Where(&APIKey{UserUUID: userUUID, Key: key}).First(&k)
	if k.UserUUID == "" {
		return nil, ErrUnauthorized
	}

	// If permission already exists, do nothing.
	for _, p := range k.Permissions {
		if p == permission {
			return &k, nil
		}
	}

	k.Permissions = append(k.Permissions, permission)
	database.DB.Save(&k)
	return &k, nil
}

// RemovePermissionToAPIKeyMutation returns APIKey!
var RemovePermissionToAPIKeyMutation = &graphql.Field{
	Type:        graphql.NewNonNull(apiKeyType),
	Description: "[Auth] Remove a permission to the api key",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "RemovePermissionToAPIKeyInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"key":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"permission": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := getUserInfo(p.Context.Value("Authorization").(string))
		if err != nil {
			return nil, err
		}

		var key APIKey
		input := p.Args["input"].(map[string]interface{})
		database.DB.Where(&APIKey{UserUUID: user.UUID, Key: input["key"].(string)}).First(&key)
		if key.UserUUID == "" {
			return nil, fmt.Errorf("bad request")
		}

		permissionToRemove := input["permission"].(string)
		var newPermissions []string
		for _, p := range key.Permissions {
			if p != permissionToRemove {
				newPermissions = append(newPermissions, p)
			}
		}

		key.Permissions = newPermissions
		database.DB.Save(&key)
		return key, nil
	},
}
