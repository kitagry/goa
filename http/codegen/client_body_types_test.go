package codegen

import (
	"bytes"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/testdata"
)

func TestBodyTypeDecl(t *testing.T) {
	const genpkg = "gen"

	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"body-user-inner", testdata.PayloadBodyUserInnerDSL, BodyUserInnerDeclCode},
		{"body-path-user-validate", testdata.PayloadBodyPathUserValidateDSL, BodyPathUserValidateDeclCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL)
			fs := clientType(genpkg, expr.Root.API.HTTP.Services[0], make(map[string]struct{}))
			section := fs.SectionTemplates[1]
			code := codegen.SectionCode(t, section)
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
			}
		})
	}
}

func TestBodyTypeInit(t *testing.T) {
	const genpkg = "gen"
	cases := []struct {
		Name         string
		DSL          func()
		SectionIndex int
		Code         string
	}{
		{"body-user-inner", testdata.PayloadBodyUserInnerDSL, 3, BodyUserInnerInitCode},
		{"body-path-user-validate", testdata.PayloadBodyPathUserValidateDSL, 2, BodyPathUserValidateInitCode},
		{"body-primitive-array-user-validate", testdata.PayloadBodyPrimitiveArrayUserValidateDSL, 2, BodyPrimitiveArrayUserValidateInitCode},
		{"result-body-user", testdata.ResultBodyObjectHeaderDSL, 2, ResultBodyObjectHeaderInitCode},
		{"result-explicit-body-primitive", testdata.ExplicitBodyPrimitiveResultMultipleViewsDSL, 1, ExplicitBodyPrimitiveResultMultipleViewsInitCode},
		{"result-explicit-body-user-type", testdata.ExplicitBodyUserResultMultipleViewsDSL, 3, ExplicitBodyUserResultMultipleViewsInitCode},
		{"result-explicit-body-object", testdata.ExplicitBodyUserResultObjectDSL, 3, ExplicitBodyObjectInitCode},
		{"result-explicit-body-object-views", testdata.ExplicitBodyUserResultObjectMultipleViewDSL, 3, ExplicitBodyObjectViewsInitCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL)
			fs := clientType(genpkg, expr.Root.API.HTTP.Services[0], make(map[string]struct{}))
			section := fs.SectionTemplates[c.SectionIndex]
			code := codegen.SectionCode(t, section)
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
			}
		})
	}
}

func TestClientTypes(t *testing.T) {
	const genpkg = "gen"
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"mixed-payload-attrs", testdata.MixedPayloadInBodyDSL, MixedPayloadInBodyClientTypesFile},
		{"multiple-methods", testdata.MultipleMethodsDSL, MultipleMethodsClientTypesFile},
		{"payload-extend-validate", testdata.PayloadExtendedValidateDSL, PayloadExtendedValidateClientTypesFile},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL)
			fs := clientType(genpkg, expr.Root.API.HTTP.Services[0], make(map[string]struct{}))
			var buf bytes.Buffer
			for _, s := range fs.SectionTemplates[1:] {
				if err := s.Write(&buf); err != nil {
					t.Fatal(err)
				}
			}
			code := codegen.FormatTestCode(t, "package foo\n"+buf.String())
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
			}
		})
	}
}

func TestClientTypeFiles(t *testing.T) {
	const genpkg = "gen"
	cases := []struct {
		Name  string
		DSL   func()
		Codes []string
	}{
		{"multiple-services-same-payload-and-result", testdata.MultipleServicesSamePayloadAndResultDSL, MultipleServicesSamePayloadAndResultClientTypesFiles},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL)
			fw := ClientTypeFiles(genpkg, expr.Root)
			for i, fs := range fw {
				var buf bytes.Buffer
				for _, s := range fs.SectionTemplates[1:] {
					if err := s.Write(&buf); err != nil {
						t.Fatal(err)
					}
				}
				code := codegen.FormatTestCode(t, "package foo\n"+buf.String())
				if code != c.Codes[i] {
					t.Errorf("invalid code at inded %d, got:\n%s\ngot vs. expected:\n%s", i, code, codegen.Diff(t, code, c.Codes[i]))
				}
			}
		})
	}
}

