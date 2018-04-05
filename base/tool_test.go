package base

import "testing"

func TestGenerateRandomString(t *testing.T) {
	t.Log(string(GenerateRandomString(32, T_RAND_ALL)))
}
