package dice

import "testing"

func TestDiceError_Error(t *testing.T) {
	t.Run("validate error string", func(t *testing.T) {
		want := "heyo error"
		got := Error("heyo error")
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
