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

package groups

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/bcmi-labs/hydrasdk/common"
	"github.com/pkg/errors"
)

// Manager provides methods to create and update ladon groups
type Manager struct {
	Endpoint *url.URL
	Client   *http.Client
}

// Group allows or denies certain Subjects to perform certain Actions on certain Resources.
// Subjects, Resources, and Actions can be strings ('user:0001') or regexes ('resource:<.+>')
type Group struct {
	ID      string   `json:"id"`
	Members []string `json:"members"`
}

// NewManager returns a Manager connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewManager(id, secret, cluster string) (*Manager, error) {
	endpoint, client, err := common.Authenticate(id, secret, cluster, "hydra")
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate ClientManager")
	}
	manager := Manager{
		Endpoint: common.JoinURL(endpoint, "warden", "groups"),
		Client:   client,
	}
	return &manager, nil
}

// List calls the hydra api to list all groups
func (m *Manager) List() ([]string, error) {
	url := m.Endpoint.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var groups []string

	err = common.Bind(m.Client, req, &groups)
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}
	return groups, nil
}

// Create calls the hydra api to create a new group
func (m *Manager) Create(group *Group) error {
	url := m.Endpoint.String()

	payload, err := json.Marshal(group)
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", group)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, &group)
	if err != nil {
		return errors.Wrap(err, "Create")
	}
	return nil
}

// OfUser calls the hydra api to return the groups of a user
func (m *Manager) OfUser(id string) ([]string, error) {
	url := common.CopyURL(m.Endpoint)
	values := url.Query()
	values.Add("member", id)
	url.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var groups []string

	err = common.Bind(m.Client, req, &groups)
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}
	return groups, nil
}

// AddMembers calls the hydra api to add members to a group
func (m *Manager) AddMembers(id string, members []string) error {
	url := common.JoinURL(m.Endpoint, id, "members").String()

	payload, err := json.Marshal(struct {
		Members []string `json:"members"`
	}{Members: members})
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", members)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, nil)
	if err != nil {
		return errors.Wrapf(err, "Update %s", id)
	}
	return nil
}

// RemoveMembers calls the hydra api to remove members to a group
func (m *Manager) RemoveMembers(id string, members []string) error {
	url := common.JoinURL(m.Endpoint, id, "members").String()

	payload, err := json.Marshal(struct {
		Members []string `json:"members"`
	}{Members: members})
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", members)
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, nil)
	if err != nil {
		return errors.Wrapf(err, "Update %s", id)
	}
	return nil
}

// Get calls the hydra api to return a specific group
func (m *Manager) Get(id string) (*Group, error) {
	url := common.JoinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var group Group

	err = common.Bind(m.Client, req, &group)
	if err != nil {
		return nil, errors.Wrapf(err, "Get %s", id)
	}
	return &group, nil
}

// Delete calls the hydra api to remove a specific Group
func (m *Manager) Delete(id string) error {
	url := common.JoinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = common.Bind(m.Client, req, nil)
	if err != nil {
		return errors.Wrapf(err, "Delete %s", id)
	}
	return nil
}
