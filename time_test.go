// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	est *time.Location
)

func init() {
	// US Eastern Time Zone
	est, _ = time.LoadLocation("America/New_York")
}

func TestTime(t *testing.T) {
	t1 := Now(time.UTC)

	// Verify time.Time methods work
	if diff := t1.Sub(t1.Time); diff != 0 {
		t.Errorf("got %v", diff)
	}
	if tt := time.Now().Add(1 * time.Second); t1.Sub(tt) == 0 {
		t.Error("expected difference in timing")
	}
}

func TestTime__NewTime(t *testing.T) {
	f := func(_ Time) {}
	f(NewTime(time.Now())) // make sure we can lift time.Time values

	start := time.Now().Add(-1 * time.Second)

	// Example from NewTime godoc
	now := Now(time.UTC)
	fmt.Println(start.Sub(now.Time))
}

func TestTime__JSON(t *testing.T) {
	// marshal and then unmarshal
	t1 := Now(time.UTC)

	bs, err := t1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var t2 Time
	if err := json.Unmarshal(bs, &t2); err != nil {
		t.Fatal(err)
	}
	if !t1.Equal(t2) {
		t.Errorf("unequal: t1=%q t2=%q", t1, t2)
	}

	in := []byte(`"2018-11-27T00:54:53Z"`)
	var t3 Time
	if err := json.Unmarshal(in, &t3); err != nil {
		t.Fatal(err)
	}
	if t3.IsZero() {
		t.Error("t3 shouldn't be zero time")
	}

	// empty should unmarshal to nothing
	in = []byte(`""`)
	var t4 Time
	if err := json.Unmarshal(in, &t4); err != nil {
		t.Errorf("empty value for base.Time is fine, but got: %v", err)
	}
}

func TestTime__jsonRFC3339(t *testing.T) {
	// Read RFC 3339 time
	in := []byte(fmt.Sprintf(`"%s"`, time.Now().Format(time.RFC3339)))
	var t1 Time
	if err := json.Unmarshal(in, &t1); err != nil {
		t.Fatal(err)
	}
	if t1.IsZero() {
		t.Error("t4 shouldn't be zero time")
	}
}

func TestTime__javascript(t *testing.T) {
	// Generated with (new Date).toISOString() in Chrome and Firefox
	in := []byte(`{"time": "2018-12-14T20:36:58.789Z"}`)

	type wrapper struct {
		When Time `json:"time"`
	}
	var wrap wrapper
	if err := json.Unmarshal(in, &wrap); err != nil {
		t.Fatal(err)
	}
	if v := wrap.When.String(); v != "2018-12-14 20:36:58 +0000 UTC" {
		t.Errorf("got %q", v)
	}
}

var quote = []byte(`"`)

// TestTime__ruby will attempt to parse an ISO 8601 time generated by this library
func TestTime__ruby(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping ruby ISO 8601 test on windows")
	}

	bin, err := exec.LookPath("ruby")
	if err != nil || bin == "" {
		if inCI := os.Getenv("TRAVIS_OS_NAME") != ""; inCI {
			t.Fatal("ruby not found")
		} else {
			t.Skip("ruby not found")
		}
	}

	tt, err := time.Parse(ISO8601Format, "2018-11-18T09:04:23-08:00")
	if err != nil {
		t.Fatal(err)
	}
	t1 := Time{
		Time: tt,
	}

	bs, err := t1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	bs = bytes.TrimPrefix(bytes.TrimSuffix(bs, quote), quote)

	// Check with ruby
	cmd := exec.Command(bin, "time.rb", string(bs))
	bs, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("err=%v\nOutput: %v", err, string(bs))
	}

	// Validate ruby output
	if !bytes.Contains(bs, []byte(`Date: 2018-11-18`)) {
		t.Errorf("no Date: %v", string(bs))
	}
	if !bytes.Contains(bs, []byte(`Time: 09:04:23`)) {
		t.Errorf("no Time: %v", string(bs))
	}
}

