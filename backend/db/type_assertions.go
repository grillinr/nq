package db

import (
	"fmt"

	"github.com/google/uuid"
)

// Type assertion helper functions to safely extract values from Neo4j records
// These functions handle type assertions with proper error checking to avoid panics

// getStringFromMap safely extracts a string value from a map
func getStringFromMap(m map[string]any, key string) (string, error) {
	val, exists := m[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in map", key)
	}
	if val == nil {
		return "", nil
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value for key %s is not a string, got %T", key, val)
	}
	return str, nil
}

// getStringPtrFromMap safely extracts a *string value from a map
func getStringPtrFromMap(m map[string]any, key string) (*string, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	if val == nil {
		return nil, nil
	}
	str, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("value for key %s is not a string, got %T", key, val)
	}
	if str == "" {
		return nil, nil
	}
	return &str, nil
}

// getIntFromMap safely extracts an int64 value from a map
func getIntFromMap(m map[string]any, key string) (int64, error) {
	val, exists := m[key]
	if !exists {
		return 0, fmt.Errorf("key %s not found in map", key)
	}
	if val == nil {
		return 0, nil
	}
	num, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("value for key %s is not an int64, got %T", key, val)
	}
	return num, nil
}

// getIntPtrFromMap safely extracts a *int value from a map
func getIntPtrFromMap(m map[string]any, key string) (*int, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	if val == nil {
		return nil, nil
	}
	num, ok := val.(int64)
	if !ok {
		return nil, fmt.Errorf("value for key %s is not an int64, got %T", key, val)
	}
	intVal := int(num)
	return &intVal, nil
}

// getFloatFromMap safely extracts a float64 value from a map
func getFloatFromMap(m map[string]any, key string) (float64, error) {
	val, exists := m[key]
	if !exists {
		return 0, fmt.Errorf("key %s not found in map", key)
	}
	if val == nil {
		return 0, nil
	}
	num, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("value for key %s is not a float64, got %T", key, val)
	}
	return num, nil
}

// getFloatPtrFromMap safely extracts a *float64 value from a map
func getFloatPtrFromMap(m map[string]any, key string) (*float64, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	if val == nil {
		return nil, nil
	}
	num, ok := val.(float64)
	if !ok {
		return nil, fmt.Errorf("value for key %s is not a float64, got %T", key, val)
	}
	return &num, nil
}

// getBoolFromMap safely extracts a bool value from a map
func getBoolFromMap(m map[string]any, key string) (bool, error) {
	val, exists := m[key]
	if !exists {
		return false, fmt.Errorf("key %s not found in map", key)
	}
	if val == nil {
		return false, nil
	}
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value for key %s is not a bool, got %T", key, val)
	}
	return b, nil
}

// getUUIDFromMap safely extracts and parses a UUID value from a map
func getUUIDFromMap(m map[string]any, key string) (uuid.UUID, error) {
	str, err := getStringFromMap(m, key)
	if err != nil {
		return uuid.Nil, err
	}
	if str == "" {
		return uuid.Nil, fmt.Errorf("UUID string is empty for key %s", key)
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID from key %s: %w", key, err)
	}
	return id, nil
}

// getSliceFromMap safely extracts a slice value from a map
func getSliceFromMap(m map[string]any, key string) ([]any, error) {
	val, exists := m[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	if val == nil {
		return []any{}, nil
	}
	slice, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("value for key %s is not a slice, got %T", key, val)
	}
	return slice, nil
}
