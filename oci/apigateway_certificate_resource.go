// Copyright (c) 2017, 2020, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	oci_apigateway "github.com/oracle/oci-go-sdk/v29/apigateway"
	oci_common "github.com/oracle/oci-go-sdk/v29/common"
)

func init() {
	RegisterResource("oci_apigateway_certificate", ApigatewayCertificateResource())
}

func ApigatewayCertificateResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: DefaultTimeout,
		Create:   createApigatewayCertificate,
		Read:     readApigatewayCertificate,
		Update:   updateApigatewayCertificate,
		Delete:   deleteApigatewayCertificate,
		Schema: map[string]*schema.Schema{
			// Required
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			// Optional
			"defined_tags": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: definedTagsDiffSuppressFunction,
				Elem:             schema.TypeString,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"freeform_tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"intermediate_certificates": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Computed
			"lifecycle_details": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_not_valid_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createApigatewayCertificate(d *schema.ResourceData, m interface{}) error {
	sync := &ApigatewayCertificateResourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).apiGatewayClient()
	sync.WorkRequestsClient = m.(*OracleClients).gatewayWorkRequestsClient

	return CreateResource(d, sync)
}

func readApigatewayCertificate(d *schema.ResourceData, m interface{}) error {
	sync := &ApigatewayCertificateResourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).apiGatewayClient()
	sync.WorkRequestsClient = m.(*OracleClients).gatewayWorkRequestsClient

	return ReadResource(sync)
}

func updateApigatewayCertificate(d *schema.ResourceData, m interface{}) error {
	sync := &ApigatewayCertificateResourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).apiGatewayClient()
	sync.WorkRequestsClient = m.(*OracleClients).gatewayWorkRequestsClient

	return UpdateResource(d, sync)
}

func deleteApigatewayCertificate(d *schema.ResourceData, m interface{}) error {
	sync := &ApigatewayCertificateResourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).apiGatewayClient()
	sync.DisableNotFoundRetries = true
	sync.WorkRequestsClient = m.(*OracleClients).gatewayWorkRequestsClient

	return DeleteResource(d, sync)
}

type ApigatewayCertificateResourceCrud struct {
	BaseCrud
	Client                 *oci_apigateway.ApiGatewayClient
	WorkRequestsClient     *oci_apigateway.WorkRequestsClient
	Res                    *oci_apigateway.Certificate
	DisableNotFoundRetries bool
}

func (s *ApigatewayCertificateResourceCrud) ID() string {
	return *s.Res.Id
}

func (s *ApigatewayCertificateResourceCrud) CreatedPending() []string {
	return []string{
		string(oci_apigateway.CertificateLifecycleStateCreating),
	}
}

func (s *ApigatewayCertificateResourceCrud) CreatedTarget() []string {
	return []string{
		string(oci_apigateway.CertificateLifecycleStateActive),
	}
}

func (s *ApigatewayCertificateResourceCrud) DeletedPending() []string {
	return []string{
		string(oci_apigateway.CertificateLifecycleStateDeleting),
	}
}

func (s *ApigatewayCertificateResourceCrud) DeletedTarget() []string {
	return []string{
		string(oci_apigateway.CertificateLifecycleStateDeleted),
	}
}

