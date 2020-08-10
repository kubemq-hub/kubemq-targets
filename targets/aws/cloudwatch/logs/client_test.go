package logs

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	awsKey       string
	awsSecretKey string
	region       string
	token        string

	logGroupName   string
	logStreamName  string
	logGroupPrefix string

	policyName     string
	policyDocument string

	limit int64
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/aws/awsKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/awsSecretKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsSecretKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = fmt.Sprintf("%s", dat)
	t.token = ""
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/logs/logGroupName.txt")
	if err != nil {
		return nil, err
	}
	t.logGroupName = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/logs/logStreamName.txt")
	if err != nil {
		return nil, err
	}
	t.logStreamName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/logs/logGroupPrefix.txt")
	if err != nil {
		return nil, err
	}
	t.logGroupPrefix = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/logs/policyName.txt")
	if err != nil {
		return nil, err
	}
	t.policyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/logs/policyDocument.txt")
	if err != nil {
		return nil, err
	}
	t.policyDocument = fmt.Sprintf("%s", dat)

	t.limit = 10

	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init ",
			cfg: config.Spec{
				Name: "aws-cloudwatch-logs",
				Kind: "aws.cloudwatch.logs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "aws-cloudwatch-logs",
				Kind: "aws.cloudwatch.logs",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
					"region":  dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing key",
			cfg: config.Spec{
				Name: "aws-cloudwatch-logs",
				Kind: "aws.cloudwatch.logs",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}






//Log Group

func TestClient_CreateLogEventsGroup(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-logs",
		Kind: "aws.cloudwatch.logs",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_log_group").
				SetMetadataKeyValue("log_group_name", dat.logGroupName),
			wantErr: false,
		},
		{
			name: "invalid create - log group already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_log_group").
				SetMetadataKeyValue("log_group_name", dat.logGroupName),
			wantErr: true,
		},
		{
			name: "invalid create- already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_log_group").
				SetMetadataKeyValue("log_group_name", dat.logGroupName),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_DescribeLogEvents(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-logs",
		Kind: "aws.cloudwatch.logs",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid describe",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "describe_log_group").
				SetMetadataKeyValue("log_group_prefix", dat.logGroupPrefix),
			wantErr: false,
		}, {
			name: "invalid describe-missing log_group_prefix",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "describe_log_group"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_DeleteLogGroup(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-logs",
		Kind: "aws.cloudwatch.logs",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_log_group").
				SetMetadataKeyValue("log_group_name", dat.logGroupName),
			wantErr: false,
		}, {
			name: "invalid delete - log group does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_log_group").
				SetMetadataKeyValue("log_group_name", dat.logGroupName),
			wantErr: true,
		}, {
			name: "invalid delete-missing log_group_name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_log_group"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
