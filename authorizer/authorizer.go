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

package authorizer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bcmi-labs/hydrasdk/common"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
)

// Authorizer uses hydra rest apis to retrieve clients
type Authorizer struct {
	AllowedEndpoint *url.URL
	Client          *http.Client
}

// NewAuthorizer returns a Warden authorizer connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewAuthorizer(id, secret, cluster string) (*Authorizer, error) {
	endpoint, client, err := common.Authenticate(id, secret, cluster, "hydra")
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate Authorizer")
	}

	manager := Authorizer{
		AllowedEndpoint: common.JoinURL(endpoint, "warden", "allowed"),
		Client:          client,
	}
	return &manager, nil
}

// IsAllowed calls the hydra endpoint to see if a subject has the permission to perform an action
func (m *Authorizer) IsAllowed(request *ladon.Request) error {
	data, err := json.Marshal(&request)
	if err != nil {
		return errors.Wrapf(err, "marshal request")
	}

	url := m.AllowedEndpoint.String()
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return errors.Wrapf(err, "new request for %s", url)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	var res struct {
		Allowed bool `json:"allowed"`
	}
	err = common.Bind(m.Client, req, &res)
	if err != nil {
		return err
	}

	if !res.Allowed {
		return errors.New("not allowed")
	}

	return nil
}
