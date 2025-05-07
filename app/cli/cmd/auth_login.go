package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

type domainCliApp struct {
	serverChan  chan error
	redirectURI string
	callBackURI string
}

// authCmd represents the auth command
func NewAuthLoginCmd() *cobra.Command {
	var headlessLogin bool
	cmd := &cobra.Command{
		Use:   pkg.LoginResource,
		Short: "Login to the VapusData platform instance using Authenticator",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(plclient.MasterGlobals.VapusCtlClient)
			interactiveAuth(headlessLogin)
		},
	}

	cmd.Flags().BoolVar(&headlessLogin, "without-browser", false, "perform a headless login without opening a browser")
	return cmd
}

func interactiveAuth(headlessLogin bool) error {
	var authApp domainCliApp
	checkContext()
	listener, localURL, err := localAuthServer()
	if err != nil {
		return err
	}
	LoginParams, err := plclient.MasterGlobals.VapusCtlClient.RetrieveLoginURL()
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to retrieve login URL")
		cobra.CheckErr(err)
		return err
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Login URL: %v", LoginParams.LoginUrl)
	instanceCallBackURL, err := url.Parse(LoginParams.CallbackUrl)
	if err != nil {
		return err
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Callback URL: %v", instanceCallBackURL)

	localURL.Path = instanceCallBackURL.Path
	loginURLParsed, err := url.Parse(LoginParams.LoginUrl)
	if err != nil {
		return err
	}
	q := loginURLParsed.Query()
	q.Set("redirect_uri", localURL.String())
	q.Set("long-lived", "true")
	loginURLParsed.RawQuery = q.Encode()

	if headlessLogin {
		return headlessAuthn(loginURLParsed)
	}

	err = openBrowser(loginURLParsed.String())
	if err != nil {
		Logger.Debug().Err(err).Msg("falling back to manual login")
		return headlessAuthn(loginURLParsed)
	}

	authApp.serverChan = make(chan error)
	authApp.redirectURI = LoginParams.RedirectUri
	authApp.callBackURI = localURL.String()

	// Run server in background
	http.HandleFunc(instanceCallBackURL.Path, authApp.callbackHandler)
	go func() {
		Logger.Info().Msg("waiting for the authentication to be completed, please check your browser")

		server := &http.Server{ReadHeaderTimeout: time.Second}

		err := server.Serve(listener)
		if err != nil {
			plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to start server")
		}
	}()

	err = <-authApp.serverChan
	if err != nil {
		return err
	}

	return nil
}

func headlessAuthn(loginURL *url.URL) error {
	// Remove cli-callback query parameter to indicate the server to show it inline
	// q := loginURL.Query()
	// q.Del("callback")
	// loginURL.RawQuery = q.Encode()
	// fmt.Printf("To authenticate, click on the following link and paste the result back here\n\n  %s\n\n", loginURL.String())

	// fmt.Print("Enter Token: ")
	// token, err := term.ReadPassword(syscall.Stdin)
	// if err != nil {
	// 	return fmt.Errorf("retrieving password from stdin: %w", err)
	// }

	// // We just want to check that it is a token, the actual verification will happen when it is sent to the server
	// // To be clear, this is just a best effort sanity check
	// if _, _, err := new(jwt.Parser).ParseUnverified(string(token), &jwt.MapClaims{}); err != nil {
	// 	return errors.New("invalid token")
	// }

	// if err := saveAuthToken(string(token)); err != nil {
	// 	return fmt.Errorf("storing token in config file: %w", err)
	// }

	// fmt.Println("")
	// Logger.Info().Msg("login successful!")
	return nil
}

func generateIdToken(code, redirectUri string) error {
	var accessToken string
	var err error
	accessToken, idToken, err := plclient.MasterGlobals.VapusCtlClient.RetrieveAccessToken(code, redirectUri)
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to retrieve id token")
		return err
	}
	viper.Set(currentIdToken, idToken)
	viper.Set(currentAccessToken, accessToken)
	return viper.WriteConfig() // may remove it because token is passed inside the context
}

func (a *domainCliApp) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	defer func() {
		a.serverChan <- nil
	}()
	if code == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	if err := generateIdToken(code, a.callBackURI); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	plclient.MasterGlobals.Logger.Info().Msg("login successful!")
	plclient.MasterGlobals.Logger.Info().Msgf("Redirecting to VapusData platform - %v", a.redirectURI)
	http.Redirect(w, r, "https://www.vapusdata.com", http.StatusSeeOther)
	fmt.Fprintln(w, "login successful, you can now close this window and go back to the terminal")
	return
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// // Create a local HTTP listener with a random available port
func localAuthServer() (net.Listener, *url.URL, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:51101")
	if err != nil {
		return nil, nil, err
	}

	// URL to where the server will post back the auth token
	URL := &url.URL{Scheme: "http", Host: listener.Addr().String()}
	return listener, URL, nil
}

func checkContext() {
	currentContext = viper.GetString(currentContextKey)
	plclient.MasterGlobals.Logger.Info().Msgf("Reading config from - %v", plclient.MasterGlobals.CfgDir)

	if currentContext == "" {
		plclient.MasterGlobals.Logger.Info().Msg("No context is set to current. Please add a context first and then set it to current or set any current context in use.")
		os.Exit(0)
	} else {
		plclient.MasterGlobals.Logger.Info().Msgf("Current context is set to %v", currentContext)
	}
}
