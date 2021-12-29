package main

import (
	"time"
	"errors"
    "net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
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
	route        string
	methods      []string // possibly restrict to certain strings (ie POST, GET)
	authRoute    bool
	bodySchema   interface{}
	handlerFn    CustomHandler
}

func (info routeInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var handlerWithMiddleware http.Handler = ErrorMiddleware(info.handlerFn)

    // validate + parse body (if applicable)
    if info.bodySchema != nil {
		handlerWithMiddleware = ValidationMiddleware(handlerWithMiddleware)
        handlerWithMiddleware = ParseBodyJSONMiddleware(info.bodySchema, handlerWithMiddleware)
    }

	// restrict auth (if applicable)
	if info.authRoute {
		handlerWithMiddleware = AuthMiddleware(handlerWithMiddleware)
	}

	// restrict method of request
	handlerWithMiddleware = RestrictMethodMiddleware(info.methods, handlerWithMiddleware)

	// delegate to handler
	handlerWithMiddleware.ServeHTTP(w, r)

}

var routeSchema = []routeInfo{
	{
		route: "/ping",
		methods: []string{"POST"},
		authRoute: false,
		handlerFn: pingHandler,
	},
	{
		route: "/register",
		methods: []string{"POST"},
		authRoute: false,
		bodySchema: &registerRequest{},
		handlerFn: registerHandler,
	},
	{
		route: "/login",
		methods: []string{"POST"},
		authRoute: false,
		bodySchema: &loginRequest{},
		handlerFn: loginHandler,
	},
	{
		route: "/vps/info",
		methods: []string{"GET"},
		authRoute: true,
		bodySchema: &vpsInfoRequest{},
		handlerFn: vpsInfoHandler,
	},
	{
		route: "/vps/create",
		methods: []string{"POST"},
		authRoute: true,
		bodySchema: &vpsCreateRequest{},
		handlerFn: vpsCreateHandler,
	},
	{
		route: "/vps/delete",
		methods: []string{"POST"},
		authRoute: true,
		bodySchema: &vpsDeleteRequest{},
		handlerFn: vpsDeleteHandler,
	},
	{
		route: "/vps/start",
		methods: []string{"POST"},
		authRoute: true,
		bodySchema: &vpsStartRequest{},
		handlerFn: vpsStartHandler,
	},
	{
		route: "/vps/stop",
		methods: []string{"POST"},
		authRoute: true,
		bodySchema: &vpsStopRequest{},
		handlerFn: vpsStopHandler,
	},
	{
		route: "/vps/snapshot",
		methods: []string{"POST"},
		authRoute: true,
		bodySchema: &vpsSnapshotRequest{},
		handlerFn: vpsSnapshotHandler,
	},
}

// ping endpoint for debug purposes
func pingHandler(w http.ResponseWriter, r *http.Request) error {

	/*
	var newVPS = VPSConfig{
		DisplayName: "pino-vps",
		Hostname:    "pino-vps",
		Username:    "pinosaur",
		Password:    "$6$rounds=4096$Z7a9LgphTzzWHJbQ$Yp8C0xPXMJhE45/Q7JLo/OoAWODjlCDGH/Zdgb7FUaX5HeGdnYH4XXP13bWZldzDlSndSKSmDWTbot88ZRuJJ1",
		SSHKey:      "ssh-rsa blah blah",
		RAM:         RAM_HIGH,
		CPU:         1,
		Disk:        25,
		OS:          "ubuntu",
	}

	err := VPSCreate(newVPS)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	err = DBMigrate(db)
	if err != nil {
		return HTTPStatusError{http.StatusInternalServerError, err}
	}
	*/

	return nil
}

