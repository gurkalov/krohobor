package tests

import (
	"github.com/gurkalov/krohobor/tests/command"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestShowListAndReadSuccess(t *testing.T) {
	convey.Convey("Show list", t, func() {
		_, err := command.Run("db", "list")
//		convey.So(out, convey.ShouldEqual, "{[{test1 8053251} {test2 8204847} {test3 8204847}]}")
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("Read first", func() {
			out, err := command.Run("--db=test1", "db", "read")
			convey.So(out, convey.ShouldEqual, "{[{account1 24576 0} {link 32768 2}]}")
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Read second", func() {
			out, err := command.Run("--db=test2", "db", "read")
			convey.So(out, convey.ShouldEqual, "{[{account2 24576 0} {link 32768 2}]}")
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Read third", func() {
			out, err := command.Run("--db=test3", "db", "read")
			convey.So(out, convey.ShouldEqual, "{[{account3 24576 0} {link 32768 2}]}")
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDumpAllSuccess(t *testing.T) {
	convey.Convey("Dump all", t, func() {
		_, err := command.Run("db", "dumpall")
		//convey.So(out, convey.ShouldEqual, "{[{test1 8204847} {test2 8204847} {test3 8204847}]}")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDumpSuccess(t *testing.T) {
	convey.Convey("Dump database", t, func() {
		_, err := command.Run("--db=test1", "db", "dump")
		//convey.So(out, convey.ShouldEqual, "{[{test1 8204847} {test2 8204847} {test3 8204847}]}")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestBackupListSuccess(t *testing.T) {
	convey.Convey("Backup all", t, func() {
		_, err := command.Run("backup", "list")
		//convey.So(out, convey.ShouldEqual, "{[{test1 8204847} {test2 8204847} {test3 8204847}]}")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRestoreSuccess(t *testing.T) {
	convey.Convey("Restore all", t, func() {
		_, err := command.Run("--name=test1.sql", "db", "restore")
		//convey.So(out, convey.ShouldEqual, "{[{test1 8204847} {test2 8204847} {test3 8204847}]}")
		convey.So(err, convey.ShouldBeNil)
	})
}
