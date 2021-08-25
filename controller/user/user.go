package user

import (
	"boframe/controller"
	"boframe/dao/giftdao"
	"boframe/dao/userdao"
	"boframe/models/giftmodels"
	"boframe/models/usermodels"
	"boframe/pkg/errcode"
	"boframe/rediselem"
	"boframe/settings/mongoI"
	"boframe/settings/redis"
	"cc"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type sentPresentRequest struct {
	From   int64 `json:"from" `
	To     int64 `json:"to"`
	GiftId int64 `json:"gift_id"`
}

func SendPresent(ctx *cc.Context) {
	requestParam := new(sentPresentRequest)

	decoder := json.NewDecoder(ctx.Req.Body)
	if err := decoder.Decode(requestParam); err != nil {
		zap.L().Info("Send Present param error.", zap.Error(err))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrInvalidParam))
		return
	}

	if requestParam.From <= 0 || requestParam.To <= 0 {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrSGiftSenderOrReceiver))
		return
	}

	if requestParam.From == requestParam.To {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrSGiftSendToSelf))
		return
	}

	if requestParam.GiftId <= 0 {
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrSGiftInvalidGift))
		return
	}

	// =======================================================================================================
	// 用户信息先查缓存，再插数据库
	userIdStrs := make([]string, 0)
	userIds := []int64{
		requestParam.From,
		requestParam.To,
	}
	userIdMap := make(map[string]int64, 2)
	for _, id := range userIds {
		sid := strconv.FormatInt(id, 10)
		userIdMap[sid] = id
		userIdStrs = append(userIdStrs, sid)
	}
	uKeys := rediselem.BuildKeys(rediselem.UserPrefix, userIdStrs...)
	cmd := redis.Cli.MGet(context.Background(), uKeys...)
	cacheMap := make(map[int64]*usermodels.User, 2)
	toSearchUserId := make([]int64, 0, 2)
	if err := cmd.Err(); err != nil {
		zap.L().Error("Send Present user from Redis error.", zap.Error(err), zap.Strings("uKeys", uKeys))
		toSearchUserId = userIds
	} else {
		for idx, cache := range cmd.Val() {
			if cache == nil {
				toSearchUserId = append(toSearchUserId, userIds[idx])
				continue
			}

			var user = usermodels.User{}
			cacheS := cache.(string)
			// TODO: byte copy
			// 转换不成功的话，从数据库取
			if err := json.Unmarshal([]byte(cacheS), &user); err != nil {
				zap.L().Error("Send Present user from Redis error.", zap.Error(err), zap.String("cache", cacheS))
				toSearchUserId = userIds
			} else {
				cacheMap[userIds[idx]] = &user
			}
		}
	}

	usersMap, err := userdao.MapByUserIds(toSearchUserId)
	if err != nil {
		zap.L().Error("Send Present userdao MapByUserIds failed", zap.Error(err), zap.Int64s("userIds", userIds))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
		return
	}
	toSetUserCache := make(map[string]string, len(toSearchUserId))
	for _, uid := range toSearchUserId {
		user, ok := usersMap[uid]
		if !ok {
			controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrSGiftSenderOrReceiver))
			return
		}
		bs, _ := json.Marshal(user)
		k := rediselem.BuildKey(rediselem.UserPrefix, strconv.FormatInt(uid, 10))
		// TODO: byte copy
		toSetUserCache[k] = string(bs)
	}

	if len(toSetUserCache) > 0 {
		redis.Cli.MSet(context.Background(), toSetUserCache)
	}

	// =======================================================================================================
	// 查询礼物信息
	var gift = &giftmodels.Gift{}
	cacheStatus := 0 // 0-正常， 1-重刷
	gKey := rediselem.BuildKey(rediselem.GiftPrefix, strconv.FormatInt(requestParam.GiftId, 10))
	cache, err := redis.Cli.Get(context.Background(), gKey).Result()
	if err != nil {
		cacheStatus = 1
		if !redis.IsErrNil(err) {
			zap.L().Error("Send Present get gift cache failed", zap.Error(err), zap.String("gKey", gKey))
		}
	} else {
		if err = json.Unmarshal([]byte(cache), gift); err != nil {
			cacheStatus = 1
		}
	}

	if cacheStatus == 1 {
		gift, err = giftdao.OneByGiftId(requestParam.GiftId)
		if mongoI.IsErrNoDocuments(err) {
			controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrSGiftInvalidGift))
			return
		}
		if err != nil {
			zap.L().Error("Send Present gift_dao OneByGiftId failed", zap.Error(err), zap.Int64("gift_id", requestParam.GiftId))
			controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
			return
		}
		gs, _ := json.Marshal(gift)
		k := rediselem.BuildKey(rediselem.GiftPrefix, strconv.FormatInt(requestParam.GiftId, 10))
		err = redis.Cli.Set(context.Background(), k, string(gs), rediselem.NoExpiration).Err()
		if err != nil {
			zap.L().Error("Send Present set gift cache failed", zap.Error(err), zap.String("key", k),
				zap.String("cache", string(gs)))
		}
	}
	// =======================================================================================================

	log := giftmodels.GiftLog{
		Id:         primitive.NewObjectID(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		GiftId:     requestParam.GiftId,
		Name:       gift.Name,
		Score:      gift.Score,
		FromUserId: requestParam.From,
		ToUserId:   requestParam.To,
		SendAt:     time.Now(),
	}
	_, err = giftdao.InsertOne(&log)
	if err != nil {
		zap.L().Error("Send Present gift_dao insert failed", zap.Error(err), zap.Any("param", log))
		controller.ResponseErrorCode(ctx, errcode.Build(errcode.ErrServer))
		return
	}

	// 处理redis的score
	err = redis.Cli.ZIncrBy(context.Background(), rediselem.RankingZSetKey, float64(gift.Score),
		strconv.FormatInt(requestParam.To, 10)).Err()
	if err != nil {
		zap.L().Error("Send Present incr score failed", zap.Error(err), zap.Int64("score", gift.Score),
			zap.Int64("anchor_id", requestParam.To))
	}

	// 返回结果
	controller.ResponseSuccess(ctx)
}
