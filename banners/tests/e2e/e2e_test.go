package e2e

import (
	"context"
	"fmt"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/headers"
	"github.com/ozontech/cute/asserts/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_Auth(t *testing.T) {
	cute.NewTestBuilder().
		Title("single auth test").
		Create().
		RequestRepeat(2).
		RequestBuilder(
			cute.WithURI("http://localhost:8080/login"),
			cute.WithMarshalBody(struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "33kjdkjj123kk2@al.ru",
				Password: "opopop111",
			}),
			cute.WithMethod(http.MethodPost),
		).
		ExpectStatus(http.StatusOK).
		AssertBody(json.Equal("message", "Successfully logged in.")).
		After(
			func(response *http.Response, errors []error) error {
				b, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				token, err := json.GetValueFromJSON(b, "token")
				if err != nil {
					return err
				}

				fmt.Println("Token from test", token)

				return nil
			},
		).
		ExecuteTest(context.Background(), t)
}

//func Test_hut(t *testing.T) {
//	cute.NewTestBuilder().Title("sdfkdskf").Create().RequestRepeat(2).RequestBuilder(
//		cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
//		cute.WithMarshalBody(struct {
//			Name string `json:"name"`
//		}{
//			Name: "vasya pupok",
//		}),
//		cute.WithQueryKV("socks", "42"),
//		cute.WithMethod(http.MethodGet),
//	).
//		ExpectExecuteTimeout(10*time.Second).
//		ExpectStatus(http.StatusOK).
//		AssertBody(json.Diff("{\"aaa\":\"bb\"}")).
//		AssertBody(
//			json.Present("$[1].name"),
//			json.Present("$[0].passport"),
//			json.Equal("$[0].email", "Eliseo@gardener.biz"),
//			examples.CustomAssertBody(),
//		).
//		ExecuteTest(context.Background(), t)
//}

//func Test_MultiSteps(t *testing.T) {
//	responseCode := 0
//	var result string
//
//	cute.NewTestBuilder().
//		Title("Test with login + post banner").
//		Tags("two_steps").
//		Create().
//		RequestBuilder(
//			cute.WithURI("http://localhost:8080/login"),
//			cute.WithMarshalBody(struct {
//				Email    string `json:"email"`
//				Password string `json:"password"`
//			}{
//				Email:    "33kjdkjj123kk2@al.ru",
//				Password: "opopop111",
//			}),
//			cute.WithMethod(http.MethodPost),
//		).
//		ExpectStatus(http.StatusOK).
//		RequireBody(json.Equal("message", "Successfully logged in.")).
//		After(
//			func(response *http.Response, errors []error) error {
//				responseCode = response.StatusCode
//
//				b, err := io.ReadAll(response.Body)
//				if err != nil {
//					return err
//				}
//
//				token, err := json.GetValueFromJSON(b, "token")
//				if err != nil {
//					return err
//				}
//
//				fmt.Println("Token from test", token)
//
//				stringSlice := make([]string, len(token))
//				for i, v := range token {
//					stringSlice[i] = fmt.Sprintf("%v", v)
//				}
//
//				result = strings.Join(stringSlice, "")
//
//				cute.NewTestBuilder().
//					Title("Test with user banner").
//					Tags("user_banner").
//					Create().
//					RequestBuilder(
//						cute.WithURI("http://localhost:8080/user_banner"),
//						cute.WithQuery(map[string][]string{
//							"tag_id":            []string{fmt.Sprint(4)},
//							"feature_id":        []string{fmt.Sprint(3)},
//							"use_last_revision": []string{fmt.Sprint(false)},
//						}),
//						cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//						cute.WithMethod(http.MethodGet),
//					).
//					ExpectStatus(http.StatusOK).
//					ExecuteTest(context.Background(), t)
//
//				return nil
//			},
//		).
//		ExpectStatus(http.StatusOK).
//		ExecuteTest(context.Background(), t)
//	fmt.Println("Response code from first request", responseCode)
//	fmt.Println("Result", result)
//}

func Test_Table_array_postbanner(t *testing.T) {
	tests := []*cute.Test{
		{
			Name: "Create banner",
			Middleware: &cute.Middleware{
				After: []cute.AfterExecute{
					func(response *http.Response, errors []error) error {
						b, err := io.ReadAll(response.Body)
						if err != nil {
							return err
						}

						token, err := json.GetValueFromJSON(b, "token")
						if err != nil {
							return err
						}

						fmt.Println("Token from test", token)

						stringSlice := make([]string, len(token))
						for i, v := range token {
							stringSlice[i] = fmt.Sprintf("%v", v)
						}

						result := strings.Join(stringSlice, "")

						cute.NewTestBuilder().
							Title("Test with user banner").
							Tags("user_banner").
							Create().
							RequestBuilder(
								cute.WithURI("http://localhost:8080/user_banner"),
								cute.WithQuery(map[string][]string{
									"tag_id":            []string{fmt.Sprint(4)},
									"feature_id":        []string{fmt.Sprint(3)},
									"use_last_revision": []string{fmt.Sprint(false)},
								}),
								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
								cute.WithMethod(http.MethodGet),
							).
							ExpectStatus(http.StatusOK).
							RequireBody(
								json.Present("content")).
							RequireHeaders(
								headers.Present("Content-Type")).
							ExecuteTest(context.Background(), t)

						return nil
					},
				},
			},
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://localhost:8080/login"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}{
						Email:    "33kjdkjj123kk2@al.ru",
						Password: "opopop111",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusOK,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully logged in."),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{},
	}

	cute.NewTestBuilder().
		Title("Table tests for banner creation").
		Tags("table_test_banner_create").
		Description("Execute array tests for banner creation").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), t)
}
