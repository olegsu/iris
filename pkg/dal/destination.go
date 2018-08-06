package dal

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"k8s.io/api/core/v1"
)

type Destination struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Secret   string `yaml:"secret"`
	Type     string `yaml:"type"`
	Branch   string `yaml:"branch"`
	Pipeline string `yaml:"pipeline"`
	CFToken  string `yaml:"cftoken"`
}

func getHmac(secret string, payload []byte) string {
	if secret != "" {
		fmt.Println("Singing payload with secret")
		key := []byte(secret)
		mac := hmac.New(sha256.New, key)
		mac.Write(payload)
		hmac := base64.URLEncoding.EncodeToString(mac.Sum(nil))
		return hmac
	}
	return ""
}

func (d *Destination) Exec(payload interface{}) {
	fmt.Printf("Executing destination %s\n", d.Name)
	if d.Type == "" {
		execDefault(d, payload)
	} else if d.Type == "Codefresh" {
		execCodefresh(d, payload)
	}
}

func execDefault(d *Destination, payload interface{}) {
	fmt.Printf("Executing default destination to %s\n", d.URL)
	mJSON, _ := json.Marshal(payload)
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", d.URL, contentReader)
	req.Header.Set("X-IRIS-HMAC", getHmac(d.Secret, mJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)
}

func execCodefresh(d *Destination, payload interface{}) {
	postBody := &codefreshPostRequestBody{
		Variables: make(map[string]string),
	}
	if d.Branch != "" {
		postBody.Branch = d.Branch
	}
	var ev *v1.Event
	b, _ := json.Marshal(payload)
	json.Unmarshal(b, &ev)

	postBody.Variables["IRIS_RESOURCE_NAME"] = ev.InvolvedObject.Name
	postBody.Variables["IRIS_NAMESPACE"] = ev.InvolvedObject.Namespace

	mJSON, _ := json.Marshal(postBody)
	contentReader := bytes.NewReader(mJSON)

	url := fmt.Sprintf("https://g.codefresh.io/api/pipelines/run/%s", url.QueryEscape(d.Pipeline))
	fmt.Printf("Executing Codefresh destination\n")
	fmt.Printf(string(mJSON))
	req, _ := http.NewRequest("POST", url, contentReader)
	req.Header.Set("authorization", d.CFToken)
	req.Header.Set("User-Agent", "IRIS")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("Build ID: %s\n", string(body))
	} else {
		fmt.Printf("Error:\nStatus Code: $d\nBody: %s\n", resp.StatusCode, string(body))
	}
}

type codefreshPostRequestBody struct {
	Options   map[string]string `json:"options"`
	Variables map[string]string `json:"variables"`
	Contexts  []string          `json:"contexts"`
	Branch    string            `json:"branch"`
}
