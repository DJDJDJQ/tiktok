package service

import (
	"go.mod/model"
	"go.mod/pkg"
)

// 关注操作
func Follow(userId int64, toUserId int64) pkg.ErrNo {
	if code := model.CreateFollow(userId, toUserId); code == -1 {
		return pkg.DataBaseErr.WithMessage("Follow action failed")
	}
	// 更新userId关注数、toUserId粉丝数
	model.UpdataUserFollowCount(userId, 1)
	model.UpdataUserFollowerCount(toUserId, 1)

	return pkg.Success
}

// 取消关注操作
func CancelFollow(temp model.Follow, userId int64, toUserId int64) pkg.ErrNo {
	if code := model.DeleteFollow(temp); code == -1 {
		return pkg.DataBaseErr.WithMessage("Cancel follow action failed")
	}
	// 更新userId关注数、toUserId粉丝数
	model.UpdataUserFollowCount(userId, -1)
	model.UpdataUserFollowerCount(toUserId, -1)

	return pkg.Success
}