func TestTime__IsBusinessDay(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected bool
	}{
		// new years day
		{time.Date(2018, time.January, 1, 1, 0, 0, 0, est), false},
		// Wednesday Canary test
		{time.Date(2018, time.January, 3, 1, 0, 0, 0, est), true},
		// saturday
		{time.Date(2018, time.January, 6, 1, 0, 0, 0, est), false},
		// sunday
		{time.Date(2018, time.January, 7, 1, 0, 0, 0, est), false},
		// Martin Luther King, JR. Day (Monday)
		{time.Date(2018, time.January, 15, 1, 0, 0, 0, est), false},
		// Memorial Day
		{time.Date(2018, time.May, 28, 1, 0, 0, 0, est), false},
		// Independence Day
		{time.Date(2018, time.July, 4, 1, 0, 0, 0, est), false},
		// Labor Day
		{time.Date(2018, time.September, 3, 1, 0, 0, 0, est), false},
		// Thanksgiving Day (Thursday)
		{time.Date(2018, time.November, 22, 1, 0, 0, 0, est), false},
		// Christmas Day (Sunday)
		{time.Date(2022, time.December, 25, 1, 0, 0, 0, est), false},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).IsBusinessDay()
		if actual != test.Expected {
			t.Fatalf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}

		actual = NewTime(test.Date).IsBusinessDay()
		if actual != test.Expected {
			t.Errorf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}
	}
}

func TestTime__IsBankingDay(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected bool
	}{
		// new years day
		{time.Date(2018, time.January, 1, 1, 0, 0, 0, est), false},
		// Wednesday Canary test
		{time.Date(2018, time.January, 3, 1, 0, 0, 0, est), true},
		// saturday
		{time.Date(2018, time.January, 6, 1, 0, 0, 0, est), false},
		// sunday
		{time.Date(2018, time.January, 7, 1, 0, 0, 0, est), false},
		// Martin Luther King, JR. Day
		{time.Date(2018, time.January, 15, 1, 0, 0, 0, est), false},
		// Presidents' Day
		{time.Date(2018, time.February, 19, 1, 0, 0, 0, est), false},
		// Memorial Day
		{time.Date(2018, time.May, 28, 1, 0, 0, 0, est), false},
		// Independence Day
		{time.Date(2018, time.July, 4, 1, 0, 0, 0, est), false},
		// Labor Day
		{time.Date(2018, time.September, 3, 1, 0, 0, 0, est), false},
		// Columbus Day
		{time.Date(2018, time.October, 8, 1, 0, 0, 0, est), false},
		// Veterans Day Observed on the monday
		{time.Date(2018, time.November, 12, 1, 0, 0, 0, est), false},
		// Thanksgiving Day
		{time.Date(2018, time.November, 22, 1, 0, 0, 0, est), false},
		// Christmas Day
		{time.Date(2018, time.December, 25, 1, 0, 0, 0, est), false},
		// Friday before Nov 2023 Veteran's Day
		{time.Date(2023, time.November, 10, 1, 0, 0, 0, est), true},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).IsBankingDay()
		if actual != test.Expected {
			t.Fatalf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}

		actual = NewTime(test.Date).IsBankingDay()
		if actual != test.Expected {
			t.Errorf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}
	}
}

func TestTime__IsWeekend(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected bool
	}{
		// saturday
		{time.Date(2018, time.January, 6, 1, 0, 0, 0, est), true},
		// sunday
		{time.Date(2018, time.January, 7, 1, 0, 0, 0, est), true},
		// monday
		{time.Date(2018, time.January, 9, 1, 0, 0, 0, est), false},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).IsWeekend()
		if actual != test.Expected {
			t.Fatalf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}

		actual = NewTime(test.Date).IsWeekend()
		if actual != test.Expected {
			t.Errorf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}
	}
}

func TestTime_AddBusinessDay(t *testing.T) {
	unchangeable := time.Date(2021, time.July, 15, 1, 0, 0, 0, est)
	tests := []struct {
		Date   time.Time
		Future time.Time
		Days   int
	}{
		// Thursday add one day needs to be friday
		{time.Date(2018, time.January, 11, 1, 0, 0, 0, est), time.Date(2018, time.January, 12, 1, 0, 0, 0, est), 1},
		// Thursday add two days over a monday holiday and needs to be following tuesday
		{time.Date(2018, time.January, 11, 1, 0, 0, 0, est), time.Date(2018, time.January, 16, 1, 0, 0, 0, est), 2},
		// Friday add two days over a monday holiday and needs to be following wednesday
		{time.Date(2018, time.January, 12, 1, 0, 0, 0, est), time.Date(2018, time.January, 17, 1, 0, 0, 0, est), 2},
		// Friday add two days over a sunday public holiday (moved to monday which is business day but not banking day) needs to be following tuesday
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2021, time.July, 6, 1, 0, 0, 0, est), 2},
		// Friday add one day over a sunday public holiday (moved to monday which is business day but not banking day) needs to be following monday
		{time.Date(2022, time.June, 17, 1, 0, 0, 0, est), time.Date(2022, time.June, 20, 1, 0, 0, 0, est), 1},
		// Negative input
		{unchangeable, unchangeable, 0},
		{unchangeable, unchangeable, -1},
		{unchangeable, unchangeable, -10},
		// Input above the max
		{unchangeable, unchangeable, 501},
		{unchangeable, unchangeable, 600},
		// Input at the max
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2023, time.June, 26, 1, 0, 0, 0, est), 500},
		// Find one year in the future
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2022, time.June, 28, 1, 0, 0, 0, est), 365 - 11 - (52 * 2)},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).AddBusinessDay(test.Days)
		if !actual.Equal(NewTime(test.Future)) {
			t.Errorf("Adding %d business days: expected %s, got %s", test.Days, NewTime(test.Future), actual)
		}
	}
}

