package models

import "github.com/pkg/errors"

var CannotMatchGameIDFromPathInDatabase = errors.New("无法通过游戏路径匹配到数据库中的游戏ID")
var GameNotFoundInDatabase = errors.New("无法在数据库中找到游戏（可能是提供的游戏id有问题）")
