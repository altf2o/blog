package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Our simple struct we use to test. By adding the various validate
// tag values, it lets us validate incoming or outgoing JSON data.
type Person struct {
	FirstName string `json:"first-name" validate:"required,alpha,lte=15"`
	LastName  string `json:"last-name" validate:"required,alpha,lte=15"`
	Age       int    `json:"age,omitempty" validate:"omitempty,number,gt=0,lt=120"`
}

// If we implement the fmt.Stringer interface, then when we
// pass an instance of our object to a "fmt.Print..."
// function, it'll automatically call it.
func (p Person) String() string {
	return fmt.Sprintf("First: %s, Last: %s, Age: %d",
		p.FirstName, p.LastName, p.Age)
}

// Just an easy way to see what JSON response we expect
const (
	ALL_INVALID int = iota
	FIRST_NAME_INVALID
	LAST_NAME_INVALID
	AGE_INVALID
	GOOD_RESPONSE
)

// This "simulates" a response JSON we can later unmarshal
func GetJSONResponse(resp int) ([]byte, error) {
	switch resp {
	case ALL_INVALID:
		// All fields invalid
		return []byte(`{"first-name":"Frank01","last-name":"^&Morales*", "age":121}`), nil
	case FIRST_NAME_INVALID:
		// FirstName too long
		return []byte(`{"first-name":"Longerthanfifteencharacters","last-name":"Morales","age":35}`), nil
	case LAST_NAME_INVALID:
		// LastName too long
		return []byte(`{"first-name":"Frank","last-name":"Longerthanfifteencharacters","age":35}`), nil
	case AGE_INVALID:
		// Age out of range
		return []byte(`{"first-name":"Frank","last-name":"Morales","age":130}`), nil
	default:
		// Default or GOOD_RESPONSE is a valid response
		return []byte(`{"first-name":"Frank","last-name":"Morales","age":35}`), nil
	}
}

// Used to simulate taking Marshaled/encoded JSON and sending
// it to an API.
// NOTE: The output shows why we'd use regular .Marshal() and not
// .MarshalIndent() as the latter would transmit more text/bytes
func SendJSONToAPI(jsonInput []byte) error {
	fmt.Printf("JSON we'd send: %s\n", string(jsonInput))
	return nil
}

func main() {
	// Get a new validator object for use later
	validate := validator.New()

	// This loops over our JSON responses and attempts to process each.
	for i := 0; i < 5; i++ {
		// "Call an API" and get a JSON response
		resp, err := GetJSONResponse(i)
		if err != nil {
			panic("should not get here")
		}

		// Unmarshal/decode the JSON response into our Person object
		// Notice that as long as the types match, it will unmarshal
		// ok even if the values are outside our expected ranges
		var decoded Person
		err = json.Unmarshal(resp, &decoded)
		if err != nil {
			fmt.Printf("error unmarshalling JSON response: %s\n", err.Error())
		} else {
			// Now we want to ensure that what we received is valid based on
			// our validate struct tags
			if err := validate.Struct(decoded); err != nil {
				fmt.Printf("error validating struct: %v\n", err.Error())
			} else {
				fmt.Println(decoded) // Will call our .String() method
				fmt.Println("struct fields validated ok!")
			}
		}
	}
	// Just a separator
	fmt.Println("-----\n")

	// Now lets simulate sending our data
	// If we have input data stored in a Person object
	// we can Marshal/encode it to be sent to an API
	p := Person{
		FirstName: "Sarah",
		LastName:  "Smith",
		Age:       28,
	}

	// We attempt to Marshal/encode our Person object into valid JSON
	encoded, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("error marshaling JSON: %s\n", err.Error())
		return
	}

	// We can then send it to a function that will send it to an API
	// for actual usage/storage
	err = SendJSONToAPI(encoded)
	if err != nil {
		fmt.Printf("error sending JSON to API: %s\n", err.Error())
		return
	}

	// What if we have no age? What does our JSON look like?
	newP := Person{
		FirstName: "Josh",
		LastName:  "Smith",
	}

	// Marshal our object into valid JSON
	encoded, err = json.Marshal(newP)
	if err != nil {
		fmt.Printf("error marshaling JSON: %s\n", err.Error())
		return
	}

	// Send our JSON to our API/endpoint
	err = SendJSONToAPI(encoded)
	if err != nil {
		fmt.Printf("error sending JSON to API: %s\n", err.Error())
	}

	// We can ensure struct values are validated before
	// being sent out as well, so they're within expected
	// API ranges/limits
	wontValidate := Person{
		FirstName: "Frank",
		LastName:  "Morales",
		Age:       130,
	}

	// It will Marshal ok
	encoded, err = json.Marshal(wontValidate)
	if err != nil {
		fmt.Printf("error marshaling JSON: %s\n", err.Error())
		return
	}

	// But validation will fail
	if err := validate.Struct(wontValidate); err != nil {
		fmt.Printf("error validating struct: %s\n", err.Error())
		return
	} else {
		// Won't get called since our struct fails validation
		_ = SendJSONToAPI(encoded)
	}
}
