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

package introspect

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bcmi-labs/hydrasdk/common"
	"github.com/codeclysm/introspector/v3"
	"github.com/pkg/errors"
)

// Introspector uses hydra rest apis to retrieve clients
type Introspector struct {
	AllowedEndpoint    *url.URL
	IntrospectEndpoint *url.URL
	Client             *http.Client
}

// NewIntrospector returns a Introspector connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewIntrospector(id, secret, cluster string) (*Introspector, error) {
	endpoint, client, err := common.Authenticate(id, secret, cluster, "hydra")
	if err != nil {
		return nil, errors.Wrap(err, "Instantiate Introspector")
	}

	manager := Introspector{
		AllowedEndpoint:    common.JoinURL(endpoint, "warden", "token", "allowed"),
		IntrospectEndpoint: common.JoinURL(endpoint, "oauth2", "introspect"),
		Client:             client,
	}
	return &manager, nil
}

// Introspect queries the endpoint with an http request. It expects that the endpoint
// implements https://tools.ietf.org/html/rfc7662
func (m *Introspector) Introspect(token string) (introspector.Introspection, error) {
	data := url.Values{
		"token": []string{token},
	}

	url := m.IntrospectEndpoint.String()
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return introspector.Introspection{}, errors.Wrapf(err, "new request for %s", url)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var i introspector.Introspection
	err = common.Bind(m.Client, req, &i)
	if err != nil {
		return introspector.Introspection{}, err
	}

	if !i.Active {
		return introspector.Introspection{}, errors.New("token not active")
	}

	return i, nil
}
