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
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// PoliciesManager provides methods to create and update ladon policies
type PoliciesManager struct {
	Endpoint *url.URL
	Client   *http.Client
}

// Policy allows or denies certain Subjects to perform certain Actions on certain Resources.
// Subjects, Resources, and Actions can be strings ('user:0001') or regexes ('resource:<.+>')
type Policy struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Subjects    []string               `json:"subjects"`
	Effect      string                 `json:"effect"`
	Resources   []string               `json:"resources"`
	Actions     []string               `json:"actions"`
	Conditions  map[string]interface{} `json:"conditions"`
}

// NewPoliciesManager returns a PoliciesManager connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewPoliciesManager(id, secret, cluster string) (*PoliciesManager, error) {
	endpoint, client, err := authenticate(id, secret, cluster)
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate ClientManager")
	}
	manager := PoliciesManager{
		Endpoint: joinURL(endpoint, "policies"),
		Client:   client,
	}
	return &manager, nil
}

// Create calls the hydra api to create a new policy
func (m *PoliciesManager) Create(policy *Policy) error {
	url := m.Endpoint.String()

	payload, err := json.Marshal(policy)
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", policy)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = bind(m.Client, req, &policy)
	if err != nil {
		return errors.Wrap(err, "Create")
	}
	return nil
}

// Update calls the hydra api to update a specific policy
func (m *PoliciesManager) Update(id string, policy *Policy) error {
	url := joinURL(m.Endpoint, id).String()

	payload, err := json.Marshal(policy)
	if err != nil {
		return errors.Wrapf(err, "json marshal of %v", policy)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = bind(m.Client, req, nil)
	if err != nil {
		return errors.Wrapf(err, "Update %s", id)
	}
	return nil
}

// GetAll calls the hydra api to return all the policies
func (m *PoliciesManager) GetAll() ([]Policy, error) {
	req, err := http.NewRequest("GET", m.Endpoint.String(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", m.Endpoint.String())
	}

	var policies []Policy

	err = bind(m.Client, req, &policies)
	if err != nil {
		return nil, errors.Wrap(err, "GetAll")
	}
	return policies, nil
}

// Get calls the hydra api to return a specific policy
func (m *PoliciesManager) Get(id string) (*Policy, error) {
	url := joinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request for %s", url)
	}

	var policy Policy

	err = bind(m.Client, req, &policy)
	if err != nil {
		return nil, errors.Wrapf(err, "Get %s", id)
	}
	return &policy, nil
}

// Delete calls the hydra api to remove a specific policy
func (m *PoliciesManager) Delete(id string) error {
	url := joinURL(m.Endpoint, id).String()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	err = bind(m.Client, req, nil)
	if err != nil {
		return errors.Wrapf(err, "Delete %s", id)
	}
	return nil
}
