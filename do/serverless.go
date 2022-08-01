/*
Copyright 2018 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package do

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/pkg/extract"
	"github.com/digitalocean/godo"
)

// ServerlessCredentials models what is stored in credentials.json for use by the plugin and nim.
// It is also the type returned by the GetServerlessNamespace function.
type ServerlessCredentials struct {
	APIHost     string                                     `json:"currentHost"`
	Namespace   string                                     `json:"currentNamespace"`
	Credentials map[string]map[string]ServerlessCredential `json:"credentials"`
}

// ServerlessCredential is the type of an individual entry in ServerlessCredentials
type ServerlessCredential struct {
	Auth string `json:"api_key"`
}

// The type of the "namespace" member of the response to /api/v2/functions/sandbox and
// /api/v2/functions/namespaces APIs.  Only relevant fields unmarshalled
type OutputNamespace struct {
	Namespace string `json:"namespace"`
	APIHost   string `json:"api_host"`
	UUID      string `json:"uuid"`
	Key       string `json:"key"`
	Label     string `json:"label"`
	Region    string `json:"Region"`
}

// FunctionParameter is the type of a parameter in the response body of action.get.  We do our
// own JSON unmarshaling of these because the go OpenWhisk client doesn't include the "init" and
// "encryption" members, of which at least "init" is needed.
type FunctionParameter struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Init       bool   `json:"init"`
	Encryption string `json:"encryption"`
}

// FunctionParameterReparse is a partial remapping of whisk.Action so that the parameters
// are declared as FunctionParameter rather than whisk.KeyValue.
type FunctionParameterReparse struct {
	Parameters []FunctionParameter `json:"parameters"`
}

// NamespaceResponse is the type of the response body for /api/v2/functions/sandbox (POST) and
// /api/v2/functions/namespaces/<nsName> (GET)
type NamespaceResponse struct {
	Namespace OutputNamespace `json:"namespace"`
}

// NamespaceListResponse is the type of the response body for /api/v2/functions/namespaces (GET)
type NamespaceListResponse struct {
	Namespaces []OutputNamespace `json:"namespaces"`
}

// newNamespaceRequest is the type of the POST body for requesting a new namespace
type newNamespaceRequest struct {
	Namespace inputNamespace `json:"namespace"`
}

// inputNamespace is the reduced representation of a namespace used when requesting a new one
type inputNamespace struct {
	Label  string `json:"label"`
	Region string `json:"Region"`
}

// ServerlessRuntime is the type of a runtime entry returned by the API host controller
// of the serverless cluster.
// Only relevant fields unmarshalled
type ServerlessRuntime struct {
	Default    bool   `json:"default"`
	Deprecated bool   `json:"deprecated"`
	Kind       string `json:"kind"`
}

// ServerlessHostInfo is the type of the host information return from the API host controller
// of the serverless cluster.
// Only relevant fields unmarshaled.
type ServerlessHostInfo struct {
	Runtimes map[string][]ServerlessRuntime `json:"runtimes"`
}

// FunctionInfo is the type of an individual function in the output
// of doctl sls fn list.  Only relevant fields are unmarshaled.
// Note: when we start replacing the sandbox plugin path with direct calls
// to backend controller operations, this will be replaced by declarations
// in the golang openwhisk client.
type FunctionInfo struct {
	Name        string       `json:"name"`
	Namespace   string       `json:"namespace"`
	Updated     int64        `json:"updated"`
	Version     string       `json:"version"`
	Annotations []Annotation `json:"annotations"`
}

// Annotation is a key/value type suitable for individual annotations
type Annotation struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ServerlessService is an interface for interacting with the sandbox plugin,
// with the namespaces service, and with the serverless cluster controller.
type ServerlessService interface {
	Cmd(string, []string) (*exec.Cmd, error)
	Exec(*exec.Cmd) (ServerlessOutput, error)
	Stream(*exec.Cmd) error
	GetServerlessNamespace(context.Context) (ServerlessCredentials, error)
	ListNamespaces(context.Context) (NamespaceListResponse, error)
	GetNamespace(context.Context, string) (ServerlessCredentials, error)
	CreateNamespace(context.Context, string, string) (ServerlessCredentials, error)
	DeleteNamespace(context.Context, string) error
	WriteCredentials(ServerlessCredentials) error
	ReadCredentials() (ServerlessCredentials, error)
	GetHostInfo(string) (ServerlessHostInfo, error)
	CheckServerlessStatus(string) error
	InstallServerless(string, bool) error
	GetFunction(string, bool) (whisk.Action, []FunctionParameter, error)
	GetConnectedAPIHost() (string, error)
}

type serverlessService struct {
	serverlessJs  string
	serverlessDir string
	credsDir      string
	node          string
	userAgent     string
	client        *godo.Client
	owClient      *whisk.Client
}

const (
	// Minimum required version of the sandbox plugin code.  The first part is
	// the version of the incorporated Nimbella CLI and the second part is the
	// version of the bridge code in the sandbox plugin repository.
	minServerlessVersion = "4.1.0-1.3.0"

	// The version of nodejs to download alongsize the plugin download.
	nodeVersion = "v16.13.0"

	// credsDir is the directory under the sandbox where all credentials are stored.
	// It in turn has a subdirectory for each access token employed (formed as a prefix of the token).
	credsDir = "creds"

	// CredentialsFile is the name of the file where the sandbox plugin stores OpenWhisk credentials.
	CredentialsFile = "credentials.json"
)

var _ ServerlessService = &serverlessService{}

var (
	// ErrServerlessNotInstalled is the error returned to users when the sandbox is not installed.
	ErrServerlessNotInstalled = errors.New("Serverless support is not installed (use `doctl serverless install`)")

	// ErrServerlessNeedsUpgrade is the error returned to users when the sandbox is at too low a version
	ErrServerlessNeedsUpgrade = errors.New("Serverless support needs to be upgraded (use `doctl serverless upgrade`)")

	// ErrServerlessNotConnected is the error returned to users when the sandbox is not connected to a namespace
	ErrServerlessNotConnected = errors.New("Serverless support is installed but not connected to a functions namespace (use `doctl serverless connect`)")
)

// ServerlessOutput contains the output returned from calls to the sandbox plugin.
type ServerlessOutput struct {
	Table     []map[string]interface{} `json:"table,omitempty"`
	Captured  []string                 `json:"captured,omitempty"`
	Formatted []string                 `json:"formatted,omitempty"`
	Entity    interface{}              `json:"entity,omitempty"`
	Error     string                   `json:"error,omitempty"`
}

// NewServerlessService returns a configured ServerlessService.
func NewServerlessService(client *godo.Client, usualServerlessDir string, credsToken string) ServerlessService {
	nodeBin := "node"
	if runtime.GOOS == "windows" {
		nodeBin = "node.exe"
	}
	// The following is needed to support snap installation.  For snap, the installation directory
	// is relocated to a snap-managed area.  That area is not user-writable, so, the credsDir location
	// is always computed relative to the normal installation area (usualServerlessDir).
	serverlessDir := os.Getenv("OVERRIDE_SANDBOX_DIR")
	if serverlessDir == "" {
		serverlessDir = usualServerlessDir
	}
	return &serverlessService{
		serverlessJs:  filepath.Join(serverlessDir, "sandbox.js"),
		serverlessDir: serverlessDir,
		credsDir:      GetCredentialDirectory(credsToken, usualServerlessDir),
		node:          filepath.Join(serverlessDir, nodeBin),
		userAgent:     fmt.Sprintf("doctl/%s serverless/%s", doctl.DoitVersion.String(), minServerlessVersion),
		client:        client,
		owClient:      nil,
	}
}

// InitWhisk is an on-demand initializer for the OpenWhisk client, called when that client
// is needed.
func initWhisk(s *serverlessService) error {
	if s.owClient != nil {
		return nil
	}
	creds, err := s.ReadCredentials()
	if err != nil {
		return err
	}
	credential := creds.Credentials[creds.APIHost][creds.Namespace]
	config := whisk.Config{Host: creds.APIHost, AuthToken: credential.Auth}
	client, err := whisk.NewClient(http.DefaultClient, &config)
	if err != nil {
		return err
	}
	s.owClient = client
	return nil
}

func (s *serverlessService) CheckServerlessStatus(leafCredsDir string) error {
	_, err := os.Stat(s.serverlessDir)
	if os.IsNotExist(err) {
		return ErrServerlessNotInstalled
	}
	if !serverlessUptodate(s.serverlessDir) {
		return ErrServerlessNeedsUpgrade
	}
	if !isServerlessConnected(leafCredsDir, s.serverlessDir) {
		return ErrServerlessNotConnected
	}
	return nil
}

// InstallServerless is the common subroutine for both serverless install and serverless upgrade
func (s *serverlessService) InstallServerless(leafCredsDir string, upgrading bool) error {
	serverlessDir := s.serverlessDir

	// Make a temporary directory for use during the install.
	// Note: we don't let this be allocated in the system temporaries area because
	// that might be on a separate file system, meaning that the final install step
	// will require an additional copy rather than a simple rename.

	tmp, err := ioutil.TempDir(filepath.Dir(serverlessDir), "sbx-install")
	if err != nil {
		return err
	}

	// Download the nodejs tarball for this os and architecture
	fmt.Print("Downloading...")

	goos := runtime.GOOS
	arch := runtime.GOARCH
	nodeBin := "node"
	if arch == "amd64" {
		arch = "x64"
	}
	if arch == "386" {
		if goos == "linux" {
			return errors.New("serverless support is not available for 32-bit linux")
		}
		arch = "x86"
	}
	if goos == "windows" {
		goos = "win"
		nodeBin = "node.exe"
	}

	var (
		nodeURL      string
		nodeFileName string
		nodeDir      string
	)

	// Download nodejs only if necessary
	if !upgrading || !canReuseNode(serverlessDir, nodeBin) {
		nodeDir = fmt.Sprintf("node-%s-%s-%s", nodeVersion, goos, arch)
		nodeURL = fmt.Sprintf("https://nodejs.org/dist/%s/%s.tar.gz", nodeVersion, nodeDir)
		nodeFileName = filepath.Join(tmp, "node-install.tar.gz")

		if goos == "win" {
			nodeURL = fmt.Sprintf("https://nodejs.org/dist/%s/%s.zip", nodeVersion, nodeDir)
			nodeFileName = filepath.Join(tmp, "node-install.zip")
		}

		err = download(nodeURL, nodeFileName)
		if err != nil {
			return err
		}
	}

	// Download the fat tarball with the nim CLI, deployer, and sandbox bridge
	// TODO do these need to be arch-specific?  Currently assuming not.
	URL := fmt.Sprintf("https://do-serverless-tools.nyc3.digitaloceanspaces.com/doctl-sandbox-%s.tar.gz",
		GetMinServerlessVersion())
	sandboxFileName := filepath.Join(tmp, "doctl-sandbox.tar.gz")
	err = download(URL, sandboxFileName)
	if err != nil {
		return err
	}

	// Exec the Extract utility at least once to unpack the fat tarball and possibly a second time if
	// node was downloaded.  If node was not downloaded, just move the existing binary into place.
	fmt.Print("Unpacking...")
	err = extract.Extract(sandboxFileName, tmp)
	if err != nil {
		return err
	}

	if nodeFileName != "" {
		err = extract.Extract(nodeFileName, tmp)
		if err != nil {
			return err
		}
	}

	// Move artifacts to final location
	fmt.Print("Installing...")
	srcPath := filepath.Join(tmp, "sandbox")
	if upgrading {
		// Preserve credentials by moving them from target (which will be replaced) to source.
		err = PreserveCreds(leafCredsDir, srcPath, serverlessDir)
		if err != nil {
			return err
		}
		// Preserve existing node if necessary
		if nodeFileName == "" {
			// Node was not downloaded
			err = moveExistingNode(serverlessDir, srcPath, nodeBin)
			if err != nil {
				return err
			}
		}
	} else {
		// Make new empty credentials directory
		emptyCreds := filepath.Join(srcPath, credsDir)
		err = os.MkdirAll(emptyCreds, 0700)
		if err != nil {
			return nil
		}

		// Create the sandbox directory if necessary.
		err := os.MkdirAll(serverlessDir, 0755)
		if err != nil {
			return err
		}
	}
	// Remove former serverlessDir before moving in the new one
	err = os.RemoveAll(serverlessDir)
	if err != nil {
		return err
	}
	err = os.Rename(srcPath, serverlessDir)
	if err != nil {
		return err
	}

	if nodeFileName != "" {
		if goos == "win" {
			srcPath = filepath.Join(tmp, nodeDir, nodeBin)
		} else {
			// Additional nesting in non-windows case
			srcPath = filepath.Join(tmp, nodeDir, "bin", nodeBin)
		}
		destPath := filepath.Join(serverlessDir, nodeBin)
		err = os.Rename(srcPath, destPath)
		if err != nil {
			return err
		}
	}
	// Clean up temp directory
	fmt.Print("Cleaning up...")
	os.RemoveAll(tmp) // Best effort, ignore error
	fmt.Println("\nDone")
	return nil
}

// Cmd builds an *exec.Cmd for calling into the sandbox plugin.
func (s *serverlessService) Cmd(command string, args []string) (*exec.Cmd, error) {
	args = append([]string{s.serverlessJs, command}, args...)
	cmd := exec.Command(s.node, args...)
	cmd.Env = append(os.Environ(), "NIMBELLA_DIR="+s.credsDir, "NIM_USER_AGENT="+s.userAgent)
	// If DEBUG is specified, we need to open up stderr for that stream.  The stdout stream
	// will continue to work for returning structured results.
	if os.Getenv("DEBUG") != "" {
		cmd.Stderr = os.Stderr
	}
	return cmd, nil
}

// Exec executes an *exec.Cmd and captures its output in a ServerlessOutput.
func (s *serverlessService) Exec(cmd *exec.Cmd) (ServerlessOutput, error) {
	output, err := cmd.Output()
	if err != nil {
		// Ignore "errors" that are just non-zero exit.  The
		// serverless uses this as a secondary indicator but the output
		// is still trustworthy (and includes error information inline)
		if _, ok := err.(*exec.ExitError); !ok {
			// Real error of some sort
			return ServerlessOutput{}, err
		}
	}
	var result ServerlessOutput
	err = json.Unmarshal(output, &result)
	if err != nil {
		return ServerlessOutput{}, err
	}
	// Result is sound JSON but also has an error field, meaning that something did
	// go wrong.  In this case we return the actual output but also the distinguished
	// error return.  Most callers will process only the error, which is fine.  Sometimes,
	// however, there is other information that can be useful as part of the error report.
	if len(result.Error) > 0 {
		return result, errors.New(result.Error)
	}
	// Result is both sound and error free
	return result, nil
}

// Stream is like Exec but assumes that output will not be captured and can be streamed.
func (s *serverlessService) Stream(cmd *exec.Cmd) error {

	return cmd.Run()
}

// GetServerlessNamespace returns the credentials of the one serverless namespace assigned to
// the invoking doctl context.
func (s *serverlessService) GetServerlessNamespace(ctx context.Context) (ServerlessCredentials, error) {
	path := "v2/functions/sandbox"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return ServerlessCredentials{}, err
	}
	return executeNamespaceRequest(ctx, s, req)
}

// executeNamespaceRequest executes a valid http.Request object where the request is expected
// to return a NamespaceResponse.  The response is converted to ServerlessCredentials.  The request
// may represent the (new) 'namespaces/<name>' GET API, the (legacy) 'sandbox' POST API, or
// a namespace creation.
// The legacy API will continue to be used by some users until feature-flipper protection is removed
// from the new one.
func executeNamespaceRequest(ctx context.Context, s *serverlessService, req *http.Request) (ServerlessCredentials, error) {
	decoded := new(NamespaceResponse)
	_, err := s.client.Do(ctx, req, decoded)
	if err != nil {
		return ServerlessCredentials{}, err
	}
	host := assignAPIHost(decoded.Namespace.APIHost, decoded.Namespace.Namespace)
	credential := ServerlessCredential{Auth: decoded.Namespace.UUID + ":" + decoded.Namespace.Key}
	namespace := decoded.Namespace.Namespace
	ans := ServerlessCredentials{
		APIHost:     host,
		Namespace:   namespace,
		Credentials: map[string]map[string]ServerlessCredential{host: {namespace: credential}},
	}
	return ans, nil
}

// ListNamespaces obtains the full list of namespaces that belong to the requesting account
func (s *serverlessService) ListNamespaces(ctx context.Context) (NamespaceListResponse, error) {
	path := "v2/functions/namespaces"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return NamespaceListResponse{}, err
	}
	decoded := new(NamespaceListResponse)
	_, err = s.client.Do(ctx, req, decoded)
	if err != nil {
		return NamespaceListResponse{}, err
	}
	return *decoded, nil
}

// GetNamespace obtains the credentials of a specific namespace, given its name
func (s *serverlessService) GetNamespace(ctx context.Context, name string) (ServerlessCredentials, error) {
	path := "v2/functions/namespaces/" + name
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ServerlessCredentials{}, err
	}
	return executeNamespaceRequest(ctx, s, req)
}

// CreateNamespace creates a new namespace and returns its credentials, given a label and region
func (s *serverlessService) CreateNamespace(ctx context.Context, label string, region string) (ServerlessCredentials, error) {
	reqBody := newNamespaceRequest{Namespace: inputNamespace{Label: label, Region: region}}
	path := "v2/functions/namespaces"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return ServerlessCredentials{}, err
	}
	return executeNamespaceRequest(ctx, s, req)
}

// DeleteNamespace deletes a namespace by name
func (s *serverlessService) DeleteNamespace(ctx context.Context, name string) error {
	path := "v2/functions/namespaces/" + name
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// GetHostInfo returns the HostInfo structure of the provided API host
func (s *serverlessService) GetHostInfo(APIHost string) (ServerlessHostInfo, error) {
	endpoint := APIHost + "/api/v1"
	resp, err := http.Get(endpoint)
	if err != nil {
		return ServerlessHostInfo{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var result ServerlessHostInfo
	err = json.Unmarshal(body, &result)
	return result, err
}

// GetFunction returns the metadata and optionally the code of a deployer function
func (s *serverlessService) GetFunction(name string, fetchCode bool) (whisk.Action, []FunctionParameter, error) {
	err := initWhisk(s)
	if err != nil {
		return whisk.Action{}, []FunctionParameter{}, nil
	}
	action, resp, err := s.owClient.Actions.Get(name, fetchCode)
	var parameters []FunctionParameter
	if resp != nil && err == nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			reparse := FunctionParameterReparse{}
			err = json.Unmarshal(body, &reparse)
			if err != nil {
				return whisk.Action{}, []FunctionParameter{}, err
			}
			parameters = reparse.Parameters
		}
	}
	return *action, parameters, err
}

// GetConnectedAPIHost retrieves the API host to which the service is currently connected
func (s *serverlessService) GetConnectedAPIHost() (string, error) {
	err := initWhisk(s)
	if err != nil {
		return "", err
	}
	return s.owClient.Config.Host, nil
}

// Assign the correct API host based on the namespace name.
// Every serverless cluster has two domain names, one ending in '.io', the other in '.co'.
// By convention, the portal only returns the '.io' one but 'doctl sbx' must start using
// only the '.co' one (the '.io' one will eventually require mtls authentication).
// During a migration period, we can continue to support reconnection to "old" namespaces
// (not prefixed by "fn-") but should make sure that all "new" namespaces (prefixed by "fn-")
// switch their API host name from '.io' to '.co'.  Eventually, reconnection to old
// namespaces will fail and they will be removed.  We will need to ensure that users are
// using a doctl containing this code so they can obtain conforming namespaces.
func assignAPIHost(origAPIHost string, namespace string) string {
	if strings.HasPrefix(namespace, "fn-") {
		hostParts := strings.Split(origAPIHost, ".")
		sansSuffix := strings.Join(hostParts[:len(hostParts)-1], ".")
		return sansSuffix + ".co"
	}
	return origAPIHost
}

// WriteCredentials writes a set of serverless credentials to the appropriate 'creds' directory
func (s *serverlessService) WriteCredentials(creds ServerlessCredentials) error {
	// Create the directory into which the file will be written.
	err := os.MkdirAll(s.credsDir, 0700)
	if err != nil {
		return err
	}
	// Write the credentials
	credsPath := filepath.Join(s.credsDir, CredentialsFile)
	bytes, err := json.MarshalIndent(&creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credsPath, bytes, 0600)
}

// ReadCredentials reads the current serverless credentials from the appropriate 'creds' diretory
func (s *serverlessService) ReadCredentials() (ServerlessCredentials, error) {
	creds := ServerlessCredentials{}
	credsPath := filepath.Join(s.credsDir, CredentialsFile)
	bytes, err := os.ReadFile(credsPath)
	if err != nil {
		return creds, err
	}
	err = json.Unmarshal(bytes, &creds)
	return creds, err
}

// Determines whether the serverlessUptodate appears to be connected.  The purpose is
// to fail fast (when feasible) on sandboxes that are clearly not connected.
// However, it is important not to add excessive overhead on each call (e.g.
// asking the plugin to validate credentials), so the test is not foolproof.
// It merely tests whether a credentials directory has been created for the
// current doctl access token and appears to have a credentials.json in it.
func isServerlessConnected(leafCredsDir string, serverlessDir string) bool {
	creds := GetCredentialDirectory(leafCredsDir, serverlessDir)
	credsFile := filepath.Join(creds, CredentialsFile)
	_, err := os.Stat(credsFile)
	return !os.IsNotExist(err)
}

// serverlessUptodate answers whether the installed version of the serverlessUptodate is at least
// what is required by doctl
func serverlessUptodate(serverlessDir string) bool {
	return GetCurrentServerlessVersion(serverlessDir) >= GetMinServerlessVersion()
}

// GetCurrentServerlessVersion gets the version of the current serverless.
// To be called only when serverless is known to exist.
// Returns "0" if the installed serverless pre-dates the versioning system
// Otherwise, returns the version string stored in the serverless directory.
func GetCurrentServerlessVersion(serverlessDir string) string {
	versionFile := filepath.Join(serverlessDir, "version")
	contents, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return "0"
	}
	return string(contents)
}

// GetMinServerlessVersion returns the minServerlessVersion (allows the constant to be overridden via an environment variable)
func GetMinServerlessVersion() string {
	fromEnv := os.Getenv("minServerlessVersion")
	if fromEnv != "" {
		return fromEnv
	}
	return minServerlessVersion
}

// GetCredentialDirectory returns the directory in which credentials should be stored for a given
// CmdConfig.  The actual leaf directory is a function of the access token being used.  This ties
// serverless credentials to DO credentials
func GetCredentialDirectory(leafDir string, serverlessDir string) string {
	return filepath.Join(serverlessDir, credsDir, leafDir)
}

// Gets the version of the node binary in the serverless.  Determine if it is
// usable or whether it has to be upgraded.
func canReuseNode(serverlessDir string, nodeBin string) bool {
	fullNodeBin := filepath.Join(serverlessDir, nodeBin)
	cmd := exec.Command(fullNodeBin, "--version")
	result, err := cmd.Output()
	if err == nil {
		installed := strings.TrimSpace(string(result))
		return installed == nodeVersion
	}
	return false
}

// Moves the existing node binary from the serverless that contains it to the new serverless being
// staged during an upgrade.  This preserves it for reuse and avoids the need to download.
func moveExistingNode(existing string, staging string, nodeBin string) error {
	srcPath := filepath.Join(existing, nodeBin)
	destPath := filepath.Join(staging, nodeBin)
	return os.Rename(srcPath, destPath)
}

// Download a network file to a local file
func download(URL, targetFile string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Received status code %d attempting to download from %s",
			response.StatusCode, URL)
	}
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// PreserveCreds preserves existing credentials in a serverless directory
// that is about to be replaced by moving them to the staging directory
// containing the replacement.
func PreserveCreds(leafDir string, stagingDir string, serverlessDir string) error {
	credPath := filepath.Join(serverlessDir, credsDir)
	relocPath := filepath.Join(stagingDir, credsDir)
	err := os.Rename(credPath, relocPath)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// There was no creds directory.  Check for legacy form and convert as part
	// of preserving.
	legacyCredPath := filepath.Join(serverlessDir, ".nimbella")
	err = os.MkdirAll(relocPath, 0700)
	if err != nil {
		return err
	}
	moveLegacyTo := GetCredentialDirectory(leafDir, stagingDir)
	return os.Rename(legacyCredPath, moveLegacyTo)
}