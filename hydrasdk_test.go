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

package hydrasdk_test

import (
	"testing"

	"github.com/bcmi-labs/hydrasdk"
)

func TestGetPolicies(t *testing.T) {
	policies, err := hydrasdk.NewPoliciesManager("admin", "demo-password", "http://localhost:4444")
	if err != nil {
		t.Error(err)
	}

	// Create some policies
	policies.Create(&hydrasdk.Policy{ID: "1"})
	policies.Create(&hydrasdk.Policy{ID: "2"})
	policies.Create(&hydrasdk.Policy{ID: "3"})
	policies.Create(&hydrasdk.Policy{ID: "4"})

	list, err := policies.GetAll()
	if err != nil {
		t.Error(err)
	}

	if len(list) != 5 { // One of those is the default one
		t.Errorf("Expected 5 policies, got %d: %s", len(list), list)
	}

	// Cleanup
	cleanPolicies(policies)
}

func TestCreatePolicy(t *testing.T) {
	policies, err := hydrasdk.NewPoliciesManager("admin", "demo-password", "http://localhost:4444")
	if err != nil {
		t.Error(err)
	}

	payload := hydrasdk.Policy{
		ID:          "example policy",
		Description: "exmaple policy",
		Subjects:    []string{"me", "my-friend"},
		Effect:      "allow",
		Actions:     []string{"eat"},
		Resources:   []string{"banana", "cake"},
	}

	err = policies.Create(&payload)
	if err != nil {
		t.Error(err)
	}

	policy, err := policies.Get(payload.ID)
	if policy.ID != payload.ID {
		t.Errorf("Expected ID='%s', got %s", payload.ID, policy.ID)
	}

	cleanPolicies(policies)

}

func cleanPolicies(policies *hydrasdk.PoliciesManager) {
	policies.Delete("1")
	policies.Delete("2")
	policies.Delete("3")
	policies.Delete("4")
	policies.Delete("example policy")
}
