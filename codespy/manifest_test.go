package analyze

import (
	"testing"
)

func TestUpToDateDependency(t *testing.T) {
	dep := Dependency{Package: "Flask", Operator: "==", Current: "0.9", Latest: "0.9"}

	if dep.Outdated() {
		t.Fatal("A dependency with the same versions shouldn't be outdated")
	}

}

func TestCurrentDependency(t *testing.T) {
	dep := Dependency{Package: "Flask", Operator: "==", Current: "0.9", Latest: "0.10.1"}

	if !dep.Outdated() {
		t.Fatal("A dependency with different versions should be outdated")
	}

}

func TestOutdatedNoLatestVersion(t *testing.T) {
	dep := Dependency{Package: "Flask", Operator: "==", Current: "0.9"}

	if dep.Outdated() {
		t.Fatal("A dependency without a Latest version shouldn't be outdated")
	}

}

// Run this test without contacting PyPI
func CheckPyPI(t *testing.T) {
	deps := []Dependency{
		Dependency{Package: "Flask", Operator: "==", Current: "0.9"},
		Dependency{Package: "Flask-Assets", Operator: "==", Current: "0.8"},
	}

	deps, _ = CheckPythonPackageIndex(deps)

	for _, dep := range deps {
		t.Log(dep.Latest)
		t.Log(dep.Outdated())
	}

	t.Fail()
}

func TestParseRequirements(t *testing.T) {
	dependencies, err := ParseRequirements("fixtures/requirements.txt")

	if err != nil {
		t.Fatal(err)
	}

	expected := []Dependency{
		Dependency{Package: "Flask", Operator: "==", Current: "0.9"},
		Dependency{Package: "Flask-Assets", Operator: "==", Current: "0.8"},
		Dependency{Package: "Flask-FlatPages", Operator: "==", Current: "0.3"},
		Dependency{Package: "Flask-Markdown", Operator: "==", Current: "0.3"},
		Dependency{Package: "Frozen-Flask", Operator: "==", Current: "0.9"},
		Dependency{Package: "Jinja2", Operator: "==", Current: "2.6"},
		Dependency{Package: "Markdown", Operator: "==", Current: "2.2.1"},
		Dependency{Package: "PyYAML", Operator: "==", Current: "3.10"},
		Dependency{Package: "Pygments", Operator: "==", Current: "1.5"},
		Dependency{Package: "Werkzeug", Operator: "==", Current: "0.8.3"},
		Dependency{Package: "argh", Operator: "==", Current: "0.17.2"},
		Dependency{Package: "argparse", Operator: "==", Current: "1.2.1"},
		Dependency{Package: "webassets", Operator: "==", Current: "0.8"},
		Dependency{Package: "wsgiref", Operator: "==", Current: "0.1.2"},
	}

	for i, dep := range dependencies {
		if dep != expected[i] {
			t.Fatalf("Package %s doesn't match %s", dep.Package, expected[i].Package)
		}
	}

}
