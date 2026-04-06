package tg_models

type BotInfo struct {
	Ok     bool           `json:"ok"`
	Result *BotInfoResult `json:"result"`
}

type BotInfoResult struct {
	Id                     int64  `json:"id"`
	IsBot                  bool   `json:"is_bot"`
	FirstName              string `json:"first_name"`
	UserName               string `json:"username"`
	CanJoinGroups          bool   `json:"can_join_groups"`
	CanReadAllGroupMessage bool   `json:"can_read_all_group_message"`
	SupportsInlineQueries  bool   `json:"supports_inline_queries"`
	CanConnectToBusiness   bool   `json:"can_connect_to_business"`
	HasMainWebApp          bool   `json:"has_main_web_app"`
}
