package test

// MockUUIDGenerator is a mock implementation of UUIDGenerator.
type MockUUIDGenerator struct {
	ReturnUUID string
}

func (g *MockUUIDGenerator) Generate() (string, error) {
	return g.ReturnUUID, nil
}