func (s *ApigatewayCertificateResourceCrud) Create() error {
	request := oci_apigateway.CreateCertificateRequest{}

	if certificate, ok := s.D.GetOkExists("certificate"); ok {
		tmp := certificate.(string)
		request.Certificate = &tmp
	}

	if compartmentId, ok := s.D.GetOkExists("compartment_id"); ok {
		tmp := compartmentId.(string)
		request.CompartmentId = &tmp
	}

	if definedTags, ok := s.D.GetOkExists("defined_tags"); ok {
		convertedDefinedTags, err := mapToDefinedTags(definedTags.(map[string]interface{}))
		if err != nil {
			return err
		}
		request.DefinedTags = convertedDefinedTags
	}

	if displayName, ok := s.D.GetOkExists("display_name"); ok {
		tmp := displayName.(string)
		request.DisplayName = &tmp
	}

	if freeformTags, ok := s.D.GetOkExists("freeform_tags"); ok {
		request.FreeformTags = objectMapToStringMap(freeformTags.(map[string]interface{}))
	}

	if intermediateCertificates, ok := s.D.GetOkExists("intermediate_certificates"); ok {
		tmp := intermediateCertificates.(string)
		request.IntermediateCertificates = &tmp
	}

	if privateKey, ok := s.D.GetOkExists("private_key"); ok {
		tmp := privateKey.(string)
		request.PrivateKey = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(s.DisableNotFoundRetries, "apigateway")

	response, err := s.Client.CreateCertificate(context.Background(), request)
	if err != nil {
		return err
	}

	workId := response.OpcWorkRequestId
	return s.getCertificateFromWorkRequest(workId, getRetryPolicy(s.DisableNotFoundRetries, "apigateway"), oci_apigateway.WorkRequestResourceActionTypeCreated, s.D.Timeout(schema.TimeoutCreate))
}

func (s *ApigatewayCertificateResourceCrud) getCertificateFromWorkRequest(workId *string, retryPolicy *oci_common.RetryPolicy,
	actionTypeEnum oci_apigateway.WorkRequestResourceActionTypeEnum, timeout time.Duration) error {

	// Wait until it finishes
	certificateId, err := certificateWaitForWorkRequest(workId, "certificate",
		actionTypeEnum, timeout, s.DisableNotFoundRetries, s.WorkRequestsClient)

	if err != nil {
		// Try to cancel the work request
		log.Printf("[DEBUG] creation failed, attempting to cancel the workrequest: %v for identifier: %v\n", workId, certificateId)
		_, cancelErr := s.WorkRequestsClient.CancelWorkRequest(context.Background(),
			oci_apigateway.CancelWorkRequestRequest{
				WorkRequestId: workId,
				RequestMetadata: oci_common.RequestMetadata{
					RetryPolicy: retryPolicy,
				},
			})
		if cancelErr != nil {
			log.Printf("[DEBUG] cleanup cancelWorkRequest failed with the error: %v\n", cancelErr)
		}
		return err
	}
	s.D.SetId(*certificateId)

	return s.Get()

}

func certificateWorkRequestShouldRetryFunc(timeout time.Duration) func(response oci_common.OCIOperationResponse) bool {
	startTime := time.Now()
	stopTime := startTime.Add(timeout)
	return func(response oci_common.OCIOperationResponse) bool {

		// Stop after timeout has elapsed
		if time.Now().After(stopTime) {
			return false
		}

		// Make sure we stop on default rules
		if shouldRetry(response, false, "apigateway", startTime) {
			return true
		}

		// Only stop if the time Finished is set
		if workRequestResponse, ok := response.Response.(oci_apigateway.GetWorkRequestResponse); ok {
			return workRequestResponse.TimeFinished == nil
		}
		return false
	}
}

func certificateWaitForWorkRequest(wId *string, entityType string, action oci_apigateway.WorkRequestResourceActionTypeEnum,
	timeout time.Duration, disableFoundRetries bool, client *oci_apigateway.WorkRequestsClient) (*string, error) {

	retryPolicy := getRetryPolicy(disableFoundRetries, "apigateway")
	retryPolicy.ShouldRetryOperation = certificateWorkRequestShouldRetryFunc(timeout)

	response := oci_apigateway.GetWorkRequestResponse{}
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(oci_apigateway.WorkRequestStatusInProgress),
			string(oci_apigateway.WorkRequestStatusAccepted),
			string(oci_apigateway.WorkRequestStatusCanceling),
		},
		Target: []string{
			string(oci_apigateway.WorkRequestStatusSucceeded),
			string(oci_apigateway.WorkRequestStatusFailed),
			string(oci_apigateway.WorkRequestStatusCanceled),
		},
		Refresh: func() (interface{}, string, error) {
			var err error
			response, err = client.GetWorkRequest(context.Background(),
				oci_apigateway.GetWorkRequestRequest{
					WorkRequestId: wId,
					RequestMetadata: oci_common.RequestMetadata{
						RetryPolicy: retryPolicy,
					},
				})
			wr := &response.WorkRequest
			return wr, string(wr.Status), err
		},
		Timeout: timeout,
	}
	if _, e := stateConf.WaitForState(); e != nil {
		return nil, e
	}

	var identifier *string
	// The work request response contains an array of objects that finished the operation
	for _, res := range response.Resources {
		if strings.Contains(strings.ToLower(*res.EntityType), entityType) {
			if res.ActionType == action {
				identifier = res.Identifier
				break
			}
		}
	}

	// The API Gateway workrequest may have failed, check for errors if identifier is not found or work failed or got cancelled
	if identifier == nil || response.Status == oci_apigateway.WorkRequestStatusFailed || response.Status == oci_apigateway.WorkRequestStatusCanceled {
		return nil, getErrorFromGatewayCertificateWorkRequest(client, wId, retryPolicy, entityType, action)
	}

	return identifier, nil
}

