package ical

import (
	"strings"
	"testing"
	"time"
)

func TestStartAndEndAtUTC(t *testing.T) {
	event := CalendarEvent{}

	if event.StartAtUTC() != nil {
		t.Error("StartAtUTC should have been nil")
	}
	if event.EndAtUTC() != nil {
		t.Error("EndAtUTC should have been nil")
	}

	tUTC := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	event.StartAt = &tUTC
	event.EndAt = &tUTC
	startTime := *(event.StartAtUTC())
	endTime := *(event.EndAtUTC())

	if startTime != tUTC {
		t.Error("StartAtUTC should have been", tUTC, ", but was", startTime)
	}
	if endTime != tUTC {
		t.Error("EndAtUTC should have been", tUTC, ", but was", endTime)
	}

	tUTC = time.Date(2010, time.March, 8, 2, 0, 0, 0, time.UTC)
	nyk, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	tNYK := tUTC.In(nyk)
	event.StartAt = &tNYK
	event.EndAt = &tNYK
	startTime = *(event.StartAtUTC())
	endTime = *(event.EndAtUTC())

	if startTime != tUTC {
		t.Error("StartAtUTC should have been", tUTC, ", but was", startTime)
	}
	if endTime != tUTC {
		t.Error("EndAtUTC should have been", tUTC, ", but was", endTime)
	}
}

func TestCalendarEventSerialize(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	createdAt := time.Date(2010, time.January, 1, 12, 0, 1, 0, time.UTC)
	modifiedAt := createdAt.Add(time.Second)
	startsAt := createdAt.Add(time.Second * 2).In(ny)
	endsAt := createdAt.Add(time.Second * 3).In(ny)

	event := CalendarEvent{
		Id:            "123",
		CreatedAtUTC:  &createdAt,
		ModifiedAtUTC: &modifiedAt,
		StartAt:       &startsAt,
		EndAt:         &endsAt,
		Summary:       "Foo Bar",
		Location:      "Berlin\nGermany",
		Description:   "Lorem\nIpsum",
		URL:           "https://www.example.com",
	}

	// expects that DTSTART and DTEND be in UTC (Z)
	// expects that string values (LOCATION for example) be escaped
	expected := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem\nIpsum
LOCATION:Berlin\nGermany
URL:https://www.example.com
END:VEVENT`

	output := event.Serialize()
	if output != strings.TrimSpace(expected) {
		t.Error("Expected calendar event serialization to be:\n", expected, "\n\nbut got:\n", output)
	}
}

func TestCalendarEventParse(t *testing.T) {
	vevent := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem\nIpsum
LOCATION:Berlin\nGermany
URL:https://www.example.com
END:VEVENT`

	output, err := ParseCalendar(vevent)
	if err != nil {
		panic(err)
	}

	node := output.ChildByName("SUMMARY")
	if node == nil {
		t.Error("Expected SUMMARY to be not nil")
	}
	if node.Value != "Foo Bar" {
		t.Error("Expected SUMMARY to be: ", "Foo Bar", "Got: ", node.Value)
	}
}

func TestCalendarEventParseMultiline(t *testing.T) {
	vevent := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem Ipsum 
 Loquitur
LOCATION:Berlin\nGermany
URL:https://www.example.com
END:VEVENT`

	output, err := ParseCalendar(vevent)
	if err != nil {
		panic(err)
	}

	node := output.ChildByName("DESCRIPTION")
	if node == nil {
		t.Error("Expected DESCRIPTION to be not nil")
	}
	if node.Value != "Lorem Ipsum Loquitur" {
		t.Error("Expected DESCRIPTION to be: ", "Lorem Ipsum Loquitur", "Got: ", node.Value)
	}
}

func TestCalendarEventParseMultiEntry(t *testing.T) {
	vevent := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem Ipsum 
 Loquitur
ATTENDEE;RSVP=TRUE;PARTSTAT=NEEDS-ACTION;ROLE=REQ-PARTICIPANT:mailto:danmo
 lik@gmail.com
ATTENDEE;RSVP=TRUE;PARTSTAT=NEEDS-ACTION;ROLE=REQ-PARTICIPANT:mailto:dan@d
 3fy.net
LOCATION:Berlin\nGermany
URL:https://www.example.com
END:VEVENT`

	output, err := ParseCalendar(vevent)
	if err != nil {
		panic(err)
	}

	node := output.ChildByName("DESCRIPTION")
	if node == nil {
		t.Error("Expected DESCRIPTION to be not nil")
	}
	if node.Value != "Lorem Ipsum Loquitur" {
		t.Error("Expected DESCRIPTION to be: ", "Lorem Ipsum Loquitur", "Got: ", node.Value)
	}

	nodes := output.ChildrenByName("ATTENDEE")
	if nodes == nil {
		t.Error("Expected ATTENDEE to be not nil")
	}
	if nodes[0].Value != "mailto:danmolik@gmail.com" {
		t.Error("Expected ATTENDEE to be: ", "mailto:danmolik@gmail.com", "Got: ", nodes[0].Value)
	}
	if nodes[1].Value != "mailto:dan@d3fy.net" {
		t.Error("Expected ATTENDEE to be: ", "mailto:dan@d3fy.net", "Got: ", nodes[1].Value)
	}
}

func TestCalendarEventParseMultiDig(t *testing.T) {
	vevent := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem Ipsum 
 Loquitur
ATTENDEE;RSVP=TRUE;PARTSTAT=NEEDS-ACTION;ROLE=REQ-PARTICIPANT:mailto:danmo
 lik@gmail.com
ATTENDEE;RSVP=TRUE;PARTSTAT=NEEDS-ACTION;ROLE=REQ-PARTICIPANT:mailto:dan@d
 3fy.net
LOCATION:Berlin\nGermany
URL:https://www.example.com
END:VEVENT`

	output, err := ParseCalendar(vevent)
	if err != nil {
		panic(err)
	}
	dig, b := output.DigProperty("ATTENDEE")
	if  dig == "" || ! b {
		t.Error("Expected dig ATTENDEE to be not empty")
	}
	digs, b := output.DigProperties("ATTENDEE")
	if len(digs) == 0 || ! b {
		t.Error("Expected digs ATTENDEE to be not empty")
	}
	if digs[0] != "mailto:danmolik@gmail.com" {
		t.Error("Expected ATTENDEE len ", len(digs), " to be: ", "mailto:danmolik@gmail.com", "Got: ", digs[0])
		for _, d := range digs {
			t.Error("Expected ATTENDEE to be: ", "mailto:danmolik@gmail.com", "Got: ", d)
		}
	}
	if digs[1] != "mailto:dan@d3fy.net" {
		t.Error("Expected ATTENDEE to be: ", "mailto:dan@d3fy.net", "Got: ", digs[1])
	}
}
