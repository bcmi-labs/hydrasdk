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

// Package clients contains methods to retrieve and save hydra clients
package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/bcmi-labs/hydrasdk/common"
	"github.com/pkg/errors"
)

// Client is an oauth2 client saved on hydra database
type Client struct {
	ID                string   `json:"id" gorethink:"id"`
	Name              string   `json:"client_name" gorethink:"client_name"`
	Secret            string   `json:"client_secret,omitempty" gorethink:"client_secret"`
	RedirectURIs      []string `json:"redirect_uris,omitempty" gorethink:"redirect_uris"`
	GrantTypes        []string `json:"grant_types" gorethink:"grant_types"`
	ResponseTypes     []string `json:"response_types,omitempty" gorethink:"response_types"`
	Scope             string   `json:"scope,omitempty" gorethink:"scope"`
	Owner             string   `json:"owner,omitempty" gorethink:"owner"`
	PolicyURI         string   `json:"policy_uri,omitempty" gorethink:"policy_uri"`
	TermsOfServiceURI string   `json:"tos_uri,omitempty" gorethink:"tos_uri"`
	ClientURI         string   `json:"client_uri,omitempty" gorethink:"client_uri"`
	LogoURI           string   `json:"logo_uri,omitempty" gorethink:"logo_uri"`
	Contacts          []string `json:"contacts,omitempty" gorethink:"contacts"`
	Public            bool     `json:"public" gorethink:"public"`
}

// ClientGetter is an abstraction that allows you to retrieve a specific client by their ID
type ClientGetter interface {
	Get(id string) (*Client, error)
}

// Manager uses hydra rest apis to retrieve clients
type Manager struct {
	Endpoint *url.URL
	Client   *http.Client
}

// NewManager returns a Manager connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewManager(id, secret, cluster string) (*Manager, error) {
	endpoint, client, err := common.Authenticate(id, secret, cluster, "hydra")
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate Manager")
	}

	manager := Manager{
		Endpoint: common.JoinURL(endpoint, "clients"),
		Client:   client,
	}
	return &manager, nil
}

// Create queries the hydra api to create a new client
func (m Manager) Create(client *Client) error {
	url := m.Endpoint.String()

	payload, err := json.Marshal(*client)
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", client)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, client)
	if err != nil {
		return err
	}

	return nil
}

// Delete queries the hydra api to retrieve a specific client by their ID.
func (m Manager) Delete(id string) error {
	url := common.JoinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// Get queries the hydra api to retrieve a specific client by their ID.
func (m Manager) Get(id string) (*Client, error) {
	url := common.JoinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var client *Client

	err = common.Bind(m.Client, req, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetAll calls the hydra api to return all the clients
func (m *Manager) GetAll() (map[string]Client, error) {
	req, err := http.NewRequest("GET", m.Endpoint.String(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", m.Endpoint.String())
	}

	var clients map[string]Client

	err = common.Bind(m.Client, req, &clients)
	if err != nil {
		return nil, errors.Wrap(err, "GetAll")
	}
	return clients, nil
}

// Update calls the hydra api to update a specific client
func (m *Manager) Update(id string, client *Client) error {
	url := common.JoinURL(m.Endpoint, id).String()

	payload, err := json.Marshal(client)
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", client)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, client)
	if err != nil {
		return errors.Wrapf(err, "Update %s", id)
	}
	return nil
}
