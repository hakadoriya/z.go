package ecsz

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
)

type credentialsFromJSONWithParamsConfig struct {
	params                   google.CredentialsParams
	tokenSourceConfigOptions []TokenSourceOption
}

type CredentialsFromJSONOption interface {
	apply(cfg *credentialsFromJSONWithParamsConfig)
}

// WithCredentialsFromJSONOptionParams sets the google.CredentialsParams for google.CredentialsFromJSONWithParams.
func WithCredentialsFromJSONOptionParams(params google.CredentialsParams) CredentialsFromJSONOption { //nolint:ireturn
	return CredentialsFromJSONOptionParams{params: params}
}

type CredentialsFromJSONOptionParams struct{ params google.CredentialsParams }

func (f CredentialsFromJSONOptionParams) apply(cfg *credentialsFromJSONWithParamsConfig) {
	cfg.params = f.params
}

// WithCredentialsFromJSONOptionTokenSourceConfigOptions sets the TokenSourceConfigOption for the credentials.
// This allows customization of the token source configuration when creating credentials from JSON.
func WithCredentialsFromJSONOptionTokenSourceConfigOptions(tokenSourceConfigOptions ...TokenSourceOption) CredentialsFromJSONOption { //nolint:ireturn
	return CredentialsFromJSONOptionTokenSourceConfigOption{tokenSourceConfigOptions: tokenSourceConfigOptions}
}

type CredentialsFromJSONOptionTokenSourceConfigOption struct {
	tokenSourceConfigOptions []TokenSourceOption
}

func (f CredentialsFromJSONOptionTokenSourceConfigOption) apply(cfg *credentialsFromJSONWithParamsConfig) {
	cfg.tokenSourceConfigOptions = f.tokenSourceConfigOptions
}

func CredentialsFromJSON(ctx context.Context, jsonData []byte, opts ...CredentialsFromJSONOption) (*google.Credentials, error) {
	cfg := &credentialsFromJSONWithParamsConfig{}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	tokenSource, errNewTokenSource := NewTokenSource(ctx, jsonData, cfg.tokenSourceConfigOptions...)
	if errNewTokenSource == nil {
		return &google.Credentials{TokenSource: tokenSource}, nil
	}

	cred, err := google.CredentialsFromJSONWithParams(ctx, jsonData, cfg.params)
	if err == nil {
		return cred, nil
	}

	return nil, fmt.Errorf("ecsz.NewTokenSource error = %w, google.CredentialsFromJSONWithParams: %w", errNewTokenSource, err)
}
