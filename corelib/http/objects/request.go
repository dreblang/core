package objects

import (
	"fmt"

	"github.com/dreblang/core/object"
	routing "github.com/qiangxue/fasthttp-routing"
)

const RequestObj = "http:Request"

type Request struct {
	ctx *routing.Context
}

func NewRequest(ctx *routing.Context) *Request {
	return &Request{ctx: ctx}
}

func (obj *Request) Type() object.ObjectType { return RequestObj }
func (obj *Request) Inspect() string {
	return fmt.Sprintf("http:Request[%p]", obj)
}

func (obj *Request) String() string { return "Request" }

func (obj *Request) GetMember(name string) object.Object {
	// switch name {
	// case "listenAndServe":
	// 	return &object.MemberFn{
	// 		Obj: obj,
	// 		Fn:  listenAndServe,
	// 	}
	// case "get":
	// 	return &object.MemberFn{
	// 		Obj: obj,
	// 		Fn:  get,
	// 	}
	// case "post":
	// 	return &object.MemberFn{
	// 		Obj: obj,
	// 		Fn:  post,
	// 	}
	// case "put":
	// 	return &object.MemberFn{
	// 		Obj: obj,
	// 		Fn:  put,
	// 	}
	// case "delete":
	// 	return &object.MemberFn{
	// 		Obj: obj,
	// 		Fn:  delete,
	// 	}
	// }
	return object.NewError("No member named [%s]", name)
}

func (obj *Request) SetMember(name string, value object.Object) object.Object {
	return object.NewError("No member named [%s]", name)
}

func (obj *Request) Equals(other object.Object) bool {
	return false
}
