package daemon

import (
	"testing"
)

func TestMinifyHTML(t *testing.T) {
	h1 := `<li>
    <a href=""></a>
          </li>`
	m1 := minifyHTML([]byte(h1))
	if m1 != "<li> <a href=\"\"></a> </li> " {
		t.Errorf("minifyHTML 1 wrong. Got: %s", m1)
	}

	h2 := `<li>

          </li>`
	m2 := minifyHTML([]byte(h2))
	if m2 != "<li> </li> " {
		t.Errorf("minifyHTML 2 wrong. Got: %s", m2)
	}
}
