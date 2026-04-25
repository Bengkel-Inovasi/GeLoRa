package adaptersoutboundclientappself

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundclient "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/client"
)

type restApiImpl struct {
	baseURL      string
	accessToken  string
	refreshToken string
	httpClient   *http.Client
}

func New(baseURL string) portsoutboundclient.Appself {
	return &restApiImpl{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (r *restApiImpl) SetAccessToken(accessToken string) {
	r.accessToken = accessToken
}

func (r *restApiImpl) SetRefreshToken(refreshToken string) {
	r.refreshToken = refreshToken
}

func (r *restApiImpl) PostAuthSignIn(ctx context.Context, username string, password string) (*domainmodel.AppselfAuthResponse, error) {
	resp, err := r.do(ctx, http.MethodPost, "/auth/sign-in", map[string]string{
		"username": username,
		"password": password,
	}, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PostAuthSignUp(ctx context.Context, name string, username string, password string) (*domainmodel.AppselfAuthResponse, error) {
	resp, err := r.do(ctx, http.MethodPost, "/auth/sign-up", map[string]string{
		"name":     name,
		"username": username,
		"password": password,
	}, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PostAuthRefresh(ctx context.Context, refreshToken string) (*domainmodel.AppselfAuthResponse, error) {
	resp, err := r.do(ctx, http.MethodPost, "/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	}, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) GetUserInfoMe(ctx context.Context) (*domainmodel.AppselfUserInfo, error) {
	resp, err := r.do(ctx, http.MethodGet, "/users/me", nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfUserInfo
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PatchUserInfoMe(ctx context.Context, name *string, bio *string) error {
	resp, err := r.do(ctx, http.MethodPatch, "/users/me", map[string]*string{"name": name, "bio": bio}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) PutUserPasswordMe(ctx context.Context, oldPassword string, newPassword string) error {
	resp, err := r.do(ctx, http.MethodPut, "/users/me/password", map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) GetUsersList(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (*domainmodel.AppselfUsersList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	if cursorId != nil {
		q.Set("cursor_id", strconv.FormatInt(*cursorId, 10))
	}
	if search != nil {
		q.Set("search", *search)
	}
	if role != nil {
		q.Set("role", string(*role))
	}

	resp, err := r.do(ctx, http.MethodGet, "/users?"+q.Encode(), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfUsersList
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) GetUserInfoById(ctx context.Context, id int64) (*domainmodel.AppselfUserInfo, error) {
	resp, err := r.do(ctx, http.MethodGet, "/users/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfUserInfo
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PatchUserInfoById(ctx context.Context, id int64, name *string, bio *string) error {
	resp, err := r.do(ctx, http.MethodPatch, "/users/"+strconv.FormatInt(id, 10), map[string]*string{"name": name, "bio": bio}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) PutUserPasswordById(ctx context.Context, id int64, newPassword string) error {
	resp, err := r.do(ctx, http.MethodPut, "/users/"+strconv.FormatInt(id, 10)+"/password", map[string]string{
		"new_password": newPassword,
	}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) PutUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) error {
	resp, err := r.do(ctx, http.MethodPut, "/users/"+strconv.FormatInt(id, 10)+"/role", map[string]string{
		"role": string(role),
	}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) DeleteUserById(ctx context.Context, id int64) error {
	resp, err := r.do(ctx, http.MethodDelete, "/users/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) GetNodesList(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (*domainmodel.AppselfNodesList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	if cursorId != nil {
		q.Set("cursor_id", strconv.FormatInt(*cursorId, 10))
	}
	if search != nil {
		q.Set("search", *search)
	}
	if isValidated != nil {
		q.Set("is_validated", strconv.FormatBool(*isValidated))
	}

	resp, err := r.do(ctx, http.MethodGet, "/nodes?"+q.Encode(), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfNodesList
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) GetNodeInfoById(ctx context.Context, id int64) (*domainmodel.AppselfNodeInfo, error) {
	resp, err := r.do(ctx, http.MethodGet, "/nodes/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfNodeInfo
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PatchNodeInfoById(ctx context.Context, id int64, name *string, description *string) error {
	resp, err := r.do(ctx, http.MethodPatch, "/nodes/"+strconv.FormatInt(id, 10), map[string]*string{"name": name, "description": description}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) PutNodeValidate(ctx context.Context, id int64, validatedBy int64) error {
	resp, err := r.do(ctx, http.MethodPut, "/nodes/"+strconv.FormatInt(id, 10)+"/validate", map[string]int64{
		"validated_by": validatedBy,
	}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) DeleteNodeById(ctx context.Context, id int64) error {
	resp, err := r.do(ctx, http.MethodDelete, "/nodes/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) PostSession(ctx context.Context, userId int64, nodeId int64) (*domainmodel.AppselfPostSessionResponse, error) {
	resp, err := r.do(ctx, http.MethodPost, "/sessions", map[string]int64{
		"user_id": userId,
		"node_id": nodeId,
	}, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfPostSessionResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) GetSessionsList(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (*domainmodel.AppselfSessionsList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	if cursorId != nil {
		q.Set("cursor_id", strconv.FormatInt(*cursorId, 10))
	}
	if userId != nil {
		q.Set("user_id", strconv.FormatInt(*userId, 10))
	}
	if nodeId != nil {
		q.Set("node_id", strconv.FormatInt(*nodeId, 10))
	}
	if active != nil {
		q.Set("active", strconv.FormatBool(*active))
	}

	resp, err := r.do(ctx, http.MethodGet, "/sessions?"+q.Encode(), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfSessionsList
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) GetSessionById(ctx context.Context, id int64) (*domainmodel.AppselfSessionInfo, error) {
	resp, err := r.do(ctx, http.MethodGet, "/sessions/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfSessionInfo
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) PutSessionEnd(ctx context.Context, id *int64, userId *int64, nodeId *int64) error {
	resp, err := r.do(ctx, http.MethodPut, "/sessions/end", map[string]*int64{"id": id, "user_id": userId, "node_id": nodeId}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) DeleteSessionById(ctx context.Context, id int64) error {
	resp, err := r.do(ctx, http.MethodDelete, "/sessions/"+strconv.FormatInt(id, 10), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkError(resp)
}

func (r *restApiImpl) GetRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (*domainmodel.AppselfRecordsList, error) {
	q := url.Values{}
	if sessionId != nil {
		q.Set("session_id", strconv.FormatInt(*sessionId, 10))
	}
	if userId != nil {
		q.Set("user_id", strconv.FormatInt(*userId, 10))
	}
	if nodeId != nil {
		q.Set("node_id", strconv.FormatInt(*nodeId, 10))
	}
	if startTime != nil {
		q.Set("start_time", startTime.Format(time.RFC3339))
	}
	if endTime != nil {
		q.Set("end_time", endTime.Format(time.RFC3339))
	}

	resp, err := r.do(ctx, http.MethodGet, "/records?"+q.Encode(), nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkError(resp); err != nil {
		return nil, err
	}

	var out domainmodel.AppselfRecordsList
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *restApiImpl) do(ctx context.Context, method, path string, body any, withAuth bool) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, r.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if withAuth {
		req.Header.Set("Authorization", "Bearer "+r.accessToken)
	}

	return r.httpClient.Do(req)
}

func checkError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	var e struct {
		Error string `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&e)
	if e.Error == "" {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return errors.New(e.Error)
}
