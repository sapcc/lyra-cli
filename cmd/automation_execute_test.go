package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetAutomationLExecute() {
	// reset automation flag vars
	ResetFlags()
}

func TestAutomationExecute(t *testing.T) {
	responseBody := `{"automation_attributes": null,"automation_id": "74","automation_name": "script_test","created_at": "2018-01-23T16:01:40.173Z","id": "109829","jobs": null,"log": null,"repository_revision": null,"selector": "@identity=\"886ea868-ba06-42f8-9bde-eb1a848938\"","state": "preparing","updated_at": "2018-01-23T16:01:40.173Z"}`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	want := `-----------------------+------------------------------------------------+
|          KEY          |                     VALUE                      |
+-----------------------+------------------------------------------------+
| automation_attributes | <nil>                                          |
| automation_id         | 74                                             |
| automation_name       | script_test                                    |
| created_at            | 2018-01-23T16:01:40.173Z                       |
| id                    | 109829                                         |
| jobs                  | <nil>                                          |
| log                   | <nil>                                          |
| repository_revision   | <nil>                                          |
| selector              | @identity="886ea868-ba06-42f8-9bde-eb1a848938" |
| state                 | preparing                                      |
| updated_at            | 2018-01-23T16:01:40.173Z                       |
+-----------------------+------------------------------------------------+`

	// reset stuff
	resetAutomationLExecute()

	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation execute --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s --selector=%s", "some_test_url", "miau", "bup", "123456789", "automation_id", "@identity=886ea868-ba06-42f8-9bde-eb1a848938"))

	if resulter.Error != nil {
		t.Error(fmt.Sprintf("Command expected to not get an error. \n \n %s", resulter.Error))
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationExecuteWatchSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		if path == "/api/v1/runs" { //create run
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `{"id":"109924","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a848938\":\nNo nodes found.\n","created_at":"2018-01-24T09:07:51.047Z","updated_at":"2018-01-24T09:07:52.854Z","repository_revision":null,"state":"executing","jobs":null,"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a848938\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
		} else if path == "/api/v1/runs/109924" { // get info run
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `{"id":"109936","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\":\n886ea868-ba06-42f8-9bde-eb1a84893b58 testito\nUsing exiting artifact for revision cc5fdd31010df5a825936c5b1edda86980007c92\nScheduled 1 job:\nca3d2c86-ff0f-41a8-8dbc-d169b1492910\n","created_at":"2018-01-24T09:38:55.294Z","updated_at":"2018-01-24T09:38:58.113Z","repository_revision":"cc5fdd31010df5a825936c5b1edda86980007c92","state":"completed","jobs":["ca3d2c86-ff0f-41a8-8dbc-d169b1492910"],"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
		} else if path == "/api/v1/jobs/cc5fdd31010df5a825936c5b1edda86980007c92" { // get job info success
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `{"version":1,"sender":"linux","request_id":"ca3d2c86-ff0f-41a8-8dbc-d169b1492910","to":"886ea868-ba06-42f8-9bde-eb1a84893b58","timeout":100,"agent":"execute","action":"tarball","payload":"{\"path\":\"script.sh\",\"url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/automation-artifacts/cc5fdd31010df5a825936c5b1edda86980007c92-script.tgz?temp_url_sig=3bf4712d59c89137c0fe59b2ac2f8c4dbc375860\\u0026temp_url_expires=1516790338\"}","status":"complete","created_at":"2018-01-24T09:38:58.105641Z","updated_at":"2018-01-24T09:39:08.142134Z","project":"abcdefghijklmnopqrstuwxyz1234567","user":{"domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName","id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123"}}`)
		}
	}))
	defer testServer.Close()

	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	// reset stuff
	resetAutomationLExecute()

	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation execute --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s --selector=%s --watch", testServer.URL, testServer.URL, "token123", "some_test_url", "miau", "bup", "123456789", "automation_id", "@identity=886ea868-ba06-42f8-9bde-eb1a848938"))

	if resulter.Error != nil {
		t.Error(fmt.Sprintf("Command expected to not get an error. \n \n %s", resulter.Error))
	}
}

