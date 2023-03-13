package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type rbac3API struct {
	auth                 string
	homePageID           string
	userRoleDatabaseID   string
	RoleAccessDatabaseID string
}

func createRbac3API() *rbac3API {
	return &rbac3API{}
}

func (api *rbac3API) homePageInit(params *rbac3InitReqModel) (*rbac3InitResModel, error) {
	url := "https://api.notion.com/v1/blocks/" + params.HomePageID + "/children"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model rbac3InitResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		api.initDatabase(&model)
		return &model, nil
	}
}

func (api *rbac3API) initDatabase(res *rbac3InitResModel) {
	for _, obj := range res.ResArr {
		if obj.Database.Title == "user-role" {
			api.userRoleDatabaseID = obj.DatabaseID
		} else if obj.Database.Title == "role-access" {
			api.RoleAccessDatabaseID = obj.DatabaseID
		}
	}
	log.Println(api.userRoleDatabaseID)
	log.Println(api.RoleAccessDatabaseID)
}

func (api *rbac3API) rbac3CheckAccess(params *checkReqModel) (*checkResModel, error) {
	var resModel *checkResModel
	roleRes, roleErr := api.checkUserRole(params)
	if roleErr != nil {
		return resModel, roleErr
	}
	accessRes, accessErr := api.checkRoleAccess(roleRes)
	if accessErr != nil {
		return resModel, accessErr
	}
	return accessRes, nil
}

func (api *rbac3API) checkUserRole(params *checkReqModel) (*checkResModel, error) {
	url := "https://api.notion.com/v1/databases/" + api.userRoleDatabaseID + "/query"

	bodyParams := &checkBody{
		Filter: &filter{
			And: []condition{
				{
					Property: "User ID",
					RichText: &textContains{
						Contains: params.UserID,
					},
				},
			},
		},
	}

	tmp, _ := json.Marshal(bodyParams)

	//log.Println(string(tmp))

	payload := strings.NewReader(string(tmp))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	//log.Println(error.Error())

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//log.Println(string(body))

	var err errorModel
	var model checkResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}
