package main

type parent interface {
	id() string
}

type pageParent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}

func (p *pageParent) id() string {
	return p.PageID
}

type databaseParent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

func (p *databaseParent) id() string {
	return p.DatabaseID
}

type workspaceParent struct {
	Type      string `json:"type"`
	Workspace bool   `json:"workspace"`
}

func (p *workspaceParent) id() string {
	return "workspace"
}

type blockID struct {
	Type    string `json:"type"`
	BlockID string `json:"block_id"`
}

func (p *blockID) id() string {
	return p.BlockID
}

type title struct {
	Title *property `json:"title"`
}

type text struct {
	RichText *property `json:"rich_text"`
}

type property struct {
}

type properties struct {
	ResourceID *title `json:"Resource ID"`
	UserID     *text  `json:"User ID"`
	Access     *text  `json:"Access"`
}

type initReqModel struct {
	DatabaseID string `json:"database_id"`
}

type initResModel struct {
	Parent *pageParent `json:"parent"`
}

type authReqModel struct {
	Token string `json:"token"`
}

type createReqModel struct {
	Parent     *pageParent `json:"parent"`
	Title      string      `json:"title"`
	Properties *properties `json:"properties"`
}

type createResModel struct {
	DatabaseID string `json:"id"`
	parent     parent `json:"parent"`
}

type errorModel struct {
	Message string `json:"message"`
}

func (err *errorModel) Error() string {
	return err.Message
}

type checkReqModel struct {
	ResourceID string `json:"rid"`
	UserID     string `json:"uid"`
}

type checkResModel struct {
	Results []result `json:"results"`
}

type result struct {
	Properties *checkProperty `json:"properties"`
	PageID     string         `json:"id"`
}

type checkProperty struct {
	Access *access `json:"Access"`
}

type access struct {
	TextArr []richText `json:"rich_text"`
}

type richText struct {
	Text *plainText `json:"text"`
}

type plainText struct {
	Content string `json:"content"`
}

type checkBody struct {
	Filter *filter `json:"filter"`
}

type filter struct {
	And []condition `json:"and"`
}

type condition struct {
	Property string        `json:"property"`
	RichText *textContains `json:"rich_text"`
}

type textContains struct {
	Contains string `json:"contains"`
}

type updateReqModel struct {
	Access string `json:"access"`
}

type updateResModel struct {
}

type grantReqModel struct {
	Parent     *databaseParent  `json:"parent"`
	Properties *grantProperties `json:"properties"`
}

type grantProperties struct {
	ResourceID *titleCol `json:"Resource ID"`
	UserID     *textCol  `json:"User ID"`
	Access     *textCol  `json:"Access"`
}

type titleCol struct {
	TitleArr []updateText `json:"title"`
}

type textCol struct {
	TextArr []updateText `json:"rich_text"`
}

type updateText struct {
	Text *plainText `json:"text"`
}

type changeReqModel struct {
	Properties *changeProperties `json:"properties"`
}

type changeProperties struct {
	Access *textCol `json:"Access"`
}

type revokeResModel struct {
}
