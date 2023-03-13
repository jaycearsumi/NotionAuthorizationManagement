package main

type rbac3InitReqModel struct {
	HomePageID string `json:"page_id"`
}

type rbac3InitResModel struct {
	ResArr []childBlock `json:"results"`
}

type childBlock struct {
	DatabaseID string         `json:"id"`
	Database   *childDatabase `json:"child_database"`
}

type childDatabase struct {
	Title string `json:"title"`
}

type checkRoleRes struct {
	Results []roleResult `json:"results"`
}

type roleResult struct {
	Properties *checkRoleProperty `json:"properties"`
	PageID     string             `json:"id"`
}

type checkRoleProperty struct {
	Role *role `json:"Role"`
}

type role struct {
}
