package anchor

import (
	"boframe/controller"
	"boframe/dao/giftdao"
	"boframe/dao/userdao"
	"boframe/pkg/errcode"
	"boframe/settings/mongoI"
	"boframe/utils"
	"cc"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type presentLogResponse struct {
	Total int64                     `json:"total"`
	List  []*presentLogResponseElem `json:"list"`
}

type presentLogResponseElem struct {
	UserId   int64     `json:"user_id"`
	UserName string    `json:"user_name"`
	GiftId   int64     `json:"gift_id"`
	GiftName string    `json:"gift_name"`
	SendAt   time.Time `json:"send_at"`
}

func PresentLog(ctx *cc.Context) {
	query := ctx.Req.URL.Query()
	anchorId, err := strconv.ParseInt(query.Get("anchor_id"), 10, 64)
	if err != nil {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidParam))
		return
	}
	page, err := strconv.ParseInt(query.Get("page"), 10, 64)
	if err != nil {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidParam))
		return
	}
	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidParam))
		return
	}

	if anchorId <= 0 {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrAnchorInvalid))
		return
	}

	if page <= 0 || limit <= 0 {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidPageInfo))
		return
	}

	_, err = userdao.OneByUserId(anchorId)
	if mongoI.IsErrNoDocuments(err) {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrAnchorNotExists))
		return
	}
	if err != nil {
		zap.L().Error("user_dao OneByUserId failed", zap.Error(err), zap.Any("anchor_id", anchorId))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
		return
	}

	condition := bson.D{
		{
			"to_user_id", anchorId,
		},
	}
	sort := bson.D{
		{
			"send_at", -1,
		},
		{
			"_id", -1,
		},
	}
	total, err := giftdao.CountByCondition(condition)
	if err != nil {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrDB))
		return
	}
	if total <= 0 {
		controller.ResponseSuccessWithData(ctx, &presentLogResponse{List: []*presentLogResponseElem{}})
		return
	}

	logs, err := giftdao.PagerByAnchorIds(page, limit, condition, sort)
	if err != nil {
		zap.L().Error("present log gift_dao pager failed", zap.Error(err))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrDB))
		return
	}

	userIds := make([]int64, 0, len(logs))
	for _, log := range logs {
		userIds = append(userIds, log.FromUserId)
	}
	userIds = utils.UniqueInt64Slice(userIds)
	userMap, err := userdao.MapByUserIds(userIds)
	if err != nil {
		zap.L().Error("present log user_dao MapByUserIds failed", zap.Error(err), zap.Int64s("userIds", userIds))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrDB))
		return
	}
	elemList := make([]*presentLogResponseElem, 0, len(logs))
	for _, log := range logs {
		elem := &presentLogResponseElem{
			UserId:   log.FromUserId,
			GiftId:   log.GiftId,
			GiftName: log.Name,
			SendAt:   utils.BackendToFrontendTime(log.SendAt),
		}
		user, ok := userMap[log.FromUserId]
		if ok {
			elem.UserName = user.Name
		}
		elemList = append(elemList, elem)
	}

	res := presentLogResponse{
		Total: total,
		List:  elemList,
	}
	controller.ResponseSuccessWithData(ctx, res)
}
