package documents

func applyBucketDefaults(bucket *Bucket) {
	if len(bucket.Versioning) == 0 {
		bucket.Versioning = VersioningDisabled
	}
}
