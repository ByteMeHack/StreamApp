package uploader

type imageUploader struct {
	presigner  Presigner
	bucketName string
}

func NewImageUploader(presigner Presigner, bucketName string) *imageUploader {
	return &imageUploader{
		presigner:  presigner,
		bucketName: bucketName,
	}
}

func (u *imageUploader) PresignPutUrl(key string, lifetime int64) (string, error) {
	req, err := u.presigner.PutObject(u.bucketName, key, lifetime)
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (u *imageUploader) PresignGetUrl(key string, contentType string) (string, error) {
	req, err := u.presigner.GetObject(u.bucketName, key, contentType)
	if err != nil {
		return "", err
	}
	return req.URL, nil
}
