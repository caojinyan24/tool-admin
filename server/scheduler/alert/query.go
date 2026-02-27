package alert

import (
	"context"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

func QueryOpenId(ctx context.Context, notifier map[string]model.Notifier) (map[string]model.Notifier, error) {
	emailList := make([]string, 0)
	for email := range notifier {
		emailList = append(emailList, email)
	}
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`open_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails(emailList).
			IncludeResigned(true).
			Build()).
		Build()

	// 发起请求
	resp, err := LarkClient.Contact.V3.User.BatchGetId(ctx, req)

	// 处理错误
	if err != nil {
		logger.Error(ctx, "QueryOpenId err: ", err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		logger.Error(ctx, "QueryOpenId err: ", resp.CodeError)
		return nil, err
	}
	logger.Info(ctx, "QueryOpenId resp: ", larkcore.Prettify(resp))
	result := make(map[string]model.Notifier)
	if resp.Data != nil && resp.Data.UserList != nil {
		for _, user := range resp.Data.UserList {
			if user.UserId != nil {
				result[*user.Email] = model.Notifier{
					ReceiveId: *user.UserId,
				}
			}
		}
	}
	return result, nil
}

func QueryWithEmail(ctx context.Context, email string) (string, error) {

	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`open_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails([]string{email}).
			IncludeResigned(true).
			Build()).
		Build()

	// 发起请求
	resp, err := LarkClient.Contact.V3.User.BatchGetId(ctx, req)

	// 处理错误
	if err != nil {
		logger.Error(ctx, "QueryOpenId err: ", err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		logger.Error(ctx, "QueryOpenId err: ", resp.CodeError)
		return "", err
	}
	logger.Info(ctx, "QueryOpenId resp: ", larkcore.Prettify(resp))
	if resp.Data != nil && resp.Data.UserList != nil {
		for _, user := range resp.Data.UserList {
			if user.UserId != nil && *user.Email == email {
				return *user.UserId, nil
			}
		}
	}

	return "", nil
}
