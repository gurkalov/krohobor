package tests

import (
	"github.com/gurkalov/krohobor/tests/command"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := command.TearDown(); err != nil {
		panic(err)
	}
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestFullStorySuccess(t *testing.T) {
	convey.Convey("Show list", t, func() {
		out, err := command.Run("db", "list")
		convey.So(err, convey.ShouldBeNil)

		response, err := command.ToStrings(out)
		convey.So(err, convey.ShouldBeNil)
		convey.So(response, convey.ShouldResemble, [][]string{
			{"test1", "8053251"},
			{"test2", "8053251"},
			{"test3", "8053251"},
		})
	})

	convey.Convey("Read first", t, func() {
		out, err := command.Run("--dbname=test1", "db", "read")
		convey.So(err, convey.ShouldBeNil)

		response, err := command.ToStrings(out)
		convey.So(err, convey.ShouldBeNil)
		convey.So(response, convey.ShouldResemble, [][]string{
			{"account1", "0", "24576"},
			{"link", "2", "32768"},
		})
	})

	convey.Convey("Read second", t, func() {
		out, err := command.Run("--dbname=test2", "db", "read")
		convey.So(err, convey.ShouldBeNil)

		response, err := command.ToStrings(out)
		convey.So(err, convey.ShouldBeNil)
		convey.So(response, convey.ShouldResemble, [][]string{
			{"account2", "0", "24576"},
			{"link", "2", "32768"},
		})
	})

	convey.Convey("Read third", t, func() {
		out, err := command.Run("--dbname=test3", "db", "read")
		convey.So(err, convey.ShouldBeNil)

		response, err := command.ToStrings(out)
		convey.So(err, convey.ShouldBeNil)
		convey.So(response, convey.ShouldResemble, [][]string{
			{"account3", "0", "24576"},
			{"link", "2", "32768"},
		})
	})

	convey.Convey("Dump all", t, func() {
		out, err := command.Run( "dump", "create")
		convey.So(err, convey.ShouldBeNil)
		convey.So(out, convey.ShouldEqual, "")
	})

	convey.Convey("Dump database", t, func() {
		out, err := command.Run("--dbname=test1", "dump", "create")
		convey.So(err, convey.ShouldBeNil)
		convey.So(out, convey.ShouldEqual, "")
	})

	convey.Convey("Dump list", t, func() {
		out, err := command.Run("dump", "list")
		convey.So(err, convey.ShouldBeNil)

		response, err := command.ToStrings(out)
		convey.So(err, convey.ShouldBeNil)
		convey.So(response, convey.ShouldResemble, [][]string{
			{response[0][0]},
			{response[1][0]},
		})

		convey.Convey("Restore all", func() {
			out, err := command.Run("--name=" + response[0][0], "--database=postgres-target", "dump", "restore")
			convey.So(out, convey.ShouldEqual, "")
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Show target list", func() {
				out, err := command.Run("--database=postgres-target", "db", "list")
				convey.So(err, convey.ShouldBeNil)

				responseList, err := command.ToStrings(out)
				convey.So(err, convey.ShouldBeNil)
				convey.So(responseList, convey.ShouldResemble, [][]string{
					{"test1", "8204847"},
					{"test2", "8204847"},
					{"test3", "8204847"},
				})

				convey.Convey("Create new db", func() {
					out, err := command.Run("--database=postgres-target", "--dbname=test1_new", "db", "create")
					convey.So(out, convey.ShouldEqual, "")
					convey.So(err, convey.ShouldBeNil)

					convey.Convey("Restore one", func() {
						out, err := command.Run("--name="+response[1][0], "--database=postgres-target", "--dbname=test1_new", "dump", "restore")
						convey.So(out, convey.ShouldEqual, "")
						convey.So(err, convey.ShouldBeNil)

						convey.Convey("Read db one", func() {
							out, err := command.Run("--dbname=test1_new", "--database=postgres-target", "db", "read")
							convey.So(err, convey.ShouldBeNil)

							response, err := command.ToStrings(out)
							convey.So(err, convey.ShouldBeNil)
							convey.So(response, convey.ShouldResemble, [][]string{
								{"account1", "0", "24576"},
								{"link", "2", "32768"},
							})

							convey.Convey("Delete db one", func() {
								out, err := command.Run("--dbname=test2", "--database=postgres-target", "db", "delete")
								convey.So(out, convey.ShouldEqual, "")
								convey.So(err, convey.ShouldBeNil)

								convey.Convey("Show target list", func() {
									out, err := command.Run("--database=postgres-target", "db", "list")
									convey.So(err, convey.ShouldBeNil)

									responseList, err := command.ToStrings(out)
									convey.So(err, convey.ShouldBeNil)
									convey.So(responseList, convey.ShouldResemble, [][]string{
										{"test1", "8204847"},
										{"test3", "8204847"},
										{"test1_new", "8204847"},
									})
								})
							})
						})
					})
				})
			})
		})
	})
}