const BodyUserInnerDeclCode = `// MethodBodyUserInnerRequestBody is the type of the "ServiceBodyUserInner"
// service "MethodBodyUserInner" endpoint HTTP request body.
type MethodBodyUserInnerRequestBody struct {
	Inner *InnerTypeRequestBody ` + "`" + `form:"inner,omitempty" json:"inner,omitempty" xml:"inner,omitempty"` + "`" + `
}
`

const BodyPathUserValidateDeclCode = `// MethodUserBodyPathValidateRequestBody is the type of the
// "ServiceBodyPathUserValidate" service "MethodUserBodyPathValidate" endpoint
// HTTP request body.
type MethodUserBodyPathValidateRequestBody struct {
	A string ` + "`" + `form:"a" json:"a" xml:"a"` + "`" + `
}
`

const BodyPrimitiveArrayUserValidateInitCode = `// NewPayloadTypeRequestBody builds the HTTP request body from the payload of
// the "MethodBodyPrimitiveArrayUserValidate" endpoint of the
// "ServiceBodyPrimitiveArrayUserValidate" service.
func NewPayloadTypeRequestBody(p []*servicebodyprimitivearrayuservalidate.PayloadType) []*PayloadTypeRequestBody {
	body := make([]*PayloadTypeRequestBody, len(p))
	for i, val := range p {
		body[i] = &PayloadTypeRequestBody{
			A: val.A,
		}
	}
	return body
}
`

const BodyUserInnerInitCode = `// NewMethodBodyUserInnerRequestBody builds the HTTP request body from the
// payload of the "MethodBodyUserInner" endpoint of the "ServiceBodyUserInner"
// service.
func NewMethodBodyUserInnerRequestBody(p *servicebodyuserinner.PayloadType) *MethodBodyUserInnerRequestBody {
	body := &MethodBodyUserInnerRequestBody{}
	if p.Inner != nil {
		body.Inner = marshalServicebodyuserinnerInnerTypeToInnerTypeRequestBody(p.Inner)
	}
	return body
}
`

const BodyPathUserValidateInitCode = `// NewMethodUserBodyPathValidateRequestBody builds the HTTP request body from
// the payload of the "MethodUserBodyPathValidate" endpoint of the
// "ServiceBodyPathUserValidate" service.
func NewMethodUserBodyPathValidateRequestBody(p *servicebodypathuservalidate.PayloadType) *MethodUserBodyPathValidateRequestBody {
	body := &MethodUserBodyPathValidateRequestBody{
		A: p.A,
	}
	return body
}
`

const ResultBodyObjectHeaderInitCode = `// NewMethodBodyObjectHeaderResultOK builds a "ServiceBodyObjectHeader" service
// "MethodBodyObjectHeader" endpoint result from a HTTP "OK" response.
func NewMethodBodyObjectHeaderResultOK(body *MethodBodyObjectHeaderResponseBody, b *string) *servicebodyobjectheader.MethodBodyObjectHeaderResult {
	v := &servicebodyobjectheader.MethodBodyObjectHeaderResult{
		A: body.A,
	}
	v.B = b
	return v
}
`

const ExplicitBodyPrimitiveResultMultipleViewsInitCode = `// NewMethodExplicitBodyPrimitiveResultMultipleViewResulttypemultipleviewsOK
// builds a "ServiceExplicitBodyPrimitiveResultMultipleView" service
// "MethodExplicitBodyPrimitiveResultMultipleView" endpoint result from a HTTP
// "OK" response.
func NewMethodExplicitBodyPrimitiveResultMultipleViewResulttypemultipleviewsOK(body string, c *string) *serviceexplicitbodyprimitiveresultmultipleviewviews.ResulttypemultipleviewsView {
	v := body
	res := &serviceexplicitbodyprimitiveresultmultipleviewviews.ResulttypemultipleviewsView{
		A: &v,
	}
	res.C = c
	return res
}
`

const ExplicitBodyUserResultMultipleViewsInitCode = `// NewMethodExplicitBodyUserResultMultipleViewResulttypemultipleviewsOK builds
// a "ServiceExplicitBodyUserResultMultipleView" service
// "MethodExplicitBodyUserResultMultipleView" endpoint result from a HTTP "OK"
// response.
func NewMethodExplicitBodyUserResultMultipleViewResulttypemultipleviewsOK(body *MethodExplicitBodyUserResultMultipleViewResponseBody, c *string) *serviceexplicitbodyuserresultmultipleviewviews.ResulttypemultipleviewsView {
	v := &serviceexplicitbodyuserresultmultipleviewviews.UserTypeView{
		X: body.X,
		Y: body.Y,
	}
	res := &serviceexplicitbodyuserresultmultipleviewviews.ResulttypemultipleviewsView{
		A: v,
	}
	res.C = c
	return res
}
`

