/*
 * This file is part of hydrasdk
 *
 * hydrasdk is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 ARDUINO AG (http://www.arduino.cc/)
 */

package hydrasdk

import (
	"crypto/rsa"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	jose "github.com/square/go-jose"
)

// KeyGetter provides functions to retrieve a key from an hydra set (tipically the first)
type KeyGetter interface {
	GetRSAPublic(set string) (*rsa.PublicKey, error)
	GetRSAPrivate(set string) (*rsa.PrivateKey, error)
}

// CachedKeyManager uses hydra rest api to retrieve keys and cache them for easy access
type CachedKeyManager struct {
	Endpoint *url.URL
	Client   *http.Client

	rsaPublics  map[string]*rsa.PublicKey
	rsaPrivates map[string]*rsa.PrivateKey
}

// NewCachedKeyManager returns a CachedKeyManager connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewCachedKeyManager(id, secret, cluster string) (*CachedKeyManager, error) {
	endpoint, client, err := authenticate(id, secret, cluster)
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate ClientManager")
	}
	manager := CachedKeyManager{
		Endpoint:    joinURL(endpoint, "keys"),
		Client:      client,
		rsaPublics:  map[string]*rsa.PublicKey{},
		rsaPrivates: map[string]*rsa.PrivateKey{},
	}
	return &manager, nil
}

// GetRSAPublic retrieves the first key of the given set. It caches them forever,
// so hope that they don't change
func (m CachedKeyManager) GetRSAPublic(set string) (*rsa.PublicKey, error) {
	// Try getting from cache
	if key, ok := m.rsaPublics[set]; ok {
		return key, nil
	}

	url := joinURL(m.Endpoint, set).String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var keyset jose.JsonWebKeySet
	err = bind(m.Client, req, &keyset)
	if err != nil {
		return nil, err
	}

	if len(keyset.Keys) == 0 {
		return nil, errors.New("The retrieved keyset is empty")
	}

	key, ok := keyset.Keys[0].Key.(*rsa.PublicKey)
	if !ok {
		return key, errors.New("Could not convert key to RSA Private Key.")
	}

	// Save on cache
	m.rsaPublics[set] = key

	return key, nil
}

// GetRSAPrivate retrieves the first key of the given set. It caches them forever,
// so hope that they don't change
func (m CachedKeyManager) GetRSAPrivate(set string) (*rsa.PrivateKey, error) {
	// Try getting from cache
	if key, ok := m.rsaPrivates[set]; ok {
		return key, nil
	}

	url := joinURL(m.Endpoint, set).String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var keyset jose.JsonWebKeySet
	err = bind(m.Client, req, &keyset)
	if err != nil {
		return nil, err
	}

	if len(keyset.Keys) == 0 {
		return nil, errors.New("The retrieved keyset is empty")
	}

	key, ok := keyset.Keys[0].Key.(*rsa.PrivateKey)
	if !ok {
		return key, errors.New("Could not convert key to RSA Private Key.")
	}

	// Save on cache
	m.rsaPrivates[set] = key

	return key, nil
}
