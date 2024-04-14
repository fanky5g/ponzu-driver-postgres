package models

import (
	"fmt"
	"github.com/fanky5g/ponzu/entities"
)

const (
	ModelNameSpace = "ponzu"
)

func wrapModelNameSpace(name string) string {
	return fmt.Sprintf("%s_%s", ModelNameSpace, name)
}

type UserDocument struct {
	Document
	entities.User
}

type UserModel struct {
	*Model
	Document UserDocument
}

func (*UserModel) Name() string {
	return wrapModelNameSpace("users")
}

type ConfigDocument struct {
	Document
	entities.Config
}

type ConfigModel struct {
	*Model
	Document ConfigDocument
}

func (*ConfigModel) Name() string {
	return wrapModelNameSpace("config")
}

type AnalyticsMetricDocument struct {
	Document
	entities.AnalyticsMetric
}

type AnalyticsMetricModel struct {
	*Model
	Document AnalyticsMetricDocument
}

func (*AnalyticsMetricModel) Name() string {
	return wrapModelNameSpace("analytics_metrics")
}

type AnalyticsHTTPRequestMetadataDocument struct {
	Document
	entities.AnalyticsHTTPRequestMetadata
}

type AnalyticsHTTPRequestMetadataModel struct {
	Model
	Document AnalyticsHTTPRequestMetadataDocument
}

func (*AnalyticsHTTPRequestMetadataModel) Name() string {
	return wrapModelNameSpace("analytics_http_request_metadata")
}

type CredentialHashDocument struct {
	Document
	entities.CredentialHash
}

type CredentialHashModel struct {
	Model
	Document CredentialHashDocument
}

func (*CredentialHashModel) Name() string {
	return wrapModelNameSpace("credential_hashes")
}

type RecoveryKeyDocument struct {
	Document
	entities.RecoveryKey
}

type RecoveryKeyModel struct {
	Model
	Document RecoveryKeyDocument
}

func (*RecoveryKeyModel) Name() string {
	return wrapModelNameSpace("recovery_keys")
}

type SlugDocument struct {
	Document
	entities.Slug
}

type SlugModel struct {
	Model
	Document SlugDocument
}

func (*SlugModel) Name() string {
	return wrapModelNameSpace("slugs")
}

type UploadDocument struct {
	Document
	entities.FileUpload
}

type UploadModel struct {
	Model
	Document UploadDocument
}

func (*UploadModel) Name() string {
	return wrapModelNameSpace("uploads")
}

func init() {
	RegisterModel(
		new(UserModel),
		new(ConfigModel),
		new(AnalyticsMetricModel),
		new(AnalyticsHTTPRequestMetadataModel),
		new(CredentialHashModel),
		new(RecoveryKeyModel),
		new(SlugModel),
		new(UploadModel),
	)
}
