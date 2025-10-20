// Test file for metadata package
package metadata

import (
	"testing"
)

func TestIsValidMediaType(t *testing.T) {
	tests := []struct {
		mediaType string
		expected  bool
	}{
		{"book", true},
		{"movie", true},
		{"tv", true},
		{"game", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidMediaType(test.mediaType)
		if result != test.expected {
			t.Errorf("IsValidMediaType(%s) = %v, expected %v", test.mediaType, result, test.expected)
		}
	}
}

func TestGetMetadataByTitleValidation(t *testing.T) {
	// Test invalid media type
	_, err := GetMetadataByTitle("invalid", "Test Title", 2023)
	if err == nil {
		t.Error("Expected error for invalid media type, got nil")
	}

	// Test empty title
	_, err = GetMetadataByTitle("movie", "", 2023)
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestGetMetadataByIDValidation(t *testing.T) {
	// Test invalid media type
	_, err := GetMetadataByID("invalid", "12345")
	if err == nil {
		t.Error("Expected error for invalid media type, got nil")
	}

	// Test empty ID
	_, err = GetMetadataByID("movie", "")
	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}
}

func TestGetMetadataValidation(t *testing.T) {
	// Test with map containing various data types
	testCases := []struct {
		name      string
		mediaType string
		mediaInfo map[string]interface{}
		expectErr bool
	}{
		{
			name:      "Valid string title",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": "Test Movie",
			},
			expectErr: false,
		},
		{
			name:      "Invalid title type",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": 123,
			},
			expectErr: true,
		},
		{
			name:      "Valid integer year",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": "Test Movie",
				"year":  2023,
			},
			expectErr: false,
		},
		{
			name:      "Valid float year",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": "Test Movie",
				"year":  2023.0,
			},
			expectErr: false,
		},
		{
			name:      "Valid string year",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": "Test Movie",
				"year":  "2023",
			},
			expectErr: false,
		},
		{
			name:      "Invalid year type",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{
				"title": "Test Movie",
				"year":  []string{"2023"},
			},
			expectErr: true,
		},
		{
			name:      "No title or ID",
			mediaType: "movie",
			mediaInfo: map[string]interface{}{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: This will fail because we need API keys, but it should at least
			// pass validation and fail on service creation or API call
			_, err := GetMetadata(tc.mediaType, tc.mediaInfo)

			if tc.expectErr && err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			} else if !tc.expectErr && err != nil {
				// Only check if it's a validation error, not API-related errors
				if err.Error() == "title must be a string" ||
					err.Error() == "year must be an integer" ||
					err.Error() == "either title or ID is required" {
					t.Errorf("Unexpected validation error for %s: %v", tc.name, err)
				}
			}
		})
	}
}
