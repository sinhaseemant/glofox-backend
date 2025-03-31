package models

import (
	"encoding/json"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomDate time.Time

const customDateFormat = "2006-01-02"

// UnmarshalJSON allows CustomDate to be parsed from "YYYY-MM-DD"
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`) // Remove quotes from JSON string
	t, err := time.Parse(customDateFormat, str)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}

// MarshalJSON ensures CustomDate is serialized as "YYYY-MM-DD"
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(cd).Format(customDateFormat))
}

// ToTime converts CustomDate to a standard time.Time object
func (cd CustomDate) ToTime() time.Time {
	return time.Time(cd)
}

// isZero checks if the CustomDate is zero (i.e., not set)
func (cd CustomDate) IsZero() bool {
	return cd.ToTime().IsZero()
}

// Class represents a class with its details.
type Class struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB ObjectID
	Name      string             `bson:"name" json:"name"`
	StartDate CustomDate         `bson:"start_date" json:"start_date"`
	EndDate   CustomDate         `bson:"end_date" json:"end_date"`
	Capacity  int                `bson:"capacity" json:"capacity"`
}

// Booking represents a member's booking for a specific class on a specific date.
type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // MongoDB ObjectID
	ClassName  string             `bson:"class_name" json:"class_name"`   // Name of the class booked
	MemberName string             `bson:"member_name" json:"member_name"` // Name of the member
	Date       CustomDate         `bson:"date" json:"date"`               // Specific date of the booking
	ClassID    primitive.ObjectID `bson:"class_id" json:"class_id"`       // Reference to the class definition
}
