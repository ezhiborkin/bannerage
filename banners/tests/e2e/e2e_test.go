package e2e

import (
	"context"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/headers"
	"github.com/ozontech/cute/asserts/json"
	"net/http"
	"testing"
)

func Test_CreateUser(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "register user",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "test1@tes22t.com",
						Password: "opopop111",
						Role:     "user",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusCreated,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully created user."),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "duplicate user",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "test1@tes22t.com",
						Password: "opopop111",
						Role:     "user",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusConflict,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "failed to create user"),
				},
			},
		},
		{
			Name:       "register admin",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "admin1@admin.com",
						Password: "opopop111",
						Role:     "admin",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusCreated,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully created user."),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "no email",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						//Email:    "admin111t@admin.com",
						Password: "opopop111",
						Role:     "admin",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "email is empty"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "no password",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email: "admin111t@admin.com",
						//Password: "opopop111",
						Role: "admin",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "password is empty"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "no role",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "admin111t@admin.com",
						Password: "opopop111",
						//Role:     "admin",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "role is empty"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "empty body",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "failed to decode request"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
	}
	cute.NewTestBuilder().
		Title("Table tests for user creation").
		Tags("table_test_user_create").
		Description("Execute array tests for user creation").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), t)
}

func Test_LoginUser(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "register user",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "test2@tes22t.com",
						Password: "opopop111",
						Role:     "user",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusCreated,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully created user."),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "register admin",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
						Role     string `json:"role"`
					}{
						Email:    "admin2@admin.com",
						Password: "opopop111",
						Role:     "admin",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusCreated,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully created user."),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
				},
			},
		},
		{
			Name:       "login ok",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/login"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}{
						Email:    "test2@tes22t.com",
						Password: "opopop111",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusOK,
				AssertBody: []cute.AssertBody{
					json.Equal("message", "Successfully logged in."),
					json.Present("token"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
					headers.Present("Authorization"),
				},
			},
		},
		{
			Name:       "no user",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/login"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}{
						Email:    "test1244435@tes22t.com",
						Password: "opopop111",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "failed to login user"),
					json.NotPresent("token"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
					headers.NotPresent("Authorization"),
				},
			},
		},
		{
			Name:       "no email",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/login"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}{
						//Email:    "test1244435@tes22t.com",
						Password: "opopop111",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "email or password is empty"),
					json.NotPresent("token"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
					headers.NotPresent("Authorization"),
				},
			},
		},
		{
			Name:       "no password",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/login"),
					cute.WithMarshalBody(struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}{
						Email: "test1244435@tes22t.com",
						//Password: "opopop111",
					}),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "email or password is empty"),
					json.NotPresent("token"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
					headers.NotPresent("Authorization"),
				},
			},
		},
		{
			Name:       "empty body",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("http://bannerage-e2e:8080/create/user"),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: http.StatusBadRequest,
				AssertBody: []cute.AssertBody{
					json.Equal("error", "failed to decode request"),
					json.NotPresent("token"),
				},
				AssertHeaders: []cute.AssertHeaders{
					headers.Present("Content-Type"),
					headers.NotPresent("Authorization"),
				},
			},
		},
	}
	cute.NewTestBuilder().
		Title("Table tests for user login").
		Tags("table_test_user_login").
		Description("Execute array tests for user login").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), t)
}

