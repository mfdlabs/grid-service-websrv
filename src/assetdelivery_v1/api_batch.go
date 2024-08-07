/*
Asset Delivery Api v1

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package assetdeliveryv1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)


// BatchAPIService BatchAPI service
type BatchAPIService service

type ApiV1AssetsBatchPostRequest struct {
	ctx context.Context
	ApiService *BatchAPIService
	robloxPlaceId *int64
	accept *string
	robloxBrowserAssetRequest *string
	assetRequestItems *[]RobloxWebAssetsBatchAssetRequestItem
}

func (r ApiV1AssetsBatchPostRequest) RobloxPlaceId(robloxPlaceId int64) ApiV1AssetsBatchPostRequest {
	r.robloxPlaceId = &robloxPlaceId
	return r
}

func (r ApiV1AssetsBatchPostRequest) Accept(accept string) ApiV1AssetsBatchPostRequest {
	r.accept = &accept
	return r
}

func (r ApiV1AssetsBatchPostRequest) RobloxBrowserAssetRequest(robloxBrowserAssetRequest string) ApiV1AssetsBatchPostRequest {
	r.robloxBrowserAssetRequest = &robloxBrowserAssetRequest
	return r
}

func (r ApiV1AssetsBatchPostRequest) AssetRequestItems(assetRequestItems []RobloxWebAssetsBatchAssetRequestItem) ApiV1AssetsBatchPostRequest {
	r.assetRequestItems = &assetRequestItems
	return r
}

func (r ApiV1AssetsBatchPostRequest) Execute() ([]RobloxWebAssetsIAssetResponseItem, *http.Response, error) {
	return r.ApiService.V1AssetsBatchPostExecute(r)
}

/*
V1AssetsBatchPost Method for V1AssetsBatchPost

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiV1AssetsBatchPostRequest
*/
func (a *BatchAPIService) V1AssetsBatchPost(ctx context.Context) ApiV1AssetsBatchPostRequest {
	return ApiV1AssetsBatchPostRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return []RobloxWebAssetsIAssetResponseItem
func (a *BatchAPIService) V1AssetsBatchPostExecute(r ApiV1AssetsBatchPostRequest) ([]RobloxWebAssetsIAssetResponseItem, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  []RobloxWebAssetsIAssetResponseItem
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BatchAPIService.V1AssetsBatchPost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/assets/batch"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.robloxPlaceId == nil {
		return localVarReturnValue, nil, reportError("robloxPlaceId is required and must be specified")
	}
	if r.accept == nil {
		return localVarReturnValue, nil, reportError("accept is required and must be specified")
	}
	if r.robloxBrowserAssetRequest == nil {
		return localVarReturnValue, nil, reportError("robloxBrowserAssetRequest is required and must be specified")
	}
	if r.assetRequestItems == nil {
		return localVarReturnValue, nil, reportError("assetRequestItems is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json", "text/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json", "text/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	parameterAddToHeaderOrQuery(localVarHeaderParams, "Roblox-Place-Id", r.robloxPlaceId, "")
	parameterAddToHeaderOrQuery(localVarHeaderParams, "Accept", r.accept, "")
	parameterAddToHeaderOrQuery(localVarHeaderParams, "Roblox-Browser-Asset-Request", r.robloxBrowserAssetRequest, "")
	// body params
	localVarPostBody = r.assetRequestItems
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