type registerRequest struct {
	Email         string         `json:"email"    validate:"required,email"`
    Password      string         `json:"password" validate:"required,min=8,max=128"`
}
type registerResponse struct {
    Token         string         `json:"token"`
}
func registerHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*registerRequest)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	// (possibly have db connection be part of the context)
	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// make sure user name and email are not already taken
	registered, err := DBUserCheckRegistered(db, parsedBody.Email)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if registered {
		// possibly return more generic error to reduce info leak
        return HTTPStatusError{http.StatusConflict, errors.New("username or email already taken")}
	}

	// hash pass
	hashed, err := bcrypt.GenerateFromPassword([]byte(parsedBody.Password),  bcrypt.MinCost)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	newUser := User{
		Email:    parsedBody.Email,
		Password: string(hashed),
	}
    userID, err := DBUserRegister(db, &newUser)
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
	Email         string         `json:"email"    validate:"required,email"`
	Password      string         `json:"password" validate:"required"`
}
type loginResponse struct {
    Token         string         `json:"token"`
}
func loginHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*loginRequest)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	userID, err := DBUserCheckCreds(db, parsedBody.Email, parsedBody.Password)
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
	VPSName      string          `json:"vps_name" validate:"required"`
}
type vpsInfoResponse struct {
	AllVPS       []VPSInfo       `json:"all_vps"`
}
func vpsInfoHandler(w http.ResponseWriter, r *http.Request) error {

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	allUserVPS, err := DBVPSGetUserAll(db, userID)
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
	DisplayName   string         `json:"display_name"           validate:"required,max=128"`
	Hostname      string         `json:"hostname"               validate:"required,max=128"`
	Username      string         `json:"username"               validate:"required,max=32"`
	Password      string         `json:"password"               validate:"required"`
	SSHKey        string         `json:"ssh_key"`
	RAM           int            `json:"ram"                    validate:"required,min=1,max=16"`
	CPU           int            `json:"cpu"                    validate:"required,min=1,max=8"`
	Disk          int            `json:"disk"                   validate:"required,min=5,max=50"`
	OS            string         `json:"os"                     validate:"required"`
	Message       string         `json:"message"`
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

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// TODO parse values for ram, cpu and disk
	config := VPSCreateRequestData{
		DisplayName: parsedBody.DisplayName,
		Hostname:    parsedBody.Hostname,
		UserID:      userID,
		Username:    parsedBody.Username,
		Password:    parsedBody.Password,
		SSHKey:      parsedBody.SSHKey,
		RAM:         parsedBody.RAM*1024,
		CPU:         parsedBody.CPU,
		Disk:        parsedBody.Disk,
		OS:          parsedBody.OS,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	newRequest := Request{
		UserID:         userID,
		RequestTime:    time.Now(),
		RequestPurpose: REQUEST_PURPOSE_VPS_CREATE,
		RequestData:    string(configJSON),
		Message:        parsedBody.Message,
	}
	err = DBRequestCreate(db, newRequest)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsDeleteRequest struct {
	VPSID    uint   `json:"vps_id" validate:"required"`
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

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := DBUserOwnsVPS(db, userID, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
        return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := DBVPSGet(db, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// issue delete commands
	err = VPSDestroy(vpsInfo.InternalName)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// remove from db
	err = DBVPSDestroy(db, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsStartRequest struct {
	VPSID    uint   `json:"vps_id" validate:"required"`
}
func vpsStartHandler(w http.ResponseWriter, r *http.Request) error {

	// TODO: this is a LOT of duplicated code, make this better
	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsStartRequest)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := DBUserOwnsVPS(db, userID, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
        return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := DBVPSGet(db, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	err = VPSStart(vpsInfo.InternalName)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsStopRequest struct {
	VPSID    uint   `json:"vps_id" validate:"required"`
}
func vpsStopHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsStopRequest)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := DBUserOwnsVPS(db, userID, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
        return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := DBVPSGet(db, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	err = VPSStop(vpsInfo.InternalName)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type vpsSnapshotRequest struct {
	VPSID    uint   `json:"vps_id" validate:"required"`
}
func vpsSnapshotHandler(w http.ResponseWriter, r *http.Request) error {

	parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(*vpsStopRequest)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	ownsVPS, err := DBUserOwnsVPS(db, userID, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsVPS {
        return HTTPStatusError{http.StatusForbidden, errors.New("no permission to access vps")}
	}

	vpsInfo, err := DBVPSGet(db, parsedBody.VPSID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	err = VPSSnapshot(vpsInfo.InternalName)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	return nil
}

type requestListResponse struct {
	RequestList      []RequestInfo      `json:"request_list"`
}
func requestListHandler(w http.ResponseWriter, r *http.Request) error {

	userID, ok := r.Context().Value(ContextKeyUserID).(uint)
	if !ok {
        return HTTPStatusError{http.StatusInternalServerError, nil}
	}

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	userRequests, err := DBRequestListUser(db, userID)
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
	RequestID     uint     `json:"uint" validate:"required"`
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

	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}

	// check that user owns the request first
	ownsRequest, err := DBRequestUserOwns(db, userID, parsedBody.RequestID)
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	if !ownsRequest {
        return HTTPStatusError{http.StatusForbidden, nil}
	}

	err = DBRequestDelete(db, parsedBody.RequestID)
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

