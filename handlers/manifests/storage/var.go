package storage

const (
	saTokenVolumeName            = "bound-sa-token"
	saTokenExpiration      int64 = 3600
	saTokenVolumeMountPath       = "/var/run/secrets/storage/serviceaccount"

	ServiceAccountTokenFilePath = saTokenVolumeMountPath + "/token"

	secretDirectory  = "/etc/storage/secrets"
	storageTLSVolume = "storage-tls"
	caDirectory      = "/etc/storage/ca"

	tokenAuthConfigVolumeName = "token-auth-config"
	tokenAuthConfigDirectory  = "/etc/storage/token-auth"

	awsDefaultAudience = "sts.amazonaws.com"
)
