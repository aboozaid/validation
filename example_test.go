package validation_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/aboozaid/validation"
	"github.com/aboozaid/validation/is"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	Name    string
	Gender  string
	Email   string
	Address Address
}

func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		// Name cannot be empty, and the length must be between 5 and 20.
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		// Gender is optional, and should be either "Female" or "Male".
		validation.Field(&c.Gender, validation.In("Female", "Male")),
		// Email cannot be empty and should be in a valid email format.
		validation.Field(&c.Email, validation.Required, is.Email),
		// Validate Address using its own validation rules
		validation.Field(&c.Address),
	)
}

func Example() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			Street: "123 Main Street",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		},
	}

	err := c.Validate()
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format.); Email: must be a valid email address.
}

func Example_second() {
	data := "example"
	err := validation.Validate(data,
		validation.Required,       // not empty
		validation.Length(5, 100), // length between 5 and 100
		is.URL,                    // is a valid URL
	)
	fmt.Println(err)
	// Output:
	// must be a valid URL
}

func Example_third() {
	addresses := []Address{
		{State: "MD", Zip: "12345"},
		{Street: "123 Main St", City: "Vienna", State: "VA", Zip: "12345"},
		{City: "Unknown", State: "NC", Zip: "123"},
	}
	err := validation.Validate(addresses)
	fmt.Println(err)
	// Output:
	// 0: (City: cannot be blank; Street: cannot be blank.); 2: (Street: cannot be blank; Zip: must be in a valid format.).
}

func Example_four() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			State: "Virginia",
		},
	}

	err := validation.Errors{
		"name":  validation.Validate(c.Name, validation.Required, validation.Length(5, 20)),
		"email": validation.Validate(c.Name, validation.Required, is.Email),
		"zip":   validation.Validate(c.Address.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	}.Filter()
	fmt.Println(err)
	// Output:
	// email: must be a valid email address; zip: cannot be blank.
}

func Example_five() {
	type Employee struct {
		Name string
	}

	type Manager struct {
		Employee
		Level int
	}

	m := Manager{}
	err := validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Level, validation.Required),
	)
	fmt.Println(err)
	// Output:
	// Level: cannot be blank; Name: cannot be blank.
}

type contextKey int

func Example_six() {
	key := contextKey(1)
	rule := validation.WithContext(func(ctx context.Context, value interface{}) error {
		s, _ := value.(string)
		if ctx.Value(key) == s {
			return nil
		}
		return errors.New("unexpected value")
	})
	ctx := context.WithValue(context.Background(), key, "good sample")

	err1 := validation.ValidateWithContext(ctx, "bad sample", rule)
	fmt.Println(err1)

	err2 := validation.ValidateWithContext(ctx, "good sample", rule)
	fmt.Println(err2)

	// Output:
	// unexpected value
	// <nil>
}

func Example_seven() {
	c := map[string]interface{}{
		"Name":  "Qiang Xue",
		"Email": "q",
		"Address": map[string]interface{}{
			"Street": "123",
			"City":   "Unknown",
			"State":  "Virginia",
			"Zip":    "12345",
		},
	}

	err := validation.Validate(c,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("Name", validation.Required, validation.Length(5, 20)),
			// Email cannot be empty and should be in a valid email format.
			validation.Key("Email", validation.Required, is.Email),
			// Validate Address using its own validation rules
			validation.Key("Address", validation.Map(
				// Street cannot be empty, and the length must between 5 and 50
				validation.Key("Street", validation.Required, validation.Length(5, 50)),
				// City cannot be empty, and the length must between 5 and 50
				validation.Key("City", validation.Required, validation.Length(5, 50)),
				// State cannot be empty, and must be a string consisting of two letters in upper case
				validation.Key("State", validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
				// State cannot be empty, and must be a string consisting of five digits
				validation.Key("Zip", validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
			)),
		),
	)
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format; Street: the length must be between 5 and 50.); Email: must be a valid email address.
}
