package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type UnixTimestamp = [4]byte

// 83 bytes
type Person struct {
	Age        uint16    // 2 byte
	Name       string    // 64 byte
	NetWorth   uint64    // 8 bytes
	Birthday   time.Time // 8 bytes
	HasLicense bool      // 1 byte
}

const AgeOffset = 0
const NameOffset = 2
const NetWorthOffset = 66
const BirthdayOffset = 74
const HasLicenseOffset = 82

type PersonByteArray = [83]byte

func main() {

	// f, _ := os.Open("blkahdl.txt")
	// f.Seek()

	bday, err := time.Parse(time.RFC3339, "2003-06-20T15:04:05Z")
	if err != nil {
		panic(err)
	}

	person := Person{
		Age:        32,
		Name:       "Raamiz Abbasi",
		NetWorth:   314054,
		Birthday:   bday,
		HasLicense: true,
	}

	b, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("person.json", b, 0766); err != nil {
		panic(err)
	}

	personByteArr := personToBytes(person)

	if err := os.WriteFile("person.bytes", personByteArr[:], 0766); err != nil {
		panic(err)
	}

	b, err = ioutil.ReadFile("person.bytes")
	if err != nil {
		panic(err)
	}
	personReadFromBytes := bytesToPerson(b)
	fmt.Printf("%+v\n", personReadFromBytes)

}

func personToBytes(person Person) PersonByteArray {
	personInBytes := PersonByteArray{}

	ageBytes := make([]byte, 4)
	binary.LittleEndian.PutUint16(ageBytes, person.Age)
	for i, b := range ageBytes {
		personInBytes[i+AgeOffset] = b
	}

	nameBytes := []byte("Raamiz Abbasi")
	for i, b := range nameBytes {
		personInBytes[i+NameOffset] = b
	}

	netWorthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(netWorthBytes, uint64(person.NetWorth))
	for i, b := range netWorthBytes {
		personInBytes[i+NetWorthOffset] = b
	}

	birthdayBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(birthdayBytes, uint64(person.Birthday.Unix()))

	for i, b := range birthdayBytes {
		personInBytes[i+BirthdayOffset] = b
	}

	if person.HasLicense {
		personInBytes[HasLicenseOffset] = 01
	} else {
		personInBytes[HasLicenseOffset] = 00
	}

	return personInBytes
}

func bytesToPerson(b []byte) Person {
	age := binary.LittleEndian.Uint16(b[:NameOffset])
	name := b[NameOffset:NetWorthOffset]
	netWorth := binary.LittleEndian.Uint64(b[NetWorthOffset:BirthdayOffset])
	bday := binary.LittleEndian.Uint64(b[BirthdayOffset:HasLicenseOffset])
	hasLicenseBit := b[HasLicenseOffset:]

	bdayTime := time.Unix(int64(bday), 0)

	hasLicense := false
	if hasLicenseBit[0] == 01 {
		hasLicense = true
	}
	return Person{
		Age:        age,
		Name:       string(name),
		NetWorth:   netWorth,
		Birthday:   bdayTime,
		HasLicense: hasLicense,
	}

}