func TestAutomationExecuteWatchRetry(t *testing.T) {
	runsCalls := 0
	runCalls := 0
	jobCalls := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		if path == "/api/v1/runs" { //create run
			runsCalls += 1
			if runsCalls == 1 {
				w.WriteHeader(503)
				fmt.Fprintln(w, `{"error":"503"}`)
			} else {
				w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
				fmt.Fprintln(w, `{"id":"109924","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a848938\":\nNo nodes found.\n","created_at":"2018-01-24T09:07:51.047Z","updated_at":"2018-01-24T09:07:52.854Z","repository_revision":null,"state":"executing","jobs":null,"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a848938\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
			}
		} else if path == "/api/v1/runs/109924" { // get info run
			runCalls += 1
			if runCalls == 1 {
				w.WriteHeader(503)
				fmt.Fprintln(w, `{"error":"503"}`)
			} else {
				w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
				fmt.Fprintln(w, `{"id":"109936","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\":\n886ea868-ba06-42f8-9bde-eb1a84893b58 testito\nUsing exiting artifact for revision cc5fdd31010df5a825936c5b1edda86980007c92\nScheduled 1 job:\nca3d2c86-ff0f-41a8-8dbc-d169b1492910\n","created_at":"2018-01-24T09:38:55.294Z","updated_at":"2018-01-24T09:38:58.113Z","repository_revision":"cc5fdd31010df5a825936c5b1edda86980007c92","state":"completed","jobs":["ca3d2c86-ff0f-41a8-8dbc-d169b1492910"],"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
			}
		} else if path == "/api/v1/jobs/cc5fdd31010df5a825936c5b1edda86980007c92" { // get job info success
			jobCalls += 1
			if jobCalls == 1 {
				w.WriteHeader(503)
				fmt.Fprintln(w, `{"error":"503"}`)
			} else {
				w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
				fmt.Fprintln(w, `{"version":1,"sender":"linux","request_id":"ca3d2c86-ff0f-41a8-8dbc-d169b1492910","to":"886ea868-ba06-42f8-9bde-eb1a84893b58","timeout":100,"agent":"execute","action":"tarball","payload":"{\"path\":\"script.sh\",\"url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/automation-artifacts/cc5fdd31010df5a825936c5b1edda86980007c92-script.tgz?temp_url_sig=3bf4712d59c89137c0fe59b2ac2f8c4dbc375860\\u0026temp_url_expires=1516790338\"}","status":"complete","created_at":"2018-01-24T09:38:58.105641Z","updated_at":"2018-01-24T09:39:08.142134Z","project":"abcdefghijklmnopqrstuwxyz1234567","user":{"domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName","id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123"}}`)
			}
		}
	}))
	defer testServer.Close()

	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	// reset stuff
	resetAutomationLExecute()

	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation execute --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s --selector=%s --watch", testServer.URL, testServer.URL, "token123", "some_test_url", "miau", "bup", "123456789", "automation_id", "@identity=886ea868-ba06-42f8-9bde-eb1a848938"))
	if resulter.Error != nil {
		t.Error(fmt.Sprintf("Command expected to not get an error. \n \n %s", resulter.Error))
	}
}

func TestAutomationExecuteWatchFailed(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		runCalls := 0
		if path == "/api/v1/runs" { //create run
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `{"id":"109924","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a848938\":\nNo nodes found.\n","created_at":"2018-01-24T09:07:51.047Z","updated_at":"2018-01-24T09:07:52.854Z","repository_revision":null,"state":"executing","jobs":null,"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a848938\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
		} else if path == "/api/v1/runs/109924" { // get info run
			runCalls += 1
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			if runCalls == 1 {
				fmt.Fprintln(w, `{"id":"109936","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\":\n886ea868-ba06-42f8-9bde-eb1a84893b58 testito\nUsing exiting artifact for revision cc5fdd31010df5a825936c5b1edda86980007c92\nScheduled 1 job:\nca3d2c86-ff0f-41a8-8dbc-d169b1492910\n","created_at":"2018-01-24T09:38:55.294Z","updated_at":"2018-01-24T09:38:58.113Z","repository_revision":"cc5fdd31010df5a825936c5b1edda86980007c92","state":"failed","jobs":["ca3d2c86-ff0f-41a8-8dbc-d169b1492910"],"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
			} else if runCalls > 1 {
				fmt.Fprintln(w, `{"id":"109936","log":"Selecting nodes using filter @identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\":\n886ea868-ba06-42f8-9bde-eb1a84893b58 testito\nUsing exiting artifact for revision cc5fdd31010df5a825936c5b1edda86980007c92\nScheduled 1 job:\nca3d2c86-ff0f-41a8-8dbc-d169b1492910\n","created_at":"2018-01-24T09:38:55.294Z","updated_at":"2018-01-24T09:38:58.113Z","repository_revision":"cc5fdd31010df5a825936c5b1edda86980007c92","state":"failed","jobs":["ca3d2c86-ff0f-41a8-8dbc-d169b1492910"],"owner":{"id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123","domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName"},"automation_id":"74","automation_name":"script_test","selector":"@identity=\"886ea868-ba06-42f8-9bde-eb1a84893b58\"","automation_attributes":{"name":"script_test","path":"script.sh","timeout":100,"repository":"https://github.com/user123/automation-test.git","repository_revision":"master"}}`)
			}
		} else if path == "/api/v1/jobs/cc5fdd31010df5a825936c5b1edda86980007c92" { // get job info success
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `{"version":1,"sender":"linux","request_id":"ca3d2c86-ff0f-41a8-8dbc-d169b1492910","to":"886ea868-ba06-42f8-9bde-eb1a84893b58","timeout":100,"agent":"execute","action":"tarball","payload":"{\"path\":\"script.sh\",\"url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/automation-artifacts/cc5fdd31010df5a825936c5b1edda86980007c92-script.tgz?temp_url_sig=3bf4712d59c89137c0fe59b2ac2f8c4dbc375860\\u0026temp_url_expires=1516790338\"}","status":"failed","created_at":"2018-01-24T09:38:58.105641Z","updated_at":"2018-01-24T09:39:08.142134Z","project":"abcdefghijklmnopqrstuwxyz1234567","user":{"domain_id":"abcdefghyjklmnopqrstuwxyz9876543","domain_name":"domeName","id":"abcdefghijklmnopqrstuwxyz1234567890abcdefghyjklmnopqrstuwxyz1234","name":"user123"}}`)
		}
	}))
	defer testServer.Close()

	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	// reset stuff
	resetAutomationLExecute()

	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation execute --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s --selector=%s --watch", testServer.URL, testServer.URL, "token123", "some_test_url", "miau", "bup", "123456789", "automation_id", "@identity=886ea868-ba06-42f8-9bde-eb1a848938"))

	if resulter.Error == nil {
		t.Error(fmt.Sprintf("Command expected to get an error. \n \n %s", resulter.Error))
	}
}