func TestTime_GetHoliday(t *testing.T) {
	now := time.Date(2023, time.November, 10, 1, 0, 0, 0, est)
	holiday := NewTime(now).GetHoliday()
	require.NotNil(t, holiday)
	require.Equal(t, "Veterans Day", holiday.Name)
}

func TestTime_AddBankingDay(t *testing.T) {
	unchangeable := time.Date(2021, time.July, 15, 1, 0, 0, 0, est)
	tests := []struct {
		Date   time.Time
		Future time.Time
		Days   int
	}{
		// Thursday add one day needs to be friday
		{time.Date(2018, time.January, 11, 1, 0, 0, 0, est), time.Date(2018, time.January, 12, 1, 0, 0, 0, est), 1},
		// Thursday add two days over a monday holiday and needs to be following tuesday
		{time.Date(2018, time.January, 11, 1, 0, 0, 0, est), time.Date(2018, time.January, 16, 1, 0, 0, 0, est), 2},
		// Friday add two days over a monday holiday abd needs to be following wednesday
		{time.Date(2018, time.January, 12, 1, 0, 0, 0, est), time.Date(2018, time.January, 17, 1, 0, 0, 0, est), 2},
		// Friday add two days over a sunday public holiday (moved to monday) needs to be following wednesday
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2021, time.July, 7, 1, 0, 0, 0, est), 2},
		// Friday add one day over a sunday public holiday (moved to monday) needs to be following tuesday
		{time.Date(2022, time.June, 17, 1, 0, 0, 0, est), time.Date(2022, time.June, 21, 1, 0, 0, 0, est), 1},
		// Negative input
		{unchangeable, unchangeable, 0},
		{unchangeable, unchangeable, -1},
		{unchangeable, unchangeable, -10},
		// Input above the max
		{unchangeable, unchangeable, 501},
		{unchangeable, unchangeable, 600},
		// Input at the max
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2023, time.June, 30, 1, 0, 0, 0, est), 500},
		// Find one year in the future
		{time.Date(2021, time.July, 2, 1, 0, 0, 0, est), time.Date(2022, time.June, 27, 1, 0, 0, 0, est), 365 - 14 - (52 * 2)},
		// Late evening conversions should still fall on a late evening
		{time.Date(2022, time.July, 6, 20, 1, 9, 0, est), time.Date(2022, time.July, 8, 20, 1, 9, 0, est), 2},
		// Thursday -> Friday before Nov 2023 Veteran's Day
		{time.Date(2023, time.November, 9, 1, 0, 0, 0, est), time.Date(2023, time.November, 10, 1, 0, 0, 0, est), 1},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).AddBankingDay(test.Days)
		if !actual.Equal(NewTime(test.Future)) {
			t.Fatalf("Adding %d days: expected %s, got %s", test.Days, NewTime(test.Future), actual)
		}

		actual = NewTime(test.Date).AddBankingDay(test.Days)
		if !actual.Equal(NewTime(test.Future)) {
			t.Errorf("Adding %d days: expected %s, got %s", test.Days, NewTime(test.Future), actual)
		}
	}
}

func TestTime__Conversions(t *testing.T) {
	eastern, _ := time.LoadLocation("America/New_York")

	// create dates that are on a different day earlier than a holiday in different time zone
	pacific, _ := time.LoadLocation("America/Los_Angeles")
	when := NewTime(time.Date(2018, time.December, 24, 23, 0, 0, 0, pacific)).In(eastern)
	if when.Day() != 25 {
		t.Errorf("%v but expected to fall on Christmas", t)
	}

	// create dates that are on a different day later than a holiday in different time zone
	madrid, _ := time.LoadLocation("Europe/Madrid")
	when = NewTime(time.Date(2018, time.December, 26, 0, 30, 0, 0, madrid)).In(eastern)
	if when.Day() != 25 {
		t.Errorf("%v but expected to fall on Christmas", t)
	}
}

