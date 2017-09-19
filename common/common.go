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

package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
)

// JoinURL creates an url from the given parts
func JoinURL(u *url.URL, args ...string) (ep *url.URL) {
	ep = CopyURL(u)
	ep.Path = path.Join(append([]string{ep.Path}, args...)...)
	return ep
}

// CopyURL returns a copy of the url
func CopyURL(u *url.URL) *url.URL {
	a := new(url.URL)
	*a = *u
	return a
}

// Bind does a get request and binds the body to the given interface
func Bind(client *http.Client, req *http.Request, o interface{}) error {
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "execute request %+v", req)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode > 299 {
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusOK, resp.StatusCode, body)
	} else if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(o); err != nil {
		return errors.Wrapf(err, "decode json %s", body)
	}
	return nil
}

// Authenticate returns the url of the cluster and an authenticated Client
func Authenticate(id, secret, cluster string) (*url.URL, *http.Client, error) {
	uri, err := url.Parse(cluster)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "parse url %s", cluster)
	}
	credentials := clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     JoinURL(uri, "oauth2/token").String(),
		Scopes:       []string{"hydra"},
	}

	ctx := context.Background()
	_, err = credentials.Token(ctx)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "connect to cluster %s", cluster)
	}
	return uri, credentials.Client(ctx), nil
}
