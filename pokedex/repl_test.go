package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input: "    hello   world",
			expected: []string{"hello", "world"},
		},
		{
			input: "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "\t\t\thello test string\t\t  ",
			expected: []string{"hello", "test", "string"},
		},
		{
			input: "hell to the o the beat",
			expected: []string{"hell", "to", "the", "o", "the", "beat"},
		},
	}

	for _, testCase := range testCases {
		actual := cleanInput(testCase.input)
		expected := testCase.expected
		if len(actual) != len(expected) {
			t.Errorf("Actual slice length does not match expected slice length")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := testCase.expected[i]

			if word != expectedWord {
				t.Errorf("Words don't match! \nActual: %s\nExpected: %s", word, expectedWord)
			}
		}
	}

	fmt.Println("All tests passed successfully!")

}
