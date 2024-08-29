package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func TestNodeInstallCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node install")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node install")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node install")
}

func TestNodeInstallCmdWithUserAuthenticationFlags(t *testing.T) {
	server := nodeInstallServer()
	defer server.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(server)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s --install-format=%s", "some_test_url", "kuak", "bup", "123456789", "node_id", "linux"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeInstallFormatDefault(t *testing.T) {
	want := `json script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	ResetFlags()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "http://somewhere.com", server.URL, "token123"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}

func TestNodeInstallJsonSuccessDefaultUrls(t *testing.T) {
	want := `json script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	ResetFlags()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "json"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}

func TestNodeInstallLinuxSuccessDefaultUrls(t *testing.T) {
	want := `shell script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	ResetFlags()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "linux"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}

func TestNodeInstallWindowsSuccessDefaultUrls(t *testing.T) {
	want := `powershell script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	ResetFlags()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "windows"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}

func TestNodeInstallCloudCinfigSuccessDefaultUrls(t *testing.T) {
	want := `cloud config script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	ResetFlags()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "cloud-config"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}

func TestNodeInstallCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		format := r.Header.Get("Accept")
		if !strings.Contains(method, "POST") {
			diffString := StringDiff(method, "POST")
			t.Errorf("Command API method doesn't match. \n \n %s", diffString)
		}
		if !strings.Contains(path.String(), "agents/init") {
			diffString := StringDiff(method, "agents/init")
			t.Errorf("Command API path doesn't match. \n \n %s", diffString)
		}
		if !strings.Contains(format, "text/x-shellscript") {
			diffString := StringDiff(format, "text/x-shellscript")
			t.Errorf("Command API format doesn't match. \n \n %s", diffString)
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s --install-format=%s", "https://somewhere.com", server.URL, "token123", "123456789", "linux"))
}

func nodeInstallServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		format := r.Header.Get("Accept")

		if format == "text/x-shellscript" {
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `shell script`)
		} else if format == "text/x-powershellscript" {
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `powershell script`)
		} else if format == "text/cloud-config" {
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `cloud config script`)
		} else {
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `json script`)
		}
	}))
	return server
}
