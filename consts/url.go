package consts

const (
	ReleaseChinaUrlPrefix  = "https://storage.flutter-io.cn/"
	ReleaseUrlPrefix       = "https://storage.googleapis.com/"
	ReleasePath            = "flutter_infra_release/releases/"

	ReleaseJsonFileName    = "releases_%s.json"
	ReleaseChinaJsonUrlFmt = ReleaseChinaUrlPrefix + ReleasePath + ReleaseJsonFileName
	ReleaseJsonUrlFmt      = ReleaseUrlPrefix + ReleasePath + ReleaseJsonFileName

	ReleaseDownloadUrlPrefix  = ReleaseUrlPrefix + ReleasePath
)
