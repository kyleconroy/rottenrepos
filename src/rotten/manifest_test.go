package rotten

import (
	"testing"
)

func TestParseRequirements(t *testing.T) {
	dependencies, err := ParseRequirements("fixtures/requirements.txt")

	if err != nil {
		t.Fatal(err)
	}

	expected := []Dependency{
		Dependency{Package: "Flask", Operator: "==", Version: "0.9"},
		Dependency{Package: "Flask-Assets", Operator: "==", Version: "0.8"},
		Dependency{Package: "Flask-FlatPages", Operator: "==", Version: "0.3"},
		Dependency{Package: "Flask-Markdown", Operator: "==", Version: "0.3"},
		Dependency{Package: "Frozen-Flask", Operator: "==", Version: "0.9"},
		Dependency{Package: "Jinja2", Operator: "==", Version: "2.6"},
		Dependency{Package: "Markdown", Operator: "==", Version: "2.2.1"},
		Dependency{Package: "PyYAML", Operator: "==", Version: "3.10"},
		Dependency{Package: "Pygments", Operator: "==", Version: "1.5"},
		Dependency{Package: "Werkzeug", Operator: "==", Version: "0.8.3"},
		Dependency{Package: "argh", Operator: "==", Version: "0.17.2"},
		Dependency{Package: "argparse", Operator: "==", Version: "1.2.1"},
		Dependency{Package: "webassets", Operator: "==", Version: "0.8"},
		Dependency{Package: "wsgiref", Operator: "==", Version: "0.1.2"},
	}

	for i, dep := range dependencies {
		if dep != expected[i] {
			t.Fatalf("Package %s doesn't match %s", dep.Package, expected[i].Package)
		}
	}

}
