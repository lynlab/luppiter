package services

import (
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

		var key APIKey
		input := p.Args["input"].(map[string]interface{})
		database.DB.Where(&APIKey{UserUUID: user.UUID, Key: input["key"].(string)}).First(&key)
		if key.UserUUID == "" {
			return nil, fmt.Errorf("bad request")
		}

		// If permission already exists, do nothing.
		permissionToAdd := input["permission"].(string)
		for _, p := range key.Permissions {
			if p == permissionToAdd {
				return key, nil
			}
		}

		key.Permissions = append(key.Permissions, permissionToAdd)
		database.DB.Save(&key)
		return key, nil
	},
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
