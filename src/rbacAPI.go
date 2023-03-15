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
	//log.Println(string(body))

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
	accessRes, accessErr := api.checkRoleAccess(roleRes, params)
	if accessErr != nil {
		return resModel, accessErr
	}
	return accessRes, nil
}

func (api *rbac3API) checkUserRole(params *checkReqModel) (*checkRoleResModel, error) {
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
	var model checkRoleResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *rbac3API) checkRoleAccess(roleRes *checkRoleResModel, params *checkReqModel) (*checkResModel, error) {
	url := "https://api.notion.com/v1/databases/" + api.RoleAccessDatabaseID + "/query"

	bodyParams := &checkBody{
		Filter: &filter{
			And: []condition{
				{
					Property: "Role",
					RichText: &textContains{
						Contains: roleRes.Results[0].Properties.Role.RoleArr[0].Name,
					},
				},
				{
					Property: "Resource ID",
					RichText: &textContains{
						Contains: params.ResourceID,
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

func (api *rbac3API) updateRole(params *rbac3UpdateRoleReqModel) (*rbac3UpdateRoleResModel, error) {
	checkRes, checkErr := api.checkUserRole(&checkReqModel{UserID: params.UserID})
	if checkErr != nil {
		return nil, checkErr
	}
	url := "https://api.notion.com/v1/pages/" + checkRes.Results[0].PageID

	reqParams := &roleUpdateBodyModel{
		Properties: &roleProperty{
			Role: &role{
				RoleArr: []multiSelect{
					{
						Name: params.Role,
					},
				},
			},
		},
	}
	for _, obj := range checkRes.Results[0].Properties.Role.RoleArr {
		reqParams.Properties.Role.RoleArr = append(reqParams.Properties.Role.RoleArr, multiSelect{Name: obj.Name})
	}

	p, _ := json.Marshal(reqParams)

	payload := strings.NewReader(string(p))

	req, _ := http.NewRequest("PATCH", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model rbac3UpdateRoleResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *rbac3API) checkRole() (*rbac3CheckRoleResModel, error) {
	url := "https://api.notion.com/v1/databases/" + api.userRoleDatabaseID

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//log.Println(string(body))

	var err errorModel
	var model rbac3CheckRoleResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *rbac3API) revokeRole(params *rbac3UpdateRoleReqModel) (*rbac3UpdateRoleResModel, error) {
	checkRes, checkErr := api.checkUserRole(&checkReqModel{UserID: params.UserID})
	if checkErr != nil {
		return nil, checkErr
	}
	url := "https://api.notion.com/v1/pages/" + checkRes.Results[0].PageID

	reqParams := &roleUpdateBodyModel{
		Properties: &roleProperty{
			Role: &role{
				RoleArr: []multiSelect{},
			},
		},
	}
	for _, obj := range checkRes.Results[0].Properties.Role.RoleArr {
		if obj.Name == params.Role {
			continue
		}
		reqParams.Properties.Role.RoleArr = append(reqParams.Properties.Role.RoleArr, multiSelect{Name: obj.Name})
	}

	p, _ := json.Marshal(reqParams)

	payload := strings.NewReader(string(p))

	req, _ := http.NewRequest("PATCH", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model rbac3UpdateRoleResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *rbac3API) updateAccess(params *rbac3UpdateAccessReqModel) (*rbac3UpdateAccessResModel, error) {
	checkRes, checkErr := api.retrieveRoleAccess(params)
	if checkErr != nil {
		return nil, checkErr
	}
	url := "https://api.notion.com/v1/pages/" + checkRes.Results[0].PageID

	reqParams := &accessUpdateBodyModel{
		Properties: &changeProperties{
			Access: &textCol{
				TextArr: []updateText{
					{
						Text: &plainText{
							Content: params.Access,
						},
					},
				},
			},
		},
	}

	p, _ := json.Marshal(reqParams)

	payload := strings.NewReader(string(p))

	req, _ := http.NewRequest("PATCH", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model rbac3UpdateAccessResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return nil, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *rbac3API) retrieveRoleAccess(params *rbac3UpdateAccessReqModel) (*checkResModel, error) {
	url := "https://api.notion.com/v1/databases/" + api.RoleAccessDatabaseID + "/query"

	bodyParams := &checkBody{
		Filter: &filter{
			And: []condition{
				{
					Property: "Resource ID",
					RichText: &textContains{
						Contains: params.ResourceID,
					},
				},
				{
					Property: "Role",
					RichText: &textContains{
						Contains: params.Role,
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

func (api *rbac3API) revokeAccess(params *rbac3UpdateAccessReqModel) (*rbac3UpdateAccessResModel, error) {
	checkRes, checkErr := api.retrieveRoleAccess(params)
	if checkErr != nil {
		return nil, checkErr
	}
	url := "https://api.notion.com/v1/blocks/" + checkRes.Results[0].PageID

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model rbac3UpdateAccessResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return nil, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}
