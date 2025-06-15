package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	issueModel "github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui/constants"
)

type GetIssueResponse struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type Issue struct {
	Expand string `json:"expand"`
	Id     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

// TODO: return with error if statusCode is not 20x
func CreateJiraWorkLog(timelog timelog.ExtendedLog) ([]byte, error) {

	jsonBody, err := constructBodyForWorklog(timelog.Comment, timelog.CreatedAt.Format("2006-01-02T15:04:05.000-0700"), timelog.GetLogDurationInSeconds())
	bodyReader := bytes.NewReader(jsonBody)

	url := JiraBaseUrl + "/rest/api/3/issue/" + timelog.IssueKey + "/worklog"
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(UserEmail, ApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	constants.LogRepo.MarkAsLogged(timelog.ID)

	return body, nil
}

func SyncIssuesFromJira(projectKey string) error {

	url := JiraBaseUrl + "/rest/api/3/search?jql=project%20%3D%20'" + projectKey + "'%20AND%20" + Jql
	fmt.Println(Jql)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(UserEmail, ApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response GetIssueResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	projectId, err := constants.ProjectRepo.GetProjectIdByKey(projectKey)
	if err != nil {
		return err
	}

	for _, issue := range response.Issues {

		issue := issueModel.Issue{
			IssueKey:   issue.Key,
			IssueTitle: issue.Fields.Summary,
			Desc:       "",
			ProjectID:  projectId,
		}

		constants.IssueRepo.CreateIssue(issue)
	}

	return nil
}

func constructBodyForWorklog(comment string, startedAt string, durationInSeconds int) ([]byte, error) {
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
					"type": "paragraph",
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
		return nil, err
	}

	return jsonData, nil
}
