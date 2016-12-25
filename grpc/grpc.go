package grpc

import (
	"path"
	"strconv"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	qscgrpcPkgPath = "code.aliyun.com/qschou/go_common/grpc/internal/grpc"
)

func init() {
	generator.RegisterPlugin(new(grpc))
}

// grpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gRPC support.
type grpc struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "grpcx".
func (g *grpc) Name() string {
	return "grpcx"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	qscgrpcPkg string
)

// Init initializes the plugin.
func (g *grpc) Init(gen *generator.Generator) {
	g.gen = gen
	qscgrpcPkg = generator.RegisterUniquePackageName("qscgrpc", nil)
}

// P forwards to g.gen.P.
func (g *grpc) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *grpc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	if len(file.FileDescriptorProto.Service) != 1 {
		panic("plugin grpcx only supports one service proto")
	}

	g.P("/****************************************  SDK BEGIN ****************************************/")
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = (*sync.WaitGroup)(nil)")
	g.P("var _ = atomic.LoadPointer")
	g.P("var _ = unsafe.Sizeof(0)")
	g.P()

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}

	g.P()
	g.P("/****************************************  SDK END ****************************************/")
}

// GenerateImports generates the import declaration for this file.
func (g *grpc) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	if len(file.FileDescriptorProto.Service) != 1 {
		panic("plugin grpcx only supports one service proto")
	}

	g.P("import (")
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, "sync")))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, "sync/atomic")))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, "unsafe")))
	g.P()
	g.P(qscgrpcPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, qscgrpcPkgPath)))
	g.P(")")
	g.P()
}

// generateService generates all the code for the named service.
func (g *grpc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	origServName := service.GetName()
	servName := generator.CamelCase(origServName)

	g.P()
	g.P("const ", servName, `ServiceName = "TODO: replace this with your service name"`)
	g.P()
	g.P("var (")
	g.P("__", servName, "ClientPtrMutex sync.Mutex")
	g.P("__", servName, "ClientPtr unsafe.Pointer")
	g.P(")")
	g.P()
	g.P("func MustClient() ", servName, "Client {")
	g.P("clt, err := Client()")
	g.P("if err != nil {")
	g.P("panic(err)")
	g.P("}")
	g.P("return clt")
	g.P("}")
	g.P()
	g.P("func Client() (", servName, "Client, error) {")
	g.P("p := (*", servName, "Client)(atomic.LoadPointer(&__", servName, "ClientPtr))")
	g.P("if p != nil {")
	g.P("return *p, nil")
	g.P("}")
	g.P()
	g.P("__", servName, "ClientPtrMutex.Lock()")
	g.P("defer __", servName, "ClientPtrMutex.Unlock()")
	g.P()
	g.P("p = (*", servName, "Client)(atomic.LoadPointer(&__", servName, "ClientPtr))")
	g.P("if p != nil {")
	g.P("return *p, nil")
	g.P("}")
	g.P()
	g.P("clt, err := newClient()")
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P("atomic.StorePointer(&__", servName, "ClientPtr, unsafe.Pointer(&clt))")
	g.P("return clt, nil")
	g.P("}")
	g.P()
	g.P("func newClient() (", servName, "Client, error) {")
	g.P("conn, err := qscgrpc.ClientConn(", servName, "ServiceName)")
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P("return New", servName, "Client(conn), nil")
	g.P("}")
	g.P()
	g.P("func Start(port int, srv ", servName, "Server) error {")
	g.P("return NewServer(port, srv).Serve()")
	g.P("}")
	g.P()
	g.P("type Server struct {")
	g.P("*qscgrpc.Server")
	g.P("}")
	g.P("func NewServer(port int, srv ", servName, "Server) Server {")
	g.P("register := func(s *grpc.Server) {")
	g.P("Register", servName, "Server(s, srv)")
	g.P("}")
	g.P("return Server{")
	g.P("Server: qscgrpc.NewServer(", servName, `ServiceName, "", port, register),`)
	g.P("}")
	g.P("}")
}
