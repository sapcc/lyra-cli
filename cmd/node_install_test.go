package cmd

import (
	"fmt"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetNodeInstall() {
	// reset automation flag vars
	ResetFlags()

	auth.AuthenticationV3 = auth.NewMockAuthenticationV3
	auth.CommonResult1 = map[string]interface{}{"token": map[string]interface{}{"project": map[string]string{"id": "test_project_id", "domain_id": "monsooniii", "name": "Arc_Test"}}}
}

func TestNodeInstallUserIdOrNameRequired(t *testing.T) {
	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --project-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestNodeInstallProjectIdOrNameRequired(t *testing.T) {
	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestNodeInstallOSRequired(t *testing.T) {
	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-identifier=%s", "http://some_test_url", "user", "project", "123456789", "identifier"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestNodeInstallIdentifierRequired(t *testing.T) {
	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --instance-os=%s", "http://some_test_url", "user", "project", "123456789", "linux"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestNodeInstallLinuxSuccessDefaultUrls(t *testing.T) {
	want := `curl --create-dirs -o /opt/arc/arc https://arc-updates.***REMOVED***/builds/latest/arc/linux/amd64
chmod +x /opt/arc/arc
/opt/arc/arc init --endpoint tls://arc-broker.***REMOVED***:8883 --update-uri https://arc-updates.***REMOVED***/updates --registration-url this_is_mock_registration_url`

	// set test server
	responseBody := `{"token":"some_nice_token", "url":"this_is_mock_registration_url"}`
	server := TestServer(200, responseBody, map[string]string{})

	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-identifier=%s --instance-os=%s --pki-service-url=%s", "http://some_test_url", "user", "project", "123456789", "identifer_test", "linux", server.URL))
	if resulter.Error != nil {
		t.Error(`Command expected to get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallWindowsSuccessDefaultUrls(t *testing.T) {
	want := `mkdir C:\monsoon\arc
powershell (new-object System.Net.WebClient).DownloadFile('https://arc-updates.***REMOVED***/builds/latest/arc/windows/amd64','C:\monsoon\arc\arc.exe')
C:\monsoon\arc\arc.exe init --endpoint tls://arc-broker.***REMOVED***:8883 --update-uri https://arc-updates.***REMOVED***/updates --registration-url this_is_mock_registration_url`

	// set test server
	responseBody := `{"token":"some_nice_token", "url":"this_is_mock_registration_url"}`
	server := TestServer(200, responseBody, map[string]string{})

	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-identifier=%s --instance-os=%s --pki-service-url=%s", "http://some_test_url", "user", "project", "123456789", "identifer_test", "windows", server.URL))
	if resulter.Error != nil {
		t.Error(`Command expected to get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}

func TestNodeInstallLinuxSuccessNotDefaultUrls(t *testing.T) {
	want := `curl --create-dirs -o /opt/arc/arc http://test_update_url/builds/latest/arc/linux/amd64
chmod +x /opt/arc/arc
/opt/arc/arc init --endpoint http://test_broker_url --update-uri http://test_update_url/updates --registration-url this_is_mock_registration_url`

	// set test server
	responseBody := `{"token":"some_nice_token", "url":"this_is_mock_registration_url"}`
	server := TestServer(200, responseBody, map[string]string{})

	// reset params
	resetNodeInstall()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node install --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-identifier=%s --instance-os=%s --update-service-url=%s --arc-broker-url=%s --pki-service-url=%s", "http://some_test_url", "user", "project", "123456789", "identifer_test", "linux", "http://test_update_url", "http://test_broker_url", server.URL))
	if resulter.Error != nil {
		t.Error(`Command expected to get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
	}
}