func getErrorFromGatewayCertificateWorkRequest(client *oci_apigateway.WorkRequestsClient, wId *string, retryPolicy *oci_common.RetryPolicy, entityType string, action oci_apigateway.WorkRequestResourceActionTypeEnum) error {
	response, err := client.ListWorkRequestErrors(context.Background(),
		oci_apigateway.ListWorkRequestErrorsRequest{
			WorkRequestId: wId,
			RequestMetadata: oci_common.RequestMetadata{
				RetryPolicy: retryPolicy,
			},
		})
	if err != nil {
		return err
	}

	allErrs := make([]string, 0)
	for _, wrkErr := range response.Items {
		allErrs = append(allErrs, *wrkErr.Message)
	}
	errorMessage := strings.Join(allErrs, "\n")

	workRequestErr := fmt.Errorf("work request did not succeed, workId: %s, entity: %s, action: %s. Message: %s", *wId, entityType, action, errorMessage)

	return workRequestErr
}

func (s *ApigatewayCertificateResourceCrud) Get() error {
	request := oci_apigateway.GetCertificateRequest{}

	tmp := s.D.Id()
	request.CertificateId = &tmp

	request.RequestMetadata.RetryPolicy = getRetryPolicy(s.DisableNotFoundRetries, "apigateway")

	response, err := s.Client.GetCertificate(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response.Certificate
	return nil
}

func (s *ApigatewayCertificateResourceCrud) Update() error {
	if compartment, ok := s.D.GetOkExists("compartment_id"); ok && s.D.HasChange("compartment_id") {
		oldRaw, newRaw := s.D.GetChange("compartment_id")
		if newRaw != "" && oldRaw != "" {
			err := s.updateCompartment(compartment)
			if err != nil {
				return err
			}
		}
	}
	request := oci_apigateway.UpdateCertificateRequest{}

	tmp := s.D.Id()
	request.CertificateId = &tmp

	if definedTags, ok := s.D.GetOkExists("defined_tags"); ok {
		convertedDefinedTags, err := mapToDefinedTags(definedTags.(map[string]interface{}))
		if err != nil {
			return err
		}
		request.DefinedTags = convertedDefinedTags
	}

	if displayName, ok := s.D.GetOkExists("display_name"); ok {
		tmp := displayName.(string)
		request.DisplayName = &tmp
	}

	if freeformTags, ok := s.D.GetOkExists("freeform_tags"); ok {
		request.FreeformTags = objectMapToStringMap(freeformTags.(map[string]interface{}))
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(s.DisableNotFoundRetries, "apigateway")

	response, err := s.Client.UpdateCertificate(context.Background(), request)
	if err != nil {
		return err
	}

	workId := response.OpcWorkRequestId
	return s.getCertificateFromWorkRequest(workId, getRetryPolicy(s.DisableNotFoundRetries, "apigateway"), oci_apigateway.WorkRequestResourceActionTypeUpdated, s.D.Timeout(schema.TimeoutUpdate))
}

func (s *ApigatewayCertificateResourceCrud) Delete() error {
	request := oci_apigateway.DeleteCertificateRequest{}

	tmp := s.D.Id()
	request.CertificateId = &tmp

	request.RequestMetadata.RetryPolicy = getRetryPolicy(s.DisableNotFoundRetries, "apigateway")

	_, err := s.Client.DeleteCertificate(context.Background(), request)
	return err
}

func (s *ApigatewayCertificateResourceCrud) SetData() error {
	if s.Res.Certificate != nil {
		s.D.Set("certificate", *s.Res.Certificate)
	}

	if s.Res.CompartmentId != nil {
		s.D.Set("compartment_id", *s.Res.CompartmentId)
	}

	if s.Res.DefinedTags != nil {
		s.D.Set("defined_tags", definedTagsToMap(s.Res.DefinedTags))
	}

	if s.Res.DisplayName != nil {
		s.D.Set("display_name", *s.Res.DisplayName)
	}

	s.D.Set("freeform_tags", s.Res.FreeformTags)

	if s.Res.IntermediateCertificates != nil {
		s.D.Set("intermediate_certificates", *s.Res.IntermediateCertificates)
	}

	if s.Res.LifecycleDetails != nil {
		s.D.Set("lifecycle_details", *s.Res.LifecycleDetails)
	}

	s.D.Set("state", s.Res.LifecycleState)

	s.D.Set("subject_names", s.Res.SubjectNames)

	if s.Res.TimeCreated != nil {
		s.D.Set("time_created", s.Res.TimeCreated.String())
	}

	if s.Res.TimeNotValidAfter != nil {
		s.D.Set("time_not_valid_after", s.Res.TimeNotValidAfter.String())
	}

	if s.Res.TimeUpdated != nil {
		s.D.Set("time_updated", s.Res.TimeUpdated.String())
	}

	return nil
}

func CertificateSummaryToMap(obj oci_apigateway.CertificateSummary) map[string]interface{} {
	result := map[string]interface{}{}

	if obj.CompartmentId != nil {
		result["compartment_id"] = string(*obj.CompartmentId)
	}

	if obj.DefinedTags != nil {
		result["defined_tags"] = definedTagsToMap(obj.DefinedTags)
	}

	if obj.DisplayName != nil {
		result["display_name"] = string(*obj.DisplayName)
	}

	result["freeform_tags"] = obj.FreeformTags

	if obj.Id != nil {
		result["id"] = string(*obj.Id)
	}

	if obj.LifecycleDetails != nil {
		result["lifecycle_details"] = string(*obj.LifecycleDetails)
	}

	result["state"] = string(obj.LifecycleState)

	result["subject_names"] = obj.SubjectNames

	if obj.TimeCreated != nil {
		result["time_created"] = obj.TimeCreated.String()
	}

	if obj.TimeNotValidAfter != nil {
		result["time_not_valid_after"] = obj.TimeNotValidAfter.String()
	}

	if obj.TimeUpdated != nil {
		result["time_updated"] = obj.TimeUpdated.String()
	}

	return result
}

func (s *ApigatewayCertificateResourceCrud) updateCompartment(compartment interface{}) error {
	changeCompartmentRequest := oci_apigateway.ChangeCertificateCompartmentRequest{}

	idTmp := s.D.Id()
	changeCompartmentRequest.CertificateId = &idTmp

	compartmentTmp := compartment.(string)
	changeCompartmentRequest.CompartmentId = &compartmentTmp

	changeCompartmentRequest.RequestMetadata.RetryPolicy = getRetryPolicy(s.DisableNotFoundRetries, "apigateway")

	_, err := s.Client.ChangeCertificateCompartment(context.Background(), changeCompartmentRequest)
	if err != nil {
		return err
	}
	return nil
}
