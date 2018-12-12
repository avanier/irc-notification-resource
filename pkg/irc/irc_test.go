package irc_test

import (
	"bytes"
	"os"

	. "github.com/flavorjones/irc-notification-resource/pkg/irc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

// {"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}

var _ = Describe("Out", func() {
	Describe("ParseAndCheckRequest()", func() {
		It("returns correct Source values", func() {
			request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret", "usetls": true, "join": false}, "params": {"message": "foo"}}`))
			Expect(error).To(BeNil())
			Expect(request.Source).To(MatchAllFields(Fields{
				"Server":   Equal("chat.freenode.net"),
				"Port":     Equal(7070),
				"Channel":  Equal("#random"),
				"User":     Equal("randobot1337"),
				"Password": Equal("secretsecret"),
				"UseTLS":   Equal(true),
				"Join":     Equal(false),
			}))
		})

		Describe("required Source property", func() {
			Describe("`server`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
					Expect(error.Error()).To(MatchRegexp(`No server was provided`))
				})
			})

			Describe("`port`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
					Expect(error.Error()).To(MatchRegexp(`No port was provided`))
				})
			})

			Describe("`channel`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "user": "randobot1337", "password": "secretsecret"}}`))
					Expect(error.Error()).To(MatchRegexp(`No channel was provided`))
				})
			})

			Describe("`user`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "password": "secretsecret"}}`))
					Expect(error.Error()).To(MatchRegexp(`No user was provided`))
				})
			})

			Describe("`password`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337"}}`))
					Expect(error.Error()).To(MatchRegexp(`No password was provided`))
				})
			})
		})

		Describe("optional Source property", func() {
			Describe("usetls", func() {
				It("defaults to true", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeTrue()}))
				})

				It("is settable to true", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret", "usetls": true}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeTrue()}))
				})

				It("is settable to false", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret", "usetls": false}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeFalse()}))
				})
			})

			Describe("join", func() {
				It("defaults to false", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeFalse()}))
				})

				It("is settable to true", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret", "join": true}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeTrue()}))
				})

				It("is settable to false", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret", "join": false}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeFalse()}))
				})
			})
		})

		It("returns correct Params values", func() {
			request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID"}}`))
			Expect(error).To(BeNil())
			Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{
				"Message": Equal("foo $BUILD_ID"),
			}))
		})

		Describe("required Params property", func() {
			Describe("`message`", func() {
				It("errors if not present", func() {
					_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
					Expect(error.Error()).To(MatchRegexp(`No message was provided`))
				})
			})
		})

		Describe("optional Params property", func() {
			Describe("`dry_run`", func() {
				It("defaults to false", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID"}}`))
					Expect(error).To(BeNil())
					Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
				})

				It("is settable to true", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID", "dry_run": true}}`))
					Expect(error).To(BeNil())
					Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeTrue()}))
				})

				It("is settable to false", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID", "dry_run": false}}`))
					Expect(error).To(BeNil())
					Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
				})
			})
		})
	})

	Describe("ExpandMessage()", func() {
		var request Request

		BeforeEach(func() {
			request = Request{
				Source: Source{
					Server:   "chat.freenode.net",
					Port:     7070,
					Channel:  "#random",
					User:     "randobot1337",
					Password: "secretsecret",
					UseTLS:   true,
					Join:     false,
				},
				Params: Params{DryRun: true},
			}

			os.Setenv("BUILD_ID", "id-123")
			os.Setenv("BUILD_NAME", "name-asdf")
			os.Setenv("BUILD_JOB_NAME", "job-name-asdf")
			os.Setenv("BUILD_PIPELINE_NAME", "pipeline-name-asdf")
			os.Setenv("BUILD_TEAM_NAME", "team-name-asdf")
			os.Setenv("ATC_EXTERNAL_URL", "https://ci.example.com")
		})

		It("expands environment variables", func() {
			request.Params.Message = ">> $BUILD_ID <<"
			message := ExpandMessage(&request)
			Expect(message).To(Equal(">> id-123 <<"))
		})

		It("expands BUILD_URL pseudo-metadata", func() {
			request.Params.Message = ">> $BUILD_URL <<"
			message := ExpandMessage(&request)
			Expect(message).To(Equal(">> https://ci.example.com/teams/team-name-asdf/pipelines/pipeline-name-asdf/jobs/job-name-asdf/builds/name-asdf <<"))
		})
	})

	Describe("BuildResponse()", func() {
		var request Request
		var message string

		BeforeEach(func() {
			request = Request{
				Source: Source{
					Server:   "chat.freenode.net",
					Port:     7070,
					Channel:  "#random",
					User:     "randobot1337",
					Password: "secretsecret",
					UseTLS:   true,
					Join:     false,
				},
				Params: Params{DryRun: true},
			}

			os.Setenv("BUILD_ID", "id-123")
			os.Setenv("BUILD_NAME", "name-asdf")
			os.Setenv("BUILD_JOB_NAME", "job-name-asdf")
			os.Setenv("BUILD_PIPELINE_NAME", "pipeline-name-asdf")
			os.Setenv("BUILD_TEAM_NAME", "team-name-asdf")
			os.Setenv("ATC_EXTERNAL_URL", "https://ci.example.com")

			message = "this is a message"
		})

		Describe("returned Response", func() {
			It("contains version", func() {
				response := BuildResponse(&request, message)
				Expect(response.Version.Ref).To(Equal("none"))
			})

			It("contains specific metadata", func() {
				response := BuildResponse(&request, message)
				Expect(response.Metadata).To(Equal([]Metadatum{
					Metadatum{"host", "chat.freenode.net:7070"},
					Metadatum{"channel", "#random"},
					Metadatum{"user", "randobot1337"},
					Metadatum{"usetls", "true"},
					Metadatum{"join", "false"},
					Metadatum{"message", "this is a message"},
					Metadatum{"dry_run", "true"},
				}))
			})

			It("does not contains metadata `password`", func() {
				response := BuildResponse(&request, message)
				for _, metadatum := range response.Metadata {
					Expect(metadatum.Name).To(Not(MatchRegexp(`password`)))
				}
			})
		})
	})
})