package cfgargs

// SyncRemoteCfg 同步远程配置中心的配置信息
func (s *SrvConfig) SyncRemoteCfg(key string) string {
	s.RemoteCfg[key] = "unreachable"
	return "unreachable"
}
