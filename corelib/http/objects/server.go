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
	}
	return object.NewError("No member named [%s]", name)
}

func listenAndServe(thisObj object.Object, args ...object.Object) object.Object {
	var err error
	this := thisObj.(*Server)
	if len(args) == 1 {
		if addr, ok := args[0].(*object.String); ok {
			err = fasthttp.ListenAndServe(addr.Value, this.router.HandleRequest)
		}
	}
	err = fasthttp.ListenAndServe(":8000", this.router.HandleRequest)
	if err != nil {
		return object.NewError("http:Server error: %s", err)
	}
	return object.NullObject
}

func get(thisObj object.Object, args ...object.Object) object.Object {
	this := thisObj.(*Server)
	path := args[0].(*object.String)
	handler := args[1]
	currentVM := vm.GetCurrentVM()
	this.router.Get(path.Value, func(ctx *routing.Context) error {
		res := currentVM.ExecClosure(handler.(*object.Closure), object.True, object.False)

		switch resp := res.(type) {
		case *object.String:
			_, err := ctx.WriteString(resp.Value)
			return err
		}
		return nil
	})
	return object.NullObject
}
