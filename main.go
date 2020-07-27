package main

import (
	"bytes"
	"errors"
	"encoding/json"
	"strings"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type RequestInvitePayload struct {
}

type Configure struct {
	ApiKey		string	`yaml:"apiKey" validate:"required,len=32,alphanum"`
	Auth		string	`yaml:"auth" validate:"required,is-auth"`
}

func printUsage() {
	exe, err := os.Executable()

	if err != nil {
		logrus.Fatal("Failed to lookup executable: ", err)
	}

	logrus.Println(
		"Usage:",
		filepath.Base(exe),
		"[vrchat:// link]",
	)
}

func validateAuth(fl validator.FieldLevel) bool {
	r := regexp.MustCompile(`^authcookie_[0-9a-f]{8}(-[0-9a-f]{4}){4}[0-9a-f]{8}$`)
	return r.MatchString(fl.Field().String())
}

func readConfigure() Configure {
	exe, err := os.Executable()

	if err != nil {
		logrus.Fatal("Failed to lookup executable: ", err)
	}

	bytes, err := ioutil.ReadFile(
		filepath.Join(filepath.Dir(exe), "configure.yml"),
	)

	if err != nil {
		logrus.Fatal("Failed to read configure file: ", err)
	}

	c := Configure{}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		logrus.Fatal("Failed to parse configure file: ", err)
	}

	v := validator.New()

	v.RegisterValidation("is-auth", validateAuth)

	if err := v.Struct(c); err != nil {
		logrus.Fatal(err)
	}

	return c
}

func parseVRChatLink(link string) (string, error) {
	if !strings.HasPrefix(link, "vrchat://launch") {
		return "", errors.New("link is not vrchat://launch link")
	}

	splittedLink := strings.SplitN(link, "?", 2)

	if len(splittedLink) != 2 {
		return "", errors.New("link is not valid vrchat://launch link")
	}

	for _, param := range strings.Split(splittedLink[1], "&") {
		parsedParam := strings.SplitN(param, "=", 2)

		if len(parsedParam) != 2 {
			continue
		}

		if parsedParam[0] != "id" {
			continue
		}

		return parsedParam[1], nil
	}

	return "", errors.New("link is not has id=")
}

// It's not "request-invite". request "invite".
func requestInvite(apiKey, auth, id string) error {
	url := fmt.Sprintf(
		"https://vrchat.com/api/1/instances/%s/invite?apiKey=%s",
		id,
		apiKey,
	)

	jsonBytes, err := json.Marshal(RequestInvitePayload{})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBytes),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Cookie", "auth=" + auth)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status code: %v", resp.StatusCode)
	}

	return nil
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	c := readConfigure()

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		printUsage()
		logrus.Fatal("Illegal argument(s) count.")
	}

	id, err := parseVRChatLink(args[0])

	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Invite to %v", id)

	if err := requestInvite(c.ApiKey, c.Auth, id); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Success")
}

