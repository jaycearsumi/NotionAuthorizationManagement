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

type checkRoleResModel struct {
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
	RoleArr []multiSelect `json:"multi_select"`
}

type multiSelect struct {
	Name string `json:"name"`
}

type rbac3UpdateRoleReqModel struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type rbac3CheckRoleResModel struct {
	Properties *roleProperty `json:"properties"`
}

type roleProperty struct {
	Role *role `json:"Role"`
}

type roleUpdateBodyModel struct {
	Properties *roleProperty `json:"properties"`
}

type rbac3UpdateRoleResModel struct {
}

type rbac3UpdateAccessReqModel struct {
	ResourceID string `json:"resource_id"`
	Role       string `json:"role"`
	Access     string `json:"access"`
}

type accessUpdateBodyModel struct {
	Properties *changeProperties `json:"properties"`
}

type rbac3UpdateAccessResModel struct {
}
