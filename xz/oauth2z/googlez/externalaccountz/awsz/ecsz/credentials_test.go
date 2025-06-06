package ecsz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2/google"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/xz/oauth2z/googlez/externalaccountz/awsz/ecsz"
)

func TestCredentialsFromJSON(t *testing.T) {
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
		cred, err := ecsz.CredentialsFromJSON(
			context.Background(),
			jsonData,
			ecsz.WithCredentialsFromJSONOptionParams(google.CredentialsParams{Scopes: ecsz.DefaultTokenSourceConfigScopes}),
			ecsz.WithCredentialsFromJSONOptionTokenSourceConfigOptions(
				ecsz.WithTokenSourceOptionScopes(ecsz.DefaultTokenSourceConfigScopes),
				ecsz.WithTokenSourceOptionDefaultAwsRegion("ap-northeast-1"),
				ecsz.WithTokenSourceOptionAwsEcsMetadataEndpointHost("http://"+metadataServerMock.Listener.Addr().String()),
				ecsz.WithTokenSourceOptionAwsContainerCredentialsRelativeURI("/v2/credentials/00000000-0000-0000-0000-000000000000"),
			),
		)
		assertz.NoError(t, err)
		assertz.NotNil(t, cred)
	})

	t.Run("success,FederationID", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		_, err := ecsz.CredentialsFromJSON(context.Background(), jsonData)
		assertz.NoError(t, err)
	})

	t.Run("success,ServiceAccountImpersonation", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		_, err := ecsz.CredentialsFromJSON(context.Background(), jsonData)
		assertz.NoError(t, err)
	})

	t.Run("success,ServiceAccountImpersonation", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{`)
		_, err := ecsz.CredentialsFromJSON(context.Background(), jsonData)
		assertz.ErrorContains(t, err, `ecs.NewTokenSource error = NewTokenSource: TokenSourceConfigFromJSON: failed to unmarshal google workload identity federation config: json.Unmarshal: unexpected end of JSON input, google.CredentialsFromJSONWithParams: unexpected end of JSON input`)
	})
}