func TestTime__FridayHoliday(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected bool
	}{
		{time.Date(2026, time.December, 24, 9, 30, 0, 0, est), false},
		{time.Date(2026, time.December, 28, 9, 30, 0, 0, est), false},

		// Friday
		{time.Date(2026, time.December, 25, 9, 30, 0, 0, est), true},
	}
	for _, test := range tests {
		actual := NewTime(test.Date).IsHoliday()
		if actual != test.Expected {
			t.Fatalf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}

		actual = NewTime(test.Date).IsHoliday()
		if actual != test.Expected {
			t.Errorf("Date %s: expected %t, got %t", test.Date, test.Expected, actual)
		}
	}
}

func TestTime__SaturdayHoliday(t *testing.T) {
	est, _ := time.LoadLocation("America/New_York")
	cases := []time.Time{
		// The following are Saturday holidays where the Fed is open on Friday
		time.Date(2023, time.November, 10, 10, 0, 0, 0, est),
		time.Date(2026, time.July, 3, 10, 0, 0, 0, est),
		time.Date(2027, time.June, 18, 10, 0, 0, 0, est),
		time.Date(2027, time.December, 24, 10, 0, 0, 0, est),
	}
	for idx := range cases {
		actual := NewTime(cases[idx]).IsBankingDay()
		if !actual {
			t.Fatalf("expected %s to be a banking day", cases[idx])
		}

		actual = NewTime(cases[idx]).IsBankingDay()
		if !actual {
			t.Fatalf("expected %s to be a banking day", cases[idx])
		}
	}

}

func TestTime__SundayHoliday(t *testing.T) {
	// "if any holiday falls on a Sunday, the next following Monday is a standard
	// Federal Reserve Bank holiday. ... process the file on the first business
	// day after the original posting date."
	//
	// 2021-07-04 (July 4th) is on a Sunday
	eastern, _ := time.LoadLocation("America/New_York")
	ts := NewTime(time.Date(2021, time.July, 4, 10, 30, 0, 0, eastern)) // 10:30am

	if ts.IsBankingDay() || !ts.IsWeekend() {
		t.Errorf("%s it not a banking day", ts)
	}

	// move ahead 1 banking day
	ts = ts.AddBankingDay(1)
	if ts.Year() != 2021 || ts.Month() != time.July || ts.Day() != 6 {
		t.Errorf("unexpected banking day: %s", ts)
	}
	if wd := ts.Weekday(); wd != time.Tuesday {
		t.Errorf("expected Tuesday, got %s", wd)
	}

	// July 4 2027 (Sunday) is observed on Monday
	ts = NewTime(time.Date(2027, time.July, 2, 10, 0, 0, 0, eastern))
	require.True(t, ts.IsBankingDay())
	ts = ts.AddBankingDay(1)
	require.Equal(t, "2027-07-06T10:00:00-04:00", ts.Format(time.RFC3339))
}

func TestTime__GetHoliday(t *testing.T) {
	eastern, _ := time.LoadLocation("America/New_York")

	when := NewTime(time.Date(2022, time.December, 25, 10, 30, 0, 0, eastern))
	require.True(t, when.IsHoliday())

	holiday := when.GetHoliday()
	require.NotNil(t, holiday)
	require.Equal(t, "Christmas Day", holiday.Name)
}

func TestTime__YearlyHolidays(t *testing.T) {
	est, _ := time.LoadLocation("America/New_York")
	cases := []struct {
		when    time.Time
		holiday bool
	}{
		{time.Date(2023, time.December, 25, 12, 0, 0, 0, est), true},
		{time.Date(2023, time.December, 26, 12, 0, 0, 0, est), false},
		{time.Date(2024, time.January, 1, 12, 0, 0, 0, est), true},
		{time.Date(2024, time.January, 2, 12, 0, 0, 0, est), false},
	}
	for i := range cases {
		description := fmt.Sprintf("%s holiday=%v", cases[i].when.Format(time.RFC3339), cases[i].holiday)

		require.Equal(t, cases[i].holiday, NewTime(cases[i].when).IsHoliday(), description)
		require.NotEqual(t, cases[i].holiday, NewTime(cases[i].when).IsBankingDay(), description)
	}
}