const ExplicitBodyObjectInitCode = `// NewMethodExplicitBodyUserResultObjectResulttypeOK builds a
// "ServiceExplicitBodyUserResultObject" service
// "MethodExplicitBodyUserResultObject" endpoint result from a HTTP "OK"
// response.
func NewMethodExplicitBodyUserResultObjectResulttypeOK(body *MethodExplicitBodyUserResultObjectResponseBody, c *string, b *string) *serviceexplicitbodyuserresultobjectviews.ResulttypeView {
	v := &serviceexplicitbodyuserresultobjectviews.ResulttypeView{}
	if body.A != nil {
		v.A = unmarshalUserTypeResponseBodyToServiceexplicitbodyuserresultobjectviewsUserTypeView(body.A)
	}
	v.C = c
	v.B = b
	return v
}
`

const ExplicitBodyObjectViewsInitCode = `// NewMethodExplicitBodyUserResultObjectMultipleViewResulttypemultipleviewsOK
// builds a "ServiceExplicitBodyUserResultObjectMultipleView" service
// "MethodExplicitBodyUserResultObjectMultipleView" endpoint result from a HTTP
// "OK" response.
func NewMethodExplicitBodyUserResultObjectMultipleViewResulttypemultipleviewsOK(body *MethodExplicitBodyUserResultObjectMultipleViewResponseBody, c *string) *serviceexplicitbodyuserresultobjectmultipleviewviews.ResulttypemultipleviewsView {
	v := &serviceexplicitbodyuserresultobjectmultipleviewviews.ResulttypemultipleviewsView{}
	if body.A != nil {
		v.A = unmarshalUserTypeResponseBodyToServiceexplicitbodyuserresultobjectmultipleviewviewsUserTypeView(body.A)
	}
	v.C = c
	return v
}
`
const MixedPayloadInBodyClientTypesFile = `// MethodARequestBody is the type of the "ServiceMixedPayloadInBody" service
// "MethodA" endpoint HTTP request body.
type MethodARequestBody struct {
	Any    interface{}          ` + "`" + `form:"any,omitempty" json:"any,omitempty" xml:"any,omitempty"` + "`" + `
	Array  []float32            ` + "`" + `form:"array" json:"array" xml:"array"` + "`" + `
	Map    map[uint]interface{} ` + "`" + `form:"map,omitempty" json:"map,omitempty" xml:"map,omitempty"` + "`" + `
	Object *BPayloadRequestBody ` + "`" + `form:"object" json:"object" xml:"object"` + "`" + `
	DupObj *BPayloadRequestBody ` + "`" + `form:"dup_obj,omitempty" json:"dup_obj,omitempty" xml:"dup_obj,omitempty"` + "`" + `
}

// BPayloadRequestBody is used to define fields on request body types.
type BPayloadRequestBody struct {
	Int   int    ` + "`" + `form:"int" json:"int" xml:"int"` + "`" + `
	Bytes []byte ` + "`" + `form:"bytes,omitempty" json:"bytes,omitempty" xml:"bytes,omitempty"` + "`" + `
}

// NewMethodARequestBody builds the HTTP request body from the payload of the
// "MethodA" endpoint of the "ServiceMixedPayloadInBody" service.
func NewMethodARequestBody(p *servicemixedpayloadinbody.APayload) *MethodARequestBody {
	body := &MethodARequestBody{
		Any: p.Any,
	}
	if p.Array != nil {
		body.Array = make([]float32, len(p.Array))
		for i, val := range p.Array {
			body.Array[i] = val
		}
	}
	if p.Map != nil {
		body.Map = make(map[uint]interface{}, len(p.Map))
		for key, val := range p.Map {
			tk := key
			tv := val
			body.Map[tk] = tv
		}
	}
	if p.Object != nil {
		body.Object = marshalServicemixedpayloadinbodyBPayloadToBPayloadRequestBody(p.Object)
	}
	if p.DupObj != nil {
		body.DupObj = marshalServicemixedpayloadinbodyBPayloadToBPayloadRequestBody(p.DupObj)
	}
	return body
}
`

