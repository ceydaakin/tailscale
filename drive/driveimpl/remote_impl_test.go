// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package driveimpl

import (
	"os/user"
	"testing"

	"tailscale.com/tstest"
)

func TestAssertCanRunAsSelf(t *testing.T) {
	tests := []struct {
		name     string
		current  user.User
		shareAs  string
		wantFail bool
	}{
		{"non-root as other user", user.User{Uid: "1000", Username: "alice"}, "bob", false},
		{"non-root as self", user.User{Uid: "1000", Username: "alice"}, "alice", false},
		{"root as other user", user.User{Uid: "0", Username: "root"}, "alice", true},
		{"root without share user", user.User{Uid: "0", Username: "root"}, "", true},
		// OpenWrt runs tailscaled as root and has neither sudo nor su, and
		// shares made by root are shared as root (#22413).
		{"root as root", user.User{Uid: "0", Username: "root"}, "root", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tstest.Replace(t, &currentUser, func() (*user.User, error) { return &tt.current, nil })
			s := &userServer{username: tt.shareAs}
			err := s.assertCanRunAsSelf()
			if gotFail := err != nil; gotFail != tt.wantFail {
				t.Errorf("assertCanRunAsSelf() = %v, want failure %v", err, tt.wantFail)
			}
		})
	}
}
