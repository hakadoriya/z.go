package ecsz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/xz/oauth2z/googlez/externalaccountz/awsz/ecsz"
)

func TestTokenSourceConfigFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("success,AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"AccessKeyId":"TestingAccessKeyId","SecretAccessKey":"TestingSecretAccessKey","Token":"TestingToken"}`))
		}))

		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecsz.NewTokenSource(
			context.Background(),
			jsonData,
			ecsz.WithTokenSourceOptionScopes(ecsz.DefaultTokenSourceConfigScopes),
			ecsz.WithTokenSourceOptionHTTPClient(http.DefaultClient),
			ecsz.WithTokenSourceOptionDefaultAwsRegion("ap-northeast-1"),
			ecsz.WithTokenSourceOptionAwsEcsMetadataEndpointHost("http://"+metadataServerMock.Listener.Addr().String()),
			ecsz.WithTokenSourceOptionAwsContainerCredentialsRelativeURI("/v2/credentials/00000000-0000-0000-0000-000000000000"),
		)
		assertz.NotNil(t, ts)
		assertz.NoError(t, err)
		tok, err := ts.Token()
		if err != nil {
			t.Logf("📝: ts.Token: %v", err)
		}

		// error
		assertz.Error(t, err) // can't get token because credentials and audience is invalid.
		assertz.Nil(t, tok)
	})

	t.Run("failure,ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet,FederationID", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecsz.NewTokenSource(context.Background(), jsonData)
		assertz.ErrorIs(t, err, ecsz.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet)
		assertz.Nil(t, ts)
	})

	t.Run("failure,ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet,ServiceAccountImpersonation", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecsz.NewTokenSource(context.Background(), jsonData)
		assertz.ErrorIs(t, err, ecsz.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet)
		assertz.Nil(t, ts)
	})

	t.Run("failure,json.Unmarshal", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
`)
		ts, err := ecsz.NewTokenSource(context.Background(), jsonData)
		assertz.ErrorContains(t, err, `failed to unmarshal google workload identity federation config: json.Unmarshal: `)
		assertz.Nil(t, ts)
	})

	t.Run("failure,audience", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecsz.NewTokenSource(context.Background(), jsonData, ecsz.WithTokenSourceOptionAwsContainerCredentialsRelativeURI("/v2/credentials/00000000-0000-0000-0000-000000000000"))
		assertz.ErrorContains(t, err, `oauth2/google/externalaccount: Audience must be set`)
		assertz.Nil(t, ts)
	})
}
