package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus/log"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func CreateKerberosClientWithPassword(principal string, password string) (*client.Client, error) {

	// Load the client krb5 config
	cfg, err := config.Load("/etc/krb5.conf")

	if err != nil {
		return nil, fmt.Errorf("failed to load krb5 cfg")

	}

	username, realm := ExtractUsernameAndRealm(principal)

	if username == "" {
		return nil, fmt.Errorf("failed to extract username and realm from principal")

	}

	cli := client.NewClientWithPassword(username, realm, password, cfg)

	// Log in the client
	err = cli.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to login krb5 client")

	}

	return cli, nil
}

func CreateKerberosClientWithKeytab(ktPath string, principal string) (*client.Client, error) {
	// https://github.com/jcmturner/gokrb5/blob/855dbc707a37a21467aef6c0245fcf3328dc39ed/USAGE.md?plain=1#L20
	kt, err := keytab.Load(ktPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load keytab file: %v", err)
	}

	krb5Conf, err := config.Load("/etc/krb5.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to load Kerberos config: %v", err)
	}

	username, realm := ExtractUsernameAndRealm(principal)

	if username == "" {
		return nil, fmt.Errorf("failed to extract username and realm from principal")
	}

	cli := client.NewClientWithKeytab(username, realm, kt, krb5Conf)

	// Log in the client
	err = cli.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to login krb5 client")

	}

	return cli, nil
}

func ExtractUsernameAndRealm(principal string) (string, string) {
	parts := strings.Split(principal, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func ExtractDomainFromURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()

	return host, nil
}

func MakeKrb5Request(client *client.Client, url string) ([]byte, error) {

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("could not create request: %v", err)
		return nil, fmt.Errorf("could not create request: %v", err)

	}

	fqdn, err := ExtractDomainFromURL(url)
	if err != nil {
		log.Errorf("could not extract fqdn from url: %v", err)
		return nil, fmt.Errorf("could not extract fqdn from url: %v", err)

	}

	spn := fmt.Sprintf("HTTP/%s", fqdn)

	spnegoCl := spnego.NewClient(client, nil, spn)

	err = spnego.SetSPNEGOHeader(client, r, spn)
	if err != nil {
		log.Errorf("error set spnego header: %v", err)
		return nil, fmt.Errorf("error set spnego header: %v", err)
	}

	// Make the request
	resp, err := spnegoCl.Client.Do(r)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	// fmt.Println(string(body))

	defer resp.Body.Close()

	return body, nil

}

func MakeKrb5RequestWithKeytab(ktPath string, principal string, url string) ([]byte, error) {

	krb5cli, err := CreateKerberosClientWithKeytab(ktPath, principal)

	if err != nil {
		log.Errorf("could not create krb5 client: %v", err)
		return nil, fmt.Errorf("could not create krb5 client: %v", err)
	}

	return MakeKrb5Request(krb5cli, url)

}

func MakeKrb5RequestWithPassword(principal string, password string, url string) ([]byte, error) {

	krb5cli, err := CreateKerberosClientWithPassword(principal, password)

	if err != nil {
		log.Errorf("could not create krb5 client: %v", err)
		return nil, fmt.Errorf("could not create krb5 client: %v", err)
	}

	return MakeKrb5Request(krb5cli, url)

}
