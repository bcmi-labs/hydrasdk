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

package groups_test

import (
	"testing"

	"github.com/bcmi-labs/hydrasdk/groups"
)

func TestCreateGroup(t *testing.T) {
	manager, err := groups.NewManager("admin", "demo-password", "http://localhost:4444")
	if err != nil {
		t.Error(err)
	}

	payload := groups.Group{
		ID:      "example group",
		Members: []string{"admin", "user1"},
	}

	err = manager.Create(&payload)
	if err != nil {
		t.Error(err)
	}

	group, err := manager.Get("example group")
	if err != nil {
		t.Error(err)
	}

	if group.ID != payload.ID {
		t.Error("group id should be equal")
	}

	if !in(group.Members, "admin") {
		t.Error("admin should be in group members")
	}

	if !in(group.Members, "user1") {
		t.Error("user1 should be in group members")
	}
	if in(group.Members, "user2") {
		t.Error("user2 shouldn't be in group members")
	}

	// Check admin groups
	groups, err := manager.OfUser("admin")
	if err != nil {
		t.Error(err)
	}

	if !in(groups, "example group") {
		t.Error("example group should be in admin's groups")
	}

	// Check user1 groups
	groups, err = manager.OfUser("user1")
	if err != nil {
		t.Error(err)
	}

	if !in(groups, "example group") {
		t.Error("example group should be in user1's groups")
	}

	// Check user2 groups
	groups, err = manager.OfUser("user2")
	if err != nil {
		t.Error(err)
	}

	if in(groups, "example group") {
		t.Error("example group shouldn't be in user2's groups")
	}

	cleangroups(manager)

}

func TestMembers(t *testing.T) {
	manager, err := groups.NewManager("admin", "demo-password", "http://localhost:4444")
	if err != nil {
		t.Error(err)
	}

	payload := groups.Group{
		ID:      "example group",
		Members: []string{"admin"},
	}

	err = manager.Create(&payload)
	if err != nil {
		t.Error(err)
	}

	// Check members group
	groups, err := manager.OfUser("user1")
	if err != nil {
		t.Error(err)
	}
	if in(groups, "example group") {
		t.Error("example group shouldn't be in user1's groups")
	}

	// Add user
	err = manager.AddMembers("example group", []string{"user1"})
	if err != nil {
		t.Error(err)
	}

	groups, err = manager.OfUser("user1")
	if err != nil {
		t.Error(err)
	}
	if !in(groups, "example group") {
		t.Error("example group should be in user1's groups")
	}

	// Remove user
	err = manager.RemoveMembers("example group", []string{"user1"})
	if err != nil {
		t.Error(err)
	}

	groups, err = manager.OfUser("user1")
	if err != nil {
		t.Error(err)
	}
	if in(groups, "example group") {
		t.Error("example group shouldn't be in user1's groups")
	}
	cleangroups(manager)

}

func in(slice []string, el string) bool {
	for i := range slice {
		if slice[i] == el {
			return true
		}
	}
	return false
}

func cleangroups(groups *groups.Manager) {
	groups.Delete("example group")
}
