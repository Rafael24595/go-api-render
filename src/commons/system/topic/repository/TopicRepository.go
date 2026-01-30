package topic_repository

import (
	core_topic_repository "github.com/Rafael24595/go-api-core/src/commons/system/topic/repository"
)

const (
	TOPIC_WEB_DATA core_topic_repository.TopicRepository = "rep_web"
)

var meta = []core_topic_repository.Extension{
	{
		Topic:       TOPIC_WEB_DATA,
		Description: "Represents the repository of user web data.",
	},
}

func init() {
	core_topic_repository.ExtendMany(meta...)
}
