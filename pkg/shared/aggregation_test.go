package shared

import "testing"

func TestInsertNth(t *testing.T) {
	hash := map[string]string{
		"This is a test":        "This @is a @test",
		"With NL  and with Q":   "With @NL  a@nd wi@th Q",
		"sddfhsejsrtjrstjrtjdr": "sddfh@sejsr@tjrst@jrtjd@r",
	}

	for original, validRes := range hash {
		res := InsertNth(original, 5, "@")
		if validRes != res {
			t.Errorf("Result is not equal to valid result. Res: %s", res)
		}
	}
}
