package groupvalidate

type GroupCreateValidate struct {
	CreatedUid int    `json:"created_uid" bind:"required"`
	GroupName  string `json:"group_name" bind:"required"`
}
