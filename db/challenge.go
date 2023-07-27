package db

// ChallengeType is a custom type for representing the challenge type as an enum.
type ChallengeType string

// String returns the string representation of the challenge type.
func (c ChallengeType) String() string {
	return string(c)
}
