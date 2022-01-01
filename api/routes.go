package api

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/KaratsubaLabs/takoyaki-backend/db"
	"github.com/KaratsubaLabs/takoyaki-backend/vps"
)

type ContextKey string

func (c ContextKey) String() string {
	return "takoyaki:contextKey:" + string(c)
}

var (
	ContextKeyUserID     = ContextKey("userid")
	ContextKeyParsedBody = ContextKey("parsedbody")
)

type CustomHandler = func(http.ResponseWriter, *http.Request) error

type routeInfo struct {
	route   string
	methods map[string]methodEndpoint
}

type methodEndpoint struct {
	authRoute  bool
	bodySchema interface{}
	handlerFn  CustomHandler
}

func (info routeInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// find which method handler to use
	endpoint, ok := info.methods[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var handlerWithMiddleware http.Handler = ErrorMiddleware(endpoint.handlerFn)

	// validate + parse body (if applicable)
	if endpoint.bodySchema != nil {
		handlerWithMiddleware = ValidationMiddleware(handlerWithMiddleware)
		handlerWithMiddleware = ParseBodyJSONMiddleware(endpoint.bodySchema, handlerWithMiddleware)
	}

	// restrict auth (if applicable)
	if endpoint.authRoute {
		handlerWithMiddleware = AuthMiddleware(handlerWithMiddleware)
	}

	// delegate to handler
	handlerWithMiddleware.ServeHTTP(w, r)

}

var routeSchema = []routeInfo{
	{
		route: "/ping",
		methods: map[string]methodEndpoint{
			"POST": methodEndpoint{
				authRoute: false,
				handlerFn: pingHandler,
			},
		},
	},
	{
		route: "/register",
		methods: map[string]methodEndpoint{
			"POST": methodEndpoint{
				authRoute:  false,
				bodySchema: &registerRequest{},
				handlerFn:  registerHandler,
			},
		},
	},
	{
		route: "/login",
		methods: map[string]methodEndpoint{
			"POST": methodEndpoint{
				authRoute:  false,
				bodySchema: &loginRequest{},
				handlerFn:  loginHandler,
			},
		},
	},
	{
		route: "/vps",
		methods: map[string]methodEndpoint{
			"GET": methodEndpoint{
				authRoute:  true,
				bodySchema: &vpsInfoRequest{},
				handlerFn:  vpsInfoHandler,
			},
			// TODO this technically should be under request POST endpoint?
			"POST": methodEndpoint{
				authRoute:  true,
				bodySchema: &vpsCreateRequest{},
				handlerFn:  vpsCreateHandler,
			},
			"DELETE": methodEndpoint{
				authRoute:  true,
				bodySchema: &vpsDeleteRequest{},
				handlerFn:  vpsDeleteHandler,
			},
		},
	},
	{
		route: "/vps/action",
		methods: map[string]methodEndpoint{
			"POST": methodEndpoint{
				authRoute:  true,
				bodySchema: &vpsActionRequest{},
				handlerFn:  vpsActionHandler,
			},
		},
	},
	{
		route: "/request",
		methods: map[string]methodEndpoint{
			"GET": methodEndpoint{
				authRoute: true,
				handlerFn: requestListHandler,
			},
			"DELETE": methodEndpoint{
				authRoute:  true,
				bodySchema: &requestDeleteRequest{},
				handlerFn:  requestDeleteHandler,
			},
		},
	},
}

// ping endpoint for debug purposes
func pingHandler(w http.ResponseWriter, r *http.Request) error {

	return nil
}

type registerRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}
type registerResponse struct {
	Token string `json:"token"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*registerRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	// (possibly have conn connection be part of the context)
	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// make sure user name and email are not already taken
	registered, err := db.UserCheckRegistered(conn, parsedBody.Email)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if registered {
		// possibly return more generic error to reduce info leak
		return HTTPStatusError{http.StatusConflict, errors.New("username or email already taken")}
	}

	// hash pass
	hashed, err := bcrypt.GenerateFromPassword([]byte(parsedBody.Password), bcrypt.MinCost)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	newUser := db.User{
		Email:    parsedBody.Email,
		Password: string(hashed),
	}
	userID, err := db.UserRegister(conn, &newUser)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	token, err := GenerateToken(userID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	json.NewEncoder(w).Encode(registerResponse{Token: token})

	return nil
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*loginRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	userID, err := db.UserCheckCreds(conn, parsedBody.Email, parsedBody.Password)
	if err != nil {
		return HTTPStatusError{http.StatusUnauthorized, err}
	}

	token, err := GenerateToken(userID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	json.NewEncoder(w).Encode(loginResponse{Token: token})

	return nil
}

type vpsInfoRequest struct {
}
type vpsInfoResponse struct {
	AllVPS []VPSInfo `json:"all_vps"`
}

func vpsInfoHandler(w http.ResponseWriter, r *http.Request) error {

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	allUserVPS, err := db.VPSGetUserAll(conn, userID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	var allUserVPSInfo []VPSInfo
	for _, vps := range allUserVPS {
		allUserVPSInfo = append(allUserVPSInfo, VPSToVPSInfo(vps))
	}

	json.NewEncoder(w).Encode(vpsInfoResponse{AllVPS: allUserVPSInfo})

	return nil
}

// or just use the VPSConfig struct directly
type vpsCreateRequest struct {
	DisplayName string `json:"display_name"           validate:"required,max=128"`
	Hostname    string `json:"hostname"               validate:"required,max=128"`
	Username    string `json:"username"               validate:"required,max=32"`
	Password    string `json:"password"               validate:"required"`
	SSHKey      string `json:"ssh_key"`
	RAM         int    `json:"ram"                    validate:"required,min=1,max=16"`
	CPU         int    `json:"cpu"                    validate:"required,min=1,max=8"`
	Disk        int    `json:"disk"                   validate:"required,min=5,max=50"`
	OS          string `json:"os"                     validate:"required"`
	Message     string `json:"message"`
}

func vpsCreateHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsCreateRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// TODO parse values for ram, cpu and disk
	config := db.VPSCreateRequestData{
		DisplayName: parsedBody.DisplayName,
		Hostname:    parsedBody.Hostname,
		UserID:      userID,
		Username:    parsedBody.Username,
		Password:    parsedBody.Password,
		SSHKey:      parsedBody.SSHKey,
		RAM:         parsedBody.RAM * 1024,
		CPU:         parsedBody.CPU,
		Disk:        parsedBody.Disk,
		OS:          parsedBody.OS,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	newRequest := db.Request{
		UserID:         userID,
		RequestTime:    time.Now(),
		RequestPurpose: db.REQUEST_PURPOSE_VPS_CREATE,
		RequestData:    string(configJSON),
		Message:        parsedBody.Message,
	}
	err = db.RequestCreate(conn, newRequest)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsDeleteRequest struct {
	VPSID uint `json:"vps_id" validate:"required"`
}

func vpsDeleteHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsDeleteRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := db.UserOwnsVPS(conn, userID, parsedBody.VPSID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
		return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := db.VPSGet(conn, parsedBody.VPSID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// issue delete commands
	err = vps.Destroy(vpsInfo.InternalName)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// remove from db
	err = db.VPSDestroy(conn, parsedBody.VPSID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsActionRequest struct {
	VPSID      uint   `json:"vps_id"      validate:"required"`
	ActionType string `json:"action_type" validate:"required"` // TODO: add validation for this
}

func vpsActionHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsActionRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := db.UserOwnsVPS(conn, userID, parsedBody.VPSID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
		return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := db.VPSGet(conn, parsedBody.VPSID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	switch parsedBody.ActionType {
	case "start":
		{

			err = vps.Start(vpsInfo.InternalName)
			if err != nil {
				return HTTPStatusError{http.StatusInternalServerError, err}
			}

		}
	case "stop":
		{

			err = vps.Stop(vpsInfo.InternalName)
			if err != nil {
				return HTTPStatusError{http.StatusInternalServerError, err}
			}

		}
	case "snapshot":
		{

			err = vps.Snapshot(vpsInfo.InternalName)
			if err != nil {
				return HTTPStatusError{http.StatusInternalServerError, err}
			}

		}
	default:
		{
			return HTTPStatusError{http.StatusBadRequest, errors.New("invalid action types")}
		}
	}

	return nil
}

type requestListResponse struct {
	RequestList []RequestInfo `json:"request_list"`
}

func requestListHandler(w http.ResponseWriter, r *http.Request) error {

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	userRequests, err := db.RequestListUser(conn, userID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	var userRequestInfos []RequestInfo
	for _, request := range userRequests {
		userRequestInfos = append(userRequestInfos, RequestToRequestInfo(request))
	}

	json.NewEncoder(w).Encode(requestListResponse{RequestList: userRequestInfos})

	return nil
}

type requestDeleteRequest struct {
	RequestID uint `json:"uint" validate:"required"`
}

func requestDeleteHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*requestDeleteRequest)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
		return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	conn, err := db.Connection()
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// check that user owns the request first
	ownsRequest, err := db.RequestUserOwns(conn, userID, parsedBody.RequestID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsRequest {
		return HTTPStatusError{http.StatusForbidden, nil}
	}

	err = db.RequestDelete(conn, parsedBody.RequestID)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

func Routes(mux *http.ServeMux) {
	for _, routeInfo := range routeSchema {
		mux.Handle(routeInfo.route, routeInfo)
	}
}