const MultipleMethodsClientTypesFile = `// MethodARequestBody is the type of the "ServiceMultipleMethods" service
// "MethodA" endpoint HTTP request body.
type MethodARequestBody struct {
	A *string ` + "`" + `form:"a,omitempty" json:"a,omitempty" xml:"a,omitempty"` + "`" + `
}

// MethodBRequestBody is the type of the "ServiceMultipleMethods" service
// "MethodB" endpoint HTTP request body.
type MethodBRequestBody struct {
	A string               ` + "`" + `form:"a" json:"a" xml:"a"` + "`" + `
	B *string              ` + "`" + `form:"b,omitempty" json:"b,omitempty" xml:"b,omitempty"` + "`" + `
	C *APayloadRequestBody ` + "`" + `form:"c" json:"c" xml:"c"` + "`" + `
}

// APayloadRequestBody is used to define fields on request body types.
type APayloadRequestBody struct {
	A *string ` + "`" + `form:"a,omitempty" json:"a,omitempty" xml:"a,omitempty"` + "`" + `
}

// NewMethodARequestBody builds the HTTP request body from the payload of the
// "MethodA" endpoint of the "ServiceMultipleMethods" service.
func NewMethodARequestBody(p *servicemultiplemethods.APayload) *MethodARequestBody {
	body := &MethodARequestBody{
		A: p.A,
	}
	return body
}

// NewMethodBRequestBody builds the HTTP request body from the payload of the
// "MethodB" endpoint of the "ServiceMultipleMethods" service.
func NewMethodBRequestBody(p *servicemultiplemethods.PayloadType) *MethodBRequestBody {
	body := &MethodBRequestBody{
		A: p.A,
		B: p.B,
	}
	if p.C != nil {
		body.C = marshalServicemultiplemethodsAPayloadToAPayloadRequestBody(p.C)
	}
	return body
}

// ValidateAPayloadRequestBody runs the validations defined on
// APayloadRequestBody
func ValidateAPayloadRequestBody(body *APayloadRequestBody) (err error) {
	if body.A != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.a", *body.A, "patterna"))
	}
	return
}
`

const PayloadExtendedValidateClientTypesFile = `// MethodQueryStringExtendedValidatePayloadRequestBody is the type of the
// "ServiceQueryStringExtendedValidatePayload" service
// "MethodQueryStringExtendedValidatePayload" endpoint HTTP request body.
type MethodQueryStringExtendedValidatePayloadRequestBody struct {
	Body string ` + "`" + `form:"body" json:"body" xml:"body"` + "`" + `
}

// NewMethodQueryStringExtendedValidatePayloadRequestBody builds the HTTP
// request body from the payload of the
// "MethodQueryStringExtendedValidatePayload" endpoint of the
// "ServiceQueryStringExtendedValidatePayload" service.
func NewMethodQueryStringExtendedValidatePayloadRequestBody(p *servicequerystringextendedvalidatepayload.MethodQueryStringExtendedValidatePayloadPayload) *MethodQueryStringExtendedValidatePayloadRequestBody {
	body := &MethodQueryStringExtendedValidatePayloadRequestBody{
		Body: p.Body,
	}
	return body
}
`

