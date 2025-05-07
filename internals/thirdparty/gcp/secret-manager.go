package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"google.golang.org/api/iterator"
	option "google.golang.org/api/option"
)

type GcpSecretManager interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
	UpdateSecret(ctx context.Context, data any, secretName string) error
}

type GcpSMStore struct {
	Client                                      *secretmanager.Client
	projectID, secretPrefix, secretNameTemplate string
}

func NewGcpSMStore(ctx context.Context, opts *GcpConfig) (*GcpSMStore, error) {
	client, err := secretmanager.NewClient(context.TODO(), option.WithCredentialsJSON(opts.ServiceAccountKey))

	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingGcpSMClient, err)
	}
	return &GcpSMStore{
		Client:    client,
		projectID: opts.ProjectID,
	}, nil
}

func (gsm *GcpSMStore) getGcpSecRes(childRes string) string {
	return fmt.Sprintf(GCP_RES_TEMP+"/"+childRes, gsm.projectID)
}

func (gsm *GcpSMStore) WriteSecret(ctx context.Context, data any, secretName string) error {
	c, err := json.Marshal(data)
	if err != nil {
		return ErrCreatingGcpSecret
	}
	log.Println("Secret data:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ", string(c))
	// first create the secret itself
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf(GCP_RES_TEMP, gsm.projectID),
		SecretId: secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}
	secret, err := gsm.Client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		return err
	}
	// once the secret is created store it as the newest version
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: c,
		},
	}
	_, err = gsm.Client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return err
	}

	return nil
}

func (gsm *GcpSMStore) UpdateSecret(ctx context.Context, data any, secretName string) error {
	c, err := json.Marshal(data)
	if err != nil {
		return ErrUpdatingGcpSecret
	}

	secretResource := gsm.getGcpSecRes(fmt.Sprintf(GCP_SM_RES, secretName))

	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretResource,
		Payload: &secretmanagerpb.SecretPayload{
			Data: c,
		},
	}

	newVersion, err := gsm.Client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return err
	}
	newVersionName := newVersion.Name

	it := gsm.Client.ListSecretVersions(ctx, &secretmanagerpb.ListSecretVersionsRequest{
		Parent: secretResource,
	})
	for {
		secretVersion, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if secretVersion.Name != newVersionName && secretVersion.State == secretmanagerpb.SecretVersion_ENABLED {
			_, err := gsm.Client.DisableSecretVersion(ctx, &secretmanagerpb.DisableSecretVersionRequest{
				Name: secretVersion.Name,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (gsm *GcpSMStore) ReadSecret(ctx context.Context, secretId string) (any, error) {
	secResId := gsm.getGcpSecRes(fmt.Sprintf(GCP_SM_RES, secretId)) + "/versions/latest"
	sec, err := gsm.Client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: secResId})
	if err != nil {
		return nil, dmerrors.DMError(ErrReadingGcpSecret, err)
	}

	return sec.Payload.Data, nil
}
func (gsm *GcpSMStore) DeleteSecret(ctx context.Context, secretName string) error {
	secResId := gsm.getGcpSecRes(fmt.Sprintf(GCP_SM_RES, secretName))
	err := gsm.Client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{Name: secResId})
	if err != nil {
		return dmerrors.DMError(ErrDeletingGcpSecret, err)
	}
	return nil
}

func (gsm *GcpSMStore) Close() {
	if gsm.Client != nil {
		gsm.Client.Close()
	}
}
