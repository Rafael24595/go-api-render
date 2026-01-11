package system

import (
	core_system "github.com/Rafael24595/go-api-core/src/commons/system"
)

const (
	SNAPSHOT_TOPIC_WEB_DATA core_system.TopicSnapshot = "snpsh_web"
)

var snapshotMeta = []core_system.SnapshotExtension{
	{
		Topic:       SNAPSHOT_TOPIC_WEB_DATA,
		Description: "Represents a snapshot of user web data.",
		CsvPath:     "./db/snapshot/web",
	},
}

func init() {
	core_system.ExtendMany(snapshotMeta...)
}
