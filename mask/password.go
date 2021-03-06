// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package mask

import (
	"fmt"
	"unicode/utf8"
)

func Password(s string) string {
	if utf8.RuneCountInString(s) < 5 {
		return "*****" // too short, we can't mask anything
	} else {
		// turn 'password' into 'p*****d'
		first, last := s[0:1], s[len(s)-1:]
		return fmt.Sprintf("%s*****%s", first, last)
	}
}
