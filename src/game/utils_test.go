package game

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultFuzzy(t *testing.T) {
	Convey("It panics when an argument is out of range", t, func() {
		So(func() { defaultFuzzy(-.01) }, ShouldPanic)
		So(func() { defaultFuzzy(1.01) }, ShouldPanic)
	})
	Convey("It returns the correct values", t, func() {
		So(defaultFuzzy(0), ShouldEqual, "--")
		So(defaultFuzzy(.09), ShouldEqual, "--")
		So(defaultFuzzy(.1), ShouldEqual, "-")
		So(defaultFuzzy(.29), ShouldEqual, "-")
		So(defaultFuzzy(.3), ShouldEqual, "o")
		So(defaultFuzzy(.69), ShouldEqual, "o")
		So(defaultFuzzy(.7), ShouldEqual, "+")
		So(defaultFuzzy(.89), ShouldEqual, "+")
		So(defaultFuzzy(.9), ShouldEqual, "++")
		So(defaultFuzzy(.99), ShouldEqual, "++")
	})
}
