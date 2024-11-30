//go:build darwin

package shims

import (
	"encoding/json"
	"errors"
	errors2 "github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/secrets/errors"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/auth"
	"github.com/keybase/go-keychain"
	"log"
)

type MacSecretsShim struct {
}

func (m *MacSecretsShim) SetSecret(client auth.ClientCredentials) error {

	// remove existing secret, if there is one
	m.RemoveSecret()

	serialized, err := json.Marshal(client)
	if err != nil {
		log.Fatalf("error serializing secret: %v", err)
	}

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(Service)
	item.SetAccount(Account)
	item.SetAccessGroup(AccessGroup)
	item.SetData(serialized)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	e := keychain.AddItem(item)
	if errors.Is(e, keychain.ErrorDuplicateItem) {
		return &errors2.DuplicateError{
			Err:  "credential exists",
			Hint: "use the logout command to remove any existing credential and try again",
		}
	}
	return e
}

func (m *MacSecretsShim) GetSecret() (auth.ClientCredentials, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(Service)
	query.SetAccount(Account)
	query.SetAccessGroup(AccessGroup)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		return auth.ClientCredentials{}, err
	} else if len(results) != 1 {
		return auth.ClientCredentials{}, errors2.NewNotFoundError()
	} else {
		var client auth.ClientCredentials
		err := json.Unmarshal(results[0].Data, &client)
		if err != nil {
			return auth.ClientCredentials{}, err
		}
		return client, nil
	}
}

func (m *MacSecretsShim) RemoveSecret() error {
	deleteItem := keychain.NewItem()
	deleteItem.SetSecClass(keychain.SecClassGenericPassword)
	deleteItem.SetService(Service)
	deleteItem.SetAccount(Account)
	deleteItem.SetAccessGroup(AccessGroup)
	if err := keychain.DeleteItem(deleteItem); err != nil {
		if !errors.Is(err, keychain.ErrorItemNotFound) {
			return err
		}
	}
	return nil
}
