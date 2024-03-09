package types

import "github.com/pkg/errors"

var GameIdNotFound = errors.New("无法通过路径找到游戏ID")
var GamePathNotFound = errors.New("无法通过ID找到游戏路径")
var GameScreenshotsNotFound = errors.New("无法找到游戏截图")
