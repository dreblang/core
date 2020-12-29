package objects

import (
	"fmt"

	"github.com/dreblang/core/object"
	"github.com/dreblang/core/vm"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const ServerObj = "http:Server"

type Server struct {
	router *routing.Router
}

func NewServer() *Server {
	return &Server{
		router: routing.New(),
	}
}

func (obj *Server) Type() object.ObjectType { return ServerObj }
func (obj *Server) Inspect() string {
	return fmt.Sprintf("http:Server[%p]", obj)
}

func (obj *Server) String() string { return "Server" }

func (obj *Server) GetMember(name string) object.Object {
	switch name {
	case "listenAndServe":
		return &object.MemberFn{
			Obj: obj,
			Fn:  listenAndServe,
		}
	case "get":
		return &object.MemberFn{
			Obj: obj,
			Fn:  get,
		}
	case "post":
		return &object.MemberFn{
			Obj: obj,
			Fn:  post,
		}
	case "put":
		return &object.MemberFn{
			Obj: obj,
			Fn:  put,
		}
	case "delete":
		return &object.MemberFn{
			Obj: obj,
			Fn:  delete,
		}
	}
	return object.NewError("No member named [%s]", name)
}

func (obj *Server) SetMember(name string, value object.Object) object.Object {
	return object.NewError("No member named [%s]", name)
}

func (obj *Server) Equals(other object.Object) bool {
	return false
}

func listenAndServe(thisObj object.Object, args ...object.Object) object.Object {
	var err error
	this := thisObj.(*Server)

	server := &fasthttp.Server{
		Handler:      this.router.HandleRequest,
		LogAllErrors: true,
	}

	if len(args) == 1 {
		if addr, ok := args[0].(*object.String); ok {
			err = server.ListenAndServe(addr.Value)
			return object.NullObject
		}
	}

	err = server.ListenAndServe(":8000")
	if err != nil {
		return object.NewError("http:Server error: %s", err)
	}
	return object.NullObject
}

func get(thisObj object.Object, args ...object.Object) object.Object {
	this := thisObj.(*Server)
	path := args[0].(*object.String)
	handler := args[1]
	this.router.Get(path.Value, getHTTPHandler(handler.(*object.Closure)))
	return object.NullObject
}

func post(thisObj object.Object, args ...object.Object) object.Object {
	this := thisObj.(*Server)
	path := args[0].(*object.String)
	handler := args[1]
	this.router.Post(path.Value, getHTTPHandler(handler.(*object.Closure)))
	return object.NullObject
}

func put(thisObj object.Object, args ...object.Object) object.Object {
	this := thisObj.(*Server)
	path := args[0].(*object.String)
	handler := args[1]
	this.router.Put(path.Value, getHTTPHandler(handler.(*object.Closure)))
	return object.NullObject
}

func delete(thisObj object.Object, args ...object.Object) object.Object {
	this := thisObj.(*Server)
	path := args[0].(*object.String)
	handler := args[1]
	this.router.Delete(path.Value, getHTTPHandler(handler.(*object.Closure)))
	return object.NullObject
}

func getHTTPHandler(handler *object.Closure) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		currentVM := vm.GetCurrentVM()
		res := currentVM.ExecClosure(handler, NewRequest(ctx), object.False)

		switch resp := res.(type) {
		case *object.String:
			_, err := ctx.WriteString(resp.Value)
			return err
		}
		return nil
	}
}
