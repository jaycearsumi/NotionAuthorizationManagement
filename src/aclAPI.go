package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type aclAPI struct {
	auth       string
	databaseID string
}

func createACLAPI() *aclAPI {
	return &aclAPI{}
}

func (api *aclAPI) retrieveDatabase(params *initReqModel) (*initResModel, error) {
	url := "https://api.notion.com/v1/databases/" + params.DatabaseID

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model initResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *aclAPI) createDatabase(params *createReqModel) (*createResModel, error) {
	url := "https://api.notion.com/v1/databases"

	params.Properties = &properties{}

	p, _ := json.Marshal(params)
	log.Println(string(p))

	payload := strings.NewReader(string(p))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model createResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *aclAPI) checkAccess(params *checkReqModel) (*checkResModel, error) {
	url := "https://api.notion.com/v1/databases/" + api.databaseID + "/query"

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
					Property: "User ID",
					RichText: &textContains{
						Contains: params.UserID,
					},
				},
			},
		},
	}

	tmp, _ := json.Marshal(bodyParams)

	log.Println(string(tmp))

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

	log.Println(string(body))

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

func (api *aclAPI) updateAccess(bodyParams *updateReqModel, checkParams *checkReqModel) (*updateResModel, error) {
	checkRes, err := api.checkAccess(checkParams)
	var resModel *updateResModel
	var error error
	if err != nil {
		return resModel, err
	} else if len(checkRes.Results) == 0 {
		resModel, error = api.grantAccess(bodyParams, checkParams)
	} else {
		resModel, error = api.changeAccess(bodyParams, checkRes)
	}
	return resModel, error
}

func (api *aclAPI) grantAccess(bodyParams *updateReqModel, checkParams *checkReqModel) (*updateResModel, error) {
	url := "https://api.notion.com/v1/pages"

	reqParams := &grantReqModel{

		Parent: &databaseParent{
			DatabaseID: api.databaseID,
			Type:       "database_id",
		},
		Properties: &grantProperties{
			ResourceID: &titleCol{
				TitleArr: []updateText{
					{
						Text: &plainText{
							Content: checkParams.ResourceID,
						},
					},
				},
			},
			UserID: &textCol{
				TextArr: []updateText{
					{
						Text: &plainText{
							Content: checkParams.UserID,
						},
					},
				},
			},
			Access: &textCol{
				TextArr: []updateText{
					{
						Text: &plainText{
							Content: bodyParams.Access,
						},
					},
				},
			},
		},
	}
	p, _ := json.Marshal(reqParams)

	payload := strings.NewReader(string(p))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var err errorModel
	var model updateResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *aclAPI) changeAccess(bodyParams *updateReqModel, checkRes *checkResModel) (*updateResModel, error) {

	url := "https://api.notion.com/v1/pages/" + checkRes.Results[0].PageID

	reqParams := &changeReqModel{
		Properties: &changeProperties{
			Access: &textCol{
				TextArr: []updateText{
					{
						Text: &plainText{
							Content: bodyParams.Access,
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
	var model updateResModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &err)
		return &model, &err
	} else {
		json.Unmarshal(body, &model)
		return &model, nil
	}
}

func (api *aclAPI) revokeAccess(checkParams *checkReqModel) (*revokeResModel, error) {
	checkRes, err := api.checkAccess(checkParams)
	var resModel *revokeResModel
	if err != nil {
		return resModel, err
	}
	url := "https://api.notion.com/v1/blocks/" + checkRes.Results[0].PageID
	req, _ := http.NewRequest("PATCH", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+api.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))

	var error errorModel
	if strings.Contains(string(body), "error") {
		json.Unmarshal(body, &error)
		return resModel, &error
	} else {
		json.Unmarshal(body, resModel)
		return resModel, nil
	}
}
