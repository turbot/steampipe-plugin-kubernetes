package kubernetes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// KOTS API response types (matching JSON from kotsadm API)

type KotsDownstreamVersion struct {
	VersionLabel       string     `json:"versionLabel"`
	UpdateCursor       string     `json:"updateCursor"`
	ChannelID          string     `json:"channelId,omitempty"`
	IsRequired         bool       `json:"isRequired"`
	Status             string     `json:"status"`
	CreatedOn          *time.Time `json:"createdOn,omitempty"`
	ParentSequence     int64      `json:"parentSequence"`
	Sequence           int64      `json:"sequence"`
	DeployedAt         *time.Time `json:"deployedAt,omitempty"`
	Source             string     `json:"source"`
	PreflightSkipped   bool       `json:"preflightSkipped"`
	CommitURL          string     `json:"commitUrl,omitempty"`
	GitDeployable      bool       `json:"gitDeployable,omitempty"`
	UpstreamReleasedAt *time.Time `json:"upstreamReleasedAt,omitempty"`
	ReleaseNotes       string     `json:"releaseNotes,omitempty"`
	IsDeployable       bool       `json:"isDeployable,omitempty"`
	NonDeployableCause string     `json:"nonDeployableCause,omitempty"`
}

type KotsVersionHistoryResponse struct {
	VersionHistory         []*KotsDownstreamVersion `json:"versionHistory"`
	TotalCount             int                      `json:"totalCount"`
	NumOfSkippedVersions   int                      `json:"numOfSkippedVersions"`
	NumOfRemainingVersions int                      `json:"numOfRemainingVersions"`
}

type KotsConfigItem struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Value      string `json:"value"`
	Default    string `json:"default"`
	Data       string `json:"data,omitempty"`
	Filename   string `json:"filename,omitempty"`
	Hidden     bool   `json:"hidden"`
	ReadOnly   bool   `json:"readonly"`
	WriteOnce  bool   `json:"write_once"`
	HelpText   string `json:"help_text,omitempty"`
	Repeatable bool   `json:"repeatable"`
}

type KotsConfigGroup struct {
	Name        string           `json:"name"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Items       []KotsConfigItem `json:"items"`
}

type KotsConfigResponse struct {
	Success      bool              `json:"success"`
	Error        string            `json:"error,omitempty"`
	ConfigGroups []KotsConfigGroup `json:"configGroups"`
}

type KotsApp struct {
	ID                string     `json:"id"`
	Slug              string     `json:"slug"`
	Name              string     `json:"name"`
	IsAirgap          bool       `json:"isAirgap"`
	CurrentSequence   int64      `json:"currentSequence"`
	UpstreamURI       string     `json:"upstreamUri"`
	IconURI           string     `json:"iconUri"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	LastUpdateCheckAt *time.Time `json:"lastUpdateCheckAt"`
	HasPreflight      bool       `json:"hasPreflight"`
	IsConfigurable    bool       `json:"isConfigurable"`
	UpdateCheckerSpec string     `json:"updateCheckerSpec"`
	AutoDeploy        string     `json:"autoDeploy"`
	AppState          string     `json:"appState"`
	LicenseType       string     `json:"licenseType"`
	AllowRollback     bool       `json:"allowRollback"`
	AllowSnapshots    bool       `json:"allowSnapshots"`
	TargetKotsVersion string     `json:"targetKotsVersion"`
	IsSemverRequired  bool       `json:"isSemverRequired"`
	Downstream        KotsAppDownstream `json:"downstream"`
}

type KotsAppDownstream struct {
	Name           string                     `json:"name"`
	CurrentVersion *KotsAppDownstreamVersion  `json:"currentVersion"`
}

type KotsAppDownstreamVersion struct {
	VersionLabel string     `json:"versionLabel"`
	Sequence     int64      `json:"sequence"`
	Status       string     `json:"status"`
	CreatedOn    *time.Time `json:"createdOn,omitempty"`
	DeployedAt   *time.Time `json:"deployedAt,omitempty"`
}

type KotsListAppsResponse struct {
	Apps []KotsApp `json:"apps"`
}

type KotsAppStatusResponse struct {
	AppStatus *KotsAppStatus `json:"appstatus"`
}

type KotsAppStatus struct {
	AppID     string    `json:"appId"`
	UpdatedAt time.Time `json:"updatedAt"`
	State     string    `json:"state"`
	Sequence  int64     `json:"sequence"`
}

// findKotsadmNamespaces discovers all namespaces that have a running kotsadm pod
func findKotsadmNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	// Search across all namespaces for pods with label app=kotsadm
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{LabelSelector: "app=kotsadm"})
	if err != nil {
		return nil, fmt.Errorf("failed to list kotsadm pods across namespaces: %w", err)
	}

	seen := map[string]bool{}
	var namespaces []string
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning && !seen[pod.Namespace] {
			seen[pod.Namespace] = true
			namespaces = append(namespaces, pod.Namespace)
		}
	}

	return namespaces, nil
}

// getKotsNamespaces returns either the single namespace from the query filter,
// or discovers all namespaces where kotsadm is installed.
func getKotsNamespaces(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	namespace := d.EqualsQualString("namespace")
	if namespace != "" {
		return []string{namespace}, nil
	}

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to get clientset: %w", err)
	}
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not available (manifest-only mode)")
	}

	return findKotsadmNamespaces(ctx, clientset)
}

// kotsPortForwardSession holds a port-forward session to kotsadm
type kotsPortForwardSession struct {
	LocalPort int
	StopChan  chan struct{}
	AuthSlug  string
}

