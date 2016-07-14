package dedicated_node 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Backups struct {

	/*Path - Descr: Path within the above bucket to which backups will be uploaded Default: 
*/
	Path interface{} `yaml:"path,omitempty"`

	/*EndpointUrl - Descr: HTTP(S) endpoint of the S3-compatible blob store that backups will be uploaded to Default: 
*/
	EndpointUrl interface{} `yaml:"endpoint_url,omitempty"`

	/*RestoreAvailable - Descr: Makes the backup restore binary available Default: true
*/
	RestoreAvailable interface{} `yaml:"restore_available,omitempty"`

	/*AccessKeyId - Descr: Access Key ID for the S3-compatible blob store that backups will be uploaded to Default: 
*/
	AccessKeyId interface{} `yaml:"access_key_id,omitempty"`

	/*SecretAccessKey - Descr: Secret Access Key for the S3-compatible blob store that backups will be uploaded to Default: 
*/
	SecretAccessKey interface{} `yaml:"secret_access_key,omitempty"`

	/*BucketName - Descr: Name of the bucket into which backups will be uploaded Default: 
*/
	BucketName interface{} `yaml:"bucket_name,omitempty"`

	/*BgSaveTimeout - Descr: Timeout in seconds for Redis background save to complete when backing up instance Default: 3600
*/
	BgSaveTimeout interface{} `yaml:"bg_save_timeout,omitempty"`

	/*BackupTmpDir - Descr: Temporary directory to use for backups. MUST be on same device (persistent disk) as redis data Default: /var/vcap/store/tmp_backup
*/
	BackupTmpDir interface{} `yaml:"backup_tmp_dir,omitempty"`

}