var MultipleServicesSamePayloadAndResultClientTypesFiles = []string{
	`// ListRequestBody is the type of the "ServiceA" service "list" endpoint HTTP
// request body.
type ListRequestBody struct {
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListStreamingBody is the type of the "ServiceA" service "list" endpoint HTTP
// request body.
type ListStreamingBody struct {
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListResponseBody is the type of the "ServiceA" service "list" endpoint HTTP
// response body.
type ListResponseBody struct {
	ID   *int    ` + "`" + `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"` + "`" + `
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListSomethingWentWrongResponseBody is the type of the "ServiceA" service
// "list" endpoint HTTP response body for the "something_went_wrong" error.
type ListSomethingWentWrongResponseBody struct {
	// Name is the name of this class of errors.
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string ` + "`" + `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"` + "`" + `
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string ` + "`" + `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"` + "`" + `
	// Is the error temporary?
	Temporary *bool ` + "`" + `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"` + "`" + `
	// Is the error a timeout?
	Timeout *bool ` + "`" + `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"` + "`" + `
	// Is the error a server-side fault?
	Fault *bool ` + "`" + `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"` + "`" + `
}

// NewListRequestBody builds the HTTP request body from the payload of the
// "list" endpoint of the "ServiceA" service.
func NewListRequestBody(p *servicea.ListPayload) *ListRequestBody {
	body := &ListRequestBody{
		Name: p.Name,
	}
	return body
}

// NewListStreamingBody builds the HTTP request body from the payload of the
// "list" endpoint of the "ServiceA" service.
func NewListStreamingBody(p *servicea.ListStreamingPayload) *ListStreamingBody {
	body := &ListStreamingBody{
		Name: p.Name,
	}
	return body
}

// NewListResultOK builds a "ServiceA" service "list" endpoint result from a
// HTTP "OK" response.
func NewListResultOK(body *ListResponseBody) *servicea.ListResult {
	v := &servicea.ListResult{
		ID:   *body.ID,
		Name: *body.Name,
	}
	return v
}

// NewListSomethingWentWrong builds a ServiceA service list endpoint
// something_went_wrong error.
func NewListSomethingWentWrong(body *ListSomethingWentWrongResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}
	return v
}

// ValidateListResponseBody runs the validations defined on ListResponseBody
func ValidateListResponseBody(body *ListResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}

// ValidateListSomethingWentWrongResponseBody runs the validations defined on
// list_something_went_wrong_response_body
func ValidateListSomethingWentWrongResponseBody(body *ListSomethingWentWrongResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}
`,
	`// ListRequestBody is the type of the "ServiceB" service "list" endpoint HTTP
// request body.
type ListRequestBody struct {
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListStreamingBody is the type of the "ServiceB" service "list" endpoint HTTP
// request body.
type ListStreamingBody struct {
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListResponseBody is the type of the "ServiceB" service "list" endpoint HTTP
// response body.
type ListResponseBody struct {
	ID   *int    ` + "`" + `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"` + "`" + `
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
}

// ListSomethingWentWrongResponseBody is the type of the "ServiceB" service
// "list" endpoint HTTP response body for the "something_went_wrong" error.
type ListSomethingWentWrongResponseBody struct {
	// Name is the name of this class of errors.
	Name *string ` + "`" + `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"` + "`" + `
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string ` + "`" + `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"` + "`" + `
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string ` + "`" + `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"` + "`" + `
	// Is the error temporary?
	Temporary *bool ` + "`" + `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"` + "`" + `
	// Is the error a timeout?
	Timeout *bool ` + "`" + `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"` + "`" + `
	// Is the error a server-side fault?
	Fault *bool ` + "`" + `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"` + "`" + `
}

// NewListRequestBody builds the HTTP request body from the payload of the
// "list" endpoint of the "ServiceB" service.
func NewListRequestBody(p *serviceb.ListPayload) *ListRequestBody {
	body := &ListRequestBody{
		Name: p.Name,
	}
	return body
}

// NewListStreamingBody builds the HTTP request body from the payload of the
// "list" endpoint of the "ServiceB" service.
func NewListStreamingBody(p *serviceb.ListStreamingPayload) *ListStreamingBody {
	body := &ListStreamingBody{
		Name: p.Name,
	}
	return body
}

// NewListResultOK builds a "ServiceB" service "list" endpoint result from a
// HTTP "OK" response.
func NewListResultOK(body *ListResponseBody) *serviceb.ListResult {
	v := &serviceb.ListResult{
		ID:   *body.ID,
		Name: *body.Name,
	}
	return v
}

// NewListSomethingWentWrong builds a ServiceB service list endpoint
// something_went_wrong error.
func NewListSomethingWentWrong(body *ListSomethingWentWrongResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}
	return v
}

// ValidateListResponseBody runs the validations defined on ListResponseBody
func ValidateListResponseBody(body *ListResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}

// ValidateListSomethingWentWrongResponseBody runs the validations defined on
// list_something_went_wrong_response_body
func ValidateListSomethingWentWrongResponseBody(body *ListSomethingWentWrongResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}
`,
}