//func Test_GetUser(t *testing.T) {
//	tests := []*cute.Test{
//		{
//			Name:       "not authenticated",
//			Middleware: nil,
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//					cute.WithQuery(map[string][]string{
//						"tag_id":            []string{fmt.Sprint(4)},
//						"feature_id":        []string{fmt.Sprint(3)},
//						"use_last_revision": []string{fmt.Sprint(false)},
//					}),
//					cute.WithMethod(http.MethodGet),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusUnauthorized,
//				AssertBody: []cute.AssertBody{
//					json.Present("error"),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name:       "register user",
//			Middleware: nil,
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/create/user"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//						Role     string `json:"role"`
//					}{
//						Email:    "test3@tes22t.com",
//						Password: "opopop111",
//						Role:     "user",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusCreated,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully created user."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name:       "register admin",
//			Middleware: nil,
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/create/user"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//						Role     string `json:"role"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//						Role:     "admin",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusCreated,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully created user."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "no tag",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									//"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(3)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusBadRequest).
//							AssertBody(
//								json.Equal("error", "tagID or featureID is not provided")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "no feature",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id": []string{fmt.Sprint(4)},
//									//"feature_id":        []string{fmt.Sprint(3)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusBadRequest).
//							AssertBody(
//								json.Equal("error", "tagID or featureID is not provided")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "empty query",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusBadRequest).
//							AssertBody(
//								json.Equal("error", "tagID or featureID is not provided")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "banner does not exist",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(2)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusNotFound).
//							AssertBody(
//								json.Equal("error", "failed to get banner")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "create banner",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/banner"),
//								cute.WithMarshalBody(struct {
//									FeatureId int64            `json:"feature_id"`
//									TagIDs    []int64          `json:"tag_ids"`
//									Content   json2.RawMessage `json:"content"`
//									IsActive  bool             `json:"is_active"`
//								}{
//									FeatureId: 3,
//									TagIDs:    []int64{1, 2, 3, 4},
//									Content:   []byte{123, 125},
//									//"{\"title\": \"memes\"}"
//									IsActive: true,
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodPost),
//							).
//							ExpectStatus(http.StatusOK).
//							AssertBody(
//								json.Equal("message", "Successfully created banner."),
//								json.Present("banner_id")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "get banner ok",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(3)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusOK).
//							AssertBody(
//								json.Present("content"),
//								json.NotPresent("error")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "create inactive banner",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/banner"),
//								cute.WithMarshalBody(struct {
//									FeatureId int64            `json:"feature_id"`
//									TagIDs    []int64          `json:"tag_ids"`
//									Content   json2.RawMessage `json:"content"`
//									IsActive  bool             `json:"is_active"`
//								}{
//									FeatureId: 4,
//									TagIDs:    []int64{1, 2, 3, 4},
//									Content:   []byte{123, 125},
//									IsActive:  false,
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodPost),
//							).
//							ExpectStatus(http.StatusOK).
//							AssertBody(
//								json.Equal("message", "Successfully created banner."),
//								json.Present("banner_id")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "get banner inactive admin",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(4)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusOK).
//							AssertBody(
//								json.Present("content"),
//								json.NotPresent("error")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "admin3@admin.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "get banner inactive user",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(4)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusUnauthorized).
//							AssertBody(
//								json.NotPresent("content"),
//								json.Equal("error", "you are not admin")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "test3@tes22t.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{
//			Name: "get cached banner",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(3)},
//									"use_last_revision": []string{fmt.Sprint(true)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithHeadersKV("Content-Type", "application/json"),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusOK).
//							AssertBody(
//								json.Present("content"),
//								json.NotPresent("error")).
//							AssertHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "test3@tes22t.com",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//	}
//	cute.NewTestBuilder().
//		Title("Table tests for user login").
//		Tags("table_test_user_login").
//		Description("Execute array tests for user login").
//		CreateTableTest().
//		PutTests(tests...).
//		ExecuteTest(context.Background(), t)
//}

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
//			cute.WithURI("http://bannerage-e2e:8080/login"),
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
//						cute.WithURI("http://bannerage-e2e:8080/user_banner"),
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

//func Test_Table_array_postbanner(t *testing.T) {
//	tests := []*cute.Test{
//		{
//			Name: "Create banner",
//			Middleware: &cute.Middleware{
//				After: []cute.AfterExecute{
//					func(response *http.Response, errors []error) error {
//						b, err := io.ReadAll(response.Body)
//						if err != nil {
//							return err
//						}
//
//						token, err := json.GetValueFromJSON(b, "token")
//						if err != nil {
//							return err
//						}
//
//						fmt.Println("Token from test", token)
//
//						stringSlice := make([]string, len(token))
//						for i, v := range token {
//							stringSlice[i] = fmt.Sprintf("%v", v)
//						}
//
//						result := strings.Join(stringSlice, "")
//
//						cute.NewTestBuilder().
//							Title("Test with user banner").
//							Tags("user_banner").
//							Create().
//							RequestBuilder(
//								cute.WithURI("http://bannerage-e2e:8080/user_banner"),
//								cute.WithQuery(map[string][]string{
//									"tag_id":            []string{fmt.Sprint(4)},
//									"feature_id":        []string{fmt.Sprint(3)},
//									"use_last_revision": []string{fmt.Sprint(false)},
//								}),
//								cute.WithHeadersKV("Authorization", fmt.Sprintf("Bearer %s", result)),
//								cute.WithMethod(http.MethodGet),
//							).
//							ExpectStatus(http.StatusOK).
//							RequireBody(
//								json.Present("content")).
//							RequireHeaders(
//								headers.Present("Content-Type")).
//							ExecuteTest(context.Background(), t)
//
//						return nil
//					},
//				},
//			},
//			Request: &cute.Request{
//				Builders: []cute.RequestBuilder{
//					cute.WithURI("http://bannerage-e2e:8080/login"),
//					cute.WithMarshalBody(struct {
//						Email    string `json:"email"`
//						Password string `json:"password"`
//					}{
//						Email:    "33kjdkjj123kk2@al.ru",
//						Password: "opopop111",
//					}),
//					cute.WithMethod(http.MethodPost),
//				},
//			},
//			Expect: &cute.Expect{
//				Code: http.StatusOK,
//				AssertBody: []cute.AssertBody{
//					json.Equal("message", "Successfully logged in."),
//				},
//				AssertHeaders: []cute.AssertHeaders{
//					headers.Present("Content-Type"),
//				},
//			},
//		},
//		{},
//	}
//
//	cute.NewTestBuilder().
//		Title("Table tests for banner creation").
//		Tags("table_test_banner_create").
//		Description("Execute array tests for banner creation").
//		CreateTableTest().
//		PutTests(tests...).
//		ExecuteTest(context.Background(), t)
//}
