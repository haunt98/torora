package comparison

import "testing"

func CompareError(t *testing.T, want, got error) {
	if want == nil && got == nil {
		return
	}

	if want != nil && got != nil {
		if got.Error() != want.Error() {
			t.Errorf("want err [%+v] got err [%+v]\n", want, got)
		}
		return
	}

	t.Errorf("want err [%+v] got err [%+v]\n", want, got)
}

func CompareInterface(t *testing.T, want, got interface{}, thing string) {
	if got != want {
		t.Errorf("want %s [%+v] got %s [%+v]", thing, want, thing, got)
	}
}
