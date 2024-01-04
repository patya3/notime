package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/patya3/notime/pkg/models/timelog"
)

// NOTE: itt majd lehet a startedAt mas lesz
func constructBody(comment string, startedAt string, durationInSeconds int) ([]byte, error) {
	body := map[string]interface{}{
		"comment": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"content": []map[string]interface{}{
						{
							"text": comment,
							"type": "text",
						},
					},
				},
			},
			"type":    "doc",
			"version": 1,
		},
		"started":          startedAt,
		"timeSpentSeconds": durationInSeconds,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("cant parse json", err)
		return nil, err
	}

	fmt.Println(body)
	return jsonData, nil
}

func CreateWorklog(timelog timelog.ExtendedLog) {
	jiraBaseUrl := os.Getenv("JIRA_BASE_URL")
	userEmail := os.Getenv("JIRA_EMAIL")
	apiToken := os.Getenv("JIRA_API_TOKEN")

	jsonBody, err := constructBody(timelog.Comment, timelog.CreatedAt.String(), timelog.GetLogDurationInSeconds())
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", jiraBaseUrl+"/rest/api/3/issue/VM-3135/worklog", bodyReader)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(userEmail, apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client: response body: %s\n", body)

}
