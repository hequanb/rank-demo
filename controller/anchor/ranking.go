package anchor

import (
	"boframe/controller"
	"boframe/dao/userdao"
	"boframe/models"
	"boframe/models/usermodels"
	"boframe/pkg/errcode"
	"boframe/rediselem"
	"boframe/settings/redis"
	"boframe/utils"
	"cc"
	"context"
	"encoding/json"
	redis2 "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

type rankingResponse struct {
	Total int64                  `json:"total"`
	List  []*rankingResponseElem `json:"list"`
}

type rankingResponseElem struct {
	UserId   int64   `json:"user_id,string"`
	UserName string  `json:"user_name"`
	Rank     int64   `json:"rank"`
	Score    float64 `json:"score"`
}

func Ranking(ctx *cc.Context) {
	query := ctx.Req.URL.Query()

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

	if page <= 0 || limit <= 0 {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidPageInfo))
		return
	}

	total, err := redis.Cli.ZCard(context.Background(), rediselem.RankingZSetKey).Result()
	if err != nil {
		zap.L().Error("ranking zcard failed", zap.Error(err))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
		return
	}

	offset := utils.CalPageOffset(page, limit)
	args := redis2.ZRangeArgs{
		Key:     rediselem.RankingZSetKey,
		Start:   "-inf",
		Stop:    "+inf",
		ByScore: true,
		ByLex:   false,
		Rev:     true,
		Offset:  offset,
		Count:   limit,
	}
	cmd := redis.Cli.ZRangeArgsWithScores(context.Background(), args)
	zs, err := cmd.Result()
	if err != nil {
		zap.L().Error("ranking zrange failed", zap.Error(err))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
		return
	}
	rank := offset + 1
	elemList := make([]*rankingResponseElem, 0, len(zs))
	userIds := make([]int64, 0)
	for _, z := range zs {
		uid, _ := strconv.ParseInt((z.Member).(string), 10, 64)
		elem := rankingResponseElem{
			UserId:   uid,
			UserName: "",
			Rank:     rank,
			Score:    z.Score,
		}
		rank++
		userIds = append(userIds, uid)
		elemList = append(elemList, &elem)
	}
	if len(elemList) <= 0 {
		resp := rankingResponse{
			Total: total,
			List:  elemList,
		}
		controller.ResponseSuccessWithData(ctx, resp)
		return
	}

	userIds = utils.UniqueInt64Slice(userIds)
	uKeys := make([]string, 0, len(userIds))
	if len(userIds) > 0 {
		uKeys = rediselem.BuildKeysInt64(rediselem.UserPrefix, userIds...)
	}

	toSearchUids := make([]int64, 0)
	cacheStatus := 0	// 0-正常，1-重刷缓存
	caches, err := redis.Cli.MGet(context.Background(), uKeys...).Result()
	if err != nil {
		zap.L().Error("ranking mget user info failed", zap.Error(err), zap.Strings("uKeys", uKeys))
		cacheStatus = 1
		toSearchUids = userIds
	}
	// 可能存在nil的情况
	cacheUserMap := make(map[int64]*usermodels.User, len(caches))
	for idx, cache := range caches {
		if cache == nil {
			cacheStatus = 1
			toSearchUids = append(toSearchUids, userIds[idx])
			continue
		}

		var temp = &usermodels.User{}
		_ = json.Unmarshal([]byte(cache.(string)), temp)
		cacheUserMap[userIds[idx]] = temp
	}

	userMap := make(map[int64]*usermodels.User, len(toSearchUids))
	if cacheStatus != 0 {
		userMap, err = userdao.MapByUserIds(toSearchUids)
		if err != nil {
			zap.L().Error("anchor ranking user_dao MapByUserIds failed", zap.Error(err), zap.Int64s("userIds", userIds))
			controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrDB))
			return
		}

		toSetUserCache := make(map[string]string, len(toSearchUids))
		for _, uid := range toSearchUids {
			user, ok := userMap[uid]
			k := rediselem.BuildKey(rediselem.UserPrefix, strconv.FormatInt(uid, 10))

			if !ok {
				toSetUserCache[k] = models.EmptyString
				continue
			}
			bs, _ := json.Marshal(user)
			toSetUserCache[k] = string(bs)
		}

		if len(toSetUserCache) > 0 {
			redis.Cli.MSet(context.Background(), toSetUserCache)
		}
	}

	for _, elem := range elemList {
		user, ok := userMap[elem.UserId]
		if ok {
			elem.UserName = user.Name
			continue
		}
		user, ok = cacheUserMap[elem.UserId]
		if !ok || user== nil {
			zap.L().Warn("anchor ranking can not find user info", zap.Int64("uid", elem.UserId))
			continue
		}
		elem.UserName = user.Name
	}

	resp := rankingResponse{
		Total: total,
		List:  elemList,
	}
	controller.ResponseSuccessWithData(ctx, resp)
}
