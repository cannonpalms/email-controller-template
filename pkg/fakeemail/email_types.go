package fakeemail

import "hash/fnv"

type Email struct {
	DestinationAddress string
	Subject            string
	Body               string
}

// EmailIdentifier is a 32-bit unsigned integer that is produced by hashing the
// destination address, subject line, and message body of an email. A given set
// of inputs will always produce the same output. For the purposes of this assignment,
// ignore the infinitisimal chance of hash collisions.
type EmailIdentifier uint32

func (e Email) hash() uint32 {
	h := fnv.New32()

	// hash the fields individually and combine with XOR (implicit)
	h.Write([]byte(e.DestinationAddress))
	h.Write([]byte(e.Subject))
	h.Write([]byte(e.Body))

	return h.Sum32()
}

// ID() returns the EmailIdentifier for an Email
// EmailIdentifier is a 32-bit unsigned integer that is produced by hashing the
// destination address, subject line, and message body of an email. A given set
// of inputs will always produce the same output. For the purposes of this assignment,
// ignore the infinitisimal chance of hash collisions.
func (e Email) ID() EmailIdentifier {
	return EmailIdentifier(e.hash())
}
