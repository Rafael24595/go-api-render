package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
)

type requestCloneCollection struct {
	CollectionName string `json:"collection_name"`
}

type requestExecuteAction struct {
	Request dto.DtoRequest `json:"request"`
	Context dto.DtoContext `json:"context"`
}

type requestImportContext struct {
	Target dto.DtoContext `json:"target"`
	Source dto.DtoContext `json:"source"`
}

type requestInsertAction struct {
	Request  dto.DtoRequest
	Response dto.DtoResponse
}

type requestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type requestVerify struct {
	OldPassword  string `json:"old_password"`
	NewPassword1 string `json:"new_password_1"`
	NewPassword2 string `json:"new_password_2"`
}

type requestSigninUser struct {
	Username  string `json:"username"`
	Password1 string `json:"password_1"`
	Password2 string `json:"password_2"`
	IsAdmin   bool   `json:"is_admin"`
}

type requestPushToCollection struct {
	SourceId    string              `json:"source_id"`
	TargetId    string              `json:"target_id"`
	TargetName  string              `json:"target_name"`
	Request     dto.DtoRequest      `json:"request"`
	RequestName string              `json:"request_name"`
	Movement    repository.Movement `json:"move"`
}

type requestSortNodes struct {
	Nodes []requestNode `json:"nodes"`
}

type requestNode struct {
	Order int    `json:"order"`
	Item  string `json:"item"`
}

func requestPushToCollectionToPayload(payload *requestPushToCollection) repository.PayloadCollectRequest {
	request := dto.ToRequest(&payload.Request)
	return repository.PayloadCollectRequest{
		SourceId:    payload.SourceId,
		TargetId:    payload.TargetId,
		TargetName:  payload.TargetName,
		Request:     *request,
		RequestName: payload.RequestName,
		Movement:    payload.Movement,
	}
}

func requestSortCollectionToPayload(payload *requestSortNodes) repository.PayloadSortNodes {
	nodes := make([]repository.PayloadCollectionNode, len(payload.Nodes))
	for i, v := range payload.Nodes {
		nodes[i] = repository.PayloadCollectionNode{
			Order: v.Order,
			Item:  v.Item,
		}
	}
	return repository.PayloadSortNodes{
		Nodes: nodes,
	}
}
