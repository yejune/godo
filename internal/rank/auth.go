package rank

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// OAuth authentication constants.
const (
	OAuthTimeout = 300 * time.Second
	PortMin      = 8080
	PortMax      = 8180
	StateBytes   = 32
)

// CallbackResult holds the result received from the OAuth callback.
type CallbackResult struct {
	Credentials *Credentials
	Error       error
}

// StartOAuthFlow initiates the OAuth authentication flow for the Rank API.
func StartOAuthFlow(baseURL string) (*Credentials, error) {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	state, err := GenerateStateToken()
	if err != nil {
		return nil, fmt.Errorf("generate state token: %w", err)
	}

	port, ln, err := FindAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("find callback port: %w", err)
	}

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/callback", port)
	resultCh := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		HandleOAuthCallback(w, r, state, resultCh)
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if serveErr := server.Serve(ln); serveErr != nil && serveErr != http.ErrServerClosed {
			resultCh <- CallbackResult{Error: fmt.Errorf("callback server: %w", serveErr)}
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	authURL := fmt.Sprintf("%s/api/auth/cli?redirect_uri=%s&state=%s", baseURL, callbackURL, state)
	if openErr := OpenBrowser(authURL); openErr != nil {
		fmt.Printf("Could not open browser automatically.\nPlease visit: %s\n", authURL)
	}

	ctx, cancel := context.WithTimeout(context.Background(), OAuthTimeout)
	defer cancel()

	select {
	case result := <-resultCh:
		if result.Error != nil {
			return nil, result.Error
		}
		return result.Credentials, nil
	case <-ctx.Done():
		return nil, &AuthError{Message: "OAuth flow timed out (5 minutes)"}
	}
}

// GenerateStateToken generates a cryptographically random state token for CSRF protection.
func GenerateStateToken() (string, error) {
	b := make([]byte, StateBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// FindAvailablePort finds an available TCP port in the range [PortMin, PortMax].
func FindAvailablePort() (int, net.Listener, error) {
	for port := PortMin; port <= PortMax; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		return port, ln, nil
	}
	return 0, nil, fmt.Errorf("no available port in range %d-%d", PortMin, PortMax)
}

// OpenBrowser opens the given URL in the default browser.
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// HandleOAuthCallback processes the OAuth callback request.
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request, expectedState string, resultCh chan<- CallbackResult) {
	query := r.URL.Query()

	receivedState := query.Get("state")
	if receivedState != expectedState {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, ErrorHTML("Authentication failed: invalid state token."))
		resultCh <- CallbackResult{Error: &AuthError{Message: "state token mismatch"}}
		return
	}

	if errMsg := query.Get("error"); errMsg != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, ErrorHTML("Authentication error: "+errMsg))
		resultCh <- CallbackResult{Error: &AuthError{Message: errMsg}}
		return
	}

	apiKey := query.Get("api_key")
	if apiKey == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, ErrorHTML("Authentication failed: missing API key."))
		resultCh <- CallbackResult{Error: &AuthError{Message: "missing api_key in callback"}}
		return
	}

	creds := &Credentials{
		APIKey:    apiKey,
		Username:  query.Get("username"),
		UserID:    query.Get("user_id"),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, SuccessHTML())
	resultCh <- CallbackResult{Credentials: creds}
}

// SuccessHTML returns an HTML page shown after successful authentication.
func SuccessHTML() string {
	return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Authentication Successful</title>
<style>body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#f0fdf4}
.card{text-align:center;padding:2rem;border-radius:12px;background:#fff;box-shadow:0 2px 8px rgba(0,0,0,0.1)}
h1{color:#16a34a;margin-bottom:0.5rem}p{color:#666}</style></head>
<body><div class="card"><h1>Authentication Successful</h1>
<p>You can close this window and return to the terminal.</p></div></body></html>`
}

// ErrorHTML returns an HTML page shown after a failed authentication.
func ErrorHTML(message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Authentication Failed</title>
<style>body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#fef2f2}
.card{text-align:center;padding:2rem;border-radius:12px;background:#fff;box-shadow:0 2px 8px rgba(0,0,0,0.1)}
h1{color:#dc2626;margin-bottom:0.5rem}p{color:#666}</style></head>
<body><div class="card"><h1>Authentication Failed</h1>
<p>%s</p></div></body></html>`, message)
}