// findKotsadmPod finds a running kotsadm pod in the given namespace
func findKotsadmPod(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: "app=kotsadm"})
	if err != nil {
		return "", fmt.Errorf("failed to list kotsadm pods: %w", err)
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return pod.Name, nil
		}
	}

	return "", fmt.Errorf("no running kotsadm pod found in namespace %s", namespace)
}

// getKotsAuthSlug retrieves the auth slug from the kotsadm-authstring secret
func getKotsAuthSlug(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, "kotsadm-authstring", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get kotsadm-authstring secret: %w", err)
	}

	authSlug, ok := secret.Data["kotsadm-authstring"]
	if !ok {
		return "", fmt.Errorf("kotsadm-authstring key not found in secret")
	}

	return string(authSlug), nil
}

// startKotsPortForward establishes a port-forward to kotsadm pod
func startKotsPortForward(ctx context.Context, restConfig *rest.Config, namespace string, podName string) (int, chan struct{}, error) {
	// Find a free port
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to find free port: %w", err)
	}
	localPort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create SPDY dialer
	roundTripper, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create roundtripper: %w", err)
	}

	u, err := url.Parse(restConfig.Host)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse host: %w", err)
	}

	scheme := u.Scheme
	hostIP := u.Host
	path := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/portforward", u.Path, namespace, podName)
	serverURL := url.URL{Scheme: scheme, Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan := make(chan struct{}, 1)
	readyChan := make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:3000", localPort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create port forwarder: %w", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- forwarder.ForwardPorts()
	}()

	// Wait for port-forward to be ready
	select {
	case <-readyChan:
		// Port forward is ready
	case err := <-errChan:
		return 0, nil, fmt.Errorf("port forward failed: %w", err)
	case <-time.After(30 * time.Second):
		close(stopChan)
		return 0, nil, fmt.Errorf("port forward timed out")
	}

	return localPort, stopChan, nil
}

// getKotsSession establishes a port-forward session to kotsadm and caches it
func getKotsSession(ctx context.Context, d *plugin.QueryData, namespace string) (*kotsPortForwardSession, error) {
	cacheKey := fmt.Sprintf("kotsSession-%s", namespace)

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*kotsPortForwardSession), nil
	}

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to get clientset: %w", err)
	}
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not available (manifest-only mode)")
	}

	// Find kotsadm pod
	podName, err := findKotsadmPod(ctx, clientset, namespace)
	if err != nil {
		return nil, err
	}

	// Get auth slug
	authSlug, err := getKotsAuthSlug(ctx, clientset, namespace)
	if err != nil {
		return nil, err
	}

	// Get rest config for port forwarding
	restConfig, err := getRestConfig(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to get rest config: %w", err)
	}

	// Start port forward
	localPort, stopChan, err := startKotsPortForward(ctx, restConfig, namespace, podName)
	if err != nil {
		return nil, err
	}

	session := &kotsPortForwardSession{
		LocalPort: localPort,
		StopChan:  stopChan,
		AuthSlug:  authSlug,
	}

	d.ConnectionManager.Cache.Set(cacheKey, session)

	return session, nil
}

// getRestConfig returns the *rest.Config for the current connection
func getRestConfig(ctx context.Context, d *plugin.QueryData) (*rest.Config, error) {
	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		return nil, err
	}
	if kubeconfig == nil {
		return nil, fmt.Errorf("kubeconfig not available")
	}
	return kubeconfig.ClientConfig()
}

// kotsAPIGet performs a GET request to the kotsadm API
func kotsAPIGet(session *kotsPortForwardSession, path string) ([]byte, error) {
	url := fmt.Sprintf("http://localhost:%d%s", session.LocalPort, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", session.AuthSlug)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d from kotsadm API", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// getKotsVersions fetches version history for an app
func getKotsVersions(session *kotsPortForwardSession, appSlug string) (*KotsVersionHistoryResponse, error) {
	path := fmt.Sprintf("/api/v1/app/%s/versions?currentPage=0&pageSize=1000&pinLatest=false&pinLatestDeployable=false", url.PathEscape(appSlug))
	body, err := kotsAPIGet(session, path)
	if err != nil {
		return nil, err
	}

	var response KotsVersionHistoryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal versions response: %w", err)
	}

	return &response, nil
}

// getKotsConfig fetches config for an app at a given sequence
func getKotsConfig(session *kotsPortForwardSession, appSlug string, sequence int64) (*KotsConfigResponse, error) {
	path := fmt.Sprintf("/api/v1/app/%s/config/%d", url.PathEscape(appSlug), sequence)
	body, err := kotsAPIGet(session, path)
	if err != nil {
		return nil, err
	}

	var response KotsConfigResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("failed to get config: %s", response.Error)
	}

	return &response, nil
}

// getKotsApps fetches the list of KOTS apps
func getKotsApps(session *kotsPortForwardSession) (*KotsListAppsResponse, error) {
	body, err := kotsAPIGet(session, "/api/v1/apps")
	if err != nil {
		return nil, err
	}

	var response KotsListAppsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal apps response: %w", err)
	}

	return &response, nil
}

// getKotsAppStatus fetches the runtime status of a KOTS app
func getKotsAppStatus(session *kotsPortForwardSession, appSlug string) (*KotsAppStatusResponse, error) {
	path := fmt.Sprintf("/api/v1/app/%s/status", url.PathEscape(appSlug))
	body, err := kotsAPIGet(session, path)
	if err != nil {
		return nil, err
	}

	var response KotsAppStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app status response: %w", err)
	}

	return &response, nil
}
