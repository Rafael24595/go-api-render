package topic_snapshot

import (
	core_topic_snapshot "github.com/Rafael24595/go-api-core/src/commons/system/topic/snapshot"
	topic_repository "github.com/Rafael24595/go-api-render/src/commons/system/topic/repository"
)

const (
	TOPIC_WEB_DATA core_topic_snapshot.TopicSnapshot = "snpsh_web"
)

var meta = []core_topic_snapshot.Extension{
	{
		Topic:       TOPIC_WEB_DATA,
		Description: "Represents a snapshot of user web data.",
		CsvPath:     "./db/snapshot/web",
		Repository:  topic_repository.TOPIC_WEB_DATA,
	},
}

func init() {
	core_topic_snapshot.ExtendMany(meta...)
}
