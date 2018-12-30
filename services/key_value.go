package services

import (
	"luppiter/components/database"
	"time"

	"github.com/graphql-go/graphql"
)

type KeyValueItem struct {
	UserUUID string `gorm:"type:varchar(255);primary_key"`
	Key      string `gorm:"type:varchar(255);primary_key"`
	Value    string `gorm:"type:varchar(1024)"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var keyValueItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "KeyValueItem",
	Fields: graphql.Fields{
		"key":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"value": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

// KeyValueItemQuery returns KeyValueItem
var KeyValueItemQuery = &graphql.Field{
	Type:        keyValueItemType,
	Description: "[KeyValue] Get a key value item",
	Args: graphql.FieldConfigArgument{
		"key": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		key, err := checkPermission(p.Context.Value("APIKey").(string), "KeyValue")
		if err != nil {
			return nil, err
		}

		var item KeyValueItem
		database.DB.Where(&KeyValueItem{UserUUID: key.UserUUID, Key: p.Args["key"].(string)}).First(&item)

		if item.Value == "" {
			return nil, nil
		}
		return item, nil
	},
}

// SetKeyValueItemMutation returns KeyValueItem!
var SetKeyValueItemMutation = &graphql.Field{
	Type:        graphql.NewNonNull(keyValueItemType),
	Description: "[KeyValue] Create or update a key value item.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "SetKeyValueItemInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"key":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"value": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		key, err := checkPermission(p.Context.Value("APIKey").(string), "KeyValue")
		if err != nil {
			return nil, err
		}

		input := p.Args["input"].(map[string]interface{})

		var item KeyValueItem
		database.DB.Where(&KeyValueItem{UserUUID: key.UserUUID, Key: input["key"].(string)}).First(&item)

		var errs []error
		if item.Value != "" {
			// If key already exists, update it.
			item.Value = input["value"].(string)
			errs = database.DB.Save(&item).GetErrors()
		} else {
			// If xnot exists, create new one.
			item = KeyValueItem{
				UserUUID: key.UserUUID,
				Key:      input["key"].(string),
				Value:    input["value"].(string),
			}
			errs = database.DB.Create(&item).GetErrors()
		}

		if len(errs) > 0 {
			return nil, errs[0]
		}
		return item, nil
	},
}
