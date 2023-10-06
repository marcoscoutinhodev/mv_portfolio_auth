package __mock__

type IDGeneratorMock struct {
}

func (m *IDGeneratorMock) Generate() string {
	return "any_id"
}
