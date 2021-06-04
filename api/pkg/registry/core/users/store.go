package users

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/auth/keycloak"
	"github.com/nrc-no/core/api/pkg/registry/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// KeycloakStore implements the storage for Users
type KeycloakStore struct {
	keycloakClient *keycloak.KeycloakClient
	realmName      string
}

func NewREST(client *keycloak.KeycloakClient, realmName string) *KeycloakStore {
	return &KeycloakStore{
		keycloakClient: client,
		realmName:      realmName,
	}
}

var _ rest.Updater = &KeycloakStore{}
var _ rest.Deleter = &KeycloakStore{}
var _ rest.Lister = &KeycloakStore{}
var _ rest.Creater = &KeycloakStore{}

func (k KeycloakStore) Create(ctx context.Context, name string, obj runtime.Object) (runtime.Object, error) {

	castObj := obj.(*core.User)

	keycloakUser := mapUserToKeycloakUser(castObj)

	token, err := k.keycloakClient.GetToken(k.realmName)
	if err != nil {
		return nil, err
	}

	result, err := k.keycloakClient.CreateUser(token, k.realmName, keycloakUser)
	if err != nil {
		return nil, err
	}

	ret := mapKeycloakUserToUser(result)

	return ret, nil
}

func (k KeycloakStore) NewList() runtime.Object {
	return &core.UserList{}
}

func (k KeycloakStore) List(ctx context.Context) (runtime.Object, error) {

	token, err := k.keycloakClient.GetToken(k.realmName)
	if err != nil {
		return nil, err
	}

	keycloakUsers, err := k.keycloakClient.ListUsers(token, k.realmName, keycloak.ListUserOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]core.User, len(keycloakUsers))
	for i, keycloakUser := range keycloakUsers {
		user := mapKeycloakUserToUser(keycloakUser)
		items[i] = *user
	}

	return &core.UserList{
		Items: items,
	}, nil

}

func (k KeycloakStore) Delete(ctx context.Context, name string) (runtime.Object, bool, error) {

	token, err := k.keycloakClient.GetToken(k.realmName)
	if err != nil {
		return nil, false, err
	}

	keycloakUser, err := k.keycloakClient.GetUserByUsername(token, k.realmName, name)
	if err != nil {
		return nil, false, err
	}

	user := mapKeycloakUserToUser(keycloakUser)

	if err := k.keycloakClient.DeleteUser(token, k.realmName, keycloakUser.Username); err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (k KeycloakStore) New() runtime.Object {
	return &core.User{}
}

func mapKeycloakUserToUser(keycloakUser *keycloak.User) *core.User {
	attributes := keycloakUser.Attributes
	if len(keycloakUser.FirstName) > 0 {
		attributes["firstName"] = []string{keycloakUser.FirstName}
	}
	if len(keycloakUser.LastName) > 0 {
		attributes["lastName"] = []string{keycloakUser.LastName}
	}
	if len(keycloakUser.Email) > 0 {
		attributes["email"] = []string{keycloakUser.Email}
	}
	user := &core.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: keycloakUser.Username,
			UID:  types.UID(keycloakUser.ID),
		},
		Attributes: attributes,
	}
	return user
}

func mapUserToKeycloakUser(user *core.User) *keycloak.User {
	keycloakUser := &keycloak.User{
		ID:         string(user.UID),
		Username:   user.Name,
		Attributes: user.Attributes,
		Enabled:    true,
	}
	return keycloakUser
}

func (k KeycloakStore) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, error) {

	token, err := k.keycloakClient.GetToken(k.realmName)
	if err != nil {
		return nil, err
	}

	keycloakUser, err := k.keycloakClient.GetUserByUsername(token, k.realmName, name)
	if err != nil {
		return nil, err
	}

	user := mapKeycloakUserToUser(keycloakUser)

	updated, err := objInfo.UpdatedObject(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser := updated.(*core.User)

	updatedKeycloakUser := mapUserToKeycloakUser(updatedUser)

	_, err = k.keycloakClient.UpdateUser(token, k.realmName, updatedKeycloakUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil

}
