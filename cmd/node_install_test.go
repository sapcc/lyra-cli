package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetNodeInstall() {
	// reset automation flag vars
	ResetFlags()
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

func TestNodeInstallFormatDefault(t *testing.T) {
	want := `json script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	resetNodeInstall()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "http://somewhere.com", server.URL, "token123"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallJsonSuccessDefaultUrls(t *testing.T) {
	want := `json script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	resetNodeInstall()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "json"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallLinuxSuccessDefaultUrls(t *testing.T) {
	want := `shell script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	resetNodeInstall()

	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "linux"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallWindowsSuccessDefaultUrls(t *testing.T) {
	want := `powershell script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "windows"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallCloudCinfigSuccessDefaultUrls(t *testing.T) {
	want := `cloud config script`
	server := nodeInstallServer()
	defer server.Close()

	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --install-format=%s", "http://somewhere.com", server.URL, "token123", "cloud-config"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}